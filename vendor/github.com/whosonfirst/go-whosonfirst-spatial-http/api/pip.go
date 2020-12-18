package api

import (
	"encoding/json"
	"fmt"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/output"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/parameters"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	// "github.com/whosonfirst/go-whosonfirst-flags/placetypes"			
	"github.com/whosonfirst/go-whosonfirst-flags/existential"		
	"github.com/whosonfirst/go-whosonfirst-flags/geometry"		
	_ "log"
	"net/http"
)

type PointInPolygonHandlerOptions struct {
	EnableGeoJSON    bool
	EnableProperties bool
	GeoJSONReader    reader.Reader
}

func PointInPolygonHandler(spatial_app *app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	spatial_db := spatial_app.SpatialDatabase
	properties_r := spatial_app.PropertiesReader
	walker := spatial_app.Walker

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()
		query := req.URL.Query()

		coord, err := parameters.Coordinate(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_format, err := sanitize.GetString(req, "format")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if str_format == "geojson" && !opts.EnableGeoJSON {
			http.Error(rsp, "GeoJSON formatting is disabled.", http.StatusBadRequest)
			return
		}

		if str_format == "properties" && !opts.EnableProperties {
			http.Error(rsp, "Properties formatting is disabled.", http.StatusBadRequest)
			return
		}

		properties_paths, err := parameters.Properties(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}
		
		filters, err := filter.NewSPRFilterFromQuery(query)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		err = appendFilterWithParameters(filters)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}
		
		results, err := spatial_db.PointInPolygon(ctx, coord, filters)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if results == nil {
			http.Error(rsp, "Unable to yield results", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		var final interface{}
		final = results

		enc := json.NewEncoder(rsp)

		switch str_format {
		case "geojson":

			err := geojson.AsFeatureCollection(ctx, results, opts.GeoJSONReader, rsp)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			return

		case "properties":

			if len(properties_paths) > 0 {

				props, err := properties_r.PropertiesResponseResultsWithStandardPlacesResults(ctx, final.(spr.StandardPlacesResults), properties_paths)

				if err != nil {
					http.Error(rsp, err.Error(), http.StatusInternalServerError)
					return
				}

				final = props
			}

		default:
			// spr (above)
		}

		err = enc.Encode(final)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func PointInPolygonCandidatesHandler(spatial_app *app.SpatialApplication) (http.Handler, error) {

	walker := spatial_app.Walker
	spatial_db := spatial_app.SpatialDatabase

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()

		coord, err := parameters.Coordinate(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		candidates, err := spatial_db.PointInPolygonCandidates(ctx, coord)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		output.AsJSON(rsp, candidates)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func appendFilterWithParameters(f filter.Filter) error {

	geometries, err := parameters.Geometries()

	if err != nil {
		return err
	}

	switch geometries {
	case "all":
		// pass
	case "alt", "alternate":

		af, err := geometry.NewIsAlternateGeometryFlag(true)

		if err != nil {
			return fmt.Errorf("Failed to create alternate geometry flag, %v", err)
		}

		f.AlternateGeometry = af
		
	case "default":

		af, err := geometry.NewIsAlternateGeometryFlag(false)

		if err != nil {
			fmt.Errorf("Failed to create alternate geometry flag, %v", err)
		}

		f.AlternateGeometry = af
		
	default:
		fmt.Errorf("Invalid -geometries flag")
	}
	
	alt_geoms, err := parameters.AlternateGeometries()

	if err != nil {
		return err
	}

	if len(alt_geoms) > 0 {

		alt_flags, err := geometry.NewAlternateGeometryFlagsWithLabelArray(alt_geoms...)

		if err != nil {
			fmt.Errorf("Failed to create alternate geometries flags, %v", err)
		}
				
		f.AlternateGeometries = alt_flags
	}

	/*
	if len(pts) > 0 {

		pt_flags, err := placetypes.NewPlacetypeFlagsArray(pts...)

		if err != nil {
			fmt.Errorf("Failed to create placetype flags, %v", err)
		}
		
		f.Placetypes = pt_flags
	}
	*/

	is_current, err := parameters.IsCurrent()

	if err != nil {
		return err
	}
	
	if len(is_current) > 0 {

		existential_flags, err := existential.NewKnownUnknownFlagsArray(is_current...)

		if err != nil {
			fmt.Errorf("Failed to create is-current flags, %v", err)
		}
		
		f.Current = existential_flags
	}
	
	is_ceased, err := parameters.IsCeased()
	
	if err != nil {
		return err
	}

	if len(is_ceased) > 0 {

		existential_flags, err := existential.NewKnownUnknownFlagsArray(is_ceased...)

		if err != nil {
			fmt.Errorf("Failed to create is-ceased flags, %v", err)
		}

		f.Ceased = existential_flags
	}

	is_deprecated, err := parameters.IsDeprecated()

	if err != nil {
		return err
	}
	
	if len(is_deprecated) > 0 {

		existential_flags, err := existential.NewKnownUnknownFlagsArray(is_deprecated...)

		if err != nil {
			fmt.Errorf("Failed to create is-deprecated flags, %v", err)
		}

		f.Deprecated = existential_flags
	}

	is_superseded, err := parameters.IsSuperseded()

	if err != nil {
		return err
	}

	if len(is_superseded) > 0 {

		existential_flags, err := existential.NewKnownUnknownFlagsArray(is_superseded...)

		if err != nil {
			fmt.Errorf("Failed to create is-superseded flags, %v", err)
		}

		f.Superseded = existential_flags
	}

	is_superseding, err := parameters.IsSuperseding()

	if err != nil {
		return err
	}
	
	if len(is_superseding) > 0 {

		existential_flags, err := existential.NewKnownUnknownFlagsArray(is_superseding...)

		if err != nil {
			fmt.Errorf("Failed to create is-superseding flags, %v", err)
		}

		f.Superseding = existential_flags
	}

	return nil
}
