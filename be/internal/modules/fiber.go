package modules

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

var _ types.ModuleFunc = ConfigureFiber

func ConfigureFiber(registry types.ServiceRegistry) error {
	registration.RegisterInstance(registry, fiber.Config{
		AppName:   "chat-app",
		Immutable: true,
	})

	registry.Register(newFiber, types.LifetimeSingleton)
	registry.RegisterModule(RegisterRouteHandlers)
	return nil
}

func newFiber(config fiber.Config) *fiber.App {
	return fiber.New(config)
}
