package helper

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sendReq(client *http.Client, targetAddr string) {
	resp, err := client.Get("http://127.0.0.1:" + targetAddr)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp == nil {
		log.Println("Response is nil")
		return
	}

	log.Printf("Response status: %s", resp.Status)
	time.Sleep(2 * time.Second)
}

func ReqSimulator(targetAddr string, oneTime bool, duration time.Duration) error {
	client := &http.Client{}

	if oneTime {
		sendReq(client, targetAddr)
	} else {
		startTime := time.Now()
		endTime := startTime.Add(duration - 3*time.Second)

		for time.Now().Before(endTime) {
			sendReq(client, targetAddr)
		}
	}

	fmt.Println("Timeout reached, stopping the requests.")
	return nil
}
