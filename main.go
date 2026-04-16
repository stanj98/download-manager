package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Download struct {
	Url           string
	TargetPath    string
	TotalSections int
}

func main() {
	startTime := time.Now()
	d := Download{
		Url:           "https://file-examples.com/wp-content/storage/2018/04/file_example_AVI_480_750kB.avi",
		TargetPath:    "final.mp4",
		TotalSections: 10,
	}
	err := d.Do()
	if err != nil {
		log.Fatalf("An error occurred while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}

func (d Download) Do() error {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		log.Fatalf("An error occurred while making HEAD request: %s\n", err)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatalf("An error occurred while sending HEAD request: %s\n", err)
	}
	fmt.Printf("Got %v\n", resp.StatusCode)
	if resp.StatusCode > 299 {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.Fatalf("An error occurred while fetching file: %s\n", err)
	}
	fmt.Printf("Size is %v in bytes\n", size)
	//8:32
	return nil
}

func (d Download) getNewRequest(method string) (*http.Request, error) {
	fmt.Println("Making connection")
	r, err := http.NewRequest(
		method,
		d.Url,
		nil,
	)
	if err != nil {
		log.Fatalf("An error occurred while performing download request: %s\n", err)
	}
	r.Header.Set("User-Agent", "Simple Download Manager v001")
	return r, nil
}
