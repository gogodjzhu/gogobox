package timefmt

import (
	"testing"
	"time"
)

func TestTimeFormatter_ParseInput(t *testing.T) {
	tf := &TimeFormatter{}

	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "Unix timestamp in seconds",
			input:    "1640995200",
			expected: time.Unix(1640995200, 0),
			wantErr:  false,
		},
		{
			name:     "Unix timestamp in milliseconds",
			input:    "1640995200000",
			expected: time.Unix(1640995200, 0),
			wantErr:  false,
		},
		{
			name:     "RFC3339 format",
			input:    "2022-01-01T00:00:00Z",
			expected: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Date only format",
			input:    "2022-01-01",
			expected: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "DateTime format",
			input:    "2022-01-01 15:04:05",
			expected: time.Date(2022, 1, 1, 15, 4, 5, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "US date format",
			input:    "01/01/2022",
			expected: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:    "Invalid input",
			input:   "invalid-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tf.ParseInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equal(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTimeFormatter_FormatTime(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		outputFormat string
		expected     string
		wantErr      bool
	}{
		{
			name:         "Format timestamp to default",
			input:        "1640995200",
			outputFormat: "2006-01-02 15:04:05",
			expected:     "2022-01-01 00:00:00",
			wantErr:      false,
		},
		{
			name:         "Format date string to different format",
			input:        "2022-01-01",
			outputFormat: "Jan 2, 2006",
			expected:     "Jan 1, 2022",
			wantErr:      false,
		},
		{
			name:         "Format with time only",
			input:        "2022-01-01 15:04:05",
			outputFormat: "15:04:05",
			expected:     "15:04:05",
			wantErr:      false,
		},
		{
			name:    "Invalid input",
			input:   "invalid-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &TimeFormatter{
				outputFormat: tt.outputFormat,
				timezone:     "UTC", // Force UTC for consistent testing
			}
			result, err := tf.FormatTime(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestTimeFormatter_ConvertToTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		unit     string
		expected int64
		wantErr  bool
	}{
		{
			name:     "Convert date to seconds timestamp",
			input:    "2022-01-01T00:00:00Z",
			unit:     "s",
			expected: 1640995200,
			wantErr:  false,
		},
		{
			name:     "Convert date to milliseconds timestamp",
			input:    "2022-01-01T00:00:00Z",
			unit:     "ms",
			expected: 1640995200000,
			wantErr:  false,
		},
		{
			name:     "Convert timestamp to different unit",
			input:    "1640995200",
			unit:     "ms",
			expected: 1640995200000,
			wantErr:  false,
		},
		{
			name:    "Invalid unit",
			input:   "2022-01-01",
			unit:    "invalid",
			wantErr: true,
		},
		{
			name:    "Invalid input",
			input:   "invalid-date",
			unit:    "s",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &TimeFormatter{}
			result, err := tf.ConvertToTimestamp(tt.input, tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestTimeFormatter_WithTimezone(t *testing.T) {
	tf := &TimeFormatter{
		outputFormat: "2006-01-02 15:04:05 MST",
		timezone:     "America/New_York",
	}

	// Test with UTC input
	result, err := tf.FormatTime("2022-01-01T12:00:00Z")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should convert UTC 12:00 to EST 07:00 (UTC-5 in January)
	expected := "2022-01-01 07:00:00 EST"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestTimeFormatter_WithInputFormat(t *testing.T) {
	tf := &TimeFormatter{
		inputFormat:  "02/01/2006", // DD/MM/YYYY
		outputFormat: "2006-01-02",
	}

	result, err := tf.FormatTime("25/12/2022")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "2022-12-25"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
