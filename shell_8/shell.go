package main

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

import (
	"bufio"
	"fmt"
  "strings"
	"syscall"

	// "log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func pwd() (string, int) {
  dir, err := syscall.Getwd()
  if err != nil {
    return dir, 1
  }
	return dir, 0
}

func cd(dir string) (string, int) {
  ret := ""
  err := syscall.Chdir(dir)
  if err != nil {
    fmt.Println(err)
    return ret, 1
  }
	return ret, 0

}

func echo() (string, int) {
	err := 0
  ret := ""
	return ret, err
}

func kill() (string, int) {
	err := 0
  ret := ""
	return ret, err
}

func ps() (string, int) {
	err := 0
  ret := ""
	return ret, err
}

func shell() {
	scan := bufio.NewScanner(os.Stdin)
	process_code := 1
	for {
		scan.Scan()
		cmd := scan.Text()
		if cmd == "\\quit" || scan.Err() != nil || cmd == "" {
			break
		}
    // if strings.Contains(cmd, "|") {
    // }
    ret := ""
    output := ""
    arg := strings.Split(cmd, " ")
		switch arg[0] {
		case "pwd":
			ret, process_code = pwd()
		case "cd":
      if len(arg) < 2 {
        ret, process_code = cd("")
      } else {
        ret, process_code = cd(arg[1])
      }
    case "echo":
      ret, process_code = echo()
    case "kill":
      ret, process_code = kill()
    case "ps":
      ret, process_code = ps()
    }
    if process_code != 0 {
      output += color.RedString(">")
    } else {
      output += color.GreenString(">")
    }
    output += ret

    fmt.Printf("%s\n", output)
    process_code = 1
  }
}

func cli_startup(c *cli.Context) {
  shell()
}

func main() {
  app := cli.NewApp()
  fmt.Printf("%s", color.RedString(">"))
  app.Action = cli_startup
  app.Name = "shell utility"
  app.Description = "a shell utility for L2 course at WB"
  app.UsageText = "shell"
  app.Author = "Artemii"
  app.Email = "vllemail@icloud.com"
  app.Run(os.Args)
}
