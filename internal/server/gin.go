package server

import (
	"strings"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/xzzpig/simple-ip-kvm/internal/config"
	"github.com/xzzpig/simple-ip-kvm/web"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
)

func Run() {
	logger := zap.L().Named("server")

	var r *gin.Engine

	if *config.GetConfig().Debug {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		r.Use(ginzap.RecoveryWithZap(logger, true))
	}

	if accounts := config.GetConfig().Web.Auth; accounts != nil {
		r.Use(gin.BasicAuth(accounts))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/api/config", apiConfig)
	r.GET("/websocket", handleWebsocket)

	if config.GetConfig().Proxy.Enable {
		r.GET(config.GetConfig().Proxy.Path+"/*any", proxyStream())
		logger.Debug("proxy enabled", zap.String("path", config.GetConfig().Proxy.Path), zap.String("url", config.GetConfig().Video.StreamUrl))
	}

	uiBase := strings.TrimSuffix(config.GetConfig().Web.BasePath, "/")
	r.Use(static.Serve(uiBase+"/", web.NewStaticFileSystem()))

	r.Run(config.GetConfig().Web.Addr)
}
