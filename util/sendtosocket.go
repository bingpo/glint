package util

import (
	"glint/logger"
	"time"
)

func SendToSocket(SocketMsg *chan map[string]interface{}, status int, Value interface{}) {
	if SocketMsg == nil {
		logger.Error("socketmsg chan is nil , status = %v", status)
		return
	}
	Element := make(map[string]interface{}, 1)
	Element["status"] = status
	Element["CrawUrl"] = Value
	select {
	case (*SocketMsg) <- Element:
	case <-time.After(time.Second * 5):
	}
}
