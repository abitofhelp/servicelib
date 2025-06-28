//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Example usage of the Time value object
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/abitofhelp/servicelib/valueobject/temporal"
)

func main() {
	fmt.Println("=== Time Value Object Example ===")

	// Create a new time from string
	timeStr := "2025-06-22T01:39:59-07:00"
	t1, err := temporal.NewTimeFromString(timeStr)
	if err != nil {
		fmt.Println("Error creating time:", err)
		return
	}

	// Create a new time from time.Time
	now := time.Now()
	t2, err := temporal.NewTime(now)
	if err != nil {
		fmt.Println("Error creating time:", err)
		return
	}

	// Format as RFC 3339 string
	fmt.Printf("Time 1: %s\n", t1.String())
	fmt.Printf("Time 2: %s\n", t2.String())

	// Compare times
	fmt.Printf("\n=== Time Comparison ===\n")
	fmt.Printf("Time 1 is before Time 2: %v\n", t1.Before(t2))
	fmt.Printf("Time 1 is after Time 2: %v\n", t1.After(t2))

	// Using the new Equals method (preferred over Equal)
	fmt.Printf("Times are equal (using Equals): %v\n", t1.Equals(t2))

	// Create a copy of t1 to demonstrate equality
	t1Copy, _ := temporal.NewTimeFromString(timeStr)
	fmt.Printf("Time 1 equals its copy: %v\n", t1.Equals(t1Copy))

	// Time arithmetic
	fmt.Printf("\n=== Time Arithmetic ===\n")
	t3 := t1.Add(24 * time.Hour)
	fmt.Printf("Time 1 + 24 hours: %s\n", t3.String())

	// Get duration between times
	duration := t2.Sub(t1)
	fmt.Printf("Duration between Time 2 and Time 1: %v\n", duration)

	// Check if time is empty
	fmt.Printf("\n=== Empty Time ===\n")
	emptyTime, _ := temporal.NewTime(time.Time{})
	fmt.Printf("Is empty time empty? %v\n", emptyTime.IsEmpty())

	// Validation
	fmt.Printf("\n=== Validation ===\n")
	err = t1.Validate()
	fmt.Printf("Time 1 validation: %v\n", err == nil)
	err = emptyTime.Validate()
	fmt.Printf("Empty time validation: %v\n", err == nil)

	// JSON marshaling and unmarshaling
	fmt.Printf("\n=== JSON Serialization ===\n")

	// Marshal to JSON
	jsonBytes, err := json.Marshal(t1)
	if err != nil {
		fmt.Println("Error marshaling time:", err)
		return
	}
	fmt.Printf("Time 1 as JSON: %s\n", string(jsonBytes))

	// Unmarshal from JSON
	var unmarshaledTime temporal.Time
	err = json.Unmarshal(jsonBytes, &unmarshaledTime)
	if err != nil {
		fmt.Println("Error unmarshaling time:", err)
		return
	}
	fmt.Printf("Unmarshaled time: %s\n", unmarshaledTime.String())
	fmt.Printf("Original equals unmarshaled: %v\n", t1.Equals(unmarshaledTime))

	// Example with a struct containing Time
	type Event struct {
		ID        string        `json:"id"`
		Name      string        `json:"name"`
		StartTime temporal.Time `json:"start_time"`
		EndTime   temporal.Time `json:"end_time"`
	}

	// Create an event
	event := Event{
		ID:        "evt-123",
		Name:      "Conference",
		StartTime: t1,
		EndTime:   t3,
	}

	// Marshal the event to JSON
	eventJSON, _ := json.Marshal(event)
	fmt.Printf("\nEvent as JSON: %s\n", string(eventJSON))

	// Unmarshal back to an event
	var unmarshaledEvent Event
	_ = json.Unmarshal(eventJSON, &unmarshaledEvent)
	fmt.Printf("Event duration: %v\n", unmarshaledEvent.EndTime.Sub(unmarshaledEvent.StartTime))
}