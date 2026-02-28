package modules

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/longtk26/chat-app/configs"
	db "github.com/longtk26/chat-app/infrastructure/db"
	"github.com/longtk26/chat-app/infrastructure/db/sqlc/out"
	repositories "github.com/longtk26/chat-app/infrastructure/repo"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureDatabase(registry types.ServiceRegistry) error {
	// Register DB config
	registration.RegisterInstance(registry, configs.LoadDBConfig())

	// Register pgxpool.Pool (singleton)
	registration.RegisterSingleton(registry, func(cfg configs.DBConfig) (*pgxpool.Pool, error) {
		return db.NewPostgresPool(cfg)
	})

	// Register sqlc Queries
	registration.RegisterSingleton(registry, func(pool *pgxpool.Pool) *out.Queries {
		return out.New(pool)
	})

	// Register repository implementations
	registration.RegisterSingleton(registry, repositories.NewUserRepository)

	return nil
}
