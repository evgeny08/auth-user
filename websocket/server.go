package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	logger log.Logger
}

type Config struct {
	Logger log.Logger
}

// 1
type message struct {
	Msg string `json:"msg"`
	ID  string `json:"id"`
}

type Client struct {
	Conn *websocket.Conn
	ID   string
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(cfg *Config) (*WebSocket, error) {
	// 2
	webSocket := &WebSocket{
		logger: cfg.Logger,
	}
	return webSocket, nil
}

func Writer(coord *message) {
	select {
	case msg := <-broadcast:
		log.Fatal("received message", msg)
	default:
	}
	broadcast <- coord
}

func (s *WebSocket) MsgHandler(w http.ResponseWriter, r *http.Request) {
	var coordinates message
	if err := json.NewDecoder(r.Body).Decode(&coordinates); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go Writer(&coordinates)
}

func (s *WebSocket) WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// register client
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		clients[&Client{
			Conn: ws,
			ID:   string(p),
		}] = true
	}
}

// 3
func (s *WebSocket) Echo() {
	for {
		val := <-broadcast
		msg := fmt.Sprintf("%s", val.Msg)
		msgID := val.ID
		for client := range clients {
			if client.ID == msgID {
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}
