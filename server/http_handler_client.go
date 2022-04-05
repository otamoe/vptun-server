package server

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	libhttpMiddleware "github.com/otamoe/go-library/http/middleware"
	libutils "github.com/otamoe/go-library/utils"
	pb "github.com/otamoe/vptun-pb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	HttpHandlerClientInput struct {
		Key          string `json:"key,omitempty"`
		RouteAddress string `json:"routeAddress,omitempty"`
		Remark       string `json:"remark,omitempty"`
		State        string `json:"state,omitempty"`
		ExpiredAt    int64  `json:"expiredAt,omitempty"`
	}
)

func (httpHandler *HttpHandler) readClient(w http.ResponseWriter, r *http.Request) (client *Client) {
	id, _ := mux.Vars(r)["client"]
	client = httpHandler.clientSystem.Get(id)

	if client == nil {
		httpHandler.writeErrorJson(ErrClientNotFound, w)
		return
	}
	return
}

func (httpHandler *HttpHandler) saveClient(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["client"]

	var body []byte
	var data *HttpHandlerClientInput
	client, err := httpHandler.clientSystem.Save(id, func(client *Client) (rClient *Client, err error) {
		data = &HttpHandlerClientInput{
			Key:          client.Client.Key,
			RouteAddress: client.Client.RouteAddress,
			Remark:       client.Client.Remark,
			State:        client.Client.State.String(),
			ExpiredAt:    client.Client.ExpiredAt,
		}
		if err = httpHandler.readJson(&body, data, r); err != nil {
			return
		}
		rClient = client.Clone()

		// Key
		if data.Key == "" {
			data.Key = string(libutils.RandByte(16, libutils.RandAlphaNumber))
		}
		if id == "" || data.Key != client.Client.Key {
			rClient = rClient.WithKey(data.Key)
		}

		// RouteAddress
		if data.RouteAddress == "" {
			data.RouteAddress = httpHandler.clientSystem.NewRouteAddress(false).String()
		}
		if id == "" || data.RouteAddress != client.Client.RouteAddress {
			ip := net.ParseIP(data.RouteAddress)
			subnet, _ := viper.Get("route.subnet").(net.IPNet)
			if len(ip) == 0 || !subnet.Contains(ip) {
				err = &ValidateError{
					Name: "routeAddress",
				}
				return
			}
			rClient = rClient.WithRouteAddress(ip)
		}

		if data.Remark != client.Client.Remark {
			if len(data.Remark) > 512 {
				err = &ValidateError{
					Name: "remark",
				}
				return
			}
			rClient = rClient.WithRemark(data.Remark)
		}

		if id == "" || data.State != client.Client.State.String() {
			val, ok := pb.State_value[data.State]
			if !ok {
				err = &ValidateError{
					Name: "state",
				}
				return
			}
			rClient = rClient.WithState(pb.State(val))
		}

		if id == "" || data.ExpiredAt != client.Client.ExpiredAt {
			mint := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
			maxt := time.Date(9000, time.January, 1, 0, 0, 0, 0, time.UTC)
			if mint.Unix() > data.ExpiredAt || maxt.Unix() < data.ExpiredAt {
				err = &ValidateError{
					Name: "expiredAt",
				}
				return
			}
			rClient = rClient.WithExpiredAt(time.Unix(data.ExpiredAt, 0).UTC())
		}

		return
	})

	if data != nil {
		dataId := id
		if client != nil {
			dataId = client.Id
		}
		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", dataId),
			zap.String("data.routeAddress", data.RouteAddress),
			zap.String("data.remark", data.Remark),
			zap.String("data.state", data.State),
			zap.Int64("data.ExpiredAt", data.ExpiredAt),
		)
	}

	if err != nil {
		httpHandler.writeErrorJson(err, w)
		return
	}

	httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: client}, w)
}

