package handler

import "github.com/gin-gonic/gin"

// WebSocketConnect 建立 WebSocket 连接（仅用于文档）
// @Summary      WebSocket 连接
// @Description  通过该接口发起 WebSocket 握手，建立长连接；后续使用 topic + JSON 消息收发。
// @Tags         WebSocket
// @Security     BearerAuth
// @Produce      json
// @Router       /ws/ [get]
func WebSocketConnect(c *gin.Context) {}
