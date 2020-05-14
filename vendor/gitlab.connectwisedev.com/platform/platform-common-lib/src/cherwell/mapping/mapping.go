// Package mapping provides helpers for mapping data to Cherwell's business objects.
package mapping

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cherwell"
)

const (
	boFieldTag = "cherwell"

	// htmlTagFlag additional html tag flag that marks this value should be mapped as HTML text
	htmlTagFlag = "isHtml"

	// omitEmptyFlag is additional flag that indicates that struct's property should be ignored
	// if destination field is missing in business object
	omitEmptyFlag = "omitempty"

	// DefaultDateFormat is default Cherwell's datetime format for Save operations.
	DefaultDateFormat = "1/2/2006 3:04 PM"
)

var errMissingMappingField = errors.New("mapping destination field is missing")

// mappingProperties is mapping properties for field
type mappingProperties struct {
	// fieldName is target business object field name
	fieldName string

	// isHTML indicates that field's value should be added to HTML BO field property
	isHTML bool

	// omitEmpty indicates that property should be ignored if destination field is missing
	omitEmpty bool
}

// MapToBusinessObject maps values from provided struct to business object.
//
// Field mappings should be defined as "cherwell" field tag.
// Fields without tag will be ignored.
//
// Example:
//
//	type Foo struct {
//		Foo string `cherwell:"BusinessObjectFieldName"`
//
// 		// "omitempty" allows to ignore props if destination field is missing in BO
//		Bar string `cherwell:"FieldName,omitempty"`
//
//		// "isHtml" copies field value to HTML property of BO's field
//		Baz string `cherwell:"SomeHTMLField,isHtml"
//	}
func MapToBusinessObject(source interface{}, destination *cherwell.BusinessObject) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	if source == nil {
		return fmt.Errorf("nil value passed")
	}

	if destination == nil {
		return fmt.Errorf("output business object value is nil")
	}

	obj := reflect.ValueOf(source)
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	if obj.Kind() != reflect.Struct {
		return fmt.Errorf("provided object should be a struct pointer or struct, got %T", source)
	}

	if !obj.IsValid() {
		return fmt.Errorf("provided object is not addressable or invalid")
	}

	oType := obj.Type()
	fieldsCount := oType.NumField()

	for i := 0; i < fieldsCount; i++ {
		field := obj.Field(i)
		fType := oType.Field(i)
		err = mapFieldToBO(fType, &field, destination)
		if err != nil {
			return fmt.Errorf("failed to map field '%s', %s", fType.Name, err)
		}
	}

	return nil
}

func resolveReference(valRef *reflect.Value) (val interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred during reference resolve: %s", r)
		}
	}()

	if valRef.IsNil() {
		return nil, nil
	}

	return valRef.Elem().Interface(), nil
}

func valueToString(val interface{}) (string, error) {
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case bool:
		return strconv.FormatBool(val.(bool)), nil
	case string:
		return val.(string), nil
	case time.Time:
		t := val.(time.Time)
		return t.Format(DefaultDateFormat), nil
	default:
		return "", fmt.Errorf("cannot convert value to string, unsupported type %T", val)
	}
}

// extractMappingProperties parses tag string value into mapping properties.
//
// Supported flags:
//
// "isHtml" - adds value to 'HTML' property of BO field
//
// "omitempty" - ignores property if destination field doesn't exists in business object
func extractMappingProperties(tagValue string) (mappingProps mappingProperties, err error) {
	tagValue = strings.TrimSpace(tagValue)
	if tagValue == "" {
		err = errMissingMappingField
		return
	}

	sections := strings.Split(tagValue, ",")
	mappingProps.fieldName = sections[0]

	if len(sections) == 1 {
		return mappingProps, nil
	}

	// parse additional tag flags
	sections = sections[1:]
	for _, flag := range sections {
		flag = strings.TrimSpace(flag)

		switch flag {
		case htmlTagFlag:
			mappingProps.isHTML = true
		case omitEmptyFlag:
			mappingProps.omitEmpty = true
		default:
		}
	}

	return mappingProps, nil
}

func mapFieldToBO(fieldType reflect.StructField, fieldRef *reflect.Value, bo *cherwell.BusinessObject) (err error) {
	tagValue, ok := fieldType.Tag.Lookup(boFieldTag)
	if !ok {
		return nil
	}

	props, err := extractMappingProperties(tagValue)
	if err != nil {
		return err
	}

	destField, ok := bo.FieldByName(props.fieldName)
	if !ok {
		// Ignore error if "omitempty" flag provided
		if props.omitEmpty {
			return nil
		}

		return fmt.Errorf("field '%s' doesn't exists in business object #%s", props.fieldName, bo.ID)
	}

	var fieldVal interface{}

	if fieldRef.Kind() == reflect.Ptr {
		// Extract pointer value from field
		fieldVal, err = resolveReference(fieldRef)
		if err != nil {
			return fmt.Errorf("cannot get value of field '%s': %s", props.fieldName, err)
		}
	} else {
		// Get field value if field type is not pointer
		fieldVal = fieldRef.Interface()
	}

	// Ignore nil fields
	if fieldVal == nil {
		return nil
	}

	strVal, err := valueToString(fieldVal)
	if err != nil {
		return err
	}

	destField.Value = strVal
	if props.isHTML {
		destField.HTML = strVal
	}

	destField.Dirty = true
	return nil
}
