package main

import (
	_ "bytes"
	_ "github.com/pelletier/go-toml/v2"
	_ "github.com/yuin/goldmark"
)

func main() {
	rootContentPath := PanicErr(GetRootContentPath())
	if err := GenerateSite(rootContentPath, "../build"); err != nil {
		panic(err)
	}
}
