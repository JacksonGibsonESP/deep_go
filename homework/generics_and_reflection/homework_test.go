package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var builder strings.Builder
	objType := reflect.TypeOf(person)
	objValue := reflect.ValueOf(person)

	for i := range objType.NumField() {
		fieldType := objType.Field(i)
		fieldTag := fieldType.Tag.Get("properties")

		fieldValue := objValue.Field(i)
		tagSplits := strings.Split(fieldTag, ",")
		if len(tagSplits) > 1 &&
			tagSplits[1] == "omitempty" &&
			isEmpty(fieldValue) {
			continue
		}

		strValue := valueToString(fieldValue)
		propertyName := tagSplits[0]

		builder.WriteString(propertyName)
		builder.WriteString("=")
		builder.WriteString(strValue)
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Int:
		return v.Int() == 0
	case reflect.Bool:
		return false
	default:
		return false
	}
}

func valueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String())
	case reflect.Int:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return ""
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
