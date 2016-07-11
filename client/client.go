package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/digitalmarc/go/udp-chat/common"
	"github.com/nu7hatch/gouuid"
)

type Client struct {
	connection          *net.UDPConn
	alive               bool
	userID              uuid.UUID
	userName            string
	sendingMessageQueue chan string
	receiveMessages     chan string
}

var scanError error

func (c *Client) packMessage(msg string, messageType common.MessageType) string {
	return strings.Join([]string{c.userID.String(), strconv.Itoa(int(messageType)), c.userName, msg, time.Now().Format("15:04:05")}, "\x01")
}

func (c *Client) funcSendMessage(msg string) {
	message := c.packMessage(msg, common.FUNC)
	_, err := c.connection.Write([]byte(message))
	checkError(err, "func_sendMessage")
}

func (c *Client) sendMessage() {
	for c.alive {

		msg := <-c.sendingMessageQueue
		message := c.packMessage(msg, common.CLASSIQUE)
		_, err := c.connection.Write([]byte(message))
		checkError(err, "sendMessage")
	}

}

func (c *Client) receiveMessage() {
	var buf [512]byte
	//var userID *uuid.UUID
	for c.alive {
		n, err := c.connection.Read(buf[0:])
		checkError(err, "receiveMessage")
		//msg := string(buf[0:n])
		//stringArray := strings.Split(msg, "\x01")

		//userID, err = uuid.ParseHex(stringArray[0])
		//checkError(err, "receiveMessage")
		//if *userID != c.userID {
		c.receiveMessages <- string(buf[0:n])
		fmt.Println("")
		//}
	}
}

func (c *Client) readInput() {
	var msg string
	for c.alive {
		fmt.Println("msg: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			msg = scanner.Text()
			if msg == ":quit" || msg == ":q" {
				c.alive = false
			}
			c.sendingMessageQueue <- msg
		}
		//_,scanError := fmt.Scanln(&msg)
		//checkError(scanError, "readInput")

	}
}

func (c *Client) printMessage() {
	for c.alive {
		msg := <-c.receiveMessages
		stringArray := strings.Split(msg, "\x01")
		var userName = stringArray[2]
		var content = stringArray[3]
		var time = stringArray[4]
		fmt.Printf("%s %s: %s", time, userName, content)
		fmt.Println("")
		// pf("MESSAGE RECEIVED: %s \n", msg)
		// pf("USER NAME: %s \n", stringArray [2])
		// pf("CONTENT: %s \n", stringArray [3])
		if strings.HasPrefix(msg, ":q") || strings.HasPrefix(msg, ":quit") {
			fmt.Printf("%s is leaving", userName)
		}
	}
}

func nowTime() string {
	return time.Now().String()
}
func checkError(err error, funcName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s-----in func:%s", err.Error(), funcName)
		os.Exit(1)
	}
}
func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage:%s host:port", os.Args[0])
	// 	os.Exit(1)
	// }
	// service := os.Args[1]
	// udpAddr, err := net.ResolveUDPAddr("udp4", service)
	udpAddr, err := net.ResolveUDPAddr("udp4", "78.242.118.156:1200")
	checkError(err, "main")

	var c Client
	c.alive = true
	c.sendingMessageQueue = make(chan string)
	c.receiveMessages = make(chan string)
	u, err := uuid.NewV4()

	c.userID = *u

	fmt.Println("input name: ")
	_, err = fmt.Scanln(&c.userName)
	checkError(err, "main")

	c.connection, err = net.DialUDP("udp", nil, udpAddr)
	checkError(err, "main")
	defer c.connection.Close()

	c.funcSendMessage("joined")

	go c.printMessage()
	go c.receiveMessage()

	go c.sendMessage()
	c.readInput()

	c.funcSendMessage("left")

	os.Exit(0)
}
