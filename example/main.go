package example

import (
	"context"
	"gitlab.falabella.com/falabella-retail/txd/boss/libraries/golang/persistence/repo"
)

func main() {
	// Config
	var prop = &repo.PostgresProp{}
	pgClient := repo.GetPostgresClient(context.Background(), prop)

	var repo1 repo.DBQuery
	repo1 = repo.GetPgRepo(pgClient)

	// Secondary
	repo1.ExecuteQuery("", nil)
}
