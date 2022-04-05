package server

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	libhttpMiddleware "github.com/otamoe/go-library/http/middleware"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/zap"
)

type (
	HttpHandlerRouteInput struct {
		Type            string `json:"type,omitempty"`
		SourceIP        string `json:"sourceIP,omitempty"`
		DestinationIP   string `json:"destinationIP,omitempty"`
		SourcePort      uint32 `json:"sourcePort,omitempty"`
		DestinationPort uint32 `json:"destinationPort,omitempty"`
		Level           int32  `json:"level,omitempty"`
		Remark          string `json:"remark,omitempty"`
		Action          string `json:"action,omitempty"`
		State           string `json:"state,omitempty"`
		ExpiredAt       int64  `json:"expiredAt,omitempty"`
	}
)

func (httpHandler *HttpHandler) readRoute(w http.ResponseWriter, r *http.Request) (route *Route) {
	id, _ := mux.Vars(r)["route"]
	route = httpHandler.routeSystem.Get(id)
	if route == nil {
		httpHandler.writeErrorJson(ErrRouteNotFound, w)
		return
	}
	return
}

func (httpHandler *HttpHandler) saveRoute(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["route"]

	var body []byte
	var data *HttpHandlerRouteInput
	route, err := httpHandler.routeSystem.Save(id, func(route *Route) (rRoute *Route, err error) {
		data = &HttpHandlerRouteInput{
			Type:            route.Route.Type.String(),
			SourceIP:        route.Route.SourceIP,
			DestinationIP:   route.Route.DestinationIP,
			SourcePort:      route.Route.SourcePort,
			DestinationPort: route.Route.DestinationPort,
			Level:           route.Route.Level,
			Remark:          route.Route.Remark,
			Action:          route.Route.Action.String(),
			State:           route.Route.State.String(),
			ExpiredAt:       route.Route.ExpiredAt,
		}
		if err = httpHandler.readJson(&body, data, r); err != nil {
			return
		}
		rRoute = route.Clone()

		if id == "" || data.Type != route.Route.Type.String() {
			val, ok := pb.Route_Type_value[data.Type]
			if !ok {
				err = &ValidateError{
					Name: "type",
				}
				return
			}
			rRoute = rRoute.WithType(pb.Route_Type(val))
		}

		// SourceIP
		if id == "" || data.SourceIP != route.Route.SourceIP {
			var ipNet *net.IPNet
			if _, ipNet, err = net.ParseCIDR(data.SourceIP); err != nil {
				err = &ValidateError{
					Name: "sourceIP",
				}
				return
			}
			if ipNet == nil {
				err = &ValidateError{
					Name: "sourceIP",
				}
				return
			}
			rRoute = rRoute.WithSourceIP(ipNet)
		}

		// DestinationIP
		if id == "" || data.DestinationIP != route.Route.DestinationIP {
			var ipNet *net.IPNet
			if _, ipNet, err = net.ParseCIDR(data.DestinationIP); err != nil || ipNet == nil {
				err = &ValidateError{
					Name: "destinationIP",
				}
				return
			}
			rRoute = rRoute.WithDestinationIP(ipNet)
		}

		if data.Remark != route.Route.Remark {
			if len(data.Remark) > 512 {
				err = &ValidateError{
					Name: "remark",
				}
				return
			}
			rRoute = rRoute.WithRemark(data.Remark)
		}

		rRoute = rRoute.WithSourcePort(data.SourcePort)
		rRoute = rRoute.WithDestinationPort(data.DestinationPort)
		rRoute = rRoute.WithLevel(data.Level)

		if id == "" || data.Action != route.Route.Action.String() {
			val, ok := pb.Route_Action_value[data.Action]
			if !ok {
				err = &ValidateError{
					Name: "action",
				}
				return
			}
			rRoute = rRoute.WithAction(pb.Route_Action(val))
		}

		if id == "" || data.State != route.Route.State.String() {
			val, ok := pb.State_value[data.State]
			if !ok {
				err = &ValidateError{
					Name: "state",
				}
				return
			}
			rRoute = rRoute.WithState(pb.State(val))
		}

		if id == "" || data.ExpiredAt != route.Route.ExpiredAt {
			mint := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
			maxt := time.Date(9000, time.January, 1, 0, 0, 0, 0, time.UTC)
			if mint.Unix() > data.ExpiredAt || maxt.Unix() < data.ExpiredAt {
				err = &ValidateError{
					Name: "expiredAt",
				}
				return
			}
			rRoute = rRoute.WithExpiredAt(time.Unix(data.ExpiredAt, 0).UTC())
		}

		return
	})

	if data != nil {
		dataId := id
		if route != nil {
			dataId = route.Id
		}
		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", dataId),
			zap.String("data.type", data.Type),
			zap.String("data.sourceIP", data.SourceIP),
			zap.String("data.destinationIP", data.DestinationIP),
			zap.Uint32("data.sourcePort", data.SourcePort),
			zap.Uint32("data.destinationPort", data.DestinationPort),
			zap.String("data.remark", data.Remark),
			zap.String("data.action", data.Action),
			zap.String("data.state", data.State),
			zap.Int64("data.ExpiredAt", data.ExpiredAt),
		)
	}

	if err != nil {
		httpHandler.writeErrorJson(err, w)
		return
	}

	httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: route}, w)
}

