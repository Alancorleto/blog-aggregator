package state

import (
	config "github.com/alancorleto/gator/internal/config"
	database "github.com/alancorleto/gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
