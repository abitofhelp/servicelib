// Copyright (c) 2025 A Bit of Help, Inc.

package measurement

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileSize(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		unit      FileSizeUnit
		wantBytes uint64
		wantErr   bool
	}{
		{
			name:      "Valid bytes",
			value:     1000,
			unit:      FileSizeBytes,
			wantBytes: 1000,
			wantErr:   false,
		},
		{
			name:      "Valid kilobytes",
			value:     1,
			unit:      Kilobytes,
			wantBytes: 1000,
			wantErr:   false,
		},
		{
			name:      "Valid megabytes",
			value:     1,
			unit:      Megabytes,
			wantBytes: 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Valid gigabytes",
			value:     1,
			unit:      Gigabytes,
			wantBytes: 1000 * 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Valid terabytes",
			value:     1,
			unit:      Terabytes,
			wantBytes: 1000 * 1000 * 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Valid petabytes",
			value:     1,
			unit:      Petabytes,
			wantBytes: 1000 * 1000 * 1000 * 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Negative value",
			value:     -1,
			unit:      FileSizeBytes,
			wantBytes: 0,
			wantErr:   true,
		},
		{
			name:      "Invalid unit",
			value:     1,
			unit:      "invalid",
			wantBytes: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileSize(tt.value, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, got.Bytes())
			}
		})
	}
}

func TestFileSize_String(t *testing.T) {
	tests := []struct {
		name  string
		bytes uint64
		want  string
	}{
		{
			name:  "Bytes",
			bytes: 500,
			want:  "500 B",
		},
		{
			name:  "Kilobytes",
			bytes: 1500,
			want:  "1.50 KB",
		},
		{
			name:  "Megabytes",
			bytes: 1500000,
			want:  "1.50 MB",
		},
		{
			name:  "Gigabytes",
			bytes: 1500000000,
			want:  "1.50 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FileSize{bytes: tt.bytes}
			assert.Equal(t, tt.want, fs.String())
		})
	}
}

func TestFileSize_Equals(t *testing.T) {
	fs1 := FileSize{bytes: 1000}
	fs2 := FileSize{bytes: 1000}
	fs3 := FileSize{bytes: 2000}

	assert.True(t, fs1.Equals(fs2))
	assert.False(t, fs1.Equals(fs3))
}

func TestFileSize_IsEmpty(t *testing.T) {
	emptyValue := FileSize{bytes: 0}
	nonEmptyValue := FileSize{bytes: 1000}

	assert.True(t, emptyValue.IsEmpty())
	assert.False(t, nonEmptyValue.IsEmpty())
}

func TestFileSize_Validate(t *testing.T) {
	// FileSize.Validate() always returns nil since there are no additional validation rules
	// beyond what's checked during creation
	fs := FileSize{bytes: 1000}
	assert.NoError(t, fs.Validate())
}

func TestFileSize_ToMap(t *testing.T) {
	fs := FileSize{bytes: 1000}

	expected := map[string]interface{}{
		"bytes": uint64(1000),
	}

	assert.Equal(t, expected, fs.ToMap())
}

func TestFileSize_MarshalJSON(t *testing.T) {
	fs := FileSize{bytes: 1000}

	data, err := fs.MarshalJSON()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"bytes": float64(1000), // JSON numbers are unmarshaled as float64
	}

	assert.Equal(t, expected, result)
}

func TestFileSize_Conversions(t *testing.T) {
	fs := FileSize{bytes: 1000 * 1000 * 1000} // 1 GB

	assert.Equal(t, uint64(1000*1000*1000), fs.Bytes())
	assert.InDelta(t, 1000*1000, fs.Kilobytes(), 0.001)
	assert.InDelta(t, 1000, fs.Megabytes(), 0.001)
	assert.InDelta(t, 1.0, fs.Gigabytes(), 0.001)
	assert.InDelta(t, 1.0/1000, fs.Terabytes(), 0.001)
	assert.InDelta(t, 1.0/(1000*1000), fs.Petabytes(), 0.001)
}

func TestFileSize_Format(t *testing.T) {
	fs := FileSize{bytes: 1000 * 1000} // 1 MB

	assert.Equal(t, "1000000 B", fs.Format("B"))
	assert.Equal(t, "1000.00 KB", fs.Format("KB"))
	assert.Equal(t, "1.00 MB", fs.Format("MB"))
	assert.Equal(t, "0.00 GB", fs.Format("GB"))
	assert.Equal(t, "0.00 TB", fs.Format("TB"))
	assert.Equal(t, "0.00 PB", fs.Format("PB"))
	assert.Equal(t, "1.00 MB", fs.Format("auto"))
	assert.Equal(t, "1.00 MB", fs.Format("invalid"))
}

func TestFileSize_Operations(t *testing.T) {
	fs1 := FileSize{bytes: 1000}
	fs2 := FileSize{bytes: 2000}

	// Add
	result := fs1.Add(fs2)
	assert.Equal(t, uint64(3000), result.Bytes())

	// Subtract
	result = fs2.Subtract(fs1)
	assert.Equal(t, uint64(1000), result.Bytes())

	// Subtract with negative result
	result = fs1.Subtract(fs2)
	assert.Equal(t, uint64(0), result.Bytes())

	// IsLargerThan
	assert.True(t, fs2.IsLargerThan(fs1))
	assert.False(t, fs1.IsLargerThan(fs2))

	// IsSmallerThan
	assert.True(t, fs1.IsSmallerThan(fs2))
	assert.False(t, fs2.IsSmallerThan(fs1))
}

func TestParseFileSize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBytes uint64
		wantErr   bool
	}{
		{
			name:      "Valid bytes",
			input:     "1000B",
			wantBytes: 1000,
			wantErr:   false,
		},
		{
			name:      "Valid kilobytes",
			input:     "1KB",
			wantBytes: 1000,
			wantErr:   false,
		},
		{
			name:      "Valid megabytes",
			input:     "1MB",
			wantBytes: 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Valid gigabytes",
			input:     "1GB",
			wantBytes: 1000 * 1000 * 1000,
			wantErr:   false,
		},
		{
			name:      "Empty string",
			input:     "",
			wantBytes: 0,
			wantErr:   true,
		},
		{
			name:      "Invalid format",
			input:     "invalid",
			wantBytes: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFileSize(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, got.Bytes())
			}
		})
	}
}
