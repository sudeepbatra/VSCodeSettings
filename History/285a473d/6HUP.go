package data

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
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

	// Check if the type of data is a slice and it has elements
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		firstItem := v.Index(0)

		// If the element is a pointer, dereference it
		if firstItem.Kind() == reflect.Ptr {
			firstItem = firstItem.Elem()
		}

		// Check if the elements of the slice are of struct type
		if firstItem.Kind() == reflect.Struct {
			numFields := firstItem.Type().NumField()
			fieldNames := make([]string, numFields)

			// Iterate through struct fields and store their names in lowercase
			for i := 0; i < numFields; i++ {
				field := firstItem.Type().Field(i)
				fieldNames[i] = strings.ToLower(field.Name)
			}

			return fieldNames
		}

		// If the element is not a struct, try dereferencing it as a pointer to struct
		if firstItem.Kind() == reflect.Ptr && firstItem.Elem().Kind() == reflect.Struct {
			numFields := firstItem.Elem().Type().NumField()
			fieldNames := make([]string, numFields)

			// Iterate through struct fields and store their names in lowercase
			for i := 0; i < numFields; i++ {
				field := firstItem.Elem().Type().Field(i)
				fieldNames[i] = strings.ToLower(field.Name)
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

	logger.Log.Info().Strs("fields", fieldNames).Msg("fields extracted from data")

	copyCount, err := AlphaHftDbConn.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		fieldNames,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	logger.Log.Info().Int64("count", copyCount).Msg("bulk insert done")

	return nil
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

	_, err := AlphaHftDbConn.Exec(context.Background(), query, values...)

	return err
}

func generatePlaceholderString(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ",")
}

func GetClient(useProxy bool) *http.Client {
	if useProxy {
		proxyURL := "http://localhost:8080"

		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			panic("Failed to parse proxy URL: " + err.Error())
		}

		proxyTransport := &http.Transport{
			Proxy:               http.ProxyURL(proxyURLParsed),
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}

		return &http.Client{
			Transport: proxyTransport,
		}
	}

	return &http.Client{}
}
