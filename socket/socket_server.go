package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	ID           string
	FirstUserID  string
	SecondUserID string
	FirstConn    *websocket.Conn
	SecondConn   *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Adjust this for production
		},
	}

	rooms      = make(map[string]*ChatRoom) // Room ID mapped to ChatRoom
	userToRoom = make(map[string]string)    // Map user ID to room ID
	roomsM     sync.Mutex                   // Mutex to protect room creation and deletion
)

func validateToken(token string) (string, error) {
	claims, err := utilities.ValidateJWT(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	return claims.UserID, nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := validateToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade: %v", err)
		return
	}
	defer ws.Close()

	room := findOrCreateRoom(userID, ws)
	if room == nil {
		log.Println("No available rooms, waiting for another user...")
		return
	}

	for {
		messageType, message, err := ws.ReadMessage()
		log.Println(message)
		log.Println(messageType)
		if err != nil {
			removeUserFromRoom(userID, room)
			break
		}
		sendMessageToPartner(room, userID, messageType, message)
	}
}

func findOrCreateRoom(userID string, ws *websocket.Conn) *ChatRoom {
	roomsM.Lock()
	defer roomsM.Unlock()

	if roomID, exists := userToRoom[userID]; exists {
		room := rooms[roomID]
		if room != nil {
			log.Printf("Reconnected user %s to existing room %s", userID, room.ID)
			return room
		}
	}

	for _, room := range rooms {
		if room.SecondUserID == "" {
			room.SecondUserID = userID
			room.SecondConn = ws
			userToRoom[userID] = room.ID
			log.Printf("User %s joined room with %s", userID, room.FirstUserID)
			return room
		}
	}

	roomID := uuid.New().String()
	newRoom := &ChatRoom{
		ID:          roomID,
		FirstUserID: userID,
		FirstConn:   ws,
	}
	rooms[roomID] = newRoom
	userToRoom[userID] = roomID
	log.Printf("Created new room %s for user %s", roomID, userID)
	return newRoom
}

func sendMessageToPartner(room *ChatRoom, senderID string, messageType int, message []byte) {
	if room.FirstUserID == senderID && room.SecondConn != nil {
		room.SecondConn.WriteMessage(messageType, message)
	} else if room.SecondUserID == senderID && room.FirstConn != nil {
		room.FirstConn.WriteMessage(messageType, message)
	}
}

func removeUserFromRoom(userID string, room *ChatRoom) {
	roomsM.Lock()
	defer roomsM.Unlock()

	if room.FirstUserID == userID {
		room.FirstConn = nil
	} else {
		room.SecondConn = nil
	}

	if room.FirstConn == nil && room.SecondConn == nil {
		delete(rooms, room.ID)
		delete(userToRoom, room.FirstUserID)
		delete(userToRoom, room.SecondUserID)
		log.Printf("Deleted room %s", room.ID)
	}
}

func main() {


	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server starting on http://localhost:8000/ws")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
