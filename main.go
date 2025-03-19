package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const port = ":42069"

// getLinesChannel reads a file 8 bytes at a time and sends complete lines to a channel
func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()  // Ensure file is closed when done
		defer close(out) // Close channel when done

		buffer := make([]byte, 8) // 8-byte buffer
		var line string           // Stores incomplete lines

		for {
			n, err := f.Read(buffer) // Read up to 8 bytes
			if err == io.EOF {
				if len(line) > 0 {
					out <- line // Send last remaining line
				}
				break
			}
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			// Convert bytes to string and split by newline
			part := string(buffer[:n])
			parts := strings.Split(part, "\n")

			// Process split parts
			for i, p := range parts {
				if i == 0 {
					line += p // Append to the current line
				} else {
					out <- line // Send completed line
					line = p    // Start new line
				}
			}
		}
	}()

	fmt.Printf("new data is push for reading %s", out)

	return out
}

func main() {
	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("erro listening for tcp traffic: %s\n", err)
	}

	defer listner.Close()

	fmt.Println("Listenin for TCP traffic on", port)

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		fmt.Println("Accepted connection from", port)

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Println("connection to ", conn.RemoteAddr(), "closed")

	}
}
