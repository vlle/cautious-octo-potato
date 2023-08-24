package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"net/http"
	"os"
)

func wget(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Usage: task <url>")
	}
	wget(os.Args[1])
}
