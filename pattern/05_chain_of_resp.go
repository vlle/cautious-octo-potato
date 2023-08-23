package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/


type ServiceWorker interface {
  execute(c *Client)
  setNext(w ServiceWorker)
}

type Client struct {
  in_restaraunt bool
  ordered bool
  eating bool
} 

type Host struct {
  next ServiceWorker
}

func (h Host) execute(c *Client) {
  if c.in_restaraunt == true {
    fmt.Println("client in restaraunt, proceed to next service worker")
  } else {
    fmt.Println("Hi, follow me to your seats.")
    c.in_restaraunt = true
  }
  h.next.execute(c)
  return
}

func (h *Host) setNext(w ServiceWorker) {
  h.next = w
}

type Waiter struct {
  next ServiceWorker
}


func (w Waiter) execute(c *Client) {
  if c.ordered == true {
    fmt.Println("client already made order, proceed to next service worker")
  } else {
    fmt.Println("Hi, can I take your order?.")
    c.ordered = true
  }
  w.next.execute(c)
  return
}

func (w *Waiter) setNext(wrk ServiceWorker) {
  w.next = wrk
}

type Cook struct {
  next ServiceWorker
}

func (c Cook) execute(client *Client) {
  if client.eating == true {
    fmt.Println("client already eating...")
  } else {
    fmt.Println("cooking and serving..")
    client.eating = true
  }
  return
}

func (c *Cook) setNext(wrk ServiceWorker) {
  c.next = wrk
}

func example2() {
  client := Client{}
  h := Host{}
  w := Waiter{}
  h.setNext(&w)
  c := Cook{}
  w.setNext(&c)
  h.execute(&client)
}

