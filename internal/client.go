package internal

import (
	"fmt"
	"net/http"
	"time"
)

func NewMovieClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func GetMovie(client *http.Client, id int) error {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("movie stand returned status: %d", resp.StatusCode)
	}

	fmt.Println("Status:", resp.StatusCode)

	return nil
}
