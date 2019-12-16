package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const readBufSize = 1024
const fileName = "./1.txt"

func main() {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("%v", fmt.Errorf("open: %w", err))
	}
	defer f.Close()

	_, err = f.Seek(0, 2)
	if err != nil {
		log.Fatalf("%v", fmt.Errorf("seek: %w", err))
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
		handle(append(acc, readBuf...))
		acc = acc[:0]
	}
}

func handle(acc []byte) {
	fmt.Print(string(acc))
}
