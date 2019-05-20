package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocket struct {
	logger log.Logger
}

type Config struct {
	Logger log.Logger
}

// 1
type Message struct {
	Msg string `json:"msg"`
	ID  string `json:"id"`
}

type Client struct {
	Conn *websocket.Conn
	ID   string
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *Message)
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

func (s *WebSocket) Writer(coord *Message) error {
	select {
	case msg := <-broadcast:
		log.Fatal("received message", msg)
	default:
	}
	broadcast <- coord
	return nil
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
		fmt.Println(Client{
			Conn: ws,
			ID:   string(p),
		})
	}
}

// 3
func (s *WebSocket) Echo() {
	for {
		val := <-broadcast
		msg := fmt.Sprintf("%s", val.Msg)
		msgID := val.ID
		for client := range clients {
			fmt.Println(client.ID)
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
