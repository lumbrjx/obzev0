package helper

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func ReqSimulator(
	proxyAddr string,
	targetURL string,
	duration time.Duration,
) error {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: proxyAddr}),
		},
		Timeout: duration,
	}

	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		resp, err := client.Get(targetURL)
		if err != nil {
			log.Printf("Failed to make request: %v", err)
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			resp.Body.Close()
			return err
		}
		resp.Body.Close()

		log.Printf("Client:")
		log.Printf("Response status: %s", resp.Status)
		log.Printf("Response body: %s", string(body))
		log.Printf("\n----------------------------------------")

		time.Sleep(1 * time.Second)
	}

	return nil
}
