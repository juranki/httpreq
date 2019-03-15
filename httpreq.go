/*
Package httpreq unmarshals values from http.Request as specified in struct tags.
*/
package httpreq

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "httpreq"

// Unmarshal parses values from req into the fields of v
func Unmarshal(req *http.Request, v interface{}) error {
	vv := reflect.ValueOf(v).Elem()
	// s := vv.Elem()

	for i := 0; i < vv.NumField(); i++ {
		field := vv.Type().Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		raw, err := getRawValueForTag(req, tag)
		if err != nil {
			return err
		}
		fmt.Printf("%+v, %s", field.Type, raw)

		switch field.Type.Kind() {
		case reflect.String:
			fieldVal := vv.Field(i)
			fieldVal.SetString(raw)
		case reflect.Int:
			intVal, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return err
			}
			fieldVal := vv.Field(i)
			fieldVal.SetInt(intVal)

		}
	}

	return nil
}

func getRawValueForTag(req *http.Request, tag string) (string, error) {
	var name, at string
	args := strings.Split(tag, ",")
	if len(args) != 2 {
		return "", errors.New("invalid httpreq")
	}
	name = args[0]
	_, err := fmt.Sscanf(args[1], "at=%s", &at)
	if err != nil {
		return "", err
	}

	switch at {
	case "query":
		return req.URL.Query().Get(name), nil
	}
	return "", fmt.Errorf("unknown location %s", at)
}
