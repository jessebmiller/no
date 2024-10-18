package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type OverlayMap[T any] struct {
	Map  map[string]T
	Next *OverlayMap[T]
}

// OverlayMap.Get gets a value out of the overlay map by key
func (o OverlayMap[T]) Get(key string) (T, error) {
	v, ok := o.Map[key]
	if ok {
		return v, nil
	}
	if o.Next == nil {
		return v, errors.New(
			fmt.Sprintf("ERROR: key %s not found", key),
		)
	}
	return o.Next.Get(key)
}

// OverlayMap.Try trys to get each key returning the first value found
func (o OverlayMap[T]) Try(keys ...string) (T, error) {
	for _, key := range keys {
		v, err := o.Get(key)
		if err == nil {
			return v, nil
		}
	}
	var zero T
	return zero, errors.New(
		fmt.Sprintf("ERROR: none of keys %s found", keys),
	)
}

func (o OverlayMap[T]) Push(m map[string]T) OverlayMap[T] {
	return OverlayMap[T]{Map: m, Next: &o}
}

func NewOverlayMap[T any](m map[string]T) OverlayMap[T] {
	return OverlayMap[T]{Map: m}
}

// PPrint prints anything as indented marshaled JSON
func PPrint(content any) {
	bytes, _ := json.MarshalIndent(content, "  ", "  ")
	fmt.Println(string(bytes))
}

// PanicErr panics if it sees an error, returns the value otherwise
// only works on functions returning a single value followed by an error
func PanicErr[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Any returns true if any element passes, false otherwise
func Any[T any](elements []T, passes func(T) bool) bool {
	for _, e := range elements {
		if passes(e) {
			return true
		}
	}
	return false
}

// GetRootContentDir returns the path of the root content directory or an error
// The root content dir is defined as being named `_content` and containing a file
// named `_content.toml`
func GetRootContentPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir != "/" {
		if filepath.Base(dir) != "_content" {
			dir = filepath.Dir(dir)
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return "", err
		}
		if Any(entries, func(e os.DirEntry) bool {
			return e.Name() == "_content.tmpl"
		}) {
			return dir, nil
		}
	}
	return "", errors.New("ERROR: Not in no content dir")
}
