package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"github.com/integrii/flaggy"
	"github.com/common-nighthawk/go-figure"
	"strings"
)

var version = "1.0"
var target string
var maxport = 1024

func init() {
	flaggy.SetName("Gmap")
	flaggy.SetDescription("A (disturbingly fast) port scanner!")
	flaggy.SetVersion(version)
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	flaggy.String(&target, "t", "target", "The target host to scan")

	flaggy.Parse()
}

func thread(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%v:%d", target, p)
		conn, err := net.Dial("tcp", address)
		//fmt.Printf("Dialing: %d\n", p)
		if p == maxport {
			os.Exit(0)
		}
		if err != nil {
			results <- 0
			continue
		}
		fmt.Printf("Open: %d\n", p)
		conn.Close()
		results <- p
	}
}

func main() {
	//BEGIN BANNER
	underline := strings.Repeat("=", 40)
	banner := figure.NewFigure("Gmap", "slant", true)
	banner.Print()
	fmt.Println(underline)
	//END BANNER

	if target == "" {
		fmt.Println("Please specify a target. Use Gmap -h for help.")
		os.Exit(0)
	}

	ports := make(chan int, 500)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go thread(ports, results)
	}

	go func() {
		for i := 1; i <= maxport; i++ {
			ports <- i
		}
	}()

	for i := 0; i < maxport; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}


	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}