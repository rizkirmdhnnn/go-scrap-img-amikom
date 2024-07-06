package main

import (
	"fmt"
	"go-scrap/config"
	"go-scrap/modules"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func newHttpProxy(addr string) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
		Timeout: 15 * time.Second,
	}
	return httpClient, nil
}

func downloadImage(url string, tor *modules.Tor) (*http.Response, error) {
	for {
		httpClient, err := newHttpProxy(config.Cfg.TORSERVER_ADDRESS)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		res, err := httpClient.Get(url)
		elapsed := time.Since(start)

		if err != nil {
			if os.IsTimeout(err) || elapsed >= 15*time.Second {
				fmt.Printf("Request timed out or took too long. Changing IP and retrying...\n")
				tor.ChangeIP()
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			if res.StatusCode == http.StatusNotFound {
				fmt.Println("Image not found")
			}
		}
		return res, nil
	}
}

func saveImage(res *http.Response, filename string) error {
	if res.StatusCode == http.StatusNotFound {
		return nil
	}
	file, err := os.Create(filename)
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

func getData(jurusan int, tahun int) {
	tor := modules.NewTor(config.Cfg.TORCONTROL_ADDRESS)
	tor.Init()
	tor.ChangeIP()
	totalRequests := 5000
	for i := 4880; i < totalRequests; i++ {
		url := fmt.Sprintf("https://fotomhs.amikom.ac.id/%d/%d_%d_%d.jpg", tahun, tahun%100, jurusan, i)
		println(url)

		// Download image (or perform any HTTP request)
		res, err := downloadImage(url, tor)
		if err != nil {
			log.Fatal(err)
		}

		// Save the image
		filename := fmt.Sprintf("./results/image_%d.jpg", i)
		err = saveImage(res, filename)
		if err != nil {
			log.Fatal(err)
		}

		// Check if it's time to change IP (every 10 requests, for example)
		if (i+1)%10 == 0 {
			fmt.Printf("Changing IP after %d requests\n", i+1)
			tor.ChangeIP()
			time.Sleep(5 * time.Second)
		}
	}
}

func menu() {
	var jurusan int
	var tahun int
	fmt.Println("===================")
	fmt.Println("Informatika (11)")
	fmt.Println("Sistem Informasi (kambingsun)")
	fmt.Println("Ilmu Komunikasi (kambingsun)")
	fmt.Println("===================")
	fmt.Print("Pilih kode : ")
	fmt.Scanln(&jurusan)
	fmt.Println("==== Tahun Angkatan ====")
	fmt.Println("2020")
	fmt.Println("2021")
	fmt.Println("2022")
	fmt.Println("2023")
	fmt.Println("=========================")
	fmt.Print("Pilih Tahun Angkatan : ")
	fmt.Scanln(&tahun)
	getData(jurusan, tahun)
}

func main() {
	config.LoadConfig()
	menu()
}
