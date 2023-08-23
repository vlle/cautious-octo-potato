package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

import "fmt"

type Worker interface {
  doWork()
  accept(v Visitor)
}

type Engineer struct {
  salary int
  Worker
}

func (e Engineer) doWork() {
  fmt.Println("I am building something..")
}

func (e Engineer) accept(v Visitor) {
  v.visitForEngineer(&e)
}

type Doctor struct {
  salary int
  Worker
}

func (d Doctor) doWork() {
  fmt.Println("I am healing someone..")
}

func (d Doctor) accept(v Visitor) {
  v.visitForDoctor(&d)
}

type Teacher struct {
  salary int
  Worker
}

func (t Teacher) doWork() {
  fmt.Println("I am teaching someone..")
}

func (t Teacher) accept(v Visitor) {
  v.visitForTeacher(&t)
}

type Visitor interface {
  visitForEngineer(e *Engineer)
  visitForTeacher(t *Teacher)
  visitForDoctor(d *Doctor)
}

type SalaryIncreaser struct {
  Visitor
}

func (h *SalaryIncreaser) visitForEngineer(e *Engineer) {
  e.salary += 5
  fmt.Println("Here is your increased salary, Engineer")
}

func (h *SalaryIncreaser) visitForTeacher(t *Teacher) {
  t.salary += 3
  fmt.Println("Here is your increased salary, Teacher")
}

func (h *SalaryIncreaser) visitForDoctor(d *Doctor) {
  d.salary += 10
  fmt.Println("Here is your increased salary, Doctor")
}
