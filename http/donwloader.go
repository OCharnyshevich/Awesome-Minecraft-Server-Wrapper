package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func DownloadJSON(url string) []byte {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Panic(readErr)
	}

	return body
}

func DownloadFile(filepath string, url string) (err error) {
	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(filepath)

	start := time.Now()

	out, err := os.Create(filepath) // Create the file
	if err != nil {
		return err
	}
	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		panic(err)
	}

	done := make(chan int64)

	go printDownloadPercent(done, path.String(), int64(size))

	// Get the data
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)

	return nil
}
func printDownloadPercent(done chan int64, path string, total int64) {
	var stop = false
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		select {
		case <-done:
			fmt.Println()
			stop = true
		default:
			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()
			if size == 0 {
				size = 1
			}

			var percent = float64(size) / float64(total) * 100
			fmt.Printf("\rProgress: %.0f%s", percent, "%")
		}

		if stop {
			break
		}
		time.Sleep(time.Second)
	}
}
