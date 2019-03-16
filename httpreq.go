/*
Package httpreq unmarshals values from http.Request as specified in struct tags.
*/
package httpreq

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const tagName = "req"

var (
	// ErrIsNil is returned when value is nil
	ErrIsNil = errors.New("v is nil")
	// ErrIsNotPointer is returned when value is not a pointer
	ErrIsNotPointer = errors.New("v is not a pointer")
	// ErrCannotSet is returned when field is not settable
	ErrCannotSet = errors.New("cannot set field value")
	// ErrInvalidTag is returned when tag cannot be parsed
	ErrInvalidTag = errors.New("invalid tag")
	// ErrInvalidAt is returned when at-attribute has invalid value
	ErrInvalidAt = errors.New("invalid value for at-attribute")
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
		fieldT := p.Type().Field(i)
		tag := fieldT.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		fieldV := p.Field(i)
		if !fieldV.CanSet() {
			return errors.Wrap(ErrCannotSet, fieldT.Name)
		}

		raw, err := getRawValueForTag(req, tag)
		if err != nil {
			return err
		}

		switch fieldT.Type.Kind() {
		case reflect.String:
			fieldV.SetString(raw)
		case reflect.Int:
			intVal, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(intVal)
		}
	}

	return nil
}

func getRawValueForTag(req *http.Request, tag string) (string, error) {
	var name, at string
	args := strings.Split(tag, ",")
	if len(args) != 2 {
		return "", errors.Wrap(ErrInvalidTag, tag)
	}
	name = args[0]
	_, err := fmt.Sscanf(args[1], "at=%s", &at)
	if err != nil {
		return "", errors.Wrap(ErrInvalidTag, tag)
	}

	switch at {
	case "query":
		return req.URL.Query().Get(name), nil
	}
	return "", errors.Wrap(ErrInvalidAt, at)
}
