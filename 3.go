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
	"log"
	"os"
  "bufio"
  "fmt"
  "sort"
)

func sort_util(filename string) {
  file, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  lines := []string{}
  // read file
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
    fmt.Println(scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
  // sort file 
  sort.Strings(lines)
  for _, line := range lines {
    fmt.Println(line)
  }
  // implement flags later
}
