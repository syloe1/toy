package main

import (
	"log/slog"
	"net/http"
	"os"

	"api/internal/service"

	"api/internal/handler"
	"api/internal/router"
)

func main() {
	//创建日志器
	//slog.NewTextHandler(...)

	//log和fmt区别，log带时间，带等级
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//创建service , handler对象
	calcService := service.NewCalculatorService()
	calcHandler := handler.NewCalculatorHandler(calcService)
	//创建路由
	r := router.NewRouter(calcHandler, logger)

	//给路由挂上中间件

	//启动HTTP服务\
	server := &http.Server{
		Addr:    ":7777",
		Handler: r,
	}

	logger.Info("server started", "addr", ":7777")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		//os.Exit(0) 正常下班，关门走人
		os.Exit(1) //出事了，紧急停工
	}
}
