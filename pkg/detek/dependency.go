package detek

import "reflect"

type DependencyInfo struct {
	Type reflect.Type
	// TODO
	// Description string
	// IsOptional  bool
}

func TypeOf(v interface{}) reflect.Type {
	return reflect.TypeOf(v)
}

// map[name] default initialized interface for type hint.
type DependencyMeta map[string]DependencyInfo
