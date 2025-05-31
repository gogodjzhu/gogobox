package timefmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// Common time format patterns
var commonFormats = []string{
	time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	time.RFC822,      // "02 Jan 06 15:04 MST"
	time.RFC822Z,     // "02 Jan 06 15:04 -0700"
	time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
	time.Kitchen,     // "3:04PM"
	time.Stamp,       // "Jan _2 15:04:05"
	time.StampMilli,  // "Jan _2 15:04:05.000"
	time.StampMicro,  // "Jan _2 15:04:05.000000"
	time.StampNano,   // "Jan _2 15:04:05.000000000"
	time.DateTime,    // "2006-01-02 15:04:05"
	time.DateOnly,    // "2006-01-02"
	time.TimeOnly,    // "15:04:05"
	// Additional common formats
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02T15:04:05",
	"2006/01/02",
	"01/02/2006",
	"02/01/2006",
	"2006-01-02 15:04",
	"2006/01/02 15:04",
	"Jan 2, 2006",
	"January 2, 2006",
	"2 Jan 2006",
	"2 January 2006",
	"02-01-2006 15:04:05",
	"02-01-2006",
}

// TimeFormatter handles various time format conversions
type TimeFormatter struct {
	inputFormat  string
	outputFormat string
	timezone     string
}

// NewCmdTimeFmt creates a new time format command
func NewCmdTimeFmt(f *cmdutil.Factory) *cobra.Command {
	var inputFormat, outputFormat, timezone string
	var outputTimestamp bool
	var timestampUnit string

	cmd := &cobra.Command{
		Use:   "timefmt [input]",
		Short: "Format and convert time between various formats",
		Long: `Format and convert time between various formats.

Supports:
- Input: time string (various patterns) or timestamp (in milliseconds or seconds)
- Output: formatted time string or timestamp (in milliseconds or seconds)

Examples:
  # Parse timestamp and format as date
  gogobox timefmt 1640995200000

  # Parse date string and convert to timestamp
  gogobox timefmt "2022-01-01 00:00:00" --timestamp

  # Convert between different string formats
  gogobox timefmt "2022-01-01" --output-format "Jan 2, 2006"

  # Parse with specific input format
  gogobox timefmt "01/01/2022" --input-format "01/02/2006"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			formatter := &TimeFormatter{
				inputFormat:  inputFormat,
				outputFormat: outputFormat,
				timezone:     timezone,
			}

			if outputTimestamp {
				timestamp, err := formatter.ConvertToTimestamp(input, timestampUnit)
				if err != nil {
					return fmt.Errorf("failed to convert to timestamp: %w", err)
				}
				fmt.Println(timestamp)
			} else {
				formatted, err := formatter.FormatTime(input)
				if err != nil {
					return fmt.Errorf("failed to format time: %w", err)
				}
				fmt.Println(formatted)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFormat, "input-format", "i", "", "Input time format (auto-detected if not specified)")
	cmd.Flags().StringVarP(&outputFormat, "output-format", "o", "2006-01-02 15:04:05", "Output time format")
	cmd.Flags().StringVarP(&timezone, "timezone", "z", "", "Timezone for output (e.g., 'UTC', 'America/New_York')")
	cmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Output as timestamp instead of formatted string")
	cmd.Flags().StringVarP(&timestampUnit, "unit", "u", "ms", "Timestamp unit: 'ms' (milliseconds) or 's' (seconds)")

	return cmd
}

// ParseInput parses the input string and returns a time.Time
func (tf *TimeFormatter) ParseInput(input string) (time.Time, error) {
	// Try to parse as timestamp first
	if timestamp, err := tf.parseTimestamp(input); err == nil {
		return timestamp, nil
	}

	// If input format is specified, use it
	if tf.inputFormat != "" {
		return time.Parse(tf.inputFormat, input)
	}

	// Try common formats
	for _, format := range commonFormats {
		if t, err := time.Parse(format, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time input: %s", input)
}

// parseTimestamp attempts to parse input as a timestamp
func (tf *TimeFormatter) parseTimestamp(input string) (time.Time, error) {
	// Remove any whitespace
	input = strings.TrimSpace(input)

	// Try to parse as integer
	timestamp, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// Determine if it's seconds or milliseconds based on magnitude
	// Timestamps after year 2001 in seconds: > 1000000000
	// Timestamps before year 2286 in milliseconds: < 10000000000000
	if timestamp > 1000000000 && timestamp < 10000000000 {
		// Likely seconds
		return time.Unix(timestamp, 0), nil
	} else if timestamp >= 10000000000 {
		// Likely milliseconds
		return time.Unix(timestamp/1000, (timestamp%1000)*1000000), nil
	}

	return time.Time{}, fmt.Errorf("timestamp out of reasonable range: %d", timestamp)
}

// FormatTime formats the input time to the specified output format
func (tf *TimeFormatter) FormatTime(input string) (string, error) {
	t, err := tf.ParseInput(input)
	if err != nil {
		return "", err
	}

	// Apply timezone if specified
	if tf.timezone != "" {
		loc, err := time.LoadLocation(tf.timezone)
		if err != nil {
			return "", fmt.Errorf("invalid timezone: %s", tf.timezone)
		}
		t = t.In(loc)
	}

	return t.Format(tf.outputFormat), nil
}

// ConvertToTimestamp converts input to timestamp in specified unit
func (tf *TimeFormatter) ConvertToTimestamp(input, unit string) (int64, error) {
	t, err := tf.ParseInput(input)
	if err != nil {
		return 0, err
	}

	// Apply timezone if specified
	if tf.timezone != "" {
		loc, err := time.LoadLocation(tf.timezone)
		if err != nil {
			return 0, fmt.Errorf("invalid timezone: %s", tf.timezone)
		}
		t = t.In(loc)
	}

	switch strings.ToLower(unit) {
	case "s", "sec", "second", "seconds":
		return t.Unix(), nil
	case "ms", "milli", "millisecond", "milliseconds":
		return t.UnixMilli(), nil
	default:
		return 0, fmt.Errorf("unsupported timestamp unit: %s (use 's' or 'ms')", unit)
	}
}
