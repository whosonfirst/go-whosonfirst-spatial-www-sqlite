package reader

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	wof_reader "github.com/whosonfirst/go-reader"
	"io"
	"io/ioutil"
	_ "log"
	"net/url"
	"regexp"
	"strings"
)

type readFunc func(string) (string, error)

type queryFunc func(string) (string, []interface{}, error)

var VALID_TABLE *regexp.Regexp
var VALID_KEY *regexp.Regexp
var VALID_VALUE *regexp.Regexp

var URI_READFUNC readFunc
var URI_QUERYFUNC queryFunc

func init() {

	VALID_TABLE = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	VALID_KEY = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	VALID_VALUE = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

	ctx := context.Background()
	err := wof_reader.RegisterReader(ctx, "sql", NewSQLReader)

	if err != nil {
		panic(err)
	}
}

type SQLReader struct {
	wof_reader.Reader
	conn  *sql.DB
	table string
	key   string
	value string
}

// sql://sqlite/geojson/id/body/?dsn=....

func NewSQLReader(ctx context.Context, uri string) (wof_reader.Reader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	driver := u.Host
	path := u.Path

	path = strings.TrimLeft(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		return nil, errors.New("Invalid path")
	}

	table := parts[0]
	key := parts[1]
	value := parts[2]

	dsn := q.Get("dsn")

	if dsn == "" {
		return nil, errors.New("Missing dsn parameter")
	}

	conn, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	if !VALID_TABLE.MatchString(table) {
		return nil, errors.New("Invalid table")
	}

	if !VALID_KEY.MatchString(key) {
		return nil, errors.New("Invalid key")
	}

	if !VALID_VALUE.MatchString(value) {
		return nil, errors.New("Invalid value")
	}

	r := &SQLReader{
		conn:  conn,
		table: table,
		key:   key,
		value: value,
	}

	return r, nil
}

func (r *SQLReader) Read(ctx context.Context, raw_uri string) (io.ReadCloser, error) {

	uri := raw_uri

	if URI_READFUNC != nil {

		new_uri, err := URI_READFUNC(raw_uri)

		if err != nil {
			return nil, err
		}

		uri = new_uri
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?", r.value, r.table, r.key)

	q_args := []interface{}{
		uri,
	}

	if URI_QUERYFUNC != nil {

		extra_where, extra_args, err := URI_QUERYFUNC(raw_uri)

		if err != nil {
			return nil, err
		}

		if extra_where != "" {

			q = fmt.Sprintf("%s AND %s", q, extra_where)

			for _, a := range extra_args {
				q_args = append(q_args, a)
			}
		}
	}

	// log.Println(q, q_args)

	row := r.conn.QueryRowContext(ctx, q, q_args...)

	var body string
	err := row.Scan(&body)

	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(body)
	fh := ioutil.NopCloser(sr)

	return fh, nil
}
