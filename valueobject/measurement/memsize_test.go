// Copyright (c) 2025 A Bit of Help, Inc.

package measurement

import (
	"testing"

	"github.com/abitofhelp/servicelib/valueobject/base"
	"github.com/stretchr/testify/assert"
)

func TestNewMemSize(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		unit      MemSizeUnit
		wantBytes uint64
		wantErr   bool
	}{
		{
			name:      "Valid bytes",
			value:     1024,
			unit:      Bytes,
			wantBytes: 1024,
			wantErr:   false,
		},
		{
			name:      "Valid kibibytes",
			value:     1,
			unit:      Kibibytes,
			wantBytes: 1024,
			wantErr:   false,
		},
		{
			name:      "Valid mebibytes",
			value:     1,
			unit:      Mebibytes,
			wantBytes: 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid gibibytes",
			value:     1,
			unit:      Gibibytes,
			wantBytes: 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid tebibytes",
			value:     1,
			unit:      Tebibytes,
			wantBytes: 1024 * 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid pebibytes",
			value:     1,
			unit:      Pebibytes,
			wantBytes: 1024 * 1024 * 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Negative value",
			value:     -1,
			unit:      Bytes,
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
			got, err := NewMemSize(tt.value, tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, got.Bytes())
			}
		})
	}
}

func TestParseMemSize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBytes uint64
		wantErr   bool
	}{
		{
			name:      "Valid bytes",
			input:     "1024B",
			wantBytes: 1024,
			wantErr:   false,
		},
		{
			name:      "Valid kibibytes",
			input:     "1KiB",
			wantBytes: 1024,
			wantErr:   false,
		},
		{
			name:      "Valid mebibytes",
			input:     "1MiB",
			wantBytes: 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid gibibytes",
			input:     "1GiB",
			wantBytes: 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid tebibytes",
			input:     "1TiB",
			wantBytes: 1024 * 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Valid pebibytes",
			input:     "1PiB",
			wantBytes: 1024 * 1024 * 1024 * 1024 * 1024,
			wantErr:   false,
		},
		{
			name:      "Legacy KB",
			input:     "1KB",
			wantBytes: 1024,
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
			got, err := ParseMemSize(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, got.Bytes())
			}
		})
	}
}

func TestMemSize_String(t *testing.T) {
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
			name:  "Kibibytes",
			bytes: 1024,
			want:  "1.00 KiB",
		},
		{
			name:  "Mebibytes",
			bytes: 1024 * 1024,
			want:  "1.00 MiB",
		},
		{
			name:  "Gibibytes",
			bytes: 1024 * 1024 * 1024,
			want:  "1.00 GiB",
		},
		{
			name:  "Tebibytes",
			bytes: 1024 * 1024 * 1024 * 1024,
			want:  "1.00 TiB",
		},
		{
			name:  "Pebibytes",
			bytes: 1024 * 1024 * 1024 * 1024 * 1024,
			want:  "1.00 PiB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemSize{
				BaseStructValueObject: base.BaseStructValueObject{},
				bytes:                 tt.bytes,
			}
			assert.Equal(t, tt.want, ms.String())
		})
	}
}

func TestMemSize_Equals(t *testing.T) {
	ms1 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024,
	}
	ms2 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024,
	}
	ms3 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 2048,
	}

	assert.True(t, ms1.Equals(ms2))
	assert.False(t, ms1.Equals(ms3))
}

func TestMemSize_IsEmpty(t *testing.T) {
	ms1 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 0,
	}
	ms2 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024,
	}

	assert.True(t, ms1.IsEmpty())
	assert.False(t, ms2.IsEmpty())
}

func TestMemSize_Conversions(t *testing.T) {
	ms := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024 * 1024 * 1024,
	} // 1 GiB

	assert.Equal(t, uint64(1024*1024*1024), ms.Bytes())
	assert.InDelta(t, 1024*1024, ms.Kibibytes(), 0.001)
	assert.InDelta(t, 1024, ms.Mebibytes(), 0.001)
	assert.InDelta(t, 1.0, ms.Gibibytes(), 0.001)
	assert.InDelta(t, 1.0/1024, ms.Tebibytes(), 0.001)
	assert.InDelta(t, 1.0/(1024*1024), ms.Pebibytes(), 0.001)
}

func TestMemSize_Format(t *testing.T) {
	ms := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024 * 1024,
	} // 1 MiB

	assert.Equal(t, "1048576 B", ms.Format("B"))
	assert.Equal(t, "1024.00 KiB", ms.Format("KiB"))
	assert.Equal(t, "1.00 MiB", ms.Format("MiB"))
	assert.Equal(t, "0.00 GiB", ms.Format("GiB"))
	assert.Equal(t, "0.00 TiB", ms.Format("TiB"))
	assert.Equal(t, "0.00 PiB", ms.Format("PiB"))
	assert.Equal(t, "1.00 MiB", ms.Format("auto"))
	assert.Equal(t, "1.00 MiB", ms.Format("invalid"))
}

func TestMemSize_Operations(t *testing.T) {
	ms1 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 1024,
	}
	ms2 := MemSize{
		BaseStructValueObject: base.BaseStructValueObject{},
		bytes:                 2048,
	}

	// Add
	result := ms1.Add(ms2)
	assert.Equal(t, uint64(3072), result.Bytes())

	// Subtract
	result = ms2.Subtract(ms1)
	assert.Equal(t, uint64(1024), result.Bytes())

	// Subtract with negative result
	result = ms1.Subtract(ms2)
	assert.Equal(t, uint64(0), result.Bytes())

	// IsLargerThan
	assert.True(t, ms2.IsLargerThan(ms1))
	assert.False(t, ms1.IsLargerThan(ms2))

	// IsSmallerThan
	assert.True(t, ms1.IsSmallerThan(ms2))
	assert.False(t, ms2.IsSmallerThan(ms1))
}
