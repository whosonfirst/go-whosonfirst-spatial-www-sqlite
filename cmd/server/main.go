package main

import ()

import (
	"context"
	sql_reader "github.com/whosonfirst/go-reader-database-sql"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/server"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"log"
	"strconv"
)

func main() {

	sql_reader.URI_READFUNC = func(uri_str string) (string, error) {

		id, _, err := uri.ParseURI(uri_str)

		if err != nil {
			return "", err
		}

		str_id := strconv.FormatInt(id, 10)
		return str_id, nil
	}

	sql_reader.URI_QUERYFUNC = func(uri_str string) (string, []interface{}, error) {

		_, uri_args, err := uri.ParseURI(uri_str)

		if err != nil {
			return "", nil, err
		}

		if !uri_args.IsAlternate {
			return "", nil, nil
		}

		alt_label, err := uri_args.AltGeom.String()

		if err != nil {
			return "", nil, err
		}

		where := "alt_label = ?"

		args := []interface{}{
			alt_label,
		}

		return where, args, nil
	}

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
