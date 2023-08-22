/*
Поиск анаграмм по словарю

Написать функцию поиска всех множеств анаграмм по словарю.

Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Требования:
Входные данные для функции: ссылка на массив, каждый элемент которого - слово на русском языке в кодировке utf8
Выходные данные: ссылка на мапу множеств анаграмм
Ключ - первое встретившееся в словаре слово из множества. Значение - ссылка на массив, каждый элемент которого,
слово из множества.
Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.
*/
package main

import (
	"sort"
	"strings"
)

func findWord(words []string, word string) int {
	for i, v := range words {
		if v == word {
			return i
		}
	}
	return -1
}

func f4(words []string) map[string][]string {
	result := make(map[string][]string)
	r := make(map[[33]int]string)
	for _, word := range words {
		if len(word) == 1 {
			continue
		}
		word = strings.ToLower(word)
		letters_array := [33]int{}
		for _, letter := range word {
			letters_array[letter-'а']++
		}
		key, ok := r[letters_array]
		if !ok {
			r[letters_array] = word
			key = word
		}
		if result[key] == nil {
			result[key] = make([]string, 0)
		}
		if findWord(result[key], word) == -1 {
			result[key] = append(result[key], word)
			sort.Sort(sort.StringSlice(result[key]))
		}
	}
	return result
}
