package modules

import (
	"context"
	"fmt"

	"github.com/longtk26/chat-app/internal/routers"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func RegisterRouteHandlers(registry types.ServiceRegistry) error {
	fmt.Println("Registering route handlers...")

	if err := features.RegisterList[routers.RouteHandler](context.Background(), registry); err != nil {
		fmt.Println(err)
		return err
	}

	registration.RegisterTransient(registry, routers.NewAuthRoute)

	return nil
}
