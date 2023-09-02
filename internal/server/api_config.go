package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xzzpig/simple-ip-kvm/internal/config"
)

func apiConfig(c *gin.Context) {
	cfg := config.GetConfig()
	apiCfg := map[string]interface{}{
		"app_title":    cfg.Web.Title,
		"stream_url":   cfg.Video.StreamUrl + cfg.Video.StreamPath,
		"snapshot_url": cfg.Video.StreamUrl + cfg.Video.SnapshotPath,
		"stream_port":  cfg.Video.Streamer.Port,
	}
	if cfg.Proxy.Enable {
		apiCfg["stream_url"] = "{URL_PROTOCOL}//{URL_HOSTPORT}" + cfg.Proxy.Path + cfg.Video.StreamPath
		apiCfg["snapshot_url"] = "{URL_PROTOCOL}//{URL_HOSTPORT}" + cfg.Video.SnapshotPath
	}
	c.JSON(http.StatusOK, apiCfg)
}
