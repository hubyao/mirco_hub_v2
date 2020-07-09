package router

import (
	"github.com/gin-gonic/gin"
	"micro_demo/commonlib/middleware"
)

// Load 加载中间件
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//router := gin.Default()

	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	Handle(g)

	return g
}

func Handle(g *gin.Engine) {

	u := g.Group("/v1/user")
	{
		u.GET("name")
	}

	setSkipUri()
}

// 跳过验证的url
func setSkipUri() {
	// middleware.SetSkipUri("/v1/user/regist", "skip")
	// middleware.SetSkipUri("/v1/user/signin", "skip")
}
