Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
потому что тип err – указатель на структуру customError. Интерфейс равен nil только если тип и значение равны nil. То есть если бы тип err был бы interface{} – напечатлось бы ok
https://go.dev/doc/faq#nil_error

```
