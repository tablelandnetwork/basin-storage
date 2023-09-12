package storage

import (
	"testing"
)

func TestExtractTableName(t *testing.T) {
	tests := []struct {
		filename   string
		expected   Pub
		shouldFail bool
	}{
		{
			filename:   "202309110419279721253170000000000-4edb10f3fc54c757-1-2-00000000-basin_staging.eddie.data6-1.parquet",
			expected:   Pub{Namespace: "eddie", Relation: "data6"},
			shouldFail: false,
		},
		{
			filename:   "invalid_filename_format.parquet",
			expected:   Pub{},
			shouldFail: true,
		},
	}

	crdbClient, _ := NewDB("")

	for _, tt := range tests {
		result, err := crdbClient.extractPubName(tt.filename)
		if tt.shouldFail {
			if err == nil {
				t.Errorf("expected %s to fail, but it didn't", tt.filename)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for %s: %s", tt.filename, err)
			}
			if result != tt.expected {
				t.Errorf("expected %s to be %s, but got %s", tt.filename, tt.expected, result)
			}
		}
	}
}
