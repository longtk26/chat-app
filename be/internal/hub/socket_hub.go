package hub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/v3/socketio"
)

// wsEvent is the wire format for all server-to-client messages.
type wsEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// SocketHub manages conversation rooms and socket connections.
type SocketHub struct {
	mu                   sync.RWMutex
	rooms                map[string]map[string]bool // conversationID -> socketUUID set
	socketToConversation map[string]string          // socketUUID -> conversationID
	socketToUser         map[string]string          // socketUUID -> userID
}

func NewSocketHub() *SocketHub {
	return &SocketHub{
		rooms:                make(map[string]map[string]bool),
		socketToConversation: make(map[string]string),
		socketToUser:         make(map[string]string),
	}
}

func (h *SocketHub) JoinRoom(socketUUID, userID, conversationID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.leaveCurrentRoom(socketUUID)

	if h.rooms[conversationID] == nil {
		h.rooms[conversationID] = make(map[string]bool)
	}
	h.rooms[conversationID][socketUUID] = true
	h.socketToConversation[socketUUID] = conversationID
	h.socketToUser[socketUUID] = userID
}

func (h *SocketHub) LeaveRoom(socketUUID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.leaveCurrentRoom(socketUUID)
	delete(h.socketToUser, socketUUID)
}

// leaveCurrentRoom must be called with lock held.
func (h *SocketHub) leaveCurrentRoom(socketUUID string) {
	if convID, ok := h.socketToConversation[socketUUID]; ok {
		if room, exists := h.rooms[convID]; exists {
			delete(room, socketUUID)
			if len(room) == 0 {
				delete(h.rooms, convID)
			}
		}
		delete(h.socketToConversation, socketUUID)
	}
}

func (h *SocketHub) GetConversationID(socketUUID string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.socketToConversation[socketUUID]
}

// GetSocketUUIDByUserID finds the socket UUID for a user currently in the given conversation room.
func (h *SocketHub) GetSocketUUIDByUserID(conversationID, userID string) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	room, exists := h.rooms[conversationID]
	if !exists {
		return ""
	}
	for uuid := range room {
		if h.socketToUser[uuid] == userID {
			return uuid
		}
	}
	return ""
}

// BroadcastToRoom sends a named event to all sockets in the conversation room except excludeUUID.
// Pass excludeUUID="" to send to everyone in the room.
func (h *SocketHub) BroadcastToRoom(conversationID, excludeUUID, event string, payload interface{}) {
	h.mu.RLock()
	room, exists := h.rooms[conversationID]
	if !exists {
		h.mu.RUnlock()
		return
	}
	targets := make([]string, 0, len(room))
	for uuid := range room {
		if uuid != excludeUUID {
			targets = append(targets, uuid)
		}
	}
	h.mu.RUnlock()

	data, err := json.Marshal(wsEvent{Event: event, Data: payload})
	if err != nil {
		fmt.Printf("hub: failed to marshal event %q: %v\n", event, err)
		return
	}

	for _, uuid := range targets {
		if err := socketio.EmitTo(uuid, data, socketio.TextMessage); err != nil {
			fmt.Printf("hub: failed to emit to %s: %v\n", uuid, err)
		}
	}
}
