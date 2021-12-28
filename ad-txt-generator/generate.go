package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/zorrokid/go-movies/reader"
	"github.com/zorrokid/go-movies/reader/model"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No path provided")
		return
	}
	path := args[1]
	fmt.Printf("Convert csv to advertisement text from path: %s\n", path)
	movies := reader.ReadMovies(path)
	printSalesText(movies)
	fmt.Println("Done.")
}

const MaxTitleLength = 60

func printSalesText(movies []model.Movie) {

	var funcMap = template.FuncMap{
		"AsCommaSeparatedList": asCommaSeparatedList,
		"inc":                  increase,
	}

	templateBody, err := template.New("body.tmpl").Funcs(funcMap).ParseFiles("templates/body.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	for _, movie := range movies {
		fmt.Print(templateBody.Execute(os.Stdout, movie))
	}
}

func asCommaSeparatedList(values []string) string {
	return strings.Join(values, ",")
}

func increase(value int, step int) int {
	return value + step
}
