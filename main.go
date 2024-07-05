package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const torProxy = "socks5://127.0.0.1:9050"

var requestCount int
var requestCountMutex sync.Mutex

func httpClientWithProxy(url string) (*http.Response, error) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat dialer: %v", err)
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}
	return httpClient.Get(url)
}

func downloadImg(url, pathfile string) error {
	res, err := httpClientWithProxy(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("error response code: %d", res.StatusCode)
	}
	filename := url[strings.LastIndex(url, "/")+1:]
	file, err := os.Create(pathfile + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	return err
}

func changeIP() error {
	conn, err := net.Dial("tcp", "127.0.0.1:9051")
	if err != nil {
		return fmt.Errorf("gagal terhubung ke Tor control port: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "AUTHENTICATE \"rizkirmdhn\"\r\n")
	fmt.Fprintf(conn, "SIGNAL NEWNYM\r\n")

	time.Sleep(3 * time.Second) // Beri waktu untuk Tor membangun sirkuit baru
	return nil
}

func getData() {
	changeIP()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50) // Batasi konkurensi menjadi 10

	for i := 4000; i <= 6000; i++ {
		url := fmt.Sprintf("https://fotomhs.amikom.ac.id/2022/22_11_%d.jpg", i)
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(url string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			err := downloadImg(url, "./results/")
			if err != nil {
				fmt.Printf("Error downloading %s: %s\n", url, err)
				return
			}
			fmt.Println("Downloaded:", url)

			requestCountMutex.Lock()
			requestCount++
			if requestCount%150 == 0 {
				fmt.Println("Changing IP...")
				err := changeIP()
				if err != nil {
					fmt.Printf("Error changing IP: %v\n", err)
				}
			}
			requestCountMutex.Unlock()

			time.Sleep(1 * time.Second)
		}(url)
	}
	wg.Wait()
	fmt.Println("All downloads completed")
}

func main() {
	cmd := exec.Command("tor")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Gagal memulai Tor:", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(10 * time.Second) // Beri waktu untuk Tor memulai dan membangun sirkuit awal

	getData()
}
