package modules

import (
	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureUsers(registry types.ServiceRegistry) error {
	registry.Register(usecases.NewUsersUseCase, types.LifetimeSingleton)
	registry.Register(presenters.NewUsersPresenter, types.LifetimeSingleton)

	return nil
}
