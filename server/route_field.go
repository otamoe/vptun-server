package server

import "bytes"

type (
	RouteField string
)

var (
	RouteFieldSource      RouteField = "src"
	RouteFieldDestination RouteField = "dst"
	RouteFieldRemark      RouteField = "rem"
	RouteFieldAction      RouteField = "act"
	RouteFieldState       RouteField = "sta"
	RouteFieldCreatedAt   RouteField = "c"
	RouteFieldUpdatedAt   RouteField = "u"
	RouteFieldExpiredAt   RouteField = "e"

	AllRouteFields = []RouteField{
		RouteFieldSource,
		RouteFieldDestination,
		RouteFieldRemark,
		RouteFieldAction,
		RouteFieldState,
		RouteFieldCreatedAt,
		RouteFieldUpdatedAt,
		RouteFieldExpiredAt,
	}
)

// route 的  dbkey
func routeFieldDBKey(id string, field RouteField) []byte {
	return bytes.Join([][]byte{[]byte("route"), []byte(field), []byte(id)}, []byte{'/'})
}
