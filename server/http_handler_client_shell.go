package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	libhttpMiddleware "github.com/otamoe/go-library/http/middleware"
	pb "github.com/otamoe/vptun-pb"
	"go.uber.org/zap"
)

type HttpHandlerClientShellInput struct {
	Input   string `json:"input,omitempty"`
	Remark  string `json:"remark,omitempty"`
	Timeout uint32 `json:"timeout,omitempty"`
}

func (httpHandler *HttpHandler) readClientShell(w http.ResponseWriter, r *http.Request) (client *Client, clientShell *ClientShell) {
	clientId, _ := mux.Vars(r)["client"]
	client = httpHandler.clientSystem.Get(clientId)
	if client == nil {
		httpHandler.writeErrorJson(ErrClientNotFound, w)
		return
	}
	id, _ := mux.Vars(r)["shell"]
	var err error
	clientShell, err = httpHandler.clientShellSystem.Get(clientId, id)
	if err != nil {
		httpHandler.writeErrorJson(err, w)
		return
	}
	if clientShell == nil {
		httpHandler.writeErrorJson(ErrClientShellNotFound, w)
		return
	}
	return
}

func (httpHandler *HttpHandler) saveClientShell(w http.ResponseWriter, r *http.Request) {
	clientId, _ := mux.Vars(r)["client"]
	id, _ := mux.Vars(r)["shell"]
	var body []byte
	var data *HttpHandlerClientShellInput

	// client 未找到
	client := httpHandler.clientSystem.Get(clientId)
	if client == nil {
		httpHandler.writeErrorJson(ErrClientNotFound, w)
		return
	}

	// 客户端不在线
	var grpcClient *GrpcClient
	if id == "" {
		grpcClient = httpHandler.grpcHandler.GetClient(clientId)
		if grpcClient == nil {
			err := &HttpError{
				Err:    errors.New("Client is offline"),
				Status: http.StatusForbidden,
			}
			httpHandler.writeErrorJson(err, w)
			return
		}

		// 客户端不支持
		if !client.Shell {
			err := &HttpError{
				Err:    errors.New("Client does not support shell"),
				Status: http.StatusForbidden,
			}
			httpHandler.writeErrorJson(err, w)
			return
		}
	}

	clientShell, err := httpHandler.clientShellSystem.Save(clientId, id, func(clientShell *ClientShell) (rClientShell *ClientShell, err error) {
		data = &HttpHandlerClientShellInput{
			Input:   clientShell.Input,
			Remark:  clientShell.Remark,
			Timeout: clientShell.Timeout,
		}
		if err = httpHandler.readJson(&body, data, r); err != nil {
			return
		}

		// 超时不能大于 1d
		if data.Timeout > 86400 {
			err = &ValidateError{
				Name: "timeout",
			}
			return
		}

		if data.Timeout < 10 {
			data.Timeout = 10
		}

		// 不能更新 输入
		if id != "" && data.Input != clientShell.Input {
			err = &ValidateError{
				Name: "input",
			}
			return
		}

		rClientShell = clientShell.Clone()
		rClientShell = rClientShell.WithInput(data.Input)
		rClientShell = rClientShell.WithRemark(data.Remark)
		rClientShell = rClientShell.WithTimeout(data.Timeout)

		return
	})

	// grpcClient 发送
	if err == nil && id == "" {
		err = grpcClient.Response(&pb.StreamResponse{Shell: &pb.ShellResponse{
			Id:      clientShell.Id,
			Data:    []byte(clientShell.Input),
			Timeout: clientShell.Timeout,
		}})
	}

	if data != nil {
		dataId := id
		if clientShell != nil {
			dataId = clientShell.Id
		}
		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", dataId),
			zap.String("data.clientId", clientId),
			zap.String("data.remark", data.Remark),
			zap.Uint32("data.timeout", data.Timeout),
		)
	}

	if err != nil {
		httpHandler.writeErrorJson(err, w)
		return
	}
	httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: clientShell}, w)
}

func (httpHandler *HttpHandler) ListClientShell() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		clientId, _ := mux.Vars(r)["client"]
		client := httpHandler.clientSystem.Get(clientId)

		if client == nil {
			httpHandler.writeErrorJson(ErrClientNotFound, w)
			return
		}

		ltId := r.URL.Query().Get("ltId")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit >= 1000 {
			limit = 50
		}

		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.clientId", clientId),
			zap.String("data.ltId", ltId),
			zap.Int("data.limit", limit),
		)

		var clientShells ClientShells
		if clientShells, err = httpHandler.clientShellSystem.List(clientId, ltId, limit); err != nil {
			httpHandler.writeErrorJson(err, w)
			return
		}

		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: clientShells}, w)
	})
}

func (httpHandler *HttpHandler) CreateClientShell() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveClientShell(w, r)
	})
}

func (httpHandler *HttpHandler) ReadClientShell() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientShell *ClientShell
		if _, clientShell = httpHandler.readClientShell(w, r); clientShell == nil {
			return
		}
		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: clientShell}, w)
		return
	})
}

func (httpHandler *HttpHandler) UpdateClientShell() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.saveClientShell(w, r)
	})
}

func (httpHandler *HttpHandler) DeleteClientShell() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId, _ := mux.Vars(r)["client"]
		id, _ := mux.Vars(r)["shell"]
		var err error
		var clientShell *ClientShell
		clientShell, err = httpHandler.clientShellSystem.Delete(clientId, id)

		libhttpMiddleware.LoggerFields(
			r.Context(),
			zap.String("data.id", id),
			zap.String("data.clientId", clientId),
		)
		if err != nil {
			httpHandler.writeErrorJson(err, w)
			return
		}
		httpHandler.writeJson(http.StatusOK, &HttpHandlerResponseSuccess{Data: clientShell}, w)
	})
}
