package main

/*
Утилита sort
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func sort_util(c *cli.Context) {
	// filename cautious-octo-potato]
	filename := c.String("filename")
	if filename == "" {
		c.App.Writer.Write([]byte("filename is empty\n"))
		os.Exit(1)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []string{}
	unique_lines_map := make(map[string]bool)
	if c.Bool("u") {
		for scanner.Scan() {
			txt := scanner.Text()
			if _, ok := unique_lines_map[txt]; !ok {
				unique_lines_map[txt] = true
				lines = append(lines, txt)
			}
		}
	} else {
		for scanner.Scan() {
			txt := scanner.Text()
			lines = append(lines, txt)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	k := c.Int("k")
	if c.Bool("n") {
		if c.Bool("r") {
			sort.Slice(lines, func(i, j int) bool {
				split_lines := strings.Split(lines[i], " ")
				split_lines2 := strings.Split(lines[j], " ")
				if k > len(split_lines) {
					return lines[i] > lines[j]
				}
				first_num, err := strconv.Atoi(split_lines[k])
				if err != nil {
					return lines[i] > lines[j]
				}
				second_num, err := strconv.Atoi(split_lines2[k])
				if err != nil {
					return lines[i] > lines[j]
				}
				return first_num > second_num
			})
		} else {
			sort.Slice(lines, func(i, j int) bool {
				split_lines := strings.Split(lines[i], " ")
				if k > len(split_lines) {
					return lines[i] < lines[j]
				}
				first_num, err := strconv.Atoi(strings.Split(lines[i], " ")[k])
				if err != nil {
					return lines[i] < lines[j]
				}
				second_num, err := strconv.Atoi(strings.Split(lines[j], " ")[k])
				if err != nil {
					return lines[i] < lines[j]
				}
				return first_num < second_num
			})
		}
	} else if c.Bool("r") {
		sort.Slice(lines, func(i, j int) bool {
			if k != 0 {
				split_lines := strings.Split(lines[i], " ")
				split_lines2 := strings.Split(lines[j], " ")
				if k > len(split_lines) {
					return lines[i] > lines[j]
				}
				return split_lines[k] > split_lines2[k]
			} else {
				return lines[i] > lines[j]
			}
		})
	} else {
		sort.Slice(lines, func(i, j int) bool {
			if k != 0 {
				split_lines := strings.Split(lines[i], " ")
				split_lines2 := strings.Split(lines[j], " ")
				if k > len(split_lines) {
					return lines[i] < lines[j]
				}
				return split_lines[k] < split_lines2[k]
			} else {
				return lines[i] < lines[j]
			}
		})
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	app := cli.NewApp()
	app.Action = sort_util
	app.Name = "sort_utility"
	app.Description = "a sort utility for L2 course at WB"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "filename", Usage: "a file to sort"},
		cli.IntFlag{Name: "k", Usage: "a column to sort"},
		cli.BoolFlag{Name: "n", Usage: "sort by number"},
		cli.BoolFlag{Name: "r", Usage: "reverse sort"},
		cli.BoolFlag{Name: "u", Usage: "sort unique lines"},
	}
	app.UsageText = "app --filename [filename]"
	app.Author = "Artemii"
	app.Email = "vllemail@icloud.com"
	app.Run(os.Args)
}
