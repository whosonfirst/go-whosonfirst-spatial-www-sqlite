package main

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/server"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-database-sqlite"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"	
	"log"
)

func main() {

	ctx := context.Background()

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	flags.Parse(fs)

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}

}
