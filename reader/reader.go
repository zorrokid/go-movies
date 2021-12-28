package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/zorrokid/go-movies/reader/model"
)

func ReadMovies() []model.Movie {

	csvFile, err := os.Open("data/myynti.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(csvFile)

	movies := []model.Movie{}

	rowCount := 0

	for {
		rowCount++

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if rowCount == 1 {
			// first row is header row
			continue
		}

		if len(strings.TrimSpace(record[1])) == 0 {
			fmt.Printf("No price for %s\n", record[1])
			continue
		}

		salePrice, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			fmt.Printf("Sales price failed for %s", record[0])
			log.Fatal(err)
		}
		year, err := strconv.ParseUint(record[5], 10, 16)
		if err != nil {
			year = 0
		}
		discs, err := strconv.ParseUint(record[9], 10, 8)
		if err != nil {
			fmt.Printf("Discs failed for %s", record[0])
			log.Fatal(err)
		}

		movie := model.Movie{
			Ean:             record[0],
			SalePrice:       float32(salePrice),
			Conditions:      textToSlice(record[2], ","),
			LocalName:       record[3],
			OriginalName:    record[4],
			Year:            uint16(year),
			Directors:       textToSlice(record[6], ","),
			Actors:          textToSlice(record[7], ","),
			Format:          record[8],
			Discs:           uint8(discs),
			PublicationArea: record[10],
			Publication:     record[11],
			Subtitles:       textToSlice(record[12], ";"),
			Languages:       textToSlice(record[13], ";"),
			Other:           record[14],
			IsRental:        getBoolValue(record[15]),
			HasSlipCover:    getBoolValue(record[16]),
			IsTwoSidedDisc:  getBoolValue(record[17]),
			IsReadTested:    getBoolValue(record[18]),
			CaseType:        record[19],
			DeliveryClass:   record[20],
		}

		movies = append(movies, movie)
		rowCount++
	}
	return movies
}

func getBoolValue(record string) bool {
	if len(strings.TrimSpace(record)) == 0 {
		return false
	}
	retValue, err := strconv.ParseBool(record)
	if err != nil {
		log.Fatal(err)
	}
	return retValue
}

func textToSlice(text string, separator string) []string {
	langStr := strings.TrimSpace(text)
	var res []string
	if len(langStr) > 0 {
		res = strings.Split(langStr, separator)
	}
	return res
}
