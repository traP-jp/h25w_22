package handler

import (
	"testing"
	"time"

	"github.com/traP-jp/h25w_22/game"
)

func TestRoomDeletionOnGameEnd(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()
	roomID := room.ID

	// Verify room exists
	if _, ok := rm.GetRoom(roomID); !ok {
		t.Fatal("Room should exist after creation")
	}

	// Simulate game end
	room.mu.Lock()
	room.Game.Phase = game.PhaseGameEnd
	room.mu.Unlock()

	// Trigger endGame
	room.endGame()

	// Give some time for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// Verify room is deleted
	if _, ok := rm.GetRoom(roomID); ok {
		t.Error("Room should be deleted after game ends")
	}
}

func TestRoomDeletionOnClientDisconnect(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()
	roomID := room.ID

	// Create mock clients
	client1 := &Client{
		ID:   "client_1",
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}
	client2 := &Client{
		ID:   "client_2",
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}

	// Register clients
	room.Register <- client1
	room.Register <- client2

	// Give some time for registration
	time.Sleep(50 * time.Millisecond)

	// Verify room exists and has 2 clients
	if _, ok := rm.GetRoom(roomID); !ok {
		t.Fatal("Room should exist with clients")
	}

	room.mu.Lock()
	numClients := len(room.Clients)
	room.mu.Unlock()

	if numClients != 2 {
		t.Fatalf("Expected 2 clients, got %d", numClients)
	}

	// Simulate client disconnect
	room.Unregister <- client1

	// Give some time for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// Verify room is deleted (because only 1 client remains)
	if _, ok := rm.GetRoom(roomID); ok {
		t.Error("Room should be deleted when client disconnects and only one remains")
	}
}

func TestClientNotificationOnOpponentDisconnect(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()

	// Create mock clients
	client1 := &Client{
		ID:   "client_1",
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}
	client2 := &Client{
		ID:   "client_2",
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}

	// Register clients
	room.Register <- client1
	room.Register <- client2

	// Give some time for registration
	time.Sleep(50 * time.Millisecond)

	// Simulate client1 disconnect
	go func() {
		room.Unregister <- client1
	}()

	// Check if client2 receives notification
	select {
	case msg := <-client2.Send:
		// Verify it's an OPPONENT_DISCONNECTED event
		if len(msg) == 0 {
			t.Error("Received empty message")
		}
		// Message should contain "OPPONENT_DISCONNECTED"
		msgStr := string(msg)
		if len(msgStr) == 0 {
			t.Error("Message string is empty")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Client should receive notification when opponent disconnects")
	}
}

func TestRoomDeletionWithNoClients(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()
	roomID := room.ID

	// Create a mock client
	client1 := &Client{
		ID:   "client_1",
		Room: room,
		Send: make(chan []byte, channelBufferSize),
	}

	// Register client
	room.Register <- client1

	// Give some time for registration
	time.Sleep(50 * time.Millisecond)

	// Simulate client disconnect (last client)
	room.Unregister <- client1

	// Give some time for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	// Verify room is deleted
	if _, ok := rm.GetRoom(roomID); ok {
		t.Error("Room should be deleted when last client disconnects")
	}
}
