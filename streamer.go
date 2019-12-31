package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const readBufSize = 1024

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("file name is required")
	}

	client := newSlackClient(os.Args[2])
	fileName := os.Args[1]

	var f *os.File

	if fileName == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(fileName)
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("open: %w", err))
		}
		defer f.Close()
		_, err = f.Seek(0, 2)
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("seek: %w", err))
		}
	}

	br := bufio.NewReader(f)

	acc := make([]byte, readBufSize)

	for {
		readBuf, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				acc = append(acc, readBuf...)
				continue
			} else {
				log.Fatalf("%v", fmt.Errorf("read: %w", err))
			}
		}
		client.postMessage(string(append(acc, readBuf...)))
		acc = acc[:0]
	}
}

type slackMsg struct {
	Text string `json:"text"`
}

type slackClient struct {
	url    string
	client http.Client
}

func newSlackClient(url string) *slackClient {
	return &slackClient{
		url:    url,
		client: http.Client{Timeout: time.Second * 10},
	}
}

func (sc *slackClient) postMessage(msg string) (err error) {
	body, err := json.Marshal(slackMsg{Text: msg})
	if err != nil {
		return fmt.Errorf("postMessage: marshal: %w", err)
	}
	req, err := http.NewRequest("POST", sc.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("postMessage: newRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = sc.client.Do(req)

	return err
}
