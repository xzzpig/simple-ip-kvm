package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xzzpig/simple-ip-kvm/internal/hid"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsLoop(ws)
}

type WSMsg struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func wsLoop(ws *websocket.Conn) {
	defer ws.Close()
	ws.WriteJSON(&WSMsg{
		Type:    "welcome",
		Payload: []byte("Open IP-KVM Server"),
	})
	for {
		//读取ws中的数据
		var message WSMsg
		err := ws.ReadJSON(&message)
		if err != nil {
			break
		}
		hid.Write(message.Payload)
	}
}
