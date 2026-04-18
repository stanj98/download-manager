package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

	//50MB
	//0: 0, 5MB
	//1: 6, 11MB

	var sections = make([][2]int, d.TotalSections)
	eachSize := size / d.TotalSections
	fmt.Printf("Each size is %v bytes\n", eachSize)

	//example: if file size is 100 bytes, our section would look like:
	//[[0, 10], [11, 21], [22, 33], [34, 44], [45, 55], [56, 66], [67, 77], [78, 88], [89, 99], [99, 99]]

	for i := range sections {
		if i == 0 {
			//starting byte of first section
			sections[i][0] = 0
		} else {
			//starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}
		if i < d.TotalSections-1 {
			//ending byte of first section
			sections[i][1] = sections[i][0] + eachSize
		} else {
			//ending byte of last section
			sections[i][1] = size - 1
		}
	}
	fmt.Println(sections)

	var wg sync.WaitGroup
	for i, s := range sections {
		wg.Add(1)
		//store current values as they will keep changing
		//i := i
		//s := s
		go func() {
			defer wg.Done()
			err := d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	err = d.mergeFiles(sections)
	if err != nil {
		return fmt.Errorf("Cannot merge files!")
	}
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

func (d Download) downloadSection(i int, s [2]int) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return fmt.Errorf("Can't get file")
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", s[0], s[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}
	fmt.Printf("Downloaded %v bytes for section %v: %v\n", resp.Header.Get("Content-Length"), i, s)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Can't retrieve body, response is %v", err)
	}
	err = os.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Can't save file, response is %v", err)
	}

	return nil
}

// merge downloaded files
func (d Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Can't open file, response is %v", err)
	}
	defer f.Close()
	for i := range sections {
		b, err := os.ReadFile(fmt.Sprintf("section-%v.tmp", i))
		if err != nil {
			return fmt.Errorf("Can't read file, response is %v", err)
		}
		n, err := f.Write(b)
		if err != nil {
			return fmt.Errorf("Can't write file, response is %v", err)
		}
		fmt.Printf("%v bytes merged\n", n)
	}
	return nil

}
