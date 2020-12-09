package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/server"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/flags"	
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"log"
)

func main() {

	ctx := context.Background()

	fs, err := spatial_flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	spatial_flags.Parse(fs)

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}

}
