package server

import (
	"context"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	libviper "github.com/otamoe/go-library/viper"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	RouteSystem struct {
		mux    sync.RWMutex
		ctx    context.Context
		db     *badger.DB
		routes Routes
	}

	Routes []*Route
)

func init() {
	libviper.SetDefault("route.subnet", net.IPNet{IP: net.IPv4(0xa, 0x80, 0, 0), Mask: net.IPv4Mask(0xff, 0x80, 0, 0)}, "Subnet of the entire router")
	libviper.SetDefault("route.create", net.IPNet{IP: net.IPv4(0xa, 0xff, 0, 0), Mask: net.IPv4Mask(0xff, 0xff, 0, 0)}, "Subnet for auto-registered clients")
}

func NewRouteSystem(ctx context.Context, db *badger.DB, lc fx.Lifecycle) (routeSystem *RouteSystem, err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	routeSystem = &RouteSystem{
		db:     db,
		ctx:    ctx,
		routes: Routes{},
	}

	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		cancel()
		return nil
	}})

	if err = routeSystem.startRunLoad(); err != nil {
		return
	}
	return
}

func (routeSystem *RouteSystem) All() (routes Routes) {
	routeSystem.mux.RLock()
	defer routeSystem.mux.RUnlock()
	routes = make(Routes, len(routeSystem.routes))
	copy(routes, routeSystem.routes)
	return
}

func (routeSystem *RouteSystem) Get(id string) (route *Route) {
	routes := routeSystem.All()
	for _, v := range routes {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (routeSystem *RouteSystem) exists(id string) bool {
	for _, v := range routeSystem.routes {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (routeSystem *RouteSystem) Delete(id string) (route *Route, err error) {
	routeSystem.mux.Lock()
	defer routeSystem.mux.Unlock()
	if !routeSystem.exists(id) {
		err = ErrRouteNotFound
		return
	}

	for i := 0; i < 512; i++ {
		err = routeSystem.db.Update(func(txn *badger.Txn) (err error) {
			route = newRoute(id)
			if err = route.load(txn); err != nil {
				return
			}
			if err = route.delete(txn); err != nil {
				return
			}
			return
		})

		// 无错误 响应
		if err == nil {
			break
		}
		// badger.ErrConflict 重试错误
		if err != badger.ErrConflict {
			break
		}
	}

	// 有错误
	if err != nil {
		route = nil
		return
	}

	routes := Routes{}
	for _, v := range routeSystem.routes {
		if v.Id != id {
			routes = append(routes, v)
		}
	}
	routeSystem.routes = routes
	return
}

func (routeSystem *RouteSystem) Save(id string, cb func(route *Route) (rRoute *Route, err error)) (route *Route, err error) {
	routeSystem.mux.Lock()
	defer routeSystem.mux.Unlock()

	if id != "" {
		if !routeSystem.exists(id) {
			err = ErrRouteNotFound
		}
	}

	for i := 0; i < 512; i++ {
		err = routeSystem.db.Update(func(txn *badger.Txn) (err error) {
			if id == "" {
				now := time.Now().UTC()
				route = newRoute(NewID(now)).
					WithAction(pb.Route_REJECT).
					WithState(pb.State_AVAILABLE).
					WithCreatedAt(now).
					WithUpdatedAt(now).
					WithExpiredAt(time.Date(9000, time.January, 1, 0, 0, 0, 0, time.UTC))
			} else {
				route = newRoute(id)
				if err = route.load(txn); err != nil {
					return
				}
			}

			// 回调修改
			if route, err = cb(route); err != nil {
				return
			}

			// source 空
			if route.Route.SourceIP == "" {
				err = &ValidateError{
					Name: "source",
				}
				return
			}
			// key 空
			if route.Route.DestinationIP == "" {
				err = &ValidateError{
					Name: "destination",
				}
				return
			}

			if id != "" {
				route = route.WithUpdatedAt(time.Now().UTC())
			}

			// 保存 route
			if err = route.save(txn); err != nil {
				return
			}
			return
		})

		// 无错误 响应
		if err == nil {
			break
		}

		// badger.ErrConflict 重试错误
		if err != badger.ErrConflict {
			break
		}
	}

	// 有错误
	if err != nil {
		route = nil
		return
	}

	if id == "" {
		// 新增
		routeSystem.routes = append(routeSystem.routes, route)
	} else {
		for k, v := range routeSystem.routes {
			if v.Id == route.Id {
				routeSystem.routes[k] = route
				break
			}
		}
	}

	// 重新排序下
	sort.Sort(routeSystem.routes)
	return
}

func (routeSystem *RouteSystem) startRunLoad() (err error) {
	// 读取全部 file
	err = routeSystem.db.Update(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := routeFieldDBKey("", RouteFieldCreatedAt)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keySplit := strings.Split(string(item.Key()), "/")
			id := string(keySplit[len(keySplit)-1])
			route := newRoute(id)
			if err = route.load(txn); err != nil {
				return
			}
			routeSystem.routes = append(routeSystem.routes, route)
		}
		return
	})
	if err != nil {
		logger.Error("Load routes", zap.Error(err))
		return
	}
	return
}
