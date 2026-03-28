package modules

import (
	"github.com/longtk26/chat-app/internal/hub"
	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureSocket(registry types.ServiceRegistry) error {
	registry.Register(hub.NewSocketHub, types.LifetimeSingleton)
	registry.Register(usecases.NewSocketUsecase, types.LifetimeSingleton)
	registry.Register(presenters.NewSocketPresenter, types.LifetimeSingleton)

	return nil
}
