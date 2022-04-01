package server

import (
	"context"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	pb "github.com/otamoe/vptun-pb"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	ClientSystem struct {
		mux sync.RWMutex
		ctx context.Context
		db  *badger.DB

		clients             map[string]*Client
		clientsRouteAddress map[string]*Client
	}
)

func NewClientSystem(ctx context.Context, db *badger.DB, lc fx.Lifecycle) (clientSystem *ClientSystem, err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	clientSystem = &ClientSystem{
		db:                  db,
		ctx:                 ctx,
		clients:             map[string]*Client{},
		clientsRouteAddress: map[string]*Client{},
	}

	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		cancel()
		return nil
	}})

	if err = clientSystem.startRunLoad(); err != nil {
		return
	}
	return
}

func (clientSystem *ClientSystem) NewRouteAddress(lock bool) (routeAddress net.IP) {
	if lock {
		clientSystem.mux.RLock()
		defer clientSystem.mux.RUnlock()
	}

	routeCreate := viper.Get("route.create")
	if routeCreate == nil {
		return
	}
	ipNet, ok := routeCreate.(net.IPNet)
	if !ok {
		return
	}

	networkIP := make(net.IP, len(ipNet.IP))
	gatewayIP := make(net.IP, len(ipNet.IP))
	ip := make(net.IP, len(ipNet.IP))
	broadcastIP := make(net.IP, len(ipNet.IP))
	copy(networkIP, ipNet.IP)
	copy(gatewayIP, ipNet.IP)
	copy(ip, ipNet.IP)
	copy(broadcastIP, ipNet.IP)

	// 网络地址+1
	incIP(gatewayIP)

	// 最后
	for i, _ := range ip {
		if i >= len(ipNet.Mask) {
			break
		}
		broadcastIP[i] = broadcastIP[i] | ^ipNet.Mask[i]
	}

	for ; ipNet.Contains(ip); incIP(ip) {
		// 网络地址 ip
		if networkIP.Equal(ip) {
			continue
		}

		// 网关 ip
		if gatewayIP.Equal(ip) {
			continue
		}

		// 网段结束 ip
		if broadcastIP.Equal(ip) {
			continue
		}

		// 存在
		if _, ok := clientSystem.clientsRouteAddress[ip.String()]; ok {
			continue
		}
		routeAddress = ip
		break
	}
	return
}

func (clientSystem *ClientSystem) All() (clients Clients) {
	clientSystem.mux.RLock()
	clientSystem.mux.RUnlock()
	clients = make(Clients, len(clientSystem.clients))
	var i = 0
	for _, v := range clientSystem.clients {
		clients[i] = v
		i++
	}
	sort.Sort(clients)
	return
}

func (clientSystem *ClientSystem) Get(id string) (client *Client) {
	clientSystem.mux.RLock()
	clientSystem.mux.RUnlock()
	client, _ = clientSystem.clients[id]
	return
}

func (clientSystem *ClientSystem) GetByRouteAddress(routeAddress net.IP) (client *Client) {
	clientSystem.mux.RLock()
	clientSystem.mux.RUnlock()
	client, _ = clientSystem.clientsRouteAddress[routeAddress.String()]
	return nil
}

func (clientSystem *ClientSystem) Delete(id string) (client *Client, err error) {
	clientSystem.mux.Lock()
	defer clientSystem.mux.Unlock()
	if _, ok := clientSystem.clients[id]; !ok {
		err = ErrClientNotFound
		return
	}

	for i := 0; i < 512; i++ {
		err = clientSystem.db.Update(func(txn *badger.Txn) (err error) {
			client = newClient(id)
			if err = client.load(txn); err != nil {
				return
			}
			if err = client.delete(txn); err != nil {
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
		client = nil
		return
	}

	delete(clientSystem.clients, client.Id)
	delete(clientSystem.clientsRouteAddress, client.Client.RouteAddress)
	return
}

func (clientSystem *ClientSystem) Save(id string, cb func(client *Client) (rClient *Client, err error)) (client *Client, err error) {
	clientSystem.mux.Lock()
	defer clientSystem.mux.Unlock()
	if id != "" {
		if _, ok := clientSystem.clients[id]; !ok {
			err = ErrClientNotFound
			return
		}
	}

	var oldRouteAddress string
	for i := 0; i < 512; i++ {
		err = clientSystem.db.Update(func(txn *badger.Txn) (err error) {
			if id == "" {
				now := time.Now().UTC()
				client = newClient(NewID(now)).
					WithState(pb.State_AVAILABLE).
					WithCreatedAt(now).
					WithConnectAt(now).
					WithUpdatedAt(now).
					WithExpiredAt(time.Date(9000, time.January, 1, 0, 0, 0, 0, time.UTC))
			} else {
				client = newClient(id)
				if err = client.load(txn); err != nil {
					return
				}
				oldRouteAddress = client.Client.RouteAddress
			}

			// 回调修改
			if client, err = cb(client); err != nil {
				return
			}

			// key 空
			if client.Client.Key == "" {
				err = &ValidateError{
					Name: "key",
				}
				return
			}

			// routeAddress 空
			if client.Client.RouteAddress == "" {
				err = &ValidateError{
					Name: "routeAddress",
				}
				return
			}

			// routeAddress 存在
			if c, ok := clientSystem.clientsRouteAddress[client.Client.RouteAddress]; ok && client.Client.Id != c.Client.Id {
				err = &ValidateError{
					Name: "routeAddress",
				}
				return
			}

			if id != "" {
				client = client.WithUpdatedAt(time.Now().UTC())
			}

			// 保存 client
			if err = client.save(txn); err != nil {
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
		client = nil
		return
	}

	clientSystem.clients[client.Id] = client
	delete(clientSystem.clientsRouteAddress, oldRouteAddress)
	clientSystem.clientsRouteAddress[client.Client.RouteAddress] = client

	return
}

func (clientSystem *ClientSystem) startRunLoad() (err error) {
	// 读取全部 file
	err = clientSystem.db.Update(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := clientFieldDBKey("", ClientFieldCreatedAt)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keySplit := strings.Split(string(item.Key()), "/")
			id := string(keySplit[len(keySplit)-1])
			client := newClient(id)
			if err = client.load(txn); err != nil {
				return
			}
			clientSystem.clients[client.Id] = client
			if client.Client.RouteAddress != "" {
				clientSystem.clientsRouteAddress[client.Client.RouteAddress] = client
			}
		}
		return
	})
	if err != nil {
		logger.Error("Load clients", zap.Error(err))
		return
	}
	return
}
