package routers

import (
	"net/http"

	"github.com/Shanmukh1609/backend/controllers"
	"github.com/Shanmukh1609/backend/websocket"
)

func GameEntry(mux *http.ServeMux) {
	mux.HandleFunc("/setCookie", controllers.HandleCookie)
	mux.HandleFunc("/enterRoomWS", websocket.HandleWebSocketConnections)
	mux.HandleFunc("/", controllers.HelloFunc)
	mux.HandleFunc("/leaderBoard", controllers.LeaderBoard)
}
