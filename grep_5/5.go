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

	"github.com/urfave/cli"
)

type flags interface {
	after() int
	before() int
	context() int
	count() int
	ignore_case() bool
	invert() bool
	fixed() bool
	line_num() bool
}

type grep_flags struct {
	_after   int
	_before  int
	_context int // possible optimization
	_count   int

	_ignore_invert_fixed_line uint8
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

func (f *grep_flags) count() int {
	return f._count
}

func (f *grep_flags) set_ignore_case() {
	f._ignore_invert_fixed_line = f._ignore_invert_fixed_line | 0x1
}

func (f *grep_flags) set_invert() {
	f._ignore_invert_fixed_line = f._ignore_invert_fixed_line | 0x2
}

func (f *grep_flags) set_fixed() {
	f._ignore_invert_fixed_line = f._ignore_invert_fixed_line | 0x4
}

func (f *grep_flags) set_line_num() {
	f._ignore_invert_fixed_line = f._ignore_invert_fixed_line | 0x8
}

func (f *grep_flags) ignore_case() bool {
	return f._ignore_invert_fixed_line&0x1 == 1
}
func (f *grep_flags) invert() bool {
	return f._ignore_invert_fixed_line&0x2 == 2
}

func (f *grep_flags) fixed() bool {
	return f._ignore_invert_fixed_line&0x4 == 4
}

func (f *grep_flags) line_num() bool {
	return (f._ignore_invert_fixed_line & 0x8) == 8
}

func cli_grep(c *cli.Context) {
	filename := c.String("filename")
	fmt.Println(c.Args().First(), c.Args().Tail())
	fmt.Println(filename)
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
		_after:                    c.Int("after"),
		_before:                   c.Int("before"),
		_context:                  c.Int("context"),
		_count:                    c.Int("count"),
		_ignore_invert_fixed_line: 0,
	}
	if c.Bool("ignore-case") {
		flags.set_ignore_case()
	}
	if c.Bool("v") {
		flags.set_invert()
	}
	if c.Bool("fixed") {
		flags.set_fixed()
	}
	if c.Bool("n") {
		flags.set_line_num()
	}
	_, _ = grep(scanner, pattern, flags)
}

func fixed_grep(scanner *bufio.Scanner, pattern string, f flags) ([]byte, error) {
	line_num := 0
	for scanner.Scan() {
		txt := scanner.Text()
		line_num++
		if txt == pattern {
			fmt.Println(txt)
		}
	}
	return []byte{}, nil
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
	line_num := 0
	for scanner.Scan() {
		txt := scanner.Text()
		line_num++
		if !search_pattern.MatchString(txt) && f.invert() {
			if f.line_num() {
				fmt.Printf("%d: ", line_num)
			}
			fmt.Printf("%s\n", txt)
		} else {
			if search_pattern.MatchString(txt) && !f.invert() {
				if f.line_num() {
					fmt.Printf("%d: ", line_num)
				}
				fmt.Printf("%s\n", txt)
			}
		}
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
	}
	app.Author = "Artemii"
	app.Email = "vllemail@icloud.com"
	app.Run(os.Args)
}
