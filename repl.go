package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const prompt = ">> "

// Start starts shell
func Start(in io.Reader, out io.Writer, redisClient *RedisClient) {
	scanner := bufio.NewScanner(in)
	for {
		_, _ = fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line == "exit" {
			_, _ = fmt.Fprintln(out, "bye.")
			return
		}

		resp, err := redisClient.SendCommand(line + "\r\n")
		if err != nil {
			_, _ = fmt.Fprintln(out, err)
			continue
		}
		_, _ = fmt.Fprintln(out, FmtResp(resp))
	}
}

// FmtResp formats resp
func FmtResp(resp string) string {
	firstChar := resp[0]
	switch {
	case firstChar == ':':
		return resp[1 : len(resp)-2]
	case firstChar == '+' || firstChar == '-':
		return `"` + resp[1:len(resp)-2] + `"`
	case firstChar == '$':
		arr := strings.Split(resp[1:], "\r\n")
		return `"` + arr[1] + `"`
	case firstChar == '*':
		split := strings.Split(resp, "\r\n")
		var result = make([]string, 0, len(split)/2)
		for i := 2; i < len(split); i += 2 {
			result = append(result, split[i])
		}
		return strings.Join(result, "\r\n")
	default:
		return resp
	}
}
