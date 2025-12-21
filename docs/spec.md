# **Project Specification: "Date vs. Jam" Online Card Game Server**

## **1. Project Overview**

Build a WebSocket-based game server in **Go (Golang)** for a 1-on-1 asymmetric card game.

  - **Roles:**
      - **Dater (Defense):** Selects a date location and tries to survive with high HP to gain affection (Score).
      - **Jammer (Attack):** Tries to reduce the Dater's HP to 0 using a limited number of cards.
  - **Game Structure:** 4 Turns total. Each turn has a **Scenario/Setup Phase** and a **Battle Phase**.
  - **Anti-Cheat Policy:** **None.** The server trusts the client's timing and parameters.

## **2. Tech Stack**

  - **Language:** Go (Latest version).
  - **Routing:** Standard `net/http` with Go 1.22+ `ServeMux`.
  - **WebSocket:** `github.com/gorilla/websocket`.
  - **Architecture:** Standard library only (no external web frameworks).

## **3. Directory Structure**

```text
.
├── main.go             # Entry point. Calls handler.SetupRoutes()
├── game/
│   ├── card.go         # Card models & Master data
│   ├── logic.go        # Hit/Damage & Score calculations
│   └── state.go        # GameState & State Machine
└── handler/
    ├── manager.go      # RoomManager (Create/Delete rooms)
    ├── room.go         # Room logic (Broadcasts, Game loop)
    ├── client.go       # WebSocket Client (Read/Write pumps)
    ├── http.go         # HTTP Handlers & Route Setup
    └── message.go      # Protocol definitions (Events/Commands)
```

-----

## **4. Data & Logic Specifications**

### **A. Card Data (`game/card.go`)**

  - **HitRate:** `int` (1-100).
  - **Structure:**
    ```go
    type Card struct {
        ID      int
        Name    string
        Power   int    // Damage amount
        HitRate int    // Percentage (1-100)
        Type    string // "ATTACK" or "DEFENSE"
    }
    ```

### **B. Game State (`game/state.go`)**

  - **TurnCount:** Integer (1 to 4).
  - **Roles:** `ATTACK` and `DEFENSE` (swapped every turn).
  - **HP:** Current HP. Set by the Defender at the start of the turn.
  - **ScoreMultiplier:** Float. Set by the Defender (based on date location).
  - **CardsPlayed:** Integer. Tracks how many cards the **Attacker** has used this turn.
  - **PendingAttack:** Pointer to `Card`.

### **C. Logic & Constants (`game/logic.go`)**

  - **Constants:**

    ```go
    const BaseScoreConstant = 50 // Defined as a constant for easy adjustment
    const MaxCardsPerTurn = 4
    ```

  - **Hit/Damage Logic:**

      - `isHit := (rand.Intn(100) + 1) <= attackCard.HitRate`
      - If Hit: `HP -= attackCard.Power`.

  - **Score Calculation (Affection):**

      - This is calculated at the **end of the turn** if the Defender survives (HP \> 0).
      - **Formula:** `Score += BaseScoreConstant + (CurrentHP * ScoreMultiplier)`

-----

## **5. Protocol Definition**

### **Client -\> Server (Commands)**

*Format: Text, Space-separated* (`COMMAND ARG1 ARG2 ...`)

| Command | Arguments | Context |
| :--- | :--- | :--- |
| `SELECT_DATE` | `InitHP` (int) `Multiplier` (float) | Sent by **Defender** in Phase 2. Sets difficulty. |
| `READY` | None | Sent when client finishes Scenario or Result screen. |
| `PLAY_CARD` | `CardID` (int) | Sent when player selects a card. |

### **Server -\> Client (Events)**

*Format: JSON* (`{"event": "NAME", "payload": ...}`)

| Event Name | Key Payload Fields | Description |
| :--- | :--- | :--- |
| `MATCHED` | None | Both players joined. |
| `GAME_START` | None | Both players are READY after matching. |
| `TURN_START` | `role` | Start of a Turn. |
| `ACTION_START` | `maxHp` | Scenario/Setup finished. Battle begins. |
| `OPPONENT_PLAYED`| `cardId` | Attacker played. Notify Defender. |
| `ACTION_RESULT` | `hit`, `currentHp` | Result of a card exchange. |
| `TURN_END` | `reason` | Turn finished (HP=0 or Card Limit). |
| `GAME_RESULT` | `finalScores` | Game Over (after 4 turns). |
| `OPPONENT_DISCONNECTED` | None | Opponent has disconnected. Room will be terminated. |

-----

