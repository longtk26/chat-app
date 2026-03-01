package main

import (
	"context"

	"github.com/longtk26/chat-app/internal"
	"github.com/longtk26/chat-app/internal/modules"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

func main() {
	ctx := context.Background()

	bootstrap.RunParsleyApplication(ctx, internal.NewApp,
		modules.ConfigureConfig,
		modules.ConfigureFiber,
		modules.ConfigureDatabase,
		modules.ConfigureAuth,
	)
}