func (httpHandler *HttpHandler) ListRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routes := httpHandler.routeSystem.All()

		// 过滤器 type
		if r.URL.Query().Has("type") {
			typ := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("type")))
			typInt, ok := pb.Route_Action_value[typ]
			if !ok {
				i, _ := strconv.ParseInt(typ, 10, 32)
				typInt = int32(i)
			}
			pbTyp := pb.Route_Type(typInt)

			oRoutes := Routes{}
			for _, route := range routes {
				if route.Type == pbTyp {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 过滤器 sourceIP
		if r.URL.Query().Has("sourceIP") {
			sourceIP := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("sourceIP")))

			oRoutes := Routes{}
			if _, ipNet, _ := net.ParseCIDR(sourceIP); ipNet != nil {
				for _, route := range routes {
					if route.iSourceIP != nil && ipNet.Contains(route.iSourceIP.IP) {
						oRoutes = append(oRoutes, route)
					}
				}
			} else {
				for _, route := range routes {
					if route.SourceIP == sourceIP {
						oRoutes = append(oRoutes, route)
					}
				}
			}
			routes = oRoutes
		}
		// 过滤器 destinationIP
		if r.URL.Query().Has("destinationIP") {
			destinationIP := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("destinationIP")))

			oRoutes := Routes{}
			if _, ipNet, _ := net.ParseCIDR(destinationIP); ipNet != nil {
				for _, route := range routes {
					if route.iDestinationIP != nil && ipNet.Contains(route.iDestinationIP.IP) {
						oRoutes = append(oRoutes, route)
					}
				}
			} else {
				for _, route := range routes {
					if route.DestinationIP == destinationIP {
						oRoutes = append(oRoutes, route)
					}
				}
			}
			routes = oRoutes
		}

		// 过滤器 action
		if r.URL.Query().Has("action") {
			action := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("action")))
			actionInt, ok := pb.Route_Action_value[action]
			if !ok {
				i, _ := strconv.ParseInt(action, 10, 32)
				actionInt = int32(i)
			}
			pbAction := pb.Route_Action(actionInt)

			oRoutes := Routes{}
			for _, route := range routes {
				if route.Action == pbAction {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 过滤器 state
		if r.URL.Query().Has("state") {
			state := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("state")))
			stateInt, ok := pb.State_value[state]
			if !ok {
				i, _ := strconv.ParseInt(state, 10, 32)
				stateInt = int32(i)
			}
			pbState := pb.State(stateInt)

			oRoutes := Routes{}
			for _, route := range routes {
				if route.State == pbState {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 搜索 search
		if search := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("search"))); search != "" {
			_, ipNet, _ := net.ParseCIDR(search)
			ip := net.ParseIP(search)
			oRoutes := Routes{}
			for _, route := range routes {
				if ipNet != nil && route.iSourceIP != nil && route.iSourceIP.Contains(ipNet.IP) {
					oRoutes = append(oRoutes, route)
				} else if ipNet != nil && route.iDestinationIP != nil && route.iDestinationIP.Contains(ipNet.IP) {
					oRoutes = append(oRoutes, route)
				} else if len(ip) != 0 && route.iSourceIP != nil && route.iSourceIP.Contains(ip) {
					oRoutes = append(oRoutes, route)
				} else if len(ip) != 0 && route.iDestinationIP != nil && route.iDestinationIP.Contains(ip) {
					oRoutes = append(oRoutes, route)
				} else if strings.Index(strings.TrimSpace(strings.ToLower(route.Remark)), search) != -1 {
					oRoutes = append(oRoutes, route)
				} else if strings.Index(route.SourceIP, search) != -1 {
					oRoutes = append(oRoutes, route)
				} else if strings.Index(route.DestinationIP, search) != -1 {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 过滤器 gteCreatedAt
		if queryGteCreatedAt := r.URL.Query().Get("gteCreatedAt"); queryGteCreatedAt != "" {
			queryGteCreatedAtTime, _ := strconv.ParseInt(queryGteCreatedAt, 10, 64)
			oRoutes := Routes{}
			for _, route := range routes {
				if route.CreatedAt >= queryGteCreatedAtTime {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 过滤器 gteUpdatedAt
		if queryGteUpdatedAt := r.URL.Query().Get("gteUpdatedAt"); queryGteUpdatedAt != "" {
			queryGteUpdatedAtTime, _ := strconv.ParseInt(queryGteUpdatedAt, 10, 64)
			oRoutes := Routes{}
			for _, route := range routes {
				if route.UpdatedAt >= queryGteUpdatedAtTime {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 过滤器 gteExpiredAt
		if queryGteExpiredAt := r.URL.Query().Get("gteExpiredAt"); queryGteExpiredAt != "" {
			queryGteExpiredAtTime, _ := strconv.ParseInt(queryGteExpiredAt, 10, 64)
			oRoutes := Routes{}
			for _, route := range routes {
				if route.ExpiredAt >= queryGteExpiredAtTime {
					oRoutes = append(oRoutes, route)
				}
			}
			routes = oRoutes
		}

		// 限制 读取数量
		if queryLimit := r.URL.Query().Get("limit"); queryLimit != "" {
			if queryLimitInt, _ := strconv.Atoi(queryLimit); queryLimitInt > 0 {
				oRoutes := Routes{}
				for _, route := range routes {
					if len(oRoutes) < queryLimitInt {
						oRoutes = append(oRoutes, route)
					} else {
						break
					}
				}
				routes = oRoutes
			}
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: routes}, w)
		return
	})
}

func (httpHandler *HttpHandler) CreateRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveRoute(w, r)
	})
}

func (httpHandler *HttpHandler) ReadRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route *Route
		if route = httpHandler.readRoute(w, r); route == nil {
			return
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: route}, w)
		return
	})
}

func (httpHandler *HttpHandler) UpdateRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveRoute(w, r)
	})
}

func (httpHandler *HttpHandler) DeleteRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["route"]
		var err error
		var route *Route
		route, err = httpHandler.routeSystem.Delete(id)

		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", id),
		)

		if err != nil {
			httpHandler.writeErrorJson(err, w)
			return
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: route}, w)
	})
}