## **6. Detailed Game Flow (State Machine)**

Implement this flow in `handler/room.go`.

### **Phase 1: Matching**

1.  User A POSTs to `/rooms` -\> gets a 4-digit room ID (e.g., `0025`, `9877`).
2.  Users connect via WS -\> **Server** sends `MATCHED`.
3.  Clients send `READY` -\> **Server** sends `GAME_START`.

### **Phase 2: Turn Start & Date Selection**

1.  **Server** resets `CardsPlayed = 0` for the new turn.
2.  **Server** sends `TURN_START` (with `turnCount` and `role`).
3.  **Defender** sends `SELECT_DATE <HP> <Multiplier>`.
      - **Server:** Updates `GameState.HP` and `GameState.ScoreMultiplier`.
4.  Clients play scenario -\> Send `READY`.
5.  **Server** waits for READYs + Date Selection -\> Proceeds to Phase 3.

### **Phase 3: Battle Loop (Action)**

1.  **Server** sends `ACTION_START` (once per turn).

2.  **Card Exchange Loop:**

      * **Attacker** sends `PLAY_CARD <ID>`.
          - **Server** increments `GameState.CardsPlayed`.
          - **Server** sends `OPPONENT_PLAYED` to Defender.
      * **Defender** sends `PLAY_CARD <ID>`.
          - **Server** calculates Hit/Damage.
          - **Server** sends `ACTION_RESULT` to **both**.

3.  **Check Condition (Loop or End):**

      * **Condition A (Game Over for Turn):** If `HP <= 0`.
          - **Result:** Defender Lost. No Score added.
          - Proceed to **Phase 4** (Reason: "HP\_ZERO").
      * **Condition B (Success for Turn):** If `HP > 0` **AND** `CardsPlayed == 4` (Limit Reached).
          - **Result:** Defender Won.
          - **Calculate Score:** `Score += BaseScoreConstant + (HP * ScoreMultiplier)`.
          - Proceed to **Phase 4** (Reason: "LIMIT\_REACHED").
      * **Condition C (Continue):** If `HP > 0` **AND** `CardsPlayed < 4`.
          - Continue waiting for Attacker's `PLAY_CARD`.

### **Phase 4: Turn End & Resolution**

1.  **Server** sends `TURN_END`.
      - Payload should include `turnScore` (how much was added this turn) so the client can display it.
2.  Clients display the Result Screen.
3.  Clients send `READY` after closing the Result Screen.
4.  **Server** waits for 2 READYs.
5.  **Check Turn Count:**
      - **If Turn \< 4:** Increment TurnCount, Swap Roles, **Go to Phase 2**.
      - **If Turn == 4:** Proceed to **Phase 5**.

### **Phase 5: Game End**

1.  **Server** sends `GAME_RESULT` with final aggregated scores.
2.  Server closes connections.
3.  **Server deletes the room** from the room manager to free resources.

### **Phase 6: Client Disconnection Handling**

1.  **If a client disconnects** at any point during the game:
    - **Server** sends `OPPONENT_DISCONNECTED` event to the remaining player.
    - **Server** closes all remaining connections.
    - **Server deletes the room** from the room manager.
2.  This ensures rooms don't remain active when a player leaves unexpectedly.

-----

## **7. Implementation Details for `handler` Package**

### **`handler/room.go` -\> `HandleMessage`**

```go
func (r *Room) HandleMessage(client *Client, cmd string, args []string) {
    r.mu.Lock()
    defer r.mu.Unlock()

    switch cmd {
    case "SELECT_DATE":
        // ... (HP/Multiplier setup) ...

    case "READY":
        // ... (Phase transition logic) ...

    case "PLAY_CARD":
        // 1. Identify Role (Attacker vs Defender)
        // 2. If Attacker:
        //    - r.Game.PendingAttack = card
        //    - r.Game.CardsPlayed++  // <--- Increment here
        //    - Send OPPONENT_PLAYED
        // 3. If Defender:
        //    - Execute Logic (Hit/Damage)
        //    - Send ACTION_RESULT
        //    - Check End Conditions:
        //      if r.Game.HP <= 0 {
        //          r.endTurn("HP_ZERO")
        //      } else if r.Game.CardsPlayed >= 4 {
        //          // Calculate Score
        //          bonus := game.BaseScoreConstant + (float64(r.Game.HP) * r.Game.ScoreMultiplier)
        //          r.Game.Scores[game.RoleDefense] += int(bonus)
        //          r.endTurn("LIMIT_REACHED")
        //      }
    }
}
```
