package server

import "bytes"

type (
	RouteField string
)

var (
	RouteFieldSourceIP        RouteField = "sip"
	RouteFieldDestinationIP   RouteField = "dip"
	RouteFieldSourcePort      RouteField = "spt"
	RouteFieldDestinationPort RouteField = "dpt"
	RouteFieldRemark          RouteField = "rem"
	RouteFieldAction          RouteField = "act"
	RouteFieldState           RouteField = "sta"
	RouteFieldLevel           RouteField = "lvl"
	RouteFieldType            RouteField = "typ"
	RouteFieldCreatedAt       RouteField = "c"
	RouteFieldUpdatedAt       RouteField = "u"
	RouteFieldExpiredAt       RouteField = "e"

	AllRouteFields = []RouteField{
		RouteFieldSourceIP,
		RouteFieldDestinationIP,
		RouteFieldSourcePort,
		RouteFieldDestinationPort,
		RouteFieldRemark,
		RouteFieldAction,
		RouteFieldState,
		RouteFieldLevel,
		RouteFieldCreatedAt,
		RouteFieldUpdatedAt,
		RouteFieldExpiredAt,
	}
)

// route 的  dbkey
func routeFieldDBKey(id string, field RouteField) []byte {
	return bytes.Join([][]byte{[]byte("route"), []byte(field), []byte(id)}, []byte{'/'})
}
