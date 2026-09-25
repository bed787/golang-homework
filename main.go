package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Movie struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func fetchMovie(id int) (*Movie, error) {
	url := fmt.Sprintf("http://homeworksite.site/%d/info.0.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Серве вернул статус %d", resp.StatusCode)
	}

	var m Movie
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func main() {

	var from int
	var to int

	flag.IntVar(&from, "from", -1, "first film ID")
	flag.IntVar(&to, "to", -1, "last film ID")

	flag.Parse()

	if (from == -1) || (to == -1) {
		fmt.Println("Error: flags 'from' and 'to' are mandatory")
		os.Exit(1)
	}

	if (from <= 0) || (to <= 0) {
		fmt.Println("Error: unavaliable flags values")
		os.Exit(1)
	}

	if from > to {
		fmt.Println("Error: 'to' must be bigger than 'from' or equal")
		os.Exit(1)
	}

	for id := from; id <= to; id++ {
		movie, err := fetchMovie(id)
		if err != nil {
			fmt.Println("Сервер вернул ошибку:", err)
			return
		}
		fmt.Println(movie)
	}

}
