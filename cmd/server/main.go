package main

import (
	"context"
	www_flags "github.com/whosonfirst/go-whosonfirst-spatial-www/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/server"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"log"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	ctx := context.Background()

	fs, err := spatial_flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = spatial_flags.AppendIndexingFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = www_flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	flagset.Parse(fs)

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}

}
