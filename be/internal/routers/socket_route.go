package routers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/contrib/v3/socketio"
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/hub"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &SocketRoute{}

type SocketRoute struct {
	p   *presenters.SocketPresenter
	hub *hub.SocketHub
}

func NewSocketRoute(p *presenters.SocketPresenter, h *hub.SocketHub) RouteHandler {
	return &SocketRoute{p: p, hub: h}
}

// clientMessage is the wire format for all client-to-server messages:
// {"event": "join_conversation", "data": "conv-id"}
type clientMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func (r *SocketRoute) Register(app *fiber.App) {
	fmt.Println("Registering SocketRoute...")

	app.Get("/socket.io/*", socketio.New(func(kws *socketio.Websocket) {
		userID := kws.Query("user_id", "anonymous")

		if userID == "anonymous" {
			fmt.Printf("socket connection rejected: missing user_id\n")
			kws.Close()
			return
		}
		username := kws.Query("username", "")
		kws.SetAttribute("user_id", userID)
		kws.SetAttribute("username", username)
	}))

	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("connected: user=%s uuid=%s\n", ep.Kws.GetStringAttribute("user_id"), ep.SocketUUID)
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		convID := r.hub.GetConversationID(ep.SocketUUID)
		r.hub.LeaveRoom(ep.SocketUUID)
		if convID != "" {
			r.hub.BroadcastToRoom(convID, "", "user_offline", map[string]string{
				"user_id": userID,
			})
		}
		fmt.Printf("disconnected: user=%s uuid=%s\n", userID, ep.SocketUUID)
	})

	// All client messages arrive as EventMessage; parse and dispatch to named handlers.
	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		var msg clientMessage
		if err := json.Unmarshal(ep.Data, &msg); err != nil {
			return
		}
		ep.Kws.Fire(msg.Event, msg.Data)
	})

	// join_conversation: data is a JSON-encoded conversation ID string
	socketio.On("join_conversation", func(ep *socketio.EventPayload) {
		var convID string
		if err := json.Unmarshal(ep.Data, &convID); err != nil {
			return
		}
		userID := ep.Kws.GetStringAttribute("user_id")
		r.hub.JoinRoom(ep.SocketUUID, userID, convID)
		r.hub.BroadcastToRoom(convID, ep.SocketUUID, "user_active", map[string]string{
			"user_id": userID,
		})
		fmt.Printf("join_conversation: user=%s conv=%s\n", userID, convID)
	})

	// leave_conversation: data is a JSON-encoded conversation ID string
	socketio.On("leave_conversation", func(ep *socketio.EventPayload) {
		userID := ep.Kws.GetStringAttribute("user_id")
		convID := r.hub.GetConversationID(ep.SocketUUID)
		r.hub.LeaveRoom(ep.SocketUUID)
		if convID != "" {
			r.hub.BroadcastToRoom(convID, "", "user_offline", map[string]string{
				"user_id": userID,
			})
		}
	})

	// typing: data is a JSON-encoded conversation ID string
	socketio.On("typing", func(ep *socketio.EventPayload) {
		var convID string
		if err := json.Unmarshal(ep.Data, &convID); err != nil {
			return
		}
		r.hub.BroadcastToRoom(convID, ep.SocketUUID, "user_typing", map[string]string{
			"user_id":  ep.Kws.GetStringAttribute("user_id"),
			"username": ep.Kws.GetStringAttribute("username"),
		})
	})

	// stop_typing: data is a JSON-encoded conversation ID string
	socketio.On("stop_typing", func(ep *socketio.EventPayload) {
		var convID string
		if err := json.Unmarshal(ep.Data, &convID); err != nil {
			return
		}
		r.hub.BroadcastToRoom(convID, ep.SocketUUID, "user_stop_typing", map[string]string{
			"user_id": ep.Kws.GetStringAttribute("user_id"),
		})
	})
}
