package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func readFile(filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when function exits

	// Create a buffer to hold 8 bytes
	buffer := make([]byte, 8)
	var line string

	for {
		// Read up to 8 bytes
		n, err := file.Read(buffer)
		if err == io.EOF {
			if len(line) > 0 {
				fmt.Printf("read: %s\n", line) // Print last remaining line
			}
			break // Exit loop when end of file is reached
		}
		if err != nil {
			fmt.Println("Error while reading file:", err)
			return
		}

		// Convert the read bytes into a string
		part := string(buffer[:n])

		// Split by newline character
		parts := strings.Split(part, "\n")

		// Process split parts
		for i, p := range parts {
			if i == 0 {
				line += p // Append to the current line
			} else {
				fmt.Printf("read: %s\n", line) // Print the completed line
				line = p                       // Start a new line
			}
		}
	}

}

func main() {
	const filePath = "./messages.txt"
	readFile(filePath)
}
