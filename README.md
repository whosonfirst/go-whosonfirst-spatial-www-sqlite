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
	-enable-properties \
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-nextzen-apikey {NEXT_APIKEY} \
	-mode directory://
```

A couple things to note:

* The SQLite databases specified in the `sqlite:///?dsn` string are expected to minimally contain the `rtree` and `geojson` tables confirming to the schemas defined in the [go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features). They are typically produced by the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package.

* Do you notice the way we are passing in a `-mode directory://` flag? This should only be necessary if we are generating, or updating, a SQLite database when the `server` tool starts up. As of this writing it is always necessary even though it doesn't do anything. There is an [open ticket](https://github.com/whosonfirst/go-whosonfirst-spatial/issues/14) to address this.

When you visit `http://localhost:8080` in your web browser you should see something like this:

![](docs/images/server.png)

If you don't need, or want, to expose a user-facing interface simply remove the `-enable-www` and `-nextzen-apikey` flags. For example:

```
$> bin/server \
	-enable-properties \
	-spatial-database-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-reader-uri 'sqlite:///?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-mode directory://
```

And then to query the point-in-polygon API you would do something like this:

```
$> curl 'http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383'

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

By default, results are returned as a list of ["standard places response"](https://github.com/whosonfirst/go-whosonfirst-spr/) (SPR) elements. You can also return results as a GeoJSON `FeatureCollection` by including a `format=geojson` query parameter. For example:


```
$> curl 'http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383&format=geojson'

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
$> http://localhost:8080/api/point-in-polygon?latitude=37.61701894316063&longitude=-122.3866653442383&format=geojson&properties=sfomuseum:*
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

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-http
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spr
* https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index
* https://developers.nextzen.org/
