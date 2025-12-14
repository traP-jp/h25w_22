package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// SetupRoutes sets up all HTTP routes
func SetupRoutes(manager *RoomManager) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /rooms", func(w http.ResponseWriter, r *http.Request) {
		handleCreateRoom(w, r, manager)
	})

	mux.HandleFunc("GET /rooms/{roomID}/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, manager)
	})

	return mux
}

func handleCreateRoom(w http.ResponseWriter, r *http.Request, manager *RoomManager) {
	room := manager.CreateRoom()

	response := map[string]string{
		"roomId": room.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("Created room: %s", room.ID)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, manager *RoomManager) {
	roomID := r.PathValue("roomID")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	room, ok := manager.GetRoom(roomID)
	if !ok {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Check if room is full
	if room.IsFull() {
		http.Error(w, "Room is full", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Generate client ID
	clientID := fmt.Sprintf("client_%d", len(room.Clients)+1)
	client := NewClient(clientID, conn, room)

	room.Register <- client

	log.Printf("Client %s joined room %s", clientID, roomID)

	// Start client's read and write pumps
	go client.WritePump()
	go client.ReadPump()
}
