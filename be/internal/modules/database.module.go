package modules

import (
	"fmt"

	"github.com/longtk26/chat-app/configs"
	db "github.com/longtk26/chat-app/infrastructure/db"
	"github.com/longtk26/chat-app/infrastructure/db/sqlc/out"
	repositories "github.com/longtk26/chat-app/infrastructure/repo"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureDatabase(registry types.ServiceRegistry) error {
	appConfig := configs.LoadAppConfig()

	pool, err := db.NewPostgresPool(appConfig)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL pool: %w", err)
	}

	sqlcQueries := out.New(pool)
	registry.Register(func() *out.Queries {
		fmt.Println("Creating sqlc Queries...")
		return sqlcQueries
	}, types.LifetimeSingleton)

	// Register repository implementations
	registry.Register(repositories.NewUserRepository, types.LifetimeSingleton)
	registry.Register(repositories.NewConversationsRepo, types.LifetimeSingleton)
	registry.Register(repositories.NewMessagesRepo, types.LifetimeSingleton)

	return nil
}
