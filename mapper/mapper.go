/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */

// Package mapper provides a collection of simple mapper functions.
package mapper

import (
	"fmt"
	"math"
)

// IntToInt32 maps an int to an int32.
func IntToInt32(r int) (int32, error) {
	if r < math.MinInt32 || r > math.MaxInt32 {
		return 0, fmt.Errorf("integer overflow, cannot map %d to int32", r)
	}
	return int32(r), nil
}

// IntToInt32Ptr maps an int to an *int32.
func IntToInt32Ptr(r int) (*int32, error) {
	rr, err := IntToInt32(r)
	return &rr, err
}

// Int32PtrToInt32 maps an *int32 to an int32,
// defaulting to 0 if the pointer is nil.
func Int32PtrToInt32(r *int32) int32 {
	if r == nil {
		return 0
	}
	return *r
}

// StringToStringPtr maps a string to a *string.
func StringToStringPtr(r string) *string {
	return &r
}

// BoolToBoolPtr maps a bool to a *bool.
func BoolToBoolPtr(r bool) *bool {
	return &r
}

// StringPtrToString maps a *string to a string,
// defaulting to "" if the pointer is nil.
func StringPtrToString(r *string) string {
	if r == nil {
		return ""
	}
	return *r
}
