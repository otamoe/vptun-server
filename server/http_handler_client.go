package server

import (
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	libhttpMiddleware "github.com/otamoe/go-library/http/middleware"
	libutils "github.com/otamoe/go-library/utils"
	pb "github.com/otamoe/vptun-pb"
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
			if len(ip) == 0 {
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
			mint := time.Date(999, time.January, 1, 0, 0, 0, 0, time.UTC)
			maxt := time.Date(9001, time.January, 1, 0, 0, 0, 0, time.UTC)
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
