package storage

import (
	"testing"
)

func TestExtractPub(t *testing.T) {
	filename := "esfbmltndstj/ksvraapqfiyf/export17860a3b03221a1b0000000000000001-n901064813195493377.0.parquet"
	tests := []struct {
		filename   string
		expected   Pub
		shouldFail bool
	}{
		{
			filename:   filename,
			expected:   Pub{Namespace: "esfbmltndstj", Relation: "ksvraapqfiyf"},
			shouldFail: false,
		},
		{
			filename:   "invalid_filename_format.parquet",
			expected:   Pub{},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		result, err := extractPub(tt.filename)
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
