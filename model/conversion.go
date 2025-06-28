// Copyright (c) 2025 A Bit of Help, Inc.

// Package model provides utilities for working with domain models and DTOs.
// It includes functions for copying fields between structs and creating deep copies of objects.
package model

import (
	"fmt"
	"reflect"
)

// CopyFields copies fields from source to destination based on field names.
// Both source and destination must be pointers to structs.
// Fields are copied only if they have the same name and compatible types.
//
// Parameters:
//   - dst: Pointer to the destination struct
//   - src: Pointer to the source struct
//
// Returns:
//   - error: An error if the operation fails
func CopyFields(dst, src interface{}) error {
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	// Check if both are pointers
	if dstValue.Kind() != reflect.Ptr || srcValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both source and destination must be pointers")
	}

	// Get the values that the pointers point to
	dstElem := dstValue.Elem()
	srcElem := srcValue.Elem()

	// Check if both are structs
	if dstElem.Kind() != reflect.Struct || srcElem.Kind() != reflect.Struct {
		return fmt.Errorf("both source and destination must be pointers to structs")
	}

	// Get the type of the destination struct
	dstType := dstElem.Type()

	// Iterate through the fields of the destination struct
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		dstFieldValue := dstElem.Field(i)

		// Skip unexported fields
		if !dstFieldValue.CanSet() {
			continue
		}

		// Find the corresponding field in the source struct
		srcFieldValue := srcElem.FieldByName(dstField.Name)
		if !srcFieldValue.IsValid() {
			continue
		}

		// Check if the types are compatible
		if srcFieldValue.Type().AssignableTo(dstFieldValue.Type()) {
			dstFieldValue.Set(srcFieldValue)
		}
	}

	return nil
}

// DeepCopy creates a deep copy of the source object.
// Both source and destination must be pointers to structs of the same type.
//
// Parameters:
//   - dst: Pointer to the destination struct
//   - src: Pointer to the source struct
//
// Returns:
//   - error: An error if the operation fails
func DeepCopy(dst, src interface{}) error {
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	// Check if both are pointers
	if dstValue.Kind() != reflect.Ptr || srcValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both source and destination must be pointers")
	}

	// Get the values that the pointers point to
	dstElem := dstValue.Elem()
	srcElem := srcValue.Elem()

	// Check if both are structs
	if dstElem.Kind() != reflect.Struct || srcElem.Kind() != reflect.Struct {
		return fmt.Errorf("both source and destination must be pointers to structs")
	}

	// Check if they are of the same type
	if dstElem.Type() != srcElem.Type() {
		return fmt.Errorf("source and destination must be of the same type")
	}

	// Perform the deep copy
	dstElem.Set(deepCopyValue(srcElem))

	return nil
}

// deepCopyValue creates a deep copy of a reflect.Value.
// It handles various types including pointers, structs, slices, and maps.
//
// Parameters:
//   - src: The source value to copy
//
// Returns:
//   - reflect.Value: A deep copy of the source value
func deepCopyValue(src reflect.Value) reflect.Value {
	// Handle nil pointers
	if !src.IsValid() {
		return reflect.Value{}
	}

	// Create a new value of the same type as src
	dst := reflect.New(src.Type()).Elem()

	switch src.Kind() {
	case reflect.Ptr:
		// Handle nil pointers
		if src.IsNil() {
			return dst
		}
		// Create a new pointer
		v := deepCopyValue(src.Elem())
		dst.Set(reflect.New(v.Type()))
		dst.Elem().Set(v)
	case reflect.Struct:
		// Copy each field
		for i := 0; i < src.NumField(); i++ {
			if dst.Field(i).CanSet() {
				dst.Field(i).Set(deepCopyValue(src.Field(i)))
			}
		}
	case reflect.Slice:
		// Handle nil slices
		if src.IsNil() {
			return dst
		}
		// Create a new slice
		dst = reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		// Copy each element
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(deepCopyValue(src.Index(i)))
		}
	case reflect.Map:
		// Handle nil maps
		if src.IsNil() {
			return dst
		}
		// Create a new map
		dst = reflect.MakeMap(src.Type())
		// Copy each key-value pair
		for _, key := range src.MapKeys() {
			dst.SetMapIndex(deepCopyValue(key), deepCopyValue(src.MapIndex(key)))
		}
	default:
		// For other types (int, string, etc.), just set the value
		dst.Set(src)
	}

	return dst
}
