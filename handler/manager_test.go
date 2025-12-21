package handler

import (
	"regexp"
	"testing"
)

var roomIDPattern = regexp.MustCompile(`^\d{4}$`)

func TestNewRoomManager(t *testing.T) {
	rm := NewRoomManager()
	if rm == nil {
		t.Fatal("NewRoomManager() returned nil")
	}
	if rm.rooms == nil {
		t.Error("RoomManager.rooms should not be nil")
	}
}

func TestCreateRoom(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()

	if room == nil {
		t.Fatal("CreateRoom() returned nil")
	}
	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}

	// Verify room ID is a 4-digit number
	if !roomIDPattern.MatchString(room.ID) {
		t.Errorf("Room ID should be a 4-digit number, got: %s", room.ID)
	}

	// Verify room is in manager
	retrievedRoom, ok := rm.GetRoom(room.ID)
	if !ok {
		t.Error("Room should be retrievable after creation")
	}
	if retrievedRoom != room {
		t.Error("Retrieved room should be the same as created room")
	}
}

func TestGetRoom(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()

	tests := []struct {
		name   string
		roomID string
		wantOK bool
	}{
		{"Existing room", room.ID, true},
		{"Non-existent room", "non_existent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := rm.GetRoom(tt.roomID)
			if ok != tt.wantOK {
				t.Errorf("GetRoom(%s) ok = %v, want %v", tt.roomID, ok, tt.wantOK)
			}
		})
	}
}

func TestDeleteRoom(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()

	// Verify room exists
	if _, ok := rm.GetRoom(room.ID); !ok {
		t.Fatal("Room should exist before deletion")
	}

	// Delete room
	rm.DeleteRoom(room.ID)

	// Verify room is deleted
	if _, ok := rm.GetRoom(room.ID); ok {
		t.Error("Room should not exist after deletion")
	}
}

func TestCreateRoomUniqueness(t *testing.T) {
	rm := NewRoomManager()

	// Create multiple rooms and verify they have unique IDs
	roomIDs := make(map[string]bool)
	const numRooms = 100

	for i := 0; i < numRooms; i++ {
		room := rm.CreateRoom()
		if roomIDs[room.ID] {
			t.Errorf("Duplicate room ID generated: %s", room.ID)
		}
		roomIDs[room.ID] = true

		// Verify it's a 4-digit number
		if !roomIDPattern.MatchString(room.ID) {
			t.Errorf("Room ID should be a 4-digit number, got: %s", room.ID)
		}
	}

	if len(roomIDs) != numRooms {
		t.Errorf("Expected %d unique room IDs, got %d", numRooms, len(roomIDs))
	}
}

func TestRoomManagerReference(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom()

	if room.Manager == nil {
		t.Error("Room should have a reference to its manager")
	}

	if room.Manager != rm {
		t.Error("Room's manager reference should point to the correct manager")
	}
}
