// Copyright (c) 2025 A Bit of Help, Inc.

package model

import (
	"reflect"
	"testing"
)

// TestCopyFields tests the CopyFields function with various scenarios
func TestCopyFields(t *testing.T) {
	// Define test structures
	type Source struct {
		Name    string
		Age     int
		Address string
		private string // unexported field
	}

	type Destination struct {
		Name    string
		Age     int
		Address string
		private string // unexported field
		Extra   string // field not in source
	}

	// Test case 1: Normal copy
	t.Run("Normal copy", func(t *testing.T) {
		src := &Source{
			Name:    "John Doe",
			Age:     30,
			Address: "123 Main St",
			private: "private data",
		}
		dst := &Destination{}

		err := CopyFields(dst, src)
		if err != nil {
			t.Errorf("CopyFields returned error: %v", err)
		}

		if dst.Name != src.Name {
			t.Errorf("Name field not copied correctly, got: %s, want: %s", dst.Name, src.Name)
		}
		if dst.Age != src.Age {
			t.Errorf("Age field not copied correctly, got: %d, want: %d", dst.Age, src.Age)
		}
		if dst.Address != src.Address {
			t.Errorf("Address field not copied correctly, got: %s, want: %s", dst.Address, src.Address)
		}
		if dst.private != "" {
			t.Errorf("private field should not be copied, got: %s", dst.private)
		}
		if dst.Extra != "" {
			t.Errorf("Extra field should remain empty, got: %s", dst.Extra)
		}
	})

	// Test case 2: Non-pointer arguments
	t.Run("Non-pointer arguments", func(t *testing.T) {
		src := Source{Name: "John Doe"}
		dst := Destination{}

		err := CopyFields(dst, src)
		if err == nil {
			t.Error("Expected error for non-pointer arguments, got nil")
		}
	})

	// Test case 3: Non-struct pointers
	t.Run("Non-struct pointers", func(t *testing.T) {
		src := "string"
		dst := 123

		err := CopyFields(&dst, &src)
		if err == nil {
			t.Error("Expected error for non-struct pointers, got nil")
		}
	})

	// Test case 4: Different field types
	t.Run("Different field types", func(t *testing.T) {
		type SourceDiffType struct {
			Name string
			Age  string // string instead of int
		}

		src := &SourceDiffType{
			Name: "John Doe",
			Age:  "30",
		}
		dst := &Destination{}

		err := CopyFields(dst, src)
		if err != nil {
			t.Errorf("CopyFields returned error: %v", err)
		}

		if dst.Name != src.Name {
			t.Errorf("Name field not copied correctly, got: %s, want: %s", dst.Name, src.Name)
		}
		if dst.Age != 0 { // Age should not be copied due to type mismatch
			t.Errorf("Age field should not be copied due to type mismatch, got: %d", dst.Age)
		}
	})
}

// TestDeepCopy tests the DeepCopy function with various scenarios
func TestDeepCopy(t *testing.T) {
	// Test case 1: Simple struct
	t.Run("Simple struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		src := &Person{Name: "John Doe", Age: 30}
		dst := &Person{}

		err := DeepCopy(dst, src)
		if err != nil {
			t.Errorf("DeepCopy returned error: %v", err)
		}

		if !reflect.DeepEqual(src, dst) {
			t.Errorf("DeepCopy failed, got: %+v, want: %+v", dst, src)
		}

		// Verify it's a deep copy by modifying source
		src.Name = "Jane Doe"
		if dst.Name == src.Name {
			t.Error("DeepCopy did not create a separate copy")
		}
	})

	// Test case 2: Struct with nested struct
	t.Run("Struct with nested struct", func(t *testing.T) {
		type Address struct {
			Street string
			City   string
		}

		type Person struct {
			Name    string
			Age     int
			Address *Address
		}

		src := &Person{
			Name: "John Doe",
			Age:  30,
			Address: &Address{
				Street: "123 Main St",
				City:   "Anytown",
			},
		}
		dst := &Person{}

		err := DeepCopy(dst, src)
		if err != nil {
			t.Errorf("DeepCopy returned error: %v", err)
		}

		if !reflect.DeepEqual(src, dst) {
			t.Errorf("DeepCopy failed, got: %+v, want: %+v", dst, src)
		}

		// Verify it's a deep copy by modifying source's nested struct
		src.Address.Street = "456 Oak Ave"
		if dst.Address.Street == src.Address.Street {
			t.Error("DeepCopy did not create a separate copy of nested struct")
		}
	})

	// Test case 3: Struct with slice
	t.Run("Struct with slice", func(t *testing.T) {
		type Person struct {
			Name    string
			Hobbies []string
		}

		src := &Person{
			Name:    "John Doe",
			Hobbies: []string{"Reading", "Hiking"},
		}
		dst := &Person{}

		err := DeepCopy(dst, src)
		if err != nil {
			t.Errorf("DeepCopy returned error: %v", err)
		}

		if !reflect.DeepEqual(src, dst) {
			t.Errorf("DeepCopy failed, got: %+v, want: %+v", dst, src)
		}

		// Verify it's a deep copy by modifying source's slice
		src.Hobbies[0] = "Swimming"
		if dst.Hobbies[0] == src.Hobbies[0] {
			t.Error("DeepCopy did not create a separate copy of slice")
		}
	})

	// Test case 4: Struct with map
	t.Run("Struct with map", func(t *testing.T) {
		type Person struct {
			Name       string
			Properties map[string]string
		}

		src := &Person{
			Name: "John Doe",
			Properties: map[string]string{
				"hair": "brown",
				"eyes": "blue",
			},
		}
		dst := &Person{}

		err := DeepCopy(dst, src)
		if err != nil {
			t.Errorf("DeepCopy returned error: %v", err)
		}

		if !reflect.DeepEqual(src, dst) {
			t.Errorf("DeepCopy failed, got: %+v, want: %+v", dst, src)
		}

		// Verify it's a deep copy by modifying source's map
		src.Properties["hair"] = "black"
		if dst.Properties["hair"] == src.Properties["hair"] {
			t.Error("DeepCopy did not create a separate copy of map")
		}
	})

	// Test case 5: Non-pointer arguments
	t.Run("Non-pointer arguments", func(t *testing.T) {
		type Person struct {
			Name string
		}

		src := Person{Name: "John Doe"}
		dst := Person{}

		err := DeepCopy(dst, src)
		if err == nil {
			t.Error("Expected error for non-pointer arguments, got nil")
		}
	})

	// Test case 6: Different struct types
	t.Run("Different struct types", func(t *testing.T) {
		type Person struct {
			Name string
		}

		type Employee struct {
			Title string
		}

		src := &Person{Name: "John Doe"}
		dst := &Employee{}

		err := DeepCopy(dst, src)
		if err == nil {
			t.Error("Expected error for different struct types, got nil")
		}
	})

	// Test case 7: Nil values
	t.Run("Nil values in structs", func(t *testing.T) {
		type Person struct {
			Name  string
			Child *Person
		}

		src := &Person{
			Name:  "John Doe",
			Child: nil,
		}
		dst := &Person{}

		err := DeepCopy(dst, src)
		if err != nil {
			t.Errorf("DeepCopy returned error: %v", err)
		}

		if !reflect.DeepEqual(src, dst) {
			t.Errorf("DeepCopy failed, got: %+v, want: %+v", dst, src)
		}
		if dst.Child != nil {
			t.Error("DeepCopy did not handle nil pointer correctly")
		}
	})
}
