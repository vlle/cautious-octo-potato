/*
Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func cut(c *cli.Context) {
  scanner := bufio.NewScanner(os.Stdin)
  delim := string(rune(9))
  if c.String("d") != "" {
    delim = c.String("d")
  }
  f := c.Int("f")
  if f == 0 {
    f = 1
  }
  for scanner.Scan() {
    txt := strings.Split(scanner.Text(), delim)
    if len(txt) > c.Int("f") {
      fmt.Println(txt[f-1])
    } else {
      fmt.Println(strings.Join(txt, ""))
    }
  }

  if scanner.Err() != nil {
    fmt.Println(scanner.Err())
  }

}

func main() {
	app := cli.NewApp()
	app.Action = cut
	app.Name = "cut utility"
	app.Description = "a cut utility for L2 course at WB"
	app.UsageText = "cut [flags]"
	app.Author = "Artemii"
	app.Email = "vllemail@icloud.com"
	app.Flags = []cli.Flag{
		cli.IntFlag{Name: "f", Usage: "a field to output"},
		cli.StringFlag{Name: "d", Usage: "specify delimiter"},
  }
  app.Run(os.Args)
}


type flags interface {
	fields() int
	delimiter() string
	separated() bool
}

type cut_flags struct {
	_fields   int
	_delimiter  string
	_separated bool // possible optimization
}

func (f *cut_flags) fields() int {
  return f._fields
}

func (f *cut_flags) delimiter() string {
  return f._delimiter
}

func (f *cut_flags) separated() bool {
  return f._separated
}
