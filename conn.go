package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type RedisClient struct {
	host   string
	port   int
	pass   string
	conn   net.Conn
	reader *bufio.Reader
}

// Connect 建立连接
func (r *RedisClient) Connect() error {
	var err error
	r.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", r.host, r.port))
	if err != nil {
		return err
	}
	r.reader = bufio.NewReader(r.conn)
	if r.pass != "" {
		return r.Auth(r.pass)
	}
	return nil
}

// Auth 发送认证
func (r *RedisClient) Auth(pass string) error {
	authCommand := fmt.Sprintf("AUTH %s\r\n", pass)
	_, err := r.SendCommand(authCommand)
	return err
}

// number 转数字
func (r *RedisClient) number(response string) int {
	number, err := strconv.Atoi(response[1 : len(response)-2])
	if err != nil {
		panic(err)
	}
	return number
}

// readData 读取响应
func (r *RedisClient) readData() (string, error) {
	response, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if response == "\r\n" {
		return r.readData()
	}
	return response, nil
}

// SendCommand 发送命令
func (r *RedisClient) SendCommand(command string) (string, error) {
	_, _ = fmt.Fprintf(r.conn, command)
	return r.handleCommandResp(isBlockCmd(command))
}

// Close 关闭连接
func (r *RedisClient) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

// appendWithLen 追加根据长度的响应
func (r *RedisClient) appendWithLen(response string, byteLen int) (str string, err error) {
	if byteLen < 0 {
		return "(nil)", nil
	}
	var b = make([]byte, byteLen)
	n, err := io.ReadFull(r.reader, b)
	if err != nil {
		return "", err
	}
	str = response + string(b[:n])
	return
}

// appendWithMultiline 追加多行响应
func (r *RedisClient) appendWithMultiline(response string, loop int, sn bool) (str string, err error) {
	buf := bytes.NewBufferString(response)
	var byteLen int
	for i := 0; i < loop; i++ {
		str, err = r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		buf.WriteString(str)
		if sn {
			buf.WriteString(fmt.Sprintf("%d) ", i+1))
		}
		if strings.HasPrefix(str, "$") {
			byteLen, err = strconv.Atoi(str[1 : len(str)-2])
			str, err = r.appendWithLen("", byteLen+2)
			if err != nil {
				return "", err
			}
			buf.WriteString(`"` + str[:len(str)-2] + `"` + "\r\n")
		} else if strings.HasPrefix(str, ":") {
			buf.WriteString(str[1:len(str)-2] + "\r\n")
		}
	}
	return buf.String(), nil
}

// handleCommandResp 处理命令响应
func (r *RedisClient) handleCommandResp(block bool) (string, error) {
	var parse = func(response string, sn bool) (string, error) {
		if strings.HasPrefix(response, "$") {
			return r.appendWithLen(response, r.number(response))
		} else if strings.HasPrefix(response, "*") {
			return r.appendWithMultiline(response, r.number(response), sn)
		} else {
			return response, nil
		}
	}
	// 读取响应
	response, _ := r.readData()
	if block {
		for {
			resp, err := parse(response, true)
			if err != nil {
				return "", err
			}
			fmt.Println(FmtResp(resp))
			response, _ = r.readData()
		}
	}
	return parse(response, true)
}

// isBlockCmd 判断是否是阻塞命令
func isBlockCmd(command string) bool {
	command = strings.ToUpper(command)
	return strings.HasPrefix(command, "SUBSCRIBE") || strings.HasPrefix(command, "PSUBSCRIBE")
}
