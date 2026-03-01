package modules

import (
	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureAuth(registry types.ServiceRegistry) error {
	registry.Register(usecases.NewAuthUseCase, types.LifetimeSingleton)
	registry.Register(presenters.NewAuthPresenter, types.LifetimeSingleton)

	return nil
}
