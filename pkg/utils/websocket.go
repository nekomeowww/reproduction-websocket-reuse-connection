package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	readLimit = int64(4 * 1024 * 1024)
)

// Websocket 封装了 Websocket 连接及其交互方式
type Websocket struct {
	Conn   *websocket.Conn
	Logger *logrus.Logger

	alreadyClosed bool
}

func NewWebsocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, logger *logrus.Logger) (*Websocket, error) {
	// 升级 HTTP 连接
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	// 创建 Websocket 连接封装
	ws := &Websocket{
		Conn:   websocketConn,
		Logger: logger,
	}

	// 配置最大消息大小，当前值为 4MB
	ws.Conn.SetReadLimit(readLimit)
	return ws, nil
}

// Close 关闭连接，在关闭前将会发送 CloseMessage（RFC 定义的控制帧 8）类型的消息来通知客户端关闭自己的连接
func (w *Websocket) Close() {
	// 如果当前连接已经关闭，则不再重复关闭：乐观锁
	if w.alreadyClosed {
		return
	}

	// 优先设定已经关闭的标记，防止重复关闭
	w.alreadyClosed = true
	// 发送关闭连接消息
	_ = w.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 关闭连接
	_ = w.Conn.Close()
}
