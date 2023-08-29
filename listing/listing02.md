Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1

defer может изменить значение named return переменной: https://go.dev/blog/defer-panic-and-recover
он работает по принципу LIFO (last in, first out), т.е. последний defer будет выполнен первым, и этот список вызывается после Return. Аргументы и функция для defer вычисляются сразу же, но исполняется только после покидания функции. 

```
