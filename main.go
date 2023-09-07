package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"regexp"
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

	// 白黒モード
	cmode := true
	if len(os.Args) > 1 && os.Args[1] == "mono" {
		cmode = false
	}
	// 白黒モード時の削除用
	prefixRe := regexp.MustCompile(`\[0m\[1m\[3[12]m`)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "[0m[1m[3") {
				if cmode {
					fmt.Println("[0m")
					fmt.Println(line)
				} else {
					fmt.Println(prefixRe.ReplaceAllString(line, ""))
				}
			} else if strings.HasSuffix(line, "<<") {
				fmt.Println(line)
				fmt.Print("[0m")
			} else {
				continue
			}
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
