package models

import "github.com/gorilla/websocket"

type Player struct {
	UserId   string          `json:"userid"`
	UserName string          `json:"username"`
	Conn     *websocket.Conn `json:"-"`
	IsBot    bool            `json:"is_bot"`
	IsOnline bool            `json:"is_online"`
}
