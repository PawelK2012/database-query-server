package utils_test

import (
	"testing"

	"exmple.com/database-query-server/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestDataToJson(t *testing.T) {
	var mtbl []map[string]interface{}
	row2 := make(map[string]interface{})
	row2["column_name"] = "customername"
	row2["data_type"] = "character varying"
	row2["character_maximum_length"] = int64(200)
	mtbl = append(mtbl, row2)

	tests := []struct {
		name    string
		data    []map[string]interface{}
		want    string
		wantErr bool
	}{
		{name: "Happy Flow - data to JSON", data: mtbl, want: "[{\"character_maximum_length\":200,\"column_name\":\"customername\",\"data_type\":\"character varying\"}]", wantErr: false},
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
