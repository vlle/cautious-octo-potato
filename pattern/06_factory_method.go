package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

type (
  DisplayFactory struct {

  }

  Displayer interface {
    ShowVideo()
  }

  Monitor struct {

  }

  TV struct {

  }

)

func (M Monitor) ShowVideo() {
  fmt.Println("Watching videogame..")
}

func (D TV) ShowVideo() {
  fmt.Println("Watching football..")
}

// Ключевая идея: фабрика инкапсулирует создание объектов (скрывает от пользователя кода), и позволяет создавать объекты разных типов, но с одинаковым интерфейсом.
func (F DisplayFactory) Create(s string) (Displayer, error) {
  if s == "Monitor" {
    return Monitor{}, nil
  } else if s == "TV" {
    return TV{}, nil
  } else {
    return nil, fmt.Errorf("Unknown display type")
  }
}

func example3() {
  factory := DisplayFactory{}
  monitor, _ := factory.Create("Monitor")
  tv, _ := factory.Create("TV")
  monitor.ShowVideo()
  tv.ShowVideo()
}
