package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type WaterHolder struct {
}

type CoffeeRoaster struct {
}

type CoffeeGrinder struct {
}

func (c *CoffeeRoaster) RoastCoffee() {
	fmt.Println("roasting coffee..")
}

func (g *CoffeeGrinder) GrindCoffee() {
	fmt.Println("grinding coffee..")
}

func (w *WaterHolder) BoilWater() {
	fmt.Println("boiling water..")
}

func (w *WaterHolder) PouringWater() {
	fmt.Println("pouring water.. ")
}

type ElectricCoffeeMaker_Facade struct {
	w *WaterHolder
	c *CoffeeRoaster
	g *CoffeeGrinder
}

func (facade *ElectricCoffeeMaker_Facade) MakeCoffee() {
	facade.c.RoastCoffee()
	facade.g.GrindCoffee()

	facade.w.BoilWater()
	facade.w.PouringWater()

	fmt.Println("Coffee ready")
}
