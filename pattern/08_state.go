package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

// Паттерн похож на стратегию:
// 1. Есть контекст, который может находиться в разных состояниях
// 2. Есть интерфейс состояния, который реализуют разные состояния
// 3. Контекст хранит текущее состояние и делегирует ему вызовы методов
// 4. Состояние может менять контекст на другое состояние

// Это простая программа, в которой робот в зависимости от своего состояния будет по-разному приветствовать пользователей. При этом состояние робота меняется в зависимости от того, сколько раз он был приветствован.

type (
	State interface {
		Greet(n string)
		Insult(n string)
	}

	Greet struct {
		robot       *WeirdRobot
		greet_count int
	}

	Ask struct {
		robot     *WeirdRobot
		ask_count int
	}

	Insult struct {
		robot        *WeirdRobot
		insult_count int
	}

	Context struct {
		Greeting  State
		Asking    State
		Insulting State

		current State
	}

	WeirdRobot struct {
		Context
		name string
	}
)

func (g *Greet) Greet(n string) {
	println("Hello, " + n + "!")
	g.greet_count++
	if g.greet_count > 2 {
		g.robot.current = g.robot.Asking
	}
}

func (g *Greet) Insult(n string) {
	println("I don't know you, " + n + "!")
	g.robot.current = g.robot.Insulting
}

func (a *Ask) Greet(n string) {
	println("Do I know you, " + n + "!")
	a.ask_count++
	if a.ask_count > 2 {
		a.robot.current = a.robot.Greeting
	}
}

func (a *Ask) Insult(n string) {
	println("I don't know you, " + n + "!!")
	a.robot.current = a.robot.Insulting
}

func (i *Insult) Greet(n string) {
	println("I don't know you, weirdo " + n + "!!!")
}

func (i *Insult) Insult(n string) {
	println("I don't know you, freakin weirdo " + n + "!!!")
}

func (w *WeirdRobot) Greet(n string) {
	w.current.Greet(n)
}

func (w *WeirdRobot) Insult(n string) {
	w.current.Insult(n)
}

func NewWeirdRobot(n string) *WeirdRobot {
	robot := &WeirdRobot{
		name: n,
	}
	context := Context{
		Greeting:  &Greet{robot: robot},
		Asking:    &Ask{robot: robot},
		Insulting: &Insult{robot: robot},
	}
	robot.Context = context
	robot.current = robot.Greeting
	return robot
}

func example5() {
	robot := NewWeirdRobot("John")
	robot.Greet("Alice")
	robot.Greet("Bob")
	robot.Greet("Alice")
	robot.Greet("Bob")
	new_robot := NewWeirdRobot("John")
	new_robot.Insult("Alice")
	new_robot.Greet("Bob")
}
