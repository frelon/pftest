package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println("TEST")

	for _, i := range interfaces {
		fmt.Println(i)
	}

	fmt.Println("Hello world")
}
