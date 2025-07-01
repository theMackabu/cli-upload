package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

type Response struct {
	ID      string
	Private bool
	Name    string
	Size    float64
	Url     string
}

func SendPostRequest(url string, filename string) (string, []byte) {
	client := &http.Client{}
	data, err := os.Open(filename)
	if err != nil {
		fmt.Printf("\033[0;31mfile \033[0;33m'%s'\033[0m\033[1;31m not found\033[0m\n", os.Args[1])
		os.Exit(0)
	}

	defer data.Close()
	info, _ := data.Stat()

	bar := pb.New(int(info.Size()))
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	bar.Start()

	req, err := http.NewRequest("POST", url, bar.NewProxyReader(data))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bar.Finish()

	return resp.Status, content
}

func main() {
	var response Response

	if len(os.Args) < 2 {
		fmt.Println("\033[1;31mmissing parameter \033[0;33m'file'\033[0m\033[1;31m to upload\033[0m")
		return
	}

	if len(os.Args) == 3 {
		if os.Args[2] == "-p" {
			pathArr := strings.Split(os.Args[1], "/")
			_, content := SendPostRequest(fmt.Sprintf("https://themackabu.dev/cdn/%s?q=private", pathArr[len(pathArr)-1]), os.Args[1])
			json.Unmarshal([]byte(string(content)), &response)
		}
	} else {
		pathArr := strings.Split(os.Args[1], "/")
		_, content := SendPostRequest(fmt.Sprintf("https://themackabu.dev/cdn/%s", pathArr[len(pathArr)-1]), os.Args[1])
		json.Unmarshal([]byte(string(content)), &response)
	}

	fmt.Printf("\033[0;36mInformation\033[0m\n - Uploaded: \033[0;32m%s\033[0m\n - Size: \033[0;32m%.2fkb\033[0m\n - ID: \033[0;33m%s\033[0m\n\n\033[0;36mImportant\033[0m\n - Private: \033[0;31m%v\033[0m\n - Access URL: \033[0;35m%s\033[0m\n", response.Name, response.Size/1000, response.ID, response.Private, response.Url)
}
