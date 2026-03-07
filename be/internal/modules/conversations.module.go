package modules

import (
	"github.com/longtk26/chat-app/internal/presenters"
	"github.com/longtk26/chat-app/internal/usecases"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureConversations(registry types.ServiceRegistry) error {
	registry.Register(usecases.NewConversationsUsecase, types.LifetimeSingleton)
	registry.Register(presenters.NewConversationsPresenter, types.LifetimeSingleton)

	return nil
}
