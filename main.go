package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Movie struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func fetchMovie(ctx context.Context, client *http.Client, id int) (*Movie, error) {
	url := fmt.Sprintf("http://homeworksite.site/%d/info.0.json", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Сервер вернул статус %d", resp.StatusCode)
	}

	var m Movie
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func worker(ctx context.Context, client *http.Client, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-jobs:
			if !ok {
				return
			}

			movie, err := fetchMovie(ctx, client, id)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Printf("Ошибка при загрузке фильма %d: %v\n", id, err)
				continue
			}

			fmt.Printf("%d - %s - %d - %s\n", movie.Id, movie.Title, movie.Year, movie.Director)
		}
	}
}

func main() {
	var from int
	var to int
	var workers int
	var timeout time.Duration

	flag.IntVar(&from, "from", -1, "first film ID")
	flag.IntVar(&to, "to", -1, "last film ID")
	flag.IntVar(&workers, "workers", 10, "workers count")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "request timeout")

	flag.Parse()

	if from == -1 || to == -1 {
		fmt.Println("Error: flags 'from' and 'to' are mandatory")
		os.Exit(1)
	}

	if from <= 0 || to <= 0 {
		fmt.Println("Error: unavaliable flags values")
		os.Exit(1)
	}

	if from > to {
		fmt.Println("Error: 'to' must be bigger than 'from' or equal")
		os.Exit(1)
	}

	if workers <= 0 {
		fmt.Println("Error: invalid workers value")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := &http.Client{
		Timeout: timeout,
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(ctx, client, jobs, &wg)
	}

sendLoop:
	for id := from; id <= to; id++ {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- id:
		}
	}

	close(jobs)
	wg.Wait()
}
