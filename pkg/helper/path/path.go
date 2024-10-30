package path_helper

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func MakeDir(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatalf("Error during Make Directory : %s \n", err.Error())
		}
	}
	return path
}

func AddWorkdirToSomePath(path ...string) string {
	workdir, err := os.Getwd()
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error during Add workdir to some path cause no exist : %s", err.Error())
			return ""
		} else if errors.Is(err, io.EOF) {
			log.Fatalf("Error during Add workdir to some path (EOF) : %s", err.Error())
			return ""
		}
		log.Fatalf("Error during Add workdir to some path : %s", err.Error())
		return ""
	}
	return filepath.Join(append([]string{workdir}, path...)...)
}

func MakedirFromFieldStruct(value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Ptr || (val.Kind() == reflect.Ptr && val.IsNil()) {
		return errors.New("value must be struct or pointer struct")
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return processFields(val)
}

func processFields(val reflect.Value) error {
	typeVal := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typeVal.Field(i)
		if err := processField(field, fieldType); err != nil {
			return err
		}
	}
	return nil
}

func processField(field reflect.Value, fieldType reflect.StructField) error {
	if field.Kind() == reflect.Struct {
		return MakedirFromFieldStruct(field.Interface())
	} else if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Kind() == reflect.Struct {
		return MakedirFromFieldStruct(field.Interface())
	}
	if strings.Contains(strings.ToUpper(fieldType.Name), "PATH") {
		if field.Kind() == reflect.String {
			MakeDir(AddWorkdirToSomePath(fmt.Sprintf("%v", field.Interface())))
			return nil
		}
	}
	return nil
}
