package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type counter struct {
	total int64
	whole bool
}

func (c *counter) Write(b []byte) (int, error) {
	c.total += int64(len(b))
	return len(b), nil
}

func main() {
	fmt.Println(
		"Sample URL:\n" +
		"~5mb https://mapbasic.ru/doc/MapBasicReference-9-0.pdf\n" +
		"~1mb https://mapbasic.ru/soft/MapBasicIDE-1.5-Setup.exe\n" +
			"~ 22mb https://mapbasic.ru/soft/mapbasic_11.zip\n" +
			"~131mb https://mapbasic.ru/soft/mapbasic_12_5.zip",
	)

	var URL string
	fmt.Print("Enter your file URL: ")
	fmt.Scan(&URL)
	splits := strings.Split(URL, "/")
	fileName := splits[len(splits)-1]
	FileDownload(fileName, URL)
}

func FileDownload(fileName, URL string) error {
	res, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error downloading:", err)
	}

	count := counter{total: 0}

	local, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Unable to create file:", fileName, "-", err)
		os.Exit(1)
	}
	defer local.Close()

	fmt.Println("Downloading:", URL, "to:", fileName)

	go func(){
		for{
			time.Sleep(time.Second)
			if count.whole{
				return
			}
			fmt.Println("Downloaded: ", float32(count.total)/(1024*1024), "of",
				float32(res.ContentLength)/(1024*1024))
		}
	}()

	if _, err := io.Copy(local, io.TeeReader(res.Body, &count)); err != nil {
		panic(err)
	}

	count.whole = true
	fmt.Println()
	fmt.Println(fileName, "is downloaded")

	return nil
}

