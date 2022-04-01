package server

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	libhttp "github.com/otamoe/go-library/http"
	libhttpMiddleware "github.com/otamoe/go-library/http/middleware"
	liblogger "github.com/otamoe/go-library/logger"
	libutils "github.com/otamoe/go-library/utils"
	"github.com/otamoe/vptun-server/assets"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

type (
	HttpHandler struct {
		routeSystem       *RouteSystem
		clientSystem      *ClientSystem
		clientShellSystem *ClientShellSystem
		grpcHandler       *GrpcHandler
	}
	HttpHandlerResponseError struct {
		Error string `json:"error"`
	}
	HttpHandlerResponseSuccess struct {
		Data interface{} `json:"data"`
	}
)

// Handler
func (httpHandler *HttpHandler) Handler(next http.Handler) http.Handler {
	router := mux.NewRouter()

	// 路由
	router.Handle("/api/route", httpHandler.ListRoute()).Methods("GET", "HEAD")
	router.Handle("/api/route", httpHandler.CreateRoute()).Methods("POST", "PUT")
	router.Handle("/api/route/{route}", httpHandler.ReadRoute()).Methods("GET", "HEAD")
	router.Handle("/api/route/{route}", httpHandler.UpdateRoute()).Methods("POST", "PUT")
	router.Handle("/api/route/{route}", httpHandler.DeleteRoute()).Methods("DELETE")

	// 客户端
	router.Handle("/api/client", httpHandler.ListClient()).Methods("GET", "HEAD")
	router.Handle("/api/client", httpHandler.CreateClient()).Methods("POST", "PUT")
	router.Handle("/api/client/{client}", httpHandler.ReadClient()).Methods("GET", "HEAD")
	router.Handle("/api/client/{client}", httpHandler.UpdateClient()).Methods("POST", "PUT")
	router.Handle("/api/client/{client}", httpHandler.DeleteClient()).Methods("DELETE")

	router.Handle("/api/client/{client}/shell", httpHandler.ListClientShell()).Methods("GET", "HEAD")
	router.Handle("/api/client/{client}/shell", httpHandler.CreateClientShell()).Methods("POST", "PUT")
	router.Handle("/api/client/{client}/shell/{shell}", httpHandler.ReadClientShell()).Methods("GET", "HEAD")
	router.Handle("/api/client/{client}/shell/{shell}", httpHandler.UpdateClientShell()).Methods("GET", "HEAD")
	router.Handle("/api/client/{client}/shell/{shell}", httpHandler.DeleteClientShell()).Methods("DELETE")

	router.NotFoundHandler = next
	return router
}

func (httpHandler *HttpHandler) writeErrorJson(err error, w http.ResponseWriter) {
	if err == nil {
		err = &HttpError{
			Status: http.StatusInternalServerError,
		}
	}
	switch terr := err.(type) {
	case *HttpError:
		httpHandler.writeJson(terr.StatusCode(), &HttpHandlerResponseError{Error: terr.Error()}, w)
	case *ValidateError:
		httpHandler.writeJson(http.StatusForbidden, &HttpHandlerResponseError{Error: terr.Error()}, w)
	case *NotFoundError:
		httpHandler.writeJson(http.StatusNotFound, &HttpHandlerResponseError{Error: terr.Error()}, w)
	default:
		if strings.Index(strings.ToLower(terr.Error()), "not found") != -1 {
			httpHandler.writeJson(http.StatusNotFound, &HttpHandlerResponseError{Error: terr.Error()}, w)
		} else {
			httpHandler.writeJson(http.StatusInternalServerError, &HttpHandlerResponseError{Error: terr.Error()}, w)
		}
	}
}
func (httpHandler *HttpHandler) writeJson(statusCode int, value interface{}, w http.ResponseWriter) {
	// 输出内容类型
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// 编码json
	data, err := json.Marshal(value)

	// 编码错误
	if err != nil {
		data, _ = json.Marshal(&HttpHandlerResponseError{Error: err.Error()})
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}

func (httpHandler *HttpHandler) readJson(body *[]byte, value interface{}, r *http.Request) (err error) {
	// 没 body 载入
	if body == nil || len(*body) == 0 {
		// 没长度
		contentLength := r.Header.Get("Content-length")
		if contentLength == "" {
			err = &HttpError{Status: http.StatusLengthRequired}
			return
		}

		// 不是 json
		if strings.Index(r.Header.Get("Content-Type"), "/json") == -1 {
			err = &HttpError{Status: http.StatusUnsupportedMediaType}
			return
		}

		// 长度是空
		length, _ := strconv.Atoi(contentLength)
		if length <= 0 {
			err = &HttpError{Status: http.StatusBadRequest}
			return
		}

		// 读取错误
		var rbody []byte
		if rbody, err = ioutil.ReadAll(io.LimitReader(r.Body, int64(length))); err != nil || len(rbody) == 0 {
			err = &HttpError{Status: http.StatusBadRequest}
			return
		}

		body = &rbody
	}

	// 解析 json
	if err = json.Unmarshal(*body, value); err != nil {
		err = &HttpError{Status: http.StatusForbidden}
		return
	}
	return
}

func HttpHandlerRegister(routeSystem *RouteSystem, clientSystem *ClientSystem, clientShellSystem *ClientShellSystem, grpcHandler *GrpcHandler) (httpHandler *HttpHandler, out libhttp.OutOption, err error) {
	httpHandler = &HttpHandler{
		routeSystem:       routeSystem,
		clientSystem:      clientSystem,
		clientShellSystem: clientShellSystem,
		grpcHandler:       grpcHandler,
	}

	out.Option = func(server *libhttp.Server) error {
		// 日志
		loggerMiddleware := &libhttpMiddleware.Logger{
			Logger:    liblogger.Get("http"),
			SlowQuery: time.Second * 30,
		}

		// 认证
		password := viper.GetString("http.password")
		if password == "" {
			password = string(libutils.RandByte(16, libutils.RandAlphaNumber))
		}
		basicAuthMiddleware := &libhttpMiddleware.BasicAuth{
			Auths: map[string]string{
				viper.GetString("http.username"): password,
			},
		}

		// 编码
		compressMiddleware := &libhttpMiddleware.Compress{
			Types:     []string{"text/", "application/json", "application/javascript", "application/atom+xml", "application/rss+xml", "application/xml"},
			Gzip:      true,
			GzipLevel: gzip.DefaultCompression,
		}

		// 指标
		prometheusHandler := func(next http.Handler) http.Handler {
			h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/metrics") {
					h.ServeHTTP(rw, r)
				} else {
					next.ServeHTTP(rw, r)
				}
			})
		}

		// 静态文件
		staticMiddleware := &libhttpMiddleware.Static{
			MaxAge:   86400 * 30,
			FSPath:   "public",
			FS:       assets.FS,
			ModTime:  time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
			Redirect: "/index.html",
		}
		server.Handlers = append(
			server.Handlers,
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   100,
				Handler: loggerMiddleware.Handler,
			},
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   200,
				Handler: basicAuthMiddleware.Handler,
			},
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   300,
				Handler: compressMiddleware.Handler,
			},
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   700,
				Handler: httpHandler.Handler,
			},
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   800,
				Handler: prometheusHandler,
			},
			libhttp.HandlerOption{
				Hosts:   []string{"*"},
				Index:   1000,
				Handler: staticMiddleware.Handler,
			},
		)
		return nil
	}
	return
}

// // 指标
// appPrometheus.New(config.Prometheus.Host, config.Prometheus.Path),

// // 静态文件
// appAssets.New(config.FS.Host),
