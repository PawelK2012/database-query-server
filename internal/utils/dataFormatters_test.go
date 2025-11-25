package utils_test

import (
	"testing"

	"exmple.com/database-query-server/internal/utils"
	"github.com/stretchr/testify/assert"
)

func testData() []map[string]interface{} {
	var mtbl []map[string]interface{}
	row2 := make(map[string]interface{})
	row2["column_name"] = "customername"
	row2["data_type"] = "character varying"
	row2["character_maximum_length"] = int64(200)
	row2["bool_data"] = true
	row2["float"] = 11.22

	mtbl = append(mtbl, row2)
	return mtbl
}

func TestDataToJson(t *testing.T) {
	tests := []struct {
		name    string
		data    []map[string]interface{}
		want    string
		wantErr bool
	}{
		{name: "Happy Flow - data to JSON", data: testData(), want: "[{\"bool_data\":true,\"character_maximum_length\":200,\"column_name\":\"customername\",\"data_type\":\"character varying\",\"float\":11.22}]", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.DataToJson(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DataToJson() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DataToJson() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestDataToCSV(t *testing.T) {
	tests := []struct {
		name    string
		data    []map[string]interface{}
		want    string
		wantErr bool
	}{
		{name: "Happy Flow - data to CSV", data: testData(), want: "bool_data,character_maximum_length,column_name,data_type,float\ntrue,200,customername,character varying,11.22\n", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.DataToCSV(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DataToCSV() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DataToCSV() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestDataToHTMLTable(t *testing.T) {
	tests := []struct {
		name    string
		rows    []map[string]interface{}
		want    string
		wantErr bool
	}{
		{name: "Happy Flow - data to HTML table", rows: testData(), want: "<table><thead><tr><th>bool_data</th><th>character_maximum_length</th><th>column_name</th><th>data_type</th><th>float</th></tr></thead><tbody><tr><td>true</td><td>200</td><td>customername</td><td>character varying</td><td>11.22</td></tr></tbody></table>", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.DataToHTMLTable(tt.rows)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DataToHTMLTable() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DataToHTMLTable() succeeded unexpectedly")
			}
			if true {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
