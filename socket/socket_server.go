package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ChatRoom struct {
	ID           string `gorm:"primaryKey"`
	FirstUserID  string
	SecondUserID string
	IsEnd        bool
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients     = make(map[string]*websocket.Conn) // Map of user ID to WebSocket connections
	clientsM    sync.Mutex                         // Mutex for handling concurrent access to clients map
	redisClient *redis.Client
	db          *gorm.DB
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
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// Add user to the active clients map
	clientsM.Lock()
	clients[userID] = ws // Store user ID with its WebSocket connection
	clientsM.Unlock()

	// Check if there's an existing room for the user in the database
	room := new(ChatRoom)
	if err := db.Where("(first_user_id = ? OR second_user_id = ?) AND is_end = false", userID, userID).First(&room).Error; err == nil {
		// If room exists, handle messaging for the room
		log.Printf("User %s joined an existing room %s", userID, room.ID)
		handleRoomMessaging(ws, room, userID)
		return
	}

	// No existing room: try to pop from Redis queue to match with another user
	ctx := context.Background()
	matchedUserID, err := redisClient.SPop(ctx, "waitingQueue").Result()
	if err != nil || matchedUserID == "" || matchedUserID == userID {
		// If no match found or the user is the only one in the queue, add them to the queue and wait
		redisClient.SAdd(ctx, "waitingQueue", userID)
		log.Printf("User %s added to waiting queue", userID)
		waitForMatch(ws, userID)
		return
	}

	// Create a new chat room if match is found
	roomID := uuid.New().String()
	newRoom := &ChatRoom{
		ID:           roomID,
		FirstUserID:  userID,
		SecondUserID: matchedUserID,
		IsEnd:        false,
	}
	if err := db.Create(newRoom).Error; err != nil {
		log.Printf("Failed to create new room in the database: %v", err)
		return
	}

	log.Printf("Created new room %s for users %s and %s", roomID, userID, matchedUserID)
	handleRoomMessaging(ws, newRoom, userID)
}

// handleRoomMessaging listens for messages from the user and forwards them to their chat partner
func handleRoomMessaging(ws *websocket.Conn, room *ChatRoom, senderID string) {
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		sendMessageToPartner(room, senderID, message)
	}
}

// sendMessageToPartner forwards a message to the chat partner
func sendMessageToPartner(room *ChatRoom, senderID string, message []byte) {
	var partnerID string
	if senderID == room.FirstUserID {
		partnerID = room.SecondUserID
	} else {
		partnerID = room.FirstUserID
	}

	clientsM.Lock()
	partnerConn, ok := clients[partnerID]
	clientsM.Unlock()

	if ok {
		partnerConn.WriteMessage(websocket.TextMessage, message)
	}
}

// endChatRoom ends the chat room by marking it as ended in the database
func endChatRoom(room *ChatRoom, userID string) {
	// If one user leaves, you might want to mark the chat room as ended or remove their connection
	if room.FirstUserID == userID || room.SecondUserID == userID {
		room.IsEnd = true
		db.Save(room) // Mark the room as ended in the database
	}
}

// waitForMatch does nothing while the user waits for a match.
func waitForMatch(ws *websocket.Conn, userID string) {
	log.Printf("User %s is waiting for a match", userID)

	// As per your logic, if a user is waiting, you don't need to pop them here, they will remain in the queue until matched.
}

func main() {
	err := godotenv.Load("../.env")
    apiEnv := os.Getenv("ENV")
    if err != nil && apiEnv == "" {
        log.Println("fail to load env", err)
    }
	config.LoadEnv()
	db = config.NewPostgresql()
	redisClient = config.NewRedis()

	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server starting on http://localhost:8000/ws")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
