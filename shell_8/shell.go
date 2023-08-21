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
  "os/exec"
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

func echo(input ...string) (string, int) {
	err := 0
  ret := ""
  if len(input) == 1 {
    return ret, err
  } else {
    fmt.Println(strings.Join(input[1:], ""))
  }
	return ret, err
}

func kill(proc_id string) (string, int) {
  exec := exec.Command("/bin/sh", "-c", "kill", proc_id)
  ret, err := exec.Output()

  if err != nil {
    fmt.Println(err)
    return string(ret), 1
  }
	return string(ret), 0
}

func ps() (string, int) {
  exec := exec.Command("/bin/sh", "-c", "ps")
  ret, err := exec.Output()

  if err != nil {
    fmt.Println(err)
    return string(ret), 1
  }
	return string(ret), 0
}

func forkExecWrapper(args ...string) int {
    possible_path := "/bin/" + args[0]
    command := exec.Command(possible_path, args[1:]...)
    stdout, err := command.StdoutPipe()
    if err != nil {
      fmt.Println(err)
      return 1
    }
    stderr, err := command.StderrPipe()
    if err != nil {
      fmt.Println(err)
      return 1
    }
    err = command.Start()
    if err != nil {
      fmt.Println(err)
      return 1
    } else {
      go func() {
        output := make([]byte, 100)
        _, err := stdout.Read([]byte(output))
        if err != nil {
          fmt.Println(err)
        }
        err_output := make([]byte, 100)
        _, err = stderr.Read([]byte(err_output))
        if err != nil {
          fmt.Println(err)
        }
        fmt.Printf("%s\n", output)
        fmt.Printf("%s\n", err_output)
        err = command.Wait()
        if err != nil {
          fmt.Println(err)
        }
      }()
    }
    return 0
}

func process_cmd(cmd string, args ...string) string {
  process_code := 1
  ret := ""
  output := ""
  arg := args[0]
  switch arg {
  case "pwd":
    ret, process_code = pwd()
  case "cd":
    if len(arg) < 2 {
      ret, process_code = cd("")
    } else {
      ret, process_code = cd(args[1])
    }
  case "echo":
    ret, process_code = echo(args[1:]...)
  case "kill":
    ret, process_code = kill(args[1])
  case "ps":
    ret, process_code = ps()
  default:
    process_code = forkExecWrapper(args...)
  }
  if process_code != 0 {
    output += color.RedString(">")
  } else {
    output += color.GreenString(">")
  }
  output += ret

  fmt.Printf("%s\n", output)
  return output
}

func shell() {
  scan := bufio.NewScanner(os.Stdin)
  for {
    scan.Scan()
    cmd := scan.Text()
    if cmd == "\\quit" || scan.Err() != nil || cmd == "" {
      break
    }
    if strings.Contains(cmd, "|") {
      // cmds := strings.Split(cmd, "|")
      // arg := strings.Split(cmd, " ")
      // for i, cmd := range cmds {
      //   args = process_cmd(cmd)
      // }
    } else {
      arg := strings.Split(cmd, " ")
      process_cmd(cmd, arg...)
    }
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
