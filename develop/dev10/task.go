package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// freechess.org 5000

import (
	"bufio"
	"github.com/urfave/cli"
	"log"
	"net"
	"os"
	"time"
)

func telnet(host string, port string, timeout int) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	telnet_size := 80 * 24
	bytes := make([]byte, telnet_size)
	shutdown := make(chan bool, 1)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
	    conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			_, err := conn.Read(bytes)
			if err != nil {
				shutdown <- true
				panic(err)
			}
			log.Println(string(bytes))
			scanner.Scan()
			err = scanner.Err()
			if err != nil {
				shutdown <- true
				panic(err)
			}
			scanner_bytes := scanner.Bytes()
			if len(scanner_bytes) == 0 {
				shutdown <- true
				return
			}
      conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			conn.Write(scanner.Bytes())
			conn.Write([]byte("\n"))
		}
	}()
	<-shutdown
}

func main() {
	app := cli.NewApp()
	timeout := cli.DurationFlag{
		Name:  "timeout",
		Value: 10 * time.Second,
		Usage: "timeout for connection",
	}
	app.Flags = []cli.Flag{
		&timeout,
	}
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			panic("Usage: go-telnet host port")
		}
		if c.Duration("timeout") < 0 {
			panic("Timeout must be positive")
		}
		if c.Args()[0] == "" {
			panic("Host must be not empty")
		}
		if c.Args()[1] == "" {
			panic("Port must be not empty")
		}
		telnet(c.Args()[0], c.Args()[1], int(c.Duration("timeout").Seconds()))
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
