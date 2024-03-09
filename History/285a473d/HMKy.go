package data

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sudeepbatra/alpha-hft/logger"
)

var (
	ErrDataNotStruct       = fmt.Errorf("data is not a struct")
	ErrUnableExtractFields = fmt.Errorf("unable to extract fields from the record")
)

func getStructFieldNames(data interface{}) []string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		firstItem := v.Index(0)
		if firstItem.Kind() == reflect.Struct {
			numFields := firstItem.Type().NumField()
			fieldNames := make([]string, numFields)

			for i := 0; i < numFields; i++ {
				fieldName := strings.ToLower(firstItem.Type().Field(i).Name)
				fieldNames[i] = fieldName
			}

			return fieldNames
		}
	}

	return nil
}

func BulkInsert(tableName string, data interface{}) error {
	var rows [][]interface{}

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("data is not a slice")
	}

	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Struct {
			var row []interface{}
			for j := 0; j < item.NumField(); j++ {
				row = append(row, item.Field(j).Interface())
			}

			rows = append(rows, row)
		}
	}

	fieldNames := getStructFieldNames(data)

	if fieldNames == nil {
		return fmt.Errorf("unable to extract fields from data")
	}

	logger.Log.Trace().Strs("fields", fieldNames).Msg("fields extracted from data")

	copyCount, err := AlphaHftDbConnPool.CopyFrom(context.Background(),
		pgx.Identifier{tableName},
		fieldNames,
		pgx.CopyFromRows(rows))
	if err != nil {
		return err
	}

	logger.Log.Trace().Int64("count", copyCount).Msg("bulk insert done")

	return nil
}

func areEqualIgnoringCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

func CopyStructExcludingFields(source, destination interface{}, excludedFields map[string]bool) {
	sourceValue := reflect.ValueOf(source)
	destValue := reflect.ValueOf(destination)

	if sourceValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		field := sourceValue.Type().Field(i)
		fieldName := strings.ToLower(field.Name)

		if excludedFields[fieldName] {
			continue
		}

		destField := destValue.FieldByNameFunc(func(f string) bool {
			return areEqualIgnoringCase(f, field.Name)
		})

		if destField.IsValid() && destField.CanSet() {
			destField.Set(sourceValue.Field(i))
		}
	}
}

// InsertRecord inserts a single record for a given struct into the specified table
func InsertRecord(tableName string, data interface{}) error {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	var columns []string
	var values []interface{}

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		columns = append(columns, field.Name)
		values = append(values, dataValue.Field(i).Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName,
		strings.Join(columns, ", "),
		generatePlaceholderString(len(columns)))

	conn, err := AlphaHftDbConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer conn.Release()

	_, err = conn.Exec(context.Background(), query, values...)

	return err
}

func generatePlaceholderString(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ",")
}
