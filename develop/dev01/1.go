package main

/*
Создать программу печатающую точное время с использованием NTP -библиотеки. Инициализировать как go module. Использовать библиотеку github.com/beevik/ntp. Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS

*/

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"os"
)

func f1() {
	l := log.New(os.Stderr, "", 0)
	time, err := ntp.Time("ntp0.ntp-servers.net")
	if err != nil {
		l.Fatal(err)
	} else {
		fmt.Println(time)
	}
}
