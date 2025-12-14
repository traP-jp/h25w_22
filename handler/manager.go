package handler

import (
	"fmt"
	"sync"
)

// RoomManager manages all game rooms
type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

// NewRoomManager creates a new room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// CreateRoom creates a new room with a unique ID
func (rm *RoomManager) CreateRoom() *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Generate a unique room ID
	id := fmt.Sprintf("room_%d", len(rm.rooms)+1)
	room := NewRoom(id)
	rm.rooms[id] = room

	// Start the room's event loop
	go room.Run()

	return room
}

// GetRoom retrieves a room by ID
func (rm *RoomManager) GetRoom(id string) (*Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, ok := rm.rooms[id]
	return room, ok
}

// DeleteRoom removes a room from the manager
func (rm *RoomManager) DeleteRoom(id string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.rooms, id)
}
