Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
77 78 79
Создается слайс с началом в 77 и концом в 79. Последний элемент в срезе (4) не берется в счет.

```