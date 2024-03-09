package data

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

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

func BulkInsertV1(tableName string, data interface{}) error {
	var rows [][]interface{}

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("data is not a slice")
	}

	int64Type := reflect.TypeOf(int64(0))
	float64Type := reflect.TypeOf(float64(0.0))

	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Struct {
			var row []interface{}
			for j := 0; j < item.NumField(); j++ {
				switch fieldType := item.Field(j).Type(); fieldType {
				case int64Type:
					// Convert int64 fields to time.Time if they are timestamps
					if fieldName := item.Type().Field(j).Name; strings.Contains(strings.ToLower(fieldName), "timestamp") {
						timestampValue := item.Field(j).Int()
						timestamp := time.Unix(timestampValue, 0)
						row = append(row, timestamp)
					} else {
						// For other int64 fields, append as is
						row = append(row, item.Field(j).Interface())
					}
				case float64Type:
					// Handle float64 fields (modify as needed)
					row = append(row, item.Field(j).Interface())
				default:
					// For other field types, append as is
					row = append(row, item.Field(j).Interface())
				}
			}

			rows = append(rows, row)
		}
	}

	fieldNames := getStructFieldNames(data)

	if fieldNames == nil {
		return fmt.Errorf("unable to extract fields from data")
	}

	logger.Log.Info().Strs("fields", fieldNames).Msg("fields extracted from data")

	copyCount, err := AlphaHftDbConnPool.CopyFrom(context.Background(),
		pgx.Identifier{tableName},
		fieldNames,
		pgx.CopyFromRows(rows))

	if err != nil {
		return err
	}

	logger.Log.Info().Int64("count", copyCount).Msg("bulk insert done")

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
				// Check if the field type is time.Time
				if item.Type().Field(j).Type == reflect.TypeOf(time.Time{}) {
					// Convert Unix timestamp to time.Time
					timestampValue := item.Field(j).Int()
					timestamp := time.Unix(timestampValue, 0)
					row = append(row, timestamp)
				} else {
					// For other fields, append as is
					row = append(row, item.Field(j).Interface())
				}
			}

			rows = append(rows, row)
		}
	}

	fieldNames := getStructFieldNames(data)

	if fieldNames == nil {
		return fmt.Errorf("unable to extract fields from data")
	}

	logger.Log.Info().Strs("fields", fieldNames).Msg("fields extracted from data")

	copyCount, err := AlphaHftDbConnPool.CopyFrom(context.Background(),
		pgx.Identifier{tableName},
		fieldNames,
		pgx.CopyFromRows(rows))

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
