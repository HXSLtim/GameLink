package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// Handler handles WebSocket connections.
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// ServeWS handles websocket requests from the peer.
// @Summary WebSocket endpoint for real-time monitoring
// @Description Establishes a WebSocket connection for real-time system monitoring
// @Tags WebSocket
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/ws/monitor [get]
func (h *Handler) ServeWS(c *gin.Context) {
	// Extract user info from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userID uint64
	switch v := userIDValue.(type) {
	case uint64:
		userID = v
	case float64:
		userID = uint64(v)
	case int:
		userID = uint64(v)
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}
		userID = parsed
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
		return
	}

	role, _ := c.Get("user_role")
	roleStr, ok := role.(string)
	if !ok {
		roleStr = "user"
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create client and register with hub
	client := NewClient(h.hub, conn, userID, roleStr)
	h.hub.register <- client

	log.Printf("WebSocket client connected: UserID=%d, Role=%s, RemoteAddr=%s",
		userID, roleStr, conn.RemoteAddr().String())

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// GetHub returns the WebSocket hub.
func (h *Handler) GetHub() *Hub {
	return h.hub
}

// RegisterRoutes registers WebSocket routes with the router.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/ws/monitor", h.ServeWS)
}
