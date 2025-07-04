package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type XojoCommand struct {
	Tag    string `json:"tag"`
	Script string `json:"script"`
}

func main() {
	commandOpt := flag.String("f", "", "The script to run")
	socketPathOpt := flag.String("s", "/tmp/XojoIDE", "The path to the Xojo IDE socket")
	timeoutOpt := flag.Int("t", 30, "The number of seconds to wait for the Xojo IDE to connect")

	flag.Parse()

	// dereference the options
	command := *commandOpt
	timeout := *timeoutOpt
	socketPath := *socketPathOpt

	if command == "" {
		fmt.Println("Usage: xojo-cli -f <command-file> [-s socketPath] [-s timeout]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Command file: ", command)
	fmt.Println("Socket path: ", socketPath)
	fmt.Println("Timeout: ", timeout)

	// open the command file
	file, err := os.Open(command)
	if err != nil {
		fmt.Println("failed to open command file")
		os.Exit(1)
	}
	defer file.Close()

	// read the command file
	scanner := bufio.NewScanner(file)
	script := ""
	for scanner.Scan() {
		script += scanner.Text() + "\n"
	}

	fmt.Println(script)

	// check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the command file")
		os.Exit(1)
	}

	// connect to the Xojo IDE
	conn := net.Conn(nil)
	for {
		conn, err = net.Dial("unix", socketPath)
		if err != nil {
			fmt.Println("Waiting to connect to Xojo IDE. Time left: ", timeout)
			time.Sleep(1 * time.Second)
			timeout--

			if timeout == 0 {
				fmt.Println("Connection timeout. Failed to connect to Xojo IDE.")
				os.Exit(1)
			} else {
				continue
			}
		}
		defer conn.Close()
		fmt.Println("Connected to Xojo IDE")
		if conn != nil {
			fmt.Println("conn != nil")
		}
		break
	}

	if conn == nil {
		fmt.Println("conn == nil. Failed to connect to Xojo IDE")
		os.Exit(1)
	}

	// send the command to switch to protocol 2
	protocol_change := `{"protocol":2}`
	sendCommand(conn, protocol_change)

	// send the user supplied command
	sendCommand(conn, script)

	// read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed to read response from Xojo")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(buf[:n]))
	os.Exit(0)
}

func sendCommand(conn net.Conn, command string) {
	message := ""
	if strings.HasPrefix(command, string("{")) {
		// JSON
		message = command + "\x00"
	} else {
		jsonObj := XojoCommand{
			Tag:    "",
			Script: command,
		}

		jsonData, err := json.Marshal(jsonObj)
		if err != nil {
			fmt.Println("failed to make JSON for command")
			fmt.Println(err)
			os.Exit(1)

		}
		message = string(jsonData) + "\x00"
	}

	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("failed to send command to Xojo")
		fmt.Println(err)
		os.Exit(1)
	}
}
