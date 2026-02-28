package modules

import (
	"context"

	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/routers"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureAuth(registry types.ServiceRegistry) error {
	registration.RegisterSingleton(registry, usecases.NewAuthUseCase)
	registration.RegisterSingleton(registry, presenters.NewAuthPresenter)
	registration.RegisterTransient(registry, routers.NewAuthRoute)

	features.RegisterList[routers.RouteHandler](context.Background(), registry)

	return nil
}

