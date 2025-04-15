package tools

import (
	"apps/ecosystem/core/models"
	"fmt"
	"reflect"
	"time"
)

// TODO understand this deeply...
func CreateDynamicStruct(fields []models.Field) reflect.Type {
	structFields := make([]reflect.StructField, len(fields) + 3)

	for i, field := range fields {
		var fieldType reflect.Type
		switch field.Type {
		case "int":
			fieldType = reflect.TypeOf(0)
		case "string":
			fieldType = reflect.TypeOf("")
		case "float64":
			fieldType = reflect.TypeOf(0.0)
		case "bool":
			fieldType = reflect.TypeOf(false)
		case "json":
			fieldType = reflect.TypeOf("")
		default:
			// Default to interface{} if the type is unknown
			fieldType = reflect.TypeOf((*interface{})(nil)).Elem()
		}

		structFields[i] = reflect.StructField{
			Name: field.Name,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s"`, field.Name)),
		}
	}
	// TODO do this better
	structFields[len(fields)] = reflect.StructField{
		Name: "id",
		Type: reflect.TypeOf(""),
		Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s"`, "id")),
	}

	structFields[len(fields) + 1] = reflect.StructField{
		Name: "created",
		Type: reflect.TypeOf(time.Time{}),
		Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s"`, "created")),
	}

	structFields[len(fields) + 2] = reflect.StructField{
		Name: "updated",
		Type: reflect.TypeOf(time.Time{}),
		Tag:  reflect.StructTag(fmt.Sprintf(`db:"%s"`, "updated")),
	}

	return reflect.StructOf(structFields)
}
