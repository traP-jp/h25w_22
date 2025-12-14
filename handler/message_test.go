package handler

import (
	"encoding/json"
	"testing"
)

func TestMarshalEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventName string
		payload   interface{}
		wantError bool
	}{
		{
			"Matched event",
			"MATCHED",
			MatchedPayload{},
			false,
		},
		{
			"Turn start event",
			"TURN_START",
			TurnStartPayload{TurnCount: 1, Role: "ATTACK"},
			false,
		},
		{
			"Action result event",
			"ACTION_RESULT",
			ActionResultPayload{Hit: true, Damage: 20, CurrentHP: 80},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalEvent(tt.eventName, tt.payload)
			if (err != nil) != tt.wantError {
				t.Errorf("MarshalEvent() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				var event Event
				if err := json.Unmarshal(data, &event); err != nil {
					t.Errorf("Failed to unmarshal event: %v", err)
				}
				if event.Event != tt.eventName {
					t.Errorf("Event name = %s, want %s", event.Event, tt.eventName)
				}
			}
		})
	}
}
