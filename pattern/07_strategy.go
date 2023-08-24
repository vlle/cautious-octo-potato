package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

import "fmt"

// Ключевая идея: объединить семейство алгоритмов под один интерфейс для легкой взаимозаменяемости в зависимости от нужды
// Также нужно использовать контекст для первичной настройки и смены алгоритма
type (
	Game struct {
		str Strategy
	}

	Strategy interface {
		Execute(n int) string
	}

	Multiplier struct {
	}

	Divider struct {
	}

	ContextStrategy interface {
		SetStrategy()
		Strategy()
	}
)

func (M Multiplier) Execute(n int) string {
	return fmt.Sprintln("I am.. MULTIPLIER.. Your answer.. is..", n*n)
}
func (D Divider) Execute(n int) string {
	return fmt.Sprintln("I am.. DIVIDER.. Your answer.. is..", n/2)
}

func (g *Game) SetStrategy(new_str Strategy) {
	g.str = new_str
}

func (g *Game) Strategy() Strategy {
	return g.str
}

func (g Game) Play(n int) string {
	return g.str.Execute(n)
}

func example4() {
	g := Game{}
	m := Multiplier{}
	d := Divider{}
	g.SetStrategy(m)
	fmt.Print(g.Play(10))
	g.SetStrategy(d)
	fmt.Print(g.Play(10))
}
