package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := "https://fotomhs.amikom.ac.id/2022/22_11_4879.jpg"
	err := downloadImg(url, "./results/")
	if err != nil {
		fmt.Println(err)
	}
}

func downloadImg(url, pathfile string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("error response code: " + string(res.StatusCode))
	}

	filename := url[strings.LastIndex(url, "/")+1:]
	file, err := os.Create(pathfile + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}
