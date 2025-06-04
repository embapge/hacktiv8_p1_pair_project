package main

import (
	"context"
	"pairproject/cli"
	"pairproject/config"
)

func main() {
	db := config.InitDB()
	defer db.Close()
	ctx := context.Background()

	cli := cli.NewCLIHandler(db, ctx)
	cli.Menu()
}
