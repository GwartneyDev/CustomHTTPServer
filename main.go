package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read from the connection using a fixed size buffer
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading: %s\n", err)
		return
	}

	// Convert buffer to string and parse it manually
	request := string(buffer[:n])
	lines := strings.Split(request, "\r\n")
	if len(lines) < 1 {
		fmt.Println("Received malformed request")
		return
	}

	// Basic parsing of the request line
	parts := strings.Split(lines[0], " ")
	if len(parts) < 2 {
		fmt.Println("Malformed request line")
		return
	}
	method, urlPath := parts[0], parts[1]

	// Handle only GET requests for "/"
	if urlPath == "/" && method == "GET" {
		htmlFile, err := os.Open("index.html")
		if err != nil {
			fmt.Printf("Error opening HTML file: %s\n", err)
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
			return
		}
		defer htmlFile.Close()

		// Send the HTTP response header first
		header := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n"
		conn.Write([]byte(header))

		// Use html tokenizer to read the HTML file and write its tokens
		z := html.NewTokenizer(htmlFile)
		for {
			tt := z.Next()
			switch tt {
			case html.ErrorToken:
				// End of file or an error.
				if err != nil && err != io.EOF {
					fmt.Printf("Error tokenizing HTML: %s\n", err)
				}
				return

			default:
				// Get the token and write it to the connection
				token := z.Token()
				conn.Write([]byte(token.String()))
			}
		}
	} else {
		// If the request path is not recognized or the method is not GET
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	net.Get("http://localhost:8080/")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}

		go handleConnection(conn)
	}
}
