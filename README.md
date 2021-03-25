# go-whosonfirst-spatial-http-sqlite

## Important

This is work in progress and not properly documented yet.

_Some of the documentation in out of date, specifically documentation for "plain old GeoJSON"._

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
    	A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-cors
    	Enable CORS headers for data-related and API handlers.
  -enable-custom-placetypes
    	Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-geojson
    	Enable GeoJSON output for point-in-polygon API calls.
  -enable-gzip
    	Enable gzip-encoding for data-related and API handlers.
  -enable-tangram
    	Use Tangram.js for rendering map tiles
  -enable-www
    	Enable the interactive /debug endpoint to query points and display results.
  -is-wof
    	Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported schemes are: directory://, featurecollection://, file://, filelist://, geojsonl://, repo://. (default "repo://")
  -leaflet-initial-latitude float
    	The initial latitude for map views to use. (default 37.616906)
  -leaflet-initial-longitude float
    	The initial longitude for map views to use. (default -122.386665)
  -leaflet-initial-zoom int
    	The initial zoom level for map views to use. (default 14)
  -leaflet-max-bounds string
    	An optional comma-separated bounding box ({MINX},{MINY},{MAXX},{MAXY}) to set the boundary for map views.
  -leaflet-tile-url string
    	A valid Leaflet (slippy map) tile template URL to use for rendering maps (if -enable-tangram is false)
  -nextzen-apikey string
    	A valid Nextzen API key
  -nextzen-style-url string
    	The URL for the style bundle file to use for maps rendered with Tangram.js (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	The URL for Nextzen tiles to use for maps rendered with Tangram.js (default "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -path-data string
    	The URL for data (GeoJSON) handler (default "/data")
  -path-ping string
    	The URL for the ping (health check) handler (default "/health/ping")
  -path-pip string
    	The URL for the point in polygon web handler (default "/point-in-polygon")
  -path-prefix string
    	Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.
  -path-root-api string
    	The root URL for all API handlers (default "/api")
  -properties-reader-uri string
    	A valid whosonfirst/go-reader.Reader URI. Available options are: [file:// fs:// null://]
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/data.SpatialDatabase URI. options are: [sqlite://]
  -verbose
    	Be chatty.
```

For example:

```
$> bin/server \
	-enable-www \
	-enable-tangram \
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-nextzen-apikey {NEXTZEN_APIKEY} 
```

A couple things to note:

* The SQLite databases specified in the `sqlite:///?dsn` string are expected to minimally contain the `rtree` and `spr` and `properties` tables confirming to the schemas defined in the [go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features). They are typically produced by the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package. See the documentation in the [go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package for details.

When you visit `http://localhost:8080` in your web browser you should see something like this:

![](docs/images/server.png)

If you don't need, or want, to expose a user-facing interface simply remove the `-enable-www` and `-nextzen-apikey` flags. For example:

```
$> bin/server \
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' 
```

And then to query the point-in-polygon API you would do something like this:

```
$> curl -X POST -s 'http://localhost:8080/api/point-in-polygon' -d '{"latitude":37.61701894316063, "longitude":-122.3866653442383}'

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
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db'
```

And then:

```
$> curl -s -XPOST -H 'Accept: application/geo+json' 'http://localhost:8080/api/point-in-polygon' -d '{"latitude":37.61701894316063,"longitude":-122.3866653442383 }'

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
