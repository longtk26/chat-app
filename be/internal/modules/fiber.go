package modules

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureFiber(registry types.ServiceRegistry) error {
	registration.RegisterInstance(registry, fiber.Config{
		AppName: "chat-app",
	})

	registry.Register(newFiber, types.LifetimeSingleton)

	return nil
}

func newFiber(config fiber.Config) *fiber.App {
	return fiber.New(config)
}
