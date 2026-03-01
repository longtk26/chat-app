package modules

import (
	"github.com/longtk26/chat-app/configs"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureConfig(registry types.ServiceRegistry) error {
	appConfig := configs.LoadAppConfig()

	registry.Register(func() configs.AppConfig {
		return appConfig
	}, types.LifetimeSingleton)

	return nil
}
