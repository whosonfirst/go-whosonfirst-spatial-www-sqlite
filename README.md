# go-whosonfirst-spatial-http-sqlite

## Important

This is work in progress and not properly documented yet.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
$> ./bin/server -h
  -custom-placetypes string
    	...
  -custom-placetypes-source string
    	...
  -enable-candidates
    	Enable the /candidates endpoint to return candidate bounding boxes (as GeoJSON) for requests.
  -enable-custom-placetypes
    	...
  -enable-geojson
    	Allow users to request GeoJSON FeatureCollection formatted responses.
  -enable-properties
    	Enable support for 'properties' parameters in queries.
  -enable-www
    	Enable the interactive /debug endpoint to query points and display results.
  -exclude value
    	Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.
  -geojson-reader-uri string
    	A valid whosonfirst/go-reader.Reader URI. Required if the -enable-geojson or -enable-www flags are set.
  -initial-latitude float
    	... (default 37.616906)
  -initial-longitude float
    	... (default -122.386665)
  -initial-zoom int
    	... (default 13)
  -is-wof
    	Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents. (default true)
  -mode string
    	Valid modes are: . (default "repo://")
  -nextzen-apikey string
    	A valid Nextzen API key
  -nextzen-style-url string
    	... (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	... (default "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -properties-reader-uri string
    	Valid options are: [sqlite://]
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -setenv
    	Set flags from environment variables.
  -spatial-database-uri string
    	Valid options are: [sqlite://] (default "sqlite://")
  -static-prefix string
    	Prepend this prefix to URLs for static assets.
  -strict
    	Be strict about flags and fail if any are missing or deprecated flags are used.
  -templates string
    	An optional string for local templates. This is anything that can be read by the 'templates.ParseGlob' method.
  -verbose
    	Be chatty.
  -www-path string
    	The URL path for the interactive debug endpoint. (default "/debug")
```

For example:

```
$> bin/server \
	-enable-www \	
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=/usr/local/data/sfomuseum-data-architecture.db' \		
	-nextzen-apikey {NEXTZEN_APIKEY} \
	-mode directory://
```

A couple things to note:

* The SQLite databases specified in the `sqlite:///?dsn` string are expected to minimally contain the `rtree` and `spr` and `properties` tables confirming to the schemas defined in the [go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features). They are typically produced by the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package. See the documentation in the [go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package for details.

* The `-geojson-reader-uri` flag, and GeoJSON output for spatial queries in general, is discussed in detail below.

* Do you notice the way we are passing in a `-mode directory://` flag? This should only be necessary if we are generating, or updating, a SQLite database when the `server` tool starts up. As of this writing it is always necessary even though it doesn't do anything. There is an [open ticket](https://github.com/whosonfirst/go-whosonfirst-spatial/issues/14) to address this.

When you visit `http://localhost:8080` in your web browser you should see something like this:

![](docs/images/server.png)

If you don't need, or want, to expose a user-facing interface simply remove the `-enable-www` and `-nextzen-apikey` flags. For example:

```
$> bin/server \
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-enable-properties \	
	-mode directory://
```

And then to query the point-in-polygon API you would do something like this:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383'

{
  "places": [
    {
      "wof:id": 1360665043,
      "wof:parent_id": -1,
      "wof:name": "Central Parking Garage",
      "wof:placetype": "wing",
      "wof:country": "US",
      "wof:repo": "sfomuseum-data-architecture",
      "wof:path": "136/066/504/3/1360665043.geojson",
      "wof:superseded_by": [],
      "wof:supersedes": [
        1360665035
      ],
      "mz:uri": "https://data.whosonfirst.org/136/066/504/3/1360665043.geojson",
      "mz:latitude": 37.616332,
      "mz:longitude": -122.386047,
      "mz:min_latitude": 37.61498599208708,
      "mz:min_longitude": -122.38779093748578,
      "mz:max_latitude": 37.61767331604971,
      "mz:max_longitude": -122.38429192207244,
      "mz:is_current": 0,
      "mz:is_ceased": 1,
      "mz:is_deprecated": 0,
      "mz:is_superseded": 0,
      "mz:is_superseding": 1,
      "wof:lastmodified": 1547232156
    }
    ... and so on
}    
```

By default, results are returned as a list of ["standard places response"](https://github.com/whosonfirst/go-whosonfirst-spr/) (SPR) elements. You can also return results as a GeoJSON `FeatureCollection` by passing the `-enable-geojson` flag to the server and including a `format=geojson` query parameter with requests. For example:


```
$> bin/server \
	-enable-geojson \
	-enable-properties \	
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=/usr/local/data/sfomuseum-data-architecture.db' \		
	-mode directory://
```

And then:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383&format=geojson'

{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "MultiPolygon",
        "coordinates": [ ...omitted for the sake of brevity ]
      },
      "properties": {
        "mz:is_ceased": 1,
        "mz:is_current": 0,
        "mz:is_deprecated": 0,
        "mz:is_superseded": 0,
        "mz:is_superseding": 1,
        "mz:latitude": 37.616332,
        "mz:longitude": -122.386047,
        "mz:max_latitude": 37.61767331604971,
        "mz:max_longitude": -122.38429192207244,
        "mz:min_latitude": 37.61498599208708,
        "mz:min_longitude": -122.38779093748578,
        "mz:uri": "https://data.whosonfirst.org/136/066/504/3/1360665043.geojson",
        "wof:country": "US",
        "wof:id": 1360665043,
        "wof:lastmodified": 1547232156,
        "wof:name": "Central Parking Garage",
        "wof:parent_id": -1,
        "wof:path": "136/066/504/3/1360665043.geojson",
        "wof:placetype": "wing",
        "wof:repo": "sfomuseum-data-architecture",
        "wof:superseded_by": [],
        "wof:supersedes": [
          1360665035
        ]
      }
    }
    ... and so on
  ]
}  
```

If you are returning results as a GeoJSON `FeatureCollection` you may also request additional properties be appended by specifying them as a comma-separated list in the `?properties=` parameter. For example:

```
$> curl -s http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383&format=geojson&properties=sfomuseum:*
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "MultiPolygon",
        "coordinates": [ ... ]
      },
      "properties": {
        "mz:is_ceased": 1,
        "mz:is_current": 0,
        "mz:is_deprecated": 0,
        "mz:is_superseded": 1,
        "mz:is_superseding": 1,
        "mz:latitude": 37.617037,
        "mz:longitude": -122.385975,
        "mz:max_latitude": 37.62120978585632,
        "mz:max_longitude": -122.38125166743595,
        "mz:min_latitude": 37.61220882045874,
        "mz:min_longitude": -122.39033463643914,
        "mz:uri": "https://data.whosonfirst.org/115/939/632/7/1159396327.geojson",
        "sfomuseum:building_id": "SFO",
        "sfomuseum:is_sfo": 1,
        "sfomuseum:placetype": "building",
        "wof:country": "US",
        "wof:id": 1159396327,
        "wof:lastmodified": 1547232162,
        "wof:name": "SFO Terminal Complex",
        "wof:parent_id": 102527513,
        "wof:path": "115/939/632/7/1159396327.geojson",
        "wof:placetype": "building",
        "wof:repo": "sfomuseum-data-architecture",
        "wof:superseded_by": [
          1159554801
        ],
        "wof:supersedes": [
          1159396331
        ]
      }
    }... and so on
  ]
}
```

To return just the `properties` dictionary for results pass along the `?format=properties` parameter. For example:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383&format=properties&properties=sfomuseum:*' | jq

{
  "properties": [
    {
      "mz:is_ceased": 1,
      "mz:is_current": 1,
      "mz:is_deprecated": 0,
      "mz:is_superseded": 0,
      "mz:is_superseding": 0,
      "mz:latitude": 37.616359,
      "mz:longitude": -122.386105,
      "mz:max_latitude": -122.38789520924678,
      "mz:max_longitude": -122.38437148963614,
      "mz:min_latitude": 37.61497255972697,
      "mz:min_longitude": 37.616359,
      "sfomuseum:placetype": "garage",
      "wof:country": "US",
      "wof:id": "1477856011",
      "wof:lastmodified": 1568838528,
      "wof:name": "Central Parking Garage",
      "wof:parent_id": "102527513",
      "wof:path": "147/785/601/1/1477856011.geojson",
      "wof:placetype": "building",
      "wof:repo": "sfomuseum-data-architecture"
    }
  ]
}
```

### GeoJSON

GeoJSON output is produced by transforming a [go-whosonfirst-spr](https://github.com/whosonfirst/go-whosonfirst-spr) `StandardPlacesResults` instance, as returned by the [go-whosonfirst-spatial](https://github.com/whosonfirst/go-whosonfirst-spatial) `database.PointInPolygon` interface method, in to a GeoJSON FeatureCollection using the [go-whosonfirst-spr-geojson](https://github.com/whosonfirst/go-whosonfirst-spr-geojson) package.

The `go-whosonfirst-spr-geojson` package in turn uses the [whosonfirst/go-reader](https://github.com/whosonfirst/go-reader) packages to retrieve WOF records. There are two readers that are available by default with this tool:

* A reader for a directory containing WOF records on the local file system. This is part of the [whosonfirst/go-reader](https://github.com/whosonfirst/go-reader) package.
* A reader for a SQLite database with a table that has indexed WOF records. This is part of the [whosonfirst/go-reader-sql-database](https://github.com/whosonfirst/go-reader-sql-database) package.

Here's how you might use the local file system reader:

```
$> bin/server \
	-enable-geojson \	
	-geojson-reader-uri 'fs:///usr/local/data/sfomuseum-data-architecture/data' \		
```

And here's how you might use the SQLite reader:

```
$> bin/server \
	-enable-geojson \	
	-geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=/usr/local/data/sfomuseum-data-architecture.db' \		
```

The URI syntax for the SQLite reader is different than that used for the `-spatial-database-uri` or `-properties-reader-uri` flags because the reader itself is not specific to SQLite but designed to work with any valid [database/sql](https://golang.org/pkg/database/sql/) driver. The semantics of the `-geojson-reader-uri` flag are:

```
sql://{SQL_DRIVER}/{DATABASE_TABLE}/{ID_KEY}/{GEOJSON_COLUMN}?dsn={DATABASE_DSN}'
```

In order to account for query results that may contain alternate geometries there is some extra work that needs to be done in code, assigning custom functions for the `go-reader-sql-database.URI_READFUNC` and `go-reader-sql-database.URI_QUERYFUNC` properties, in order to extend the default `{ID_KEY}={VALUE}` query. This is done in [cmd/server/main.go](cmd/server/main.go).

In the example above we are assuming a table called `geojson` as defined by the [go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#geojson) package.

Here's an example of the creating a compatible SQLite database, with support for GeoJSON results, for all the [administative data in Canada](https://github.com/whosonfirst-data/whosonfirst-data-admin-ca) using the `wof-sqlite-index-features` tool which is part of the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package:

```
$> ./bin/wof-sqlite-index-features \
	-index-alt-files \
	-rtree \
	-spr \
	-properties \
	-geojson \	
	-dsn /usr/local/data/ca.db \
	-mode repo:// \
	/usr/local/data/whosonfirst-data-admin-ca/
```

And then to start the `server` tool using this new database you would do:

```
$> bin/server \
	-enable-www \	
	-geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=/usr/local/data/ca.db' \	
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/ca.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/ca.db' \
	-nextzen-apikey {NEXTZEN_APIKEY} \
	-mode directory://
```

GeoJSON results are noticeably slower than SPR results, particularly for features with large geometries like countries. GeoJSON output is provided as a convenience but is not recommended for public-facing scenarios where speed is a critical factor.

### Indexing "plain old" GeoJSON

There is early support for indexing "plain old" GeoJSON, as in GeoJSON documents that do not following the naming conventions for properties that Who's On First documents use. It is very likely there are still bugs or subtle gotchas.

For example, here's how we could index and serve a GeoJSON FeatureCollection of building footprints, using an in-memory SQLite database:

```
$> ./bin/server \
	-spatial-database-uri 'sqlite:///?dsn=:memory:' \
	-mode featurecollection://
	/usr/local/data/footprint.geojson
```

And then:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61686957521345&longitude=-122.3903158758416' \

| jq '.["places"][]["wof:id"]'

"1031"
"1015"
"1014"
```

If you want to enable the `properties` output format you would do this:

```
$> ./bin/server \
   -enable-properties \
   -index-properties \   
   -spatial-database-uri 'sqlite:///?dsn=:memory:' \
   -properties-reader-uri 'sqlite:///?dsn=:memory:' \
   -mode featurecollection://
   /usr/local/data/footprint.geojson
```

And then:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61686957521345&longitude=-122.3903158758416&format=properties&properties=BUILDING' \

| jq '.["properties"][]["BUILDING"]'

"400"
"400"
"100"
"100"
"100"
```

If you want to enable the `geojson` output format you will need to create a local SQLite database, rather than an in-memory database. That's because the in-memory database created by the `-spatial-database-uri` flag is different from the in-memory database created by the `-geojson-reader-uri` flag. For example:

```
$> ./bin/server \
	-enable-geojson \
	-spatial-database-uri 'sqlite:///?dsn=test4.db&index-geojson=true' \
	-geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=test4.db' \
	-mode featurecollection://
	/usr/local/data/footprint.geojson
```

And then:

```
$> curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.61686957521345&longitude=-122.3903158758416&format=geojson' \

| jq '.["features"][]["properties"]["NAME"]'

"Terminal 3"
"International Terminal"
"International Terminal"
"International Terminal"
"International Terminal"
```

Under the hood the code is using the [go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features) package to index the "plain old" GeoJSON documents. You can also index your "plain old" GeoJSON documents ahead of time (using the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package) to speed up start up times, as demonstrated in the examples at the top of this document.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-http
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spr
* https://github.com/whosonfirst/go-whosonfirst-spr-geojson
* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-reader-database-sql
* https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index
* https://developers.nextzen.org/