func (httpHandler *HttpHandler) ListClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clients := httpHandler.clientSystem.All()
		for i := 0; i < len(clients); i++ {
			clients[i] = clients[i].WithStatus(nil)
		}

		// 过滤器 hostname
		if r.URL.Query().Has("hostname") {
			hostname := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("hostname")))
			oClients := Clients{}
			for _, client := range clients {
				if strings.TrimSpace(strings.ToLower(client.Hostname)) == hostname {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 connectAddress
		if r.URL.Query().Has("connectAddress") {
			connectAddress := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("connectAddress")))
			oClients := Clients{}

			if _, ipNet, _ := net.ParseCIDR(connectAddress); ipNet != nil {
				for _, client := range clients {
					if client.ConnectAddress == "" {
						continue
					}

					var ip net.IP
					if index := strings.LastIndex(client.ConnectAddress, ":"); index > 0 {
						ip = net.ParseIP(client.ConnectAddress[0:index])
					} else {
						ip = net.ParseIP(client.ConnectAddress)
					}
					if len(ip) != 0 && ipNet.Contains(ip) {
						oClients = append(oClients, client)
					}
				}
			} else {
				for _, client := range clients {
					if client.ConnectAddress == connectAddress {
						oClients = append(oClients, client)
					}
				}
			}
			clients = oClients
		}

		// 过滤器 routeAddress
		if r.URL.Query().Has("routeAddress") {
			routeAddress := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("routeAddress")))

			oClients := Clients{}
			if _, ipNet, _ := net.ParseCIDR(routeAddress); ipNet != nil {
				for _, client := range clients {
					if len(client.iRouteAddress) != 0 && ipNet.Contains(client.iRouteAddress) {
						oClients = append(oClients, client)
					}
				}
			} else {
				for _, client := range clients {
					if client.RouteAddress == routeAddress {
						oClients = append(oClients, client)
					}
				}
			}
			clients = oClients
		}

		// 过滤器 online
		if r.URL.Query().Has("online") {
			online := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("online")))
			var onlineVal bool
			if online == "true" || online == "1" {
				onlineVal = true
			}
			oClients := Clients{}
			for _, client := range clients {
				if client.Online == onlineVal {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 shell
		if r.URL.Query().Has("shell") {
			shell := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("shell")))
			var shellVal bool
			if shell == "true" || shell == "1" {
				shellVal = true
			}
			oClients := Clients{}
			for _, client := range clients {
				if client.Shell == shellVal {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
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

			oClients := Clients{}
			for _, client := range clients {
				if client.State == pbState {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 搜索 search
		if search := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("search"))); search != "" {
			_, ipNet, _ := net.ParseCIDR(search)
			oClients := Clients{}
			for _, client := range clients {
				if ipNet != nil && len(client.iRouteAddress) != 0 && ipNet.Contains(client.iRouteAddress) {
					oClients = append(oClients, client)
				} else if strings.Index(strings.TrimSpace(strings.ToLower(client.Hostname)), search) != -1 {
					oClients = append(oClients, client)
				} else if strings.Index(strings.TrimSpace(strings.ToLower(client.Remark)), search) != -1 {
					oClients = append(oClients, client)
				} else if strings.Index(client.ConnectAddress, search) != -1 {
					oClients = append(oClients, client)
				} else if strings.Index(client.RouteAddress, search) != -1 {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 gteCreatedAt
		if queryGteCreatedAt := r.URL.Query().Get("gteCreatedAt"); queryGteCreatedAt != "" {
			queryGteCreatedAtTime, _ := strconv.ParseInt(queryGteCreatedAt, 10, 64)
			oClients := Clients{}
			for _, client := range clients {
				if client.CreatedAt >= queryGteCreatedAtTime {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 gteUpdatedAt
		if queryGteUpdatedAt := r.URL.Query().Get("gteUpdatedAt"); queryGteUpdatedAt != "" {
			queryGteUpdatedAtTime, _ := strconv.ParseInt(queryGteUpdatedAt, 10, 64)
			oClients := Clients{}
			for _, client := range clients {
				if client.UpdatedAt >= queryGteUpdatedAtTime {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 gteConnectAt
		if queryGteConnectAt := r.URL.Query().Get("gteConnectAt"); queryGteConnectAt != "" {
			queryGteConnectAtTime, _ := strconv.ParseInt(queryGteConnectAt, 10, 64)
			oClients := Clients{}
			for _, client := range clients {
				if client.ConnectAt >= queryGteConnectAtTime {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 过滤器 gteExpiredAt
		if queryGteExpiredAt := r.URL.Query().Get("gteExpiredAt"); queryGteExpiredAt != "" {
			queryGteExpiredAtTime, _ := strconv.ParseInt(queryGteExpiredAt, 10, 64)
			oClients := Clients{}
			for _, client := range clients {
				if client.ExpiredAt >= queryGteExpiredAtTime {
					oClients = append(oClients, client)
				}
			}
			clients = oClients
		}

		// 限制 读取数量
		if queryLimit := r.URL.Query().Get("limit"); queryLimit != "" {
			if queryLimitInt, _ := strconv.Atoi(queryLimit); queryLimitInt > 0 {
				oClients := Clients{}
				for _, client := range clients {
					if len(oClients) < queryLimitInt {
						oClients = append(oClients, client)
					} else {
						break
					}
				}
				clients = oClients
			}
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: clients}, w)
		return
	})
}

func (httpHandler *HttpHandler) CreateClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveClient(w, r)
	})
}

func (httpHandler *HttpHandler) ReadClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var client *Client
		if client = httpHandler.readClient(w, r); client == nil {
			return
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: client}, w)
		return
	})
}

func (httpHandler *HttpHandler) UpdateClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveClient(w, r)
	})
}

func (httpHandler *HttpHandler) DeleteClient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["client"]
		var err error
		var client *Client
		client, err = httpHandler.clientSystem.Delete(id)

		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", id),
		)
		if err != nil {
			httpHandler.writeErrorJson(err, w)
			return
		}
		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: client}, w)
	})
}
