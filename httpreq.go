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

var (
	// ErrIsNil is returned when value is nil
	ErrIsNil = errors.New("v is nil")
	// ErrIsNotPointer is returned when value is not a pointer
	ErrIsNotPointer = errors.New("v is not a pointer")
)

// Unmarshal parses values from req into the fields of v
func Unmarshal(req *http.Request, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return ErrIsNotPointer
	}
	if reflect.ValueOf(v).IsNil() {
		return ErrIsNil
	}
	p := reflect.ValueOf(v).Elem()

	for i := 0; i < p.NumField(); i++ {
		field := p.Type().Field(i)
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
			fieldVal := p.Field(i)
			fieldVal.SetString(raw)
		case reflect.Int:
			intVal, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return err
			}
			fieldVal := p.Field(i)
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
