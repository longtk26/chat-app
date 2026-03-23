package modules

import (
	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureMessages(registry types.ServiceRegistry) error {
	registry.Register(usecases.NewMessagesUseCase, types.LifetimeSingleton)
	registry.Register(presenters.NewMessagesPresenter, types.LifetimeSingleton)

	return nil
}
