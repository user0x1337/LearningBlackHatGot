package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p) //Change for real word tests
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		err2 := conn.Close()
		results <- p
		if err2 != nil {
			continue
		}
	}
}

const MaxOpenPorts = 1024

func main() {
	fmt.Print("Starting Portscanning (TCP-Scan)....")
	ports := make(chan int, 256)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= MaxOpenPorts; i++ {
			ports <- i
		}
	}()

	for i := 0; i < MaxOpenPorts; i++ {
		port := <-results
		if port > 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d / TCP\topen\n", port)
	}
}
