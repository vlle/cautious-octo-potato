package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/


type Command interface {
  execute()
}

type Button struct {
  cmd Command
}


type Heater interface {
  heatRoom()
  coolRoom()
}

type AC struct {
  is_heating bool
}

func (a* AC) heatRoom() {
  fmt.Println("heating..")
  a.is_heating = true
}

func (a* AC) coolRoom() {
  fmt.Println("cooling..")
  a.is_heating = false
}

type HeatCommand struct {
  ac *AC
}

func (h HeatCommand) execute() {
  h.ac.heatRoom()
}

type CoolCommand struct {
  ac *AC
}

func (h CoolCommand) execute() {
  h.ac.coolRoom()
}

func example() {
  ac := &AC{}
  hc := HeatCommand{ac:ac}
  heat_button := Button{cmd:hc}
  cc := CoolCommand{ac:ac}
  cool_button := Button{cmd:cc}
  heat_button.cmd.execute()
  cool_button.cmd.execute()
}
