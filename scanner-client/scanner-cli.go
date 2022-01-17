package main

import (
	"fmt"
	"os"

	"github.com/zorrokid/go-movies/scanner"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Need file path as argument")
		return
	}
	text := scanner.Scan(args[1])
	fmt.Print(text)
}
