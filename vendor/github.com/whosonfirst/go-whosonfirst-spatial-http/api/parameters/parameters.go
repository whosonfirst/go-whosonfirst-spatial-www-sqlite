package parameters

import (
	"errors"
	"github.com/aaronland/go-http-sanitize"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"net/http"
	"strconv"
	"strings"
)

func Latitude(req *http.Request) (float64, error) {

	str_lat, err := sanitize.GetString(req, "latitude")

	if err != nil {
		return 0, err
	}

	if str_lat == "" {
		return 0, errors.New("Missing 'latitude' parameter")
	}

	lat, err := strconv.ParseFloat(str_lat, 64)

	if err != nil {
		return 0, err
	}

	if !geo.IsValidLatitude(lat) {
		return 0, errors.New("Invalid 'latitude' parameter")
	}

	return lat, nil
}

func Longitude(req *http.Request) (float64, error) {

	str_lon, err := sanitize.GetString(req, "longitude")

	if err != nil {
		return 0, err
	}

	if str_lon == "" {
		return 0, errors.New("Missing 'longitude' parameter")
	}

	lon, err := strconv.ParseFloat(str_lon, 64)

	if err != nil {
		return 0, err
	}

	if !geo.IsValidLongitude(lon) {
		return 0, errors.New("Invalid 'longitude' parameter")
	}

	return lon, nil
}

func Coordinate(req *http.Request) (*geom.Coord, error) {

	lat, err := Latitude(req)

	if err != nil {
		return nil, err
	}

	lon, err := Longitude(req)

	if err != nil {
		return nil, err
	}

	return geo.NewCoordinate(lon, lat)
}

func Properties(req *http.Request) ([]string, error) {

	str_properties, err := sanitize.GetString(req, "properties")

	if err != nil {
		return nil, err
	}

	properties := listWithString(str_properties, ",")
	return properties, nil
}

// as in ?geometries=all or ?geometries=default

func Geometries(req *http.Request) (string, error) {

	var geoms string

	geoms, err := sanitize.GetString(req, "geometries")

	if err != nil {
		return "", err
	}

	geoms = strings.Trim(geoms, " ")

	return geoms, nil
}

// as in ?alternate-geometries=quattroshapes&alternate-geometries=reversegeo

func AlternateGeometries(req *http.Request) ([]string, error) {

	str_geoms, err := sanitize.GetString(req, "alternate-geometries")

	if err != nil {
		return nil, err
	}

	alt_geoms := listWithString(str_geoms, ",")
	return alt_geoms, nil
}

func listWithString(raw string, sep string) []string {

	list := make([]string, 0)	
	trimmed := strings.Trim(raw, "")

	for _, str := range strings.Split(trimmed, sep) {

		str = strings.Trim(str, "")

		if str == "" {
			continue
		}

		list = append(list, str)
	}

	return list
}
