package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/models"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/services"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan CryptoCoinResponse
}

type CryptoCoinResponseWebsocket interface {
	SendMessage(responce CryptoCoinResponse)
	HasConnectedClients() bool
}

type CryptoCoinResponse struct {
	Coins []models.Coin `json:"coins"`
}

func newHub() *hub {
	return &hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan CryptoCoinResponse),
	}
}

// Start ws http server and return interface with broadcast methods
func StartWebSocket(dataService services.DataService) CryptoCoinResponseWebsocket {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Logger.Print(err)
			return
		}

		defer func() {
			delete(hub.clients, ws)
			ws.Close()
		}()

		hub.clients[ws] = true

		log.Logger.Print("Connected!")

		coins, err := dataService.GetTop50CoinMarketCurrency()
		if err == nil {
			hub.broadcast <- CryptoCoinResponse{Coins: coins}
		}

		for {
			// Read message from the WebSocket connection
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Logger.Print("WebSocket read error:", err)
				break
			}

			// Handle the received message as needed
		}
	})

	go func() {
		// Start http server with ws
		cfg := configs.GetConfig().Server
		fmt.Printf("Server listening on port %s", cfg.WebSocketPort)
		log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebSocketPort), nil))
	}()

	return hub
}

// Send message to clients of ws
func (h *hub) SendMessage(responce CryptoCoinResponse) {
	h.broadcast <- responce
}

// Checking do you have connected clients
func (h *hub) HasConnectedClients() bool {
	return len(h.clients) > 0
}

func (h *hub) run() {
	for {
		select {
		case responce := <-h.broadcast:
			for client := range h.clients {
				if err := client.WriteJSON(responce); err != nil {
					log.Logger.Printf("error occurred: %v", err)
				}
			}
		}
	}
}
