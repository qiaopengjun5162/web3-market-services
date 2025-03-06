package routes

import (
	"github.com/go-chi/chi/v5" // https://github.com/go-chi/chi

	"github.com/qiaopengjun5162/web3-market-services/services/rest/service"
)

type Routes struct {
	router *chi.Mux
	srv    service.RestService
}

// NewRoutes 创建并初始化一个新的 Routes 实例。
// 参数:
//
//	r *chi.Mux: 一个指向 chi.Mux 路由复用器的指针，用于定义和管理HTTP路由。
//	srv service.RestService: 一个实现了 RestService 接口的服务实例，用于处理HTTP请求。
//
// 返回值:
//
//	*Routes: 一个指向初始化后的 Routes 结构的指针，通过它可以在应用程序中管理路由。
func NewRoutes(r *chi.Mux, srv service.RestService) *Routes {
	return &Routes{
		router: r,
		srv:    srv,
	}
}
