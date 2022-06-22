package example

import (
	"context"
	"gitlab.falabella.com/falabella-retail/txd/boss/libraries/golang/persistence/client"
	"gitlab.falabella.com/falabella-retail/txd/boss/libraries/golang/persistence/prop"
	"gitlab.falabella.com/falabella-retail/txd/boss/libraries/golang/persistence/repo"
)

func main() {
	// Config
	var prop = &prop.PostgresProp{}
	pgClient := client.GetPostgresClient(context.Background(), prop)

	var repo1 repo.DBQuery
	repo1 = repo.GetPgRepo(pgClient)

	// Secondary
	repo1.ExecuteQuery("", nil)
}
