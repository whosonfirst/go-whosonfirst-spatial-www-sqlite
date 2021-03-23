package main

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/server"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"log"
)

func main() {

	ctx := context.Background()

	app, err := server.NewHTTPServerApplication(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
