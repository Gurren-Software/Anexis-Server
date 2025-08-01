package main

import (
	"fmt"

	"github.com/Treefle-labs/anexis-server/builder/build"
)

func main() {
	err := build.BuildAllTSFiles("./src")
	if err != nil {
		fmt.Printf("failed to build files: %v", err)
	}
}
