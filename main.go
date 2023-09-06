package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	config := tls.Config{Certificates: []tls.Certificate{}, InsecureSkipVerify: false}
	conn, err := tls.Dial("tcp", "koukoku.shadan.open.ad.jp:992", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, ">>") {
				continue
			}
			fmt.Println(line)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Fprintln(conn, scanner.Text())
		}
	}()

	wg.Wait()
}
