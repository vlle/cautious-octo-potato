package main

/*
Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).


Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

  "github.com/fatih/color"
	"github.com/urfave/cli"
)

type flags interface {
	after() int
	before() int
	context() int
	count() bool
	ignore_case() bool
	invert() bool
	fixed() bool
	line_num() bool
}

type grep_flags struct {
	_after   int
	_before  int
	_context int // possible optimization

	_civFn uint8
}

func (f *grep_flags) after() int {
	return f._after
}

func (f *grep_flags) before() int {
	return f._before
}

func (f *grep_flags) context() int {
	return f._context
}

func (f *grep_flags) count() bool {
	return (f._civFn & 0x16) == 16
}

func (f *grep_flags) set_ignore_case() {
	f._civFn = f._civFn | 1
}

func (f *grep_flags) set_invert() {
	f._civFn = f._civFn | 2
}

func (f *grep_flags) set_fixed() {
	f._civFn = f._civFn | 4
}

func (f *grep_flags) set_line_num() {
	f._civFn = f._civFn | 8
}

func (f *grep_flags) set_count() {
	f._civFn = f._civFn | 16
}

func (f *grep_flags) ignore_case() bool {
	return f._civFn&1 == 1
}
func (f *grep_flags) invert() bool {
	return f._civFn&2 == 2
}

func (f *grep_flags) fixed() bool {
	return f._civFn&4 == 4
}

func (f *grep_flags) line_num() bool {
	return f._civFn&8 == 8
}

func cli_grep(c *cli.Context) {
	filename := c.String("filename")
	if filename == "" {
		// c.App.Writer.Write([]byte("filename is empty\n"))
		os.Exit(1)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	pattern := c.String("pattern")

	flags := &grep_flags{
		_after:   c.Int("A"),
		_before:  c.Int("B"),
		_context: c.Int("C"),
		_civFn:   0,
	}

	if c.Int("C") != 0 {
		flags._after = flags._context
		flags._before = flags._context
	}

	if c.Bool("c") {
		flags.set_count()
	}
	if c.Bool("i") {
		flags.set_ignore_case()
		pattern = "(?i)" + pattern
	}
	if c.Bool("v") {
		flags.set_invert()
	}
	if c.Bool("F") {
		flags.set_fixed()
	}
	if c.Bool("n") {
		flags.set_line_num()
	}

	_, _ = grep(scanner, pattern, flags)
}

func load_all_file(scanner *bufio.Scanner) []string {
  res := []string{}
	for scanner.Scan() {
		res = append(res, scanner.Text())
  }
  return res
}

func fixed_grep(scanner *bufio.Scanner, pattern string, f flags) ([]byte, error) {
	line_num := 0
	for scanner.Scan() {
		txt := scanner.Text()
		line_num++
		if txt == pattern {
			if f.line_num() {
				fmt.Printf("%d: ", line_num)
			}
			fmt.Printf("%s\n", txt)
		}
	}
	return []byte{}, nil
}

func print_after(i, idx int, file[]string) {
  for j := 1; j <= i; j++ {
    if idx+j > len(file) {
      break
    }
    fmt.Println(file[idx+j])
  }
}

func print_before(i, idx int, file[]string) {
  for j := i; j > 0; j-- {
    if idx-j < 0 {
      break
    }
    fmt.Println(file[idx-j])
  }
}

func grep_context(file[]string, pattern *regexp.Regexp, f flags) {
  for i := range file {
    if pattern.MatchString(file[i]) {
      print_before(f.before(), i, file)
      fmt.Println(file[i])
      print_after(f.after(), i, file)
    }
  }
}

func grep(scanner *bufio.Scanner, pattern string, f flags) ([]byte, error) {
	if pattern == "" {
		return []byte{}, nil
	}
	if f.fixed() {
		fixed_grep(scanner, pattern, f)
		return []byte{}, nil
	}
	var search_pattern = regexp.MustCompile(pattern)
  if f.after() > 0 || f.before() > 0 || f.context() > 0 {
    strings := load_all_file(scanner)
    grep_context(strings, search_pattern, f)
		return []byte{}, nil
  }
	line_num := 0
	count := 0
	for scanner.Scan() {
		txt := scanner.Text()
		line_num++
		if !search_pattern.MatchString(txt) && f.invert() {
			if f.count() {
				count++
				continue
			}
			if f.line_num() {
				fmt.Printf("%d:", line_num)
			}
      color.Red(txt)
		} else {
			if search_pattern.MatchString(txt) && !f.invert() {
				if f.count() {
					count++
					continue
				}
				if f.line_num() {
					fmt.Printf("%d:", line_num)
				}
        color.Red(txt)
			}
		}
	}
	if f.count() {
		fmt.Println(count)
	}
	return []byte{}, nil
}

func main() {
	app := cli.NewApp()
	app.Action = cli_grep
	app.Name = "grep utility"
	app.Description = "a grep utility for L2 course at WB"
	app.UsageText = "grep --filename [filename]"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "filename", Usage: "a file to search"},
		cli.StringFlag{Name: "pattern", Usage: "a pattern to search"},
		cli.BoolFlag{Name: "n", Usage: "line num, напечатать номер строки"},
		cli.BoolFlag{Name: "v", Usage: "invert search"},
		cli.BoolFlag{Name: "F", Usage: "fixed search"},
		cli.BoolFlag{Name: "i", Usage: "case insensitive search"},
		cli.BoolFlag{Name: "c", Usage: "count matching strings"},
		cli.IntFlag{Name: "A", Usage: "print N strings after match"},
		cli.IntFlag{Name: "B", Usage: "print N strings before match"},
		cli.IntFlag{Name: "C", Usage: "print N strings before and after match"},
	}
	app.Author = "Artemii"
	app.Email = "vllemail@icloud.com"
	app.Run(os.Args)
}
