package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "db"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	addr := net.JoinHostPort(host, port)

	timeout := 60 // seconds
	for i := 0; i < timeout; i++ {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			conn.Close()
			fmt.Printf("mysql reachable at %s\n", addr)
			os.Exit(0)
		}
		fmt.Printf("waiting for mysql at %s (%d/%d)\n", addr, i+1, timeout)
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintf(os.Stderr, "timeout waiting for mysql at %s\n", addr)
	os.Exit(1)
}
