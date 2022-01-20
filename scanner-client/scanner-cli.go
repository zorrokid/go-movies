package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zorrokid/go-movies/scanner"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Need file path as argument")
		return
	}
	if bbs, err := scanner.Scan(args[1], "fin"); err != nil {
		log.Fatal(err)
	} else {
		for _, bb := range bbs {
			fmt.Printf("%s, %d, %d\n", bb.Word, bb.Box.Dx(), bb.Box.Dy())
		}
	}
}
