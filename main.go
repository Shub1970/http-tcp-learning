package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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

	return out
}

func readFile(filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while opening file:", err)
		return
	}

	// Get lines channel and print each line
	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func main() {
	const filePath = "./messages.txt"
	readFile(filePath)
}
