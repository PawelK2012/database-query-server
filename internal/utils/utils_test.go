package utils_test

import (
	"testing"

	"exmple.com/database-query-server/internal/utils"
)

func TestCheckFirstWord(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   string
		wantErr bool
	}{
		{name: "Happy flow", input: "SELECT * FROM mydb", wantErr: false},
		{name: "Sad flow", input: "SELECTA SELECT * FROM mydb", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := utils.CheckFirstWord(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckFirstWord() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckFirstWord() succeeded unexpectedly")
			}
		})
	}
}
