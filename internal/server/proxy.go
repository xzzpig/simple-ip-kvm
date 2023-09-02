package server

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xzzpig/simple-ip-kvm/internal/config"
	"go.uber.org/zap"
)

func proxyStream() gin.HandlerFunc {
	cfg := config.GetConfig()
	url, err := url.Parse(cfg.Video.StreamUrl)
	if err != nil {
		zap.L().Named("proxy").Panic("stream url parse error", zap.Error(err), zap.String("url", cfg.Video.StreamUrl))
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		if cfg.Proxy.Rewrite {
			c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, cfg.Proxy.Path)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
