package router

import (
	"log/slog"
	"net/http"

	"api/internal/handler"
	"api/internal/middleware"
)

func NewRouter(calcHandler *handler.CalculatorHandler, logger *slog.Logger) http.Handler {
	//创建路由器, 相当于 gin.Default()
	mux := http.NewServeMux()
	//注册路由，相当于r.POST("/calculate", calcHandler.Calculate)
	mux.HandleFunc("POST /calculate", calcHandler.Calculate)
	//先有 mux，再让日志中间件把它包起来，最后返回一个带日志能力的总路由。
	return middleware.Logging(logger)(mux)
}
