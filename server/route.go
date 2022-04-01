package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"time"

	"github.com/dgraph-io/badger/v3"
	pb "github.com/otamoe/vptun-pb"
)

type (
	Route struct {
		*pb.Route
		iSource      *net.IPNet `json:"-"`
		iDestination *net.IPNet `json:"-"`
	}
)

var ErrRouteNotFound = &NotFoundError{
	Name: "Route",
}

func (route *Route) Clone() (rRoute *Route) {
	if route == nil {
		return nil
	}
	rRoute = &Route{}
	*rRoute = *route
	rRoute.Route = &pb.Route{}

	if route.Route != nil {
		*rRoute.Route = *route.Route
	}

	if route.iSource != nil {
		rRoute.iSource = &net.IPNet{}
		*rRoute.iSource = *route.iSource
	}
	if route.iDestination != nil {
		rRoute.iDestination = &net.IPNet{}
		*rRoute.iDestination = *route.iDestination
	}
	return
}

func (route *Route) WithSource(ipnet *net.IPNet) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.iSource = ipnet
	if rRoute.iSource == nil {
		rRoute.Route.Source = ""
	} else {
		rRoute.Route.Source = rRoute.iSource.String()
	}
	return rRoute
}

func (route *Route) WithDestination(ipnet *net.IPNet) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.iDestination = ipnet
	if rRoute.iDestination == nil {
		rRoute.Route.Destination = ""
	} else {
		rRoute.Route.Destination = rRoute.iDestination.String()
	}
	return rRoute
}

func (route *Route) WithRemark(remark string) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Remark = remark
	return rRoute
}

func (route *Route) WithAction(action pb.Route_Action) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Action = action
	return rRoute
}

func (route *Route) WithState(state pb.State) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.State = state
	return rRoute
}
func (route *Route) WithCreatedAt(t time.Time) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.CreatedAt = t.Unix()
	return rRoute
}
func (route *Route) WithUpdatedAt(t time.Time) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.UpdatedAt = t.Unix()
	return rRoute
}

func (route *Route) WithExpiredAt(t time.Time) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.ExpiredAt = t.Unix()
	return rRoute
}

func (route *Route) UnmarshalJSON(data []byte) (err error) {
	route.Route = &pb.Route{}
	if err = json.Unmarshal(data, route.Route); err != nil {
		return
	}
	if route.Route.Source != "" {
		if _, route.iSource, err = net.ParseCIDR(route.Route.Source); err != nil {
			return
		}
	}
	if route.Route.Destination != "" {
		if _, route.iDestination, err = net.ParseCIDR(route.Route.Destination); err != nil {
			return
		}
	}
	return
}

func (route *Route) load(txn *badger.Txn) (err error) {
	var item *badger.Item
	var value []byte

	// RouteFieldSource
	{
		route.Route.Source = ""
		route.iSource = nil
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldSource)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				route.Route.Source = string(value)
				if _, route.iSource, err = net.ParseCIDR(route.Route.Source); err != nil {
					return
				}
			}
		}
	}

	// RouteFieldDestination
	{
		route.Route.Destination = ""
		route.iDestination = nil
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldDestination)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				route.Route.Destination = string(value)
				if _, route.iDestination, err = net.ParseCIDR(route.Route.Destination); err != nil {
					return
				}
			}
		}
	}
	// RouteFieldRemark
	{
		route.Route.Remark = ""
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldRemark)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				route.Remark = string(value)
			}
		}
	}

	// RouteFieldAction
	{
		route.Action = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldAction)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.Action)
			}
		}
	}
	// RouteFieldState
	{
		route.State = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldState)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.State)
			}
		}
	}

	// RouteFieldCreatedAt
	{
		route.CreatedAt = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldCreatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}

		// 创建时间是必须的否则未找到
		if err == badger.ErrKeyNotFound {
			err = ErrRouteNotFound
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.CreatedAt)
			}
		}
	}
	// RouteFieldUpdatedAt
	{
		route.UpdatedAt = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldUpdatedAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.UpdatedAt)
			}
		}
	}
	// RouteFieldExpiredAt
	{
		route.ExpiredAt = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldExpiredAt)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}

			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.ExpiredAt)
			}
		}
	}

	err = nil
	return
}

func (route *Route) save(txn *badger.Txn) (err error) {
	if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldSource), []byte(route.Route.Source)); err != nil {
		return
	}
	if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldDestination), []byte(route.Route.Destination)); err != nil {
		return
	}
	if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldRemark), []byte(route.Route.Remark)); err != nil {
		return
	}

	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.Action)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldAction), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.State)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldState), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.CreatedAt)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldCreatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.UpdatedAt)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldUpdatedAt), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.ExpiredAt)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldExpiredAt), buf.Bytes()); err != nil {
			return
		}
	}
	return
}

func (route *Route) delete(txn *badger.Txn) (err error) {
	for _, field := range AllRouteFields {
		if err = txn.Delete(routeFieldDBKey(route.Id, field)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
	}
	if err == badger.ErrKeyNotFound {
		err = nil
	}
	return
}

func newRoute(id string) *Route {
	return &Route{
		Route: &pb.Route{
			Id: id,
		},
	}
}

func (routes Routes) Len() int {
	return len(routes)
}
func (routes Routes) Less(i, j int) bool {
	var iMask []byte
	var jMask []byte
	if routes[i].iSource != nil {
		iMask = routes[i].iSource.Mask
	}
	if routes[j].iSource != nil {
		iMask = routes[j].iSource.Mask
	}
	if c := bytes.Compare(iMask, jMask); c != 0 {
		return c != 1
	}
	var iAddr []byte
	var jAddr []byte
	if routes[i].iSource != nil {
		iAddr = routes[i].iSource.IP
	}
	if routes[j].iSource != nil {
		jAddr = routes[j].iSource.IP
	}
	return bytes.Compare(iAddr, jAddr) != 1
}

func (routes Routes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}
