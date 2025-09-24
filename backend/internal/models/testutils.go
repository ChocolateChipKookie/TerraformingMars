package models

// Generic utility function to convert values to pointers for tests
func Ptr[T any](v T) *T { return &v }