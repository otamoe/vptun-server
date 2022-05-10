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
		iSourceIP      *net.IPNet `json:"-"`
		iDestinationIP *net.IPNet `json:"-"`
	}
)

var ErrRouteNotFound = &NotFoundError{
	Name: "Route",
}

func (route Route) Clone() (rRoute *Route) {
	rRoute = &Route{}
	*rRoute = route
	rRoute.Route = &pb.Route{}

	if route.Route != nil {
		*rRoute.Route = *route.Route
	}

	if route.iSourceIP != nil {
		rRoute.iSourceIP = &net.IPNet{}
		*rRoute.iSourceIP = *route.iSourceIP
	}
	if route.iDestinationIP != nil {
		rRoute.iDestinationIP = &net.IPNet{}
		*rRoute.iDestinationIP = *route.iDestinationIP
	}
	return
}

func (route *Route) WithType(typ pb.Route_Type) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Type = typ
	return rRoute
}

func (route *Route) WithSourceIP(ipnet *net.IPNet) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.iSourceIP = ipnet
	if rRoute.iSourceIP == nil {
		rRoute.Route.SourceIP = ""
	} else {
		rRoute.Route.SourceIP = rRoute.iSourceIP.String()
	}
	return rRoute
}

func (route *Route) WithDestinationIP(ipnet *net.IPNet) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.iDestinationIP = ipnet
	if rRoute.iDestinationIP == nil {
		rRoute.Route.DestinationIP = ""
	} else {
		rRoute.Route.DestinationIP = rRoute.iDestinationIP.String()
	}
	return rRoute
}

func (route *Route) WithSourcePort(port uint32) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Route.SourcePort = port
	return rRoute
}

func (route *Route) WithDestinationPort(port uint32) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Route.DestinationPort = port
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
func (route *Route) WithLevel(level int32) (rRoute *Route) {
	rRoute = route.Clone()
	rRoute.Level = level
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
	if route.Route.SourceIP != "" {
		if _, route.iSourceIP, err = net.ParseCIDR(route.Route.SourceIP); err != nil {
			return
		}
	}
	if route.Route.DestinationIP != "" {
		if _, route.iDestinationIP, err = net.ParseCIDR(route.Route.DestinationIP); err != nil {
			return
		}
	}
	return
}

func (route *Route) load(txn *badger.Txn) (err error) {
	var item *badger.Item
	var value []byte

	// RouteFieldType
	{
		route.Route.Type = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldType)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.Type)
			}
		}
	}

	// RouteFieldSourceIP
	{
		route.Route.SourceIP = ""
		route.iSourceIP = nil
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldSourceIP)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				route.Route.SourceIP = string(value)
				if _, route.iSourceIP, err = net.ParseCIDR(route.Route.SourceIP); err != nil {
					return
				}
			}
		}
	}

	// RouteFieldDestinationIP
	{
		route.Route.DestinationIP = ""
		route.iDestinationIP = nil
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldDestinationIP)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				route.Route.DestinationIP = string(value)
				if _, route.iDestinationIP, err = net.ParseCIDR(route.Route.DestinationIP); err != nil {
					return
				}
			}
		}
	}
	// RouteFieldSourcePort
	{
		route.Route.SourcePort = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldSourcePort)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.SourcePort)
			}
		}
	}
	// RouteFieldDestinationPort
	{
		route.Route.DestinationPort = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldDestinationPort)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.DestinationPort)
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

	// RouteFieldLevel
	{
		route.Route.Level = 0
		if item, err = txn.Get(routeFieldDBKey(route.Id, RouteFieldLevel)); err != nil && err != badger.ErrKeyNotFound {
			return
		}
		if item != nil {
			if value, err = item.ValueCopy(nil); err != nil {
				return
			}
			if len(value) != 0 {
				buf := bytes.NewBuffer(value)
				binary.Read(buf, binary.BigEndian, &route.Level)
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
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.Type)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldType), buf.Bytes()); err != nil {
			return
		}
	}
	if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldSourceIP), []byte(route.Route.SourceIP)); err != nil {
		return
	}
	if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldDestinationIP), []byte(route.Route.DestinationIP)); err != nil {
		return
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.SourcePort)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldSourcePort), buf.Bytes()); err != nil {
			return
		}
	}
	{
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, route.DestinationPort)
		if err = txn.Set(routeFieldDBKey(route.Id, RouteFieldDestinationPort), buf.Bytes()); err != nil {
			return
		}
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
	if routes[i].Level < routes[j].Level {
		return true
	}

	if routes[i].Level > routes[j].Level {
		return false
	}

	var iMask []byte
	var jMask []byte
	if routes[i].iSourceIP != nil {
		iMask = routes[i].iSourceIP.Mask
	}
	if routes[j].iSourceIP != nil {
		iMask = routes[j].iSourceIP.Mask
	}
	if c := bytes.Compare(iMask, jMask); c != 0 {
		return c != 1
	}
	var iAddr []byte
	var jAddr []byte
	if routes[i].iSourceIP != nil {
		iAddr = routes[i].iSourceIP.IP
	}
	if routes[j].iSourceIP != nil {
		jAddr = routes[j].iSourceIP.IP
	}
	return bytes.Compare(iAddr, jAddr) != 1
}

func (routes Routes) Swap(i, j int) {
	routes[i], routes[j] = routes[j], routes[i]
}
