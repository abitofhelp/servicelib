// Copyright (c) 2025 A Bit of Help, Inc.

package temporal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTime(t *testing.T) {
	// Test with valid time
	now := time.Now()
	timeObj, err := NewTime(now)
	assert.NoError(t, err)
	assert.Equal(t, now, timeObj.Value())

	// Test with zero time
	zeroTime := time.Time{}
	timeObj, err = NewTime(zeroTime)
	assert.NoError(t, err)
	assert.True(t, timeObj.IsEmpty())
}

func TestNewTimeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid RFC3339",
			input:   "2025-06-22T01:43:35-07:00",
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Invalid format",
			input:   "2025-06-22",
			wantErr: true,
		},
		{
			name:    "Invalid time",
			input:   "2025-13-35T25:61:61-07:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeObj, err := NewTimeFromString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, timeObj.String())
			}
		})
	}
}

func TestTime_String(t *testing.T) {
	timeStr := "2025-06-22T01:43:35-07:00"
	timeObj, err := NewTimeFromString(timeStr)
	assert.NoError(t, err)
	assert.Equal(t, timeStr, timeObj.String())
}

func TestTime_IsEmpty(t *testing.T) {
	// Test non-empty time
	nonEmptyTime, err := NewTimeFromString("2025-06-22T01:43:35-07:00")
	assert.NoError(t, err)
	assert.False(t, nonEmptyTime.IsEmpty())

	// Test empty time
	emptyTime, err := NewTime(time.Time{})
	assert.NoError(t, err)
	assert.True(t, emptyTime.IsEmpty())
}

func TestTime_Equals(t *testing.T) {
	time1, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	time2, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	time3, _ := NewTimeFromString("2025-06-22T01:43:36-07:00")

	assert.True(t, time1.Equals(time2))
	assert.False(t, time1.Equals(time3))
}

func TestTime_Equal(t *testing.T) {
	// Test that Equal delegates to Equals
	time1, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	time2, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	time3, _ := NewTimeFromString("2025-06-22T01:43:36-07:00")

	assert.True(t, time1.Equal(time2))
	assert.False(t, time1.Equal(time3))
}

func TestTime_BeforeAfter(t *testing.T) {
	earlier, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	later, _ := NewTimeFromString("2025-06-22T01:43:36-07:00")

	assert.True(t, earlier.Before(later))
	assert.False(t, earlier.After(later))
	assert.True(t, later.After(earlier))
	assert.False(t, later.Before(earlier))
}

func TestTime_Add(t *testing.T) {
	start, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")

	// Add one hour
	plusHour := start.Add(time.Hour)
	expected, _ := NewTimeFromString("2025-06-22T02:43:35-07:00")
	assert.True(t, plusHour.Equals(expected))

	// Add negative duration (subtract)
	minusHour := start.Add(-time.Hour)
	expected, _ = NewTimeFromString("2025-06-22T00:43:35-07:00")
	assert.True(t, minusHour.Equals(expected))
}

func TestTime_Sub(t *testing.T) {
	time1, _ := NewTimeFromString("2025-06-22T01:43:35-07:00")
	time2, _ := NewTimeFromString("2025-06-22T02:43:35-07:00")

	// Positive duration
	duration := time2.Sub(time1)
	assert.Equal(t, time.Hour, duration)

	// Negative duration
	duration = time1.Sub(time2)
	assert.Equal(t, -time.Hour, duration)

	// Zero duration
	duration = time1.Sub(time1)
	assert.Equal(t, time.Duration(0), duration)
}

func TestTime_ToMap(t *testing.T) {
	timeStr := "2025-06-22T01:43:35-07:00"
	timeObj, err := NewTimeFromString(timeStr)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"value": timeStr,
	}
	assert.Equal(t, expected, timeObj.ToMap())
}

func TestTime_MarshalJSON(t *testing.T) {
	timeStr := "2025-06-22T01:43:35-07:00"
	timeObj, err := NewTimeFromString(timeStr)
	assert.NoError(t, err)

	// Test marshaling
	jsonBytes, err := json.Marshal(timeObj)
	assert.NoError(t, err)
	assert.Equal(t, `"`+timeStr+`"`, string(jsonBytes))

	// Test in a struct
	type TestStruct struct {
		Time Time `json:"time"`
	}
	testStruct := TestStruct{Time: timeObj}

	jsonBytes, err = json.Marshal(testStruct)
	assert.NoError(t, err)
	assert.Equal(t, `{"time":"`+timeStr+`"}`, string(jsonBytes))
}

func TestTime_UnmarshalJSON(t *testing.T) {
	timeStr := "2025-06-22T01:43:35-07:00"
	jsonStr := `"` + timeStr + `"`

	// Test unmarshaling directly
	var timeObj Time
	err := json.Unmarshal([]byte(jsonStr), &timeObj)
	assert.NoError(t, err)
	assert.Equal(t, timeStr, timeObj.String())

	// Test with invalid JSON
	err = json.Unmarshal([]byte(`"invalid time"`), &timeObj)
	assert.Error(t, err)

	// Test in a struct
	jsonStructStr := `{"time":"` + timeStr + `"}`
	type TestStruct struct {
		Time Time `json:"time"`
	}

	var testStruct TestStruct
	err = json.Unmarshal([]byte(jsonStructStr), &testStruct)
	assert.NoError(t, err)
	assert.Equal(t, timeStr, testStruct.Time.String())
}

func TestTime_Validate(t *testing.T) {
	// Test with valid time
	timeObj, err := NewTimeFromString("2025-06-22T01:43:35-07:00")
	assert.NoError(t, err)
	assert.NoError(t, timeObj.Validate())

	// Test with zero time
	zeroTime, err := NewTime(time.Time{})
	assert.NoError(t, err)
	assert.NoError(t, zeroTime.Validate())
}
