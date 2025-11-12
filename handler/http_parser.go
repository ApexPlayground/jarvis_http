package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    map[string]string
}

func ParseHttp(conn net.Conn) (*HttpRequest, error) {
	// using bufio to reduce constant reading
	reader := bufio.NewReader(conn)

	// read reques linw
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request: %v", err)
	}

	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")

	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	req := &HttpRequest{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
		Headers: make(map[string]string),
		Body:    make(map[string]string),
	}

	// read headers

	// at this point we can read the header
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header line %v", err)
		}

		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break
		}

		headerPart := strings.SplitN(headerLine, ":", 2)
		if len(headerPart) == 2 {
			key := strings.TrimSpace(headerPart[0])
			value := strings.TrimSpace(headerPart[1])
			req.Headers[key] = value
		}

	}

	// parsing body

	// check if key exist
	if val, ok := req.Headers["Content-Length"]; ok {
		length, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length: %v", err)
		}

		bodyBytes := make([]byte, length)
		_, err = io.ReadFull(reader, bodyBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}

		bodyStr := string(bodyBytes)

		fmt.Println("Raw body:", bodyStr)

		// Check content type before parsing
		if ctype, ok := req.Headers["Content-Type"]; ok {
			switch {
			case strings.Contains(ctype, "application/x-www-form-urlencoded"):
				pairs := strings.Split(bodyStr, "&")
				for _, p := range pairs {
					kv := strings.SplitN(p, "=", 2)
					if len(kv) == 2 {
						key := kv[0]
						value := kv[1]
						req.Body[key] = value
					}
				}

			case strings.Contains(ctype, "application/json"):
				var jsonData map[string]string
				err := json.Unmarshal(bodyBytes, &jsonData)
				if err != nil {
					return nil, fmt.Errorf("error parsing JSON body: %v", err)
				}
				maps.Copy(req.Body, jsonData)
			}
		}
	}

	return req, nil

}
