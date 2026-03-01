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
	var configs = configs.LoadDBConfig()

	pool, err := db.NewPostgresPool(configs)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL connection pool: %w", err)
	}

	// Register sqlc Queries
	queries := out.New(pool)
	registry.Register(func() *out.Queries {
		fmt.Println("Creating sqlc Queries...")
		return queries
	}, types.LifetimeSingleton)

	// Register repository implementations
	registry.Register(repositories.NewUserRepository, types.LifetimeSingleton)

	return nil
}
