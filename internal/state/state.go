package state

import (
	config "github.com/alancorleto/blog-aggregator/internal/config"
	database "github.com/alancorleto/blog-aggregator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
