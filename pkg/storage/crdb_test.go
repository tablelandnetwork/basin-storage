package storage

import (
	"testing"
)

func TestExtractTableName(t *testing.T) {
	tests := []struct {
		filename   string
		expected   string
		shouldFail bool
	}{
		{
			filename:   "202308291525552525242120000000000-3ab461ed932d5f1c-1-2-00000001-office_dogs-2.parquet",
			expected:   "office_dogs",
			shouldFail: false,
		},
		{
			filename:   "2020-04-02/202004022058072107140000000000000-56087568dba1e6b8-1-72-00000000-test_table-1.ndjson",
			expected:   "test_table",
			shouldFail: false,
		},
		{
			filename:   "invalid_filename_format.parquet",
			expected:   "",
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
