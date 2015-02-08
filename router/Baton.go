package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// A Baton represents a web request, with all session based tools
// which might be required.
// The Baton contains
//  - the gbxproto.GbxBuffer to use for requests
type Baton struct {
	w     http.ResponseWriter
	r     *http.Request
	path  string
	route *route
	Token interface{}
}

// Scan matches up the request url with the provided dests.
// e.g. pattern /api/site/%d/%s with request /api/site/1/hello
// Scan(&anInteger, &aString)
//
// Currently only supports:
// %d -> uint64, uint32
// %s -> string
func (b *Baton) Scan(dest ...interface{}) error {
	patternParts := strings.Split(b.route.format[1:], "/")
	urlParts := strings.Split(b.r.URL.Path[1:], "/")
	if len(urlParts) != len(patternParts) {
		fmt.Println(urlParts)
		fmt.Println(patternParts)

		return fmt.Errorf("URL had %d parts, pattern had %d", len(urlParts), len(patternParts))
	}

	// Add the url parts which are represented by placeholders in the pattern.
	// i.e., in a/%b/c/%d matched to g/h/i/j, return an array of just 'h' and 'j'
	matchList := make([]string, 0, len(dest))
	for i, patt := range patternParts {
		if !strings.HasPrefix(patt, "%") {
			continue
		}
		matchList = append(matchList, urlParts[i])
	}

	if len(matchList) != len(dest) {
		return fmt.Errorf("URL had %d parameters, expected %d", len(matchList), len(dest))
	}

	// Parse the url parts into their placeholders
	for i, src := range matchList {
		dst := dest[i]
		switch t := dst.(type) {
		case *string:
			*t = src
		case *uint64:
			srcInt, err := strconv.ParseUint(src, 10, 64)
			if err != nil {
				return fmt.Errorf("URL Parameter %d could not be converted to an unsigned integer")
			}
			*t = srcInt
		case *uint32:
			srcInt, err := strconv.ParseUint(src, 10, 32)
			if err != nil {
				return fmt.Errorf("URL Parameter %d could not be converted to an unsigned integer")
			}
			*t = uint32(srcInt)

		default:
			return fmt.Errorf("URL Parameter %d could not be converted to a %T",
				i+1, t)

		}
	}
	return nil
}

// Method returns the request method (POST, GET etc)
func (b *Baton) Method() string {
	return b.r.Method
}

func (b *Baton) Accept(t string) bool {
	a := b.r.Header.Get("accept")
	return strings.Contains(a, t)
}

// Raw returns the underlying writer and request
func (b *Baton) Raw() (http.ResponseWriter, *http.Request) {
	return b.w, b.r
}

// Path returns the path component of the URL
func (b *Baton) Path() string {
	return b.r.URL.Path
}

func (b *Baton) FormValue(name string) string {
	return b.r.FormValue(name)
}

// QueryString grabs a parameter from the querystring as a string
func (b *Baton) QueryString(key string) (string, bool) {
	str := b.r.URL.Query().Get(key)
	if len(str) < 1 {
		return "", false
	}
	return str, true
}

// QueryStringArray grabs a comma seperated value from the querystring
func (b *Baton) QueryStringArray(key string) ([]string, bool) {
	setVal, ok := b.QueryString(key)
	if !ok {
		return []string{}, false
	}
	return strings.Split(setVal, ","), true
}

func (b *Baton) QueryBool(key string) (bool, bool) {
	strVal, ok := b.QueryString(key)
	if !ok {
		return false, false
	}
	switch strings.ToLower(strVal) {
	case "true", "1", "yes":
		return true, true
		return true, true
		return true, true
	default:
		return false, true
	}
}

func (b *Baton) QueryIntArray(key string) ([]int, bool) {
	strs, ok := b.QueryStringArray(key)
	if !ok {
		return []int{}, false
	}
	ints := make([]int, len(strs), len(strs))
	for i, str := range strs {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return []int{}, false
		}
		ints[i] = int(val)
	}
	return ints, true
}

func (b *Baton) QueryUIntArray(key string) ([]uint64, bool) {
	strs, ok := b.QueryStringArray(key)
	if !ok {
		return []uint64{}, false
	}
	ints := make([]uint64, len(strs), len(strs))
	for i, str := range strs {
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return []uint64{}, false
		}
		ints[i] = val
	}
	return ints, true
}

// QueryUInt grabs a parameter from the querystring as an unsigned, 64 bit integer
func (b *Baton) QueryUInt(key string) (uint64, bool) {
	str, ok := b.QueryString(key)
	if !ok {
		return 0, false
	}
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, false
	}
	return val, true
}

// QueryUIntDefault grabs a parameter from the querystring if it is set.
func (b *Baton) QueryUIntDefault(key string, defaultValue uint64) uint64 {
	realValue, ok := b.QueryUInt(key)
	if ok {
		return realValue
	}
	return defaultValue
}

func (b *Baton) SendError(err error, status int) {
	log.Printf("SEND ERROR %s\n", err.Error())
	b.w.WriteHeader(status)
	b.w.Write([]byte(err.Error()))
}

func (b *Baton) SendJSON(obj interface{}) {
	enc := json.NewEncoder(b.w)
	enc.Encode(obj)
}
