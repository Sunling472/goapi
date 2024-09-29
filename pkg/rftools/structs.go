package rftools

import (
	"errors"
	rf "reflect"
	"strings"
)

type testStruct struct {
	Name string
	Age  int
}

type (
	structTags   map[string][]string
	allFieldTags map[string]structTags
)

func SetAttr(s any, name string, value any) {
	v := rf.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == name {
			valueType := rf.TypeOf(value)
			if valueType == field.Type && v.Field(i).CanSet() {
				v.Field(i).Set(rf.ValueOf(value))
			}
		} else {
			continue
		}
	}
}

func GetAttr[T any](s any, name string) (T, error) {
	v := rf.ValueOf(s).Elem()
	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Type().Field(i)
		if f.Name == name {
			res := v.Field(i).Interface()
			if data, ok := res.(T); ok {
				return data, nil
			} else if !ok {
				return *new(T), errors.New("Incorrect type")
			}
			continue
		}
	}
	return *new(T), errors.New("Attr is not found")
}

func TagsByName(s any, tagKey string) structTags {
	t := rf.TypeOf(s).Elem()
	result := structTags{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := f.Tag.Get(tagKey)
		if tags != "" {
			resTags := strings.Split(tags, ",")
			result[f.Name] = resTags
		}
	}
	return result
}

func AllTags(s any) allFieldTags {
	t := rf.TypeOf(s).Elem()
	result := allFieldTags{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := strings.Split(string(f.Tag), " ")
		for _, tag := range tags {
			parseTag := strings.Split(tag, ":")
			key := parseTag[0]
			values := parseTag[1]
			values = values[1 : len(values)-1]
			resultVal := strings.Split(values, ",")
			if result[f.Name] == nil {
				result[f.Name] = structTags{}
			}
			result[f.Name][key] = resultVal
		}
	}
	return result
}
