Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
[3 2 3]
после i[0] = "3" в массиве внутри слайса изменится значение 0 элемента на 3, но после первого append-a в i будет лежать уже новый слайс и новая структура, также будет выделен новый массив, и запись элементов в нем не будет отражаться в памяти оригинального слайса s
```
