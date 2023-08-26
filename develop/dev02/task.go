/*
Задача на распаковку

Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую
повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

*/

package main


func unpack(s string) string {
	s2 := ""
	var letter rune = 0
	for _, val := range s {
		if int(val)-'0' >= 0 && int(val)-'0' <= 9 {
			if letter == 0 {
				return ""
			}
			for i := 0; i < (int(val) - '0'); i++ {
				s2 += string(letter)
			}
			letter = 1
		} else {
			if letter != 1 && letter != 0 {
				s2 += string(letter)
			}
			letter = val
		}
	}
	if letter == 0 {
		return ""
	}
	if letter != 1 {
		s2 += string(letter)
	}

	return s2
}
