package routers

import (
	"fmt"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v3"
	"github.com/longtk26/chat-app/internal/presenters"
)

var _ RouteHandler = &SocketRoute{}

type SocketRoute struct {
	p *presenters.SocketPresenter
}

func NewSocketRoute(p *presenters.SocketPresenter) RouteHandler {
	return &SocketRoute{
		p: p,
	}
}

func (r *SocketRoute) Register(app *fiber.App) {
	fmt.Println("Registering SocketRoute...")
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Println("connected:")
		r.p.HandleConnect()
	})
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		fmt.Println("disconnected:")
	})
}
