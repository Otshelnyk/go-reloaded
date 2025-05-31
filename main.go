package main

import (
	"encoding/json"
	"fmt" // Импортируем пакет fmt для вывода на экран
	"io/ioutil"
	"net/http"
	"os"      // Импортируем пакет os для работы с файловой системой
	"regexp"  // Импортируем пакет regexp для работы с регулярными выражениями
	"strconv" // Импортируем пакет strconv для конвертации строк в числа
	"strings" // Импортируем пакет strings для работы со строками
	"unicode" // Импортируем пакет unicode для работы с юникодными символами
)

func main() {
	if len(os.Args) == 3 {
		// Читаем содержимое входного файла
		content := readFile(os.Args[1])
		// Разделяем содержимое файла на строки
		lines := strings.Split(content, "\n")
		// Проверяем, что весь текст состоит только из ASCII-символов
		if !isASCII(content) {
			fmt.Println("Error: Text contains non-ASCII characters") // Выводим ошибку, если есть не-ASCII символы
			return
		}
		// Создаем срез для хранения обработанных строк
		var processedLines []string
		// Обрабатываем каждую строку из входного файла
		for _, line := range lines {
			// Если строка пустая, просто добавляем её в результат
			if strings.TrimSpace(line) == "" {
				processedLines = append(processedLines, "")
				continue
			}
			str := " " + line // Копируем строку для обработки
			// Добавляем пробелы перед открывающимися скобками, если их нет
			str = addSpacesBeforeBrackets(str)
			iterations := 0                         // Счетчик итераций
			maxIterations := 10                     // Максимальное количество итераций
			processedBrackets := make(map[int]bool) // Словарь для отслеживания обработанных скобок
			// Цикл, который будет работать максимум 10 раз
			for iterations < maxIterations {
				iterations++
				// Находим все позиции открывающих скобок
				startArr := findStartBrackets(str)
				// Находим все позиции закрывающих скобок
				endArr := findEndBrackets(str)
				// Если нет скобок, выходим из цикла
				if len(startArr) == 0 || len(endArr) == 0 {
					break
				}
				startPos := -1 // Начальная позиция скобки
				endPos := -1   // Конечная позиция скобки
				// Ищем пару скобок (открывающую и закрывающую)
				for i := 0; i < len(endArr); i++ {
					// Если скобка уже была обработана, пропускаем её
					if processedBrackets[endArr[i]] {
						continue
					}
					currentEnd := endArr[i]
					// Ищем соответствующую открывающую скобку
					for j := len(startArr) - 1; j >= 0; j-- {
						if startArr[j] < currentEnd {
							startPos = startArr[j]
							endPos = currentEnd
							break
						}
					}
					// Если нашли пару, выходим из цикла
					if startPos != -1 {
						break
					}
				}
				// Если не нашли пару скобок, выходим из цикла
				if startPos == -1 || endPos == -1 {
					break
				}
				// Извлекаем строку между скобками
				newStr := strings.TrimSpace(str[startPos+1 : endPos])
				// Проверяем, является ли команда между скобками валидной
				if !isValidCommandFormat(newStr) {
					processedBrackets[startPos] = true // Отмечаем скобку как обработанную
					processedBrackets[endPos] = true
					continue
				}
				processedBrackets[startPos] = true // Отмечаем скобку как обработанную
				// Получаем команду и число (если оно есть)
				newStr2, num := getSubStrAndNum(newStr)
				// В зависимости от команды выполняем соответствующие действия
				switch strings.TrimSpace(newStr2) {
				case "hex": // Если команда "hex", конвертируем из шестнадцатеричной в десятичную
					str = hexToDec(str)
				case "bin": // Если команда "bin", конвертируем из бинарной в десятичную
					str = binToDec(str)
				case "up", "cap", "low": // Если команда "up", "cap", "low", преобразуем слова
					beforeBracket := str[:startPos] // Часть строки до скобки
					afterBracket := str[endPos+1:]  // Часть строки после скобки
					// Разбиваем строку на отдельные слова
					words := strings.Split(beforeBracket, " ")
					wordCount := 1 // Количество слов для изменения
					if num > 0 {
						wordCount = num // Если есть число, меняем указанное количество слов
					}
					lastModifiedIndex := -1 // Индекс последнего измененного слова
					// Обрабатываем слова с конца
					for i := len(words) - 1; i >= 0 && wordCount > 0; i-- {
						if words[i] != "" {
							// В зависимости от команды изменяем слово
							switch strings.TrimSpace(newStr2) {
							case "up":
								words[i] = strings.ToUpper(words[i]) // Переводим в верхний регистр
							case "cap":
								if len(words[i]) > 0 {
									words[i] = strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:]) // Делаем первую букву заглавной
								}
							case "low":
								words[i] = strings.ToLower(words[i]) // Переводим в нижний регистр
							}
							lastModifiedIndex = i // Обновляем индекс последнего измененного слова
							wordCount--           // Уменьшаем количество оставшихся слов для изменения
						}
					}
					if lastModifiedIndex == -1 {
						lastModifiedIndex = 0 // Или пропустите этот блок обработки
					}
					// Собираем строку обратно
					beforeWords := strings.Join(words[:lastModifiedIndex], " ")
					afterWords := strings.Join(words[lastModifiedIndex+1:], " ")
					// Добавляем пробелы вокруг слов
					if lastModifiedIndex > 0 && len(beforeWords) > 0 {
						beforeWords += " "
					}
					if lastModifiedIndex < len(words)-1 && len(afterWords) > 0 {
						afterWords = " " + afterWords
					}
					str = beforeWords + words[lastModifiedIndex] + afterWords + afterBracket // Обновляем строку
				}
			}
			// Преобразуем артикли "a" и "an" в зависимости от следующего слова
			str = transformAtoAn(str)
			// Форматируем пунктуацию
			str = formatPunctuation(str)
			// Настроить пробелы вокруг двойных кавычек
			str = adjustSpacesAroundDoubleQuotes(str)
			// Настроить пробелы вокруг одинарных кавычек
			str = adjustSpacesAroundSingleQuotes(str)
			// Добавляем обработанную строку в результат
			str = removeDoubleSpaces(str)
			processedLines = append(processedLines, str)
		}
		// Объединяем обработанные строки в одну
		result := strings.Join(processedLines, "\n")
		// Выводим результат
		fmt.Println(result)
		// Записываем результат в выходной файл
		writeFile(os.Args[2], result)
	} else {
		http.HandleFunc("/api/process", processHandler)
		http.Handle("/", http.FileServer(http.Dir(".")))
		fmt.Println("Server started at http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	}
}

// Проверка на наличие только ASCII символов в строке
func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 { // Если символ больше 127, то это не ASCII
			return false
		}
	}
	return true // Все символы ASCII
}

func isValidCommandFormat(cmd string) bool {
	// Удаляем пробелы
	cmd = strings.TrimSpace(cmd)
	// Проверяем базовые команды
	if cmd == "hex" || cmd == "bin" || cmd == "up" || cmd == "low" || cmd == "cap" {
		return true
	}
	// Проверяем команды с числом
	if strings.Contains(cmd, ",") {
		parts := strings.SplitN(cmd, ",", 2)
		baseCmd := strings.TrimSpace(parts[0])
		numStr := strings.TrimSpace(parts[1])
		// Проверяем базовую часть команды
		if baseCmd != "up" && baseCmd != "low" && baseCmd != "cap" {
			return false
		}
		// Проверяем, что после запятой идет число и оно не отрицательное
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 0 {
			return false
		}
		return true
	}
	return false
}

// Функция для добавления пробела перед открывающимися скобками
func addSpacesBeforeBrackets(s string) string {
	var result string
	for i, v := range s {
		// Если символ - открывающая скобка, а перед ним нет пробела
		if v == '(' && i > 0 && !unicode.IsSpace(rune(s[i-1])) {
			result += " " // Добавляем пробел
		}
		result += string(v) // Добавляем сам символ
	}
	return result // Возвращаем строку с добавленными пробелами
}

// Функция для преобразования бинарного числа в десятичное
func binToDec(sentence string) string {
	// Создаем регулярное выражение для поиска "( bin )" в строке
	re := regexp.MustCompile(`\(\s*bin\s*\)`)

	// Бесконечный цикл для поиска всех вхождений бинарных чисел
	for {
		// Ищем индекс вхождения шаблона
		loc := re.FindStringIndex(sentence)
		if loc == nil {
			break // Если больше нет вхождений, выходим из цикла
		}

		// Разделяем строку на части до и после найденного выражения
		before := sentence[:loc[0]]
		after := sentence[loc[1]:]

		// Убираем пробелы с конца части до скобки
		before = strings.TrimRight(before, " ")

		// Находим последнее слово перед скобкой
		lastSpace := strings.LastIndex(before, " ")
		var word, beforeWord string
		if lastSpace == -1 {
			word = before // Если пробела нет, все это одно слово
			beforeWord = ""
		} else {
			// Разделяем строку на слово и переднее пространство
			word = before[lastSpace+1:]
			beforeWord = before[:lastSpace]
		}

		// Преобразуем бинарное число в десятичное
		decimalNum, err := strconv.ParseInt(word, 2, 64)
		if err != nil {
			// Если произошла ошибка при конвертации, возвращаем строку без изменений
			sentence = beforeWord + " " + word + after
			fmt.Printf("Error bin number: %s", word) // Выводим ошибку
			continue                                 // Переходим к следующей итерации
		}

		// Преобразуем десятичное число обратно в строку
		newWord := fmt.Sprintf("%d", decimalNum)

		// Собираем строку обратно, добавляем преобразованное число
		if beforeWord == "" {
			sentence = newWord + after
		} else {
			sentence = beforeWord + " " + newWord + after
		}
	}

	// Возвращаем строку без лишних пробелов
	return strings.TrimSpace(sentence)
}

// Функция для преобразования шестнадцатеричного числа в десятичное
func hexToDec(sentence string) string {
	// Создаем регулярное выражение для поиска "( hex )" в строке
	re := regexp.MustCompile(`\(\s*hex\s*\)`)

	// Бесконечный цикл для поиска всех вхождений шестнадцатеричных чисел
	for {
		// Ищем индекс вхождения шаблона
		loc := re.FindStringIndex(sentence)
		if loc == nil {
			break // Если больше нет вхождений, выходим из цикла
		}

		// Разделяем строку на части до и после найденного выражения
		before := sentence[:loc[0]]
		after := sentence[loc[1]:]

		// Убираем пробелы с конца части до скобки
		before = strings.TrimRight(before, " ")

		// Если нет слов перед выражением, пропускаем его
		if len(before) == 0 {
			sentence = after
			continue
		}

		// Находим последнее слово перед скобкой
		lastSpace := strings.LastIndex(before, " ")
		var word, beforeWord string
		if lastSpace == -1 {
			word = before // Если пробела нет, все это одно слово
			beforeWord = ""
		} else {
			// Разделяем строку на слово и переднее пространство
			word = before[lastSpace+1:]
			beforeWord = before[:lastSpace]
		}

		// Преобразуем шестнадцатеричное число в десятичное
		decimalNum, err := strconv.ParseInt(word, 16, 64)
		if err != nil {
			// Если произошла ошибка при конвертации, возвращаем строку без изменений
			sentence = beforeWord + " " + word + after
			fmt.Printf("Error hex number: %s", word) // Выводим ошибку
			continue                                 // Переходим к следующей итерации
		}

		// Преобразуем десятичное число обратно в строку
		newWord := fmt.Sprintf("%d", decimalNum)

		// Собираем строку обратно, добавляем преобразованное число
		if beforeWord == "" {
			sentence = newWord + after
		} else {
			sentence = beforeWord + " " + newWord + after
		}
	}

	// Возвращаем строку без лишних пробелов
	return strings.TrimSpace(sentence)
}

// Функция для форматирования пунктуации в строке
func formatPunctuation(input string) string {
	// Убираем пробелы вокруг знаков пунктуации
	re := regexp.MustCompile(`\s*([.,!?:;])\s*`)
	input = re.ReplaceAllString(input, "$1")
	// Добавляем пробел после знаков пунктуации, если следом идет буква или цифра
	re = regexp.MustCompile(`([.,!?:;])([a-zA-Z0-9-])`)
	input = re.ReplaceAllString(input, "$1 $2")
	// Убираем двойные пробелы
	input = strings.Join(strings.Fields(input), " ")
	// Убираем начальные и конечные пробелы
	return strings.TrimSpace(input)
}

// Функция для извлечения команды и числа из строки
func getSubStrAndNum(newStr string) (string, int) {
	// Убираем пробелы с начала и конца строки
	newStr = strings.TrimSpace(newStr)

	// Если строка начинается с "(", и заканчивается на ")", удаляем эти символы
	if strings.HasPrefix(newStr, "(") && strings.HasSuffix(newStr, ")") {
		newStr = newStr[1 : len(newStr)-1]
	}

	// Проверяем, есть ли запятая в строке (для команд с числом)
	if hasComma(newStr) {
		// Разделяем строку на команду и число
		parts := strings.SplitN(newStr, ",", 2)
		if len(parts) != 2 {
			// Если формат неверный, выводим ошибку
			fmt.Println("Неверный формат команды:", newStr)
			return newStr, 1
		}

		// Получаем команду и число
		cmd := strings.TrimSpace(parts[0])
		numStr := strings.TrimSpace(parts[1])

		// Убираем все пробелы в числе
		numStr = strings.ReplaceAll(numStr, " ", "")

		// Преобразуем строку в число
		num, err := strconv.Atoi(numStr)
		if err != nil {
			// Если произошла ошибка при конвертации, выводим ошибку
			fmt.Println("Ошибка конвертации числа:", err)
			return cmd, 1
		}

		// Возвращаем команду и число
		return cmd, num
	} else {
		// Если запятой нет, возвращаем команду и 1
		return strings.TrimSpace(newStr), 1
	}
}

// Функция для чтения файла и возврата его содержимого как строки
func readFile(s string) string {
	// Читаем файл
	b, err := os.ReadFile(s)
	if err != nil {
		// Если ошибка при чтении, выводим её
		fmt.Print(err)
	}
	// Возвращаем содержимое файла в виде строки
	return string(b)
}

// Функция для записи строки в файл
func writeFile(filename, myString string) {
	// Открываем файл для записи
	f, err := os.Create(filename)
	if err != nil {
		// Если ошибка при открытии, выводим её
		fmt.Println(err)
	}
	defer f.Close() // Закрываем файл по завершению

	// Пишем строку в файл
	_, err2 := f.WriteString(myString)
	if err2 != nil {
		// Если ошибка при записи, выводим её
		fmt.Println(err2)
	}
}

// Функция для поиска всех открывающих скобок в строке и возврата их позиций
func findStartBrackets(s string) []int {
	var mArr []int
	// Проходим по всем символам строки
	for index, v := range s {
		// Если символ - открывающая скобка, добавляем её позицию в массив
		if v == '(' {
			mArr = append(mArr, index)
		}
	}
	// Возвращаем массив с позициями открывающих скобок
	return mArr
}

// Функция для поиска всех закрывающих скобок в строке и возврата их позиций
func findEndBrackets(s string) []int {
	var mArr []int
	// Проходим по всем символам строки
	for index, v := range s {
		// Если символ - закрывающая скобка, добавляем её позицию в массив
		if v == ')' {
			mArr = append(mArr, index)
		}
	}
	// Возвращаем массив с позициями закрывающих скобок
	return mArr
}

// Функция для проверки наличия запятой в строке
func hasComma(s string) bool {
	// Проходим по всем символам строки
	for _, v := range s {
		// Если встречаем запятую, возвращаем true
		if v == ',' {
			return true
		}
	}
	// Если запятая не найдена, возвращаем false
	return false
}

// Функция для преобразования "a" в "an" или наоборот в зависимости от контекста
func transformAtoAn(input string) string {
	words := strings.Fields(input) // Разбиваем строку на слова
	result := make([]string, len(words))
	copy(result, words) // Копируем слова в новый массив

	// Проходим по всем словам
	for i := 0; i < len(words)-1; i++ {
		// Пропускаем слова "and"
		if words[i+1] == "and" {
			continue
		}
		// Преобразуем текущее слово в нижний регистр
		currentWordLower := strings.ToLower(words[i])

		// Проверяем, являются ли текущие и следующие слова заглавными
		isCurrentWordUpper := words[i] == strings.ToUpper(words[i])
		isNextWordUpper := words[i+1] == strings.ToUpper(words[i+1])

		// Если текущее слово "a" и следующее начинается с гласной или "h"
		if currentWordLower == "a" {
			if len(words[i+1]) > 1 {
				firstLetter := strings.ToLower(string(words[i+1][0]))
				if strings.Contains("aeiouh", firstLetter) {
					// Преобразуем "a" в "an" или наоборот в зависимости от контекста
					if isCurrentWordUpper && isNextWordUpper {
						result[i] = "AN"
					} else if words[i+1] == strings.ToUpper(words[i+1]) && words[i] == "A" {
						result[i] = "AN"
					} else if words[i] == "A" {
						result[i] = "An"
					} else {
						result[i] = "an"
					}
				}
			}
		} else if currentWordLower == "an" {
			if len(words[i+1]) > 1 {
				firstLetter := strings.ToLower(string(words[i+1][0]))
				if !strings.Contains("aeiouh", firstLetter) {
					// Преобразуем "an" в "a" или наоборот в зависимости от контекста
					if isCurrentWordUpper && isNextWordUpper {
						result[i] = "A"
					} else if words[i] == "An" {
						result[i] = "A"
					} else {
						result[i] = "a"
					}
				}
			}
		}
	}
	// Возвращаем строку с преобразованными словами
	return strings.Join(result, " ")
}

// Функция для корректировки пробелов вокруг двойных кавычек в строке
func adjustSpacesAroundDoubleQuotes(input string) string {
	// Преобразуем строку в срез рун для удобства обработки
	runes := []rune(input)
	var output []rune // Срез для хранения результирующей строки
	i := 0            // Индекс текущего символа

	// Проходим по всем рунам в строке
	for i < len(runes) {
		// Если встречаем двойные кавычки
		if runes[i] == '"' {
			nextQuote := -1 // Переменная для хранения индекса следующей кавычки

			// Находим индекс следующей двойной кавычки
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '"' {
					nextQuote = j
					break
				}
			}

			// Если нашли пару кавычек
			if nextQuote != -1 {
				// Добавляем пробел перед кавычками, если это необходимо
				if len(output) > 0 && !unicode.IsSpace(output[len(output)-1]) {
					output = append(output, ' ')
				}

				// Добавляем саму кавычку
				output = append(output, '"')

				// Очищаем пробелы внутри кавычек и добавляем содержимое
				content := strings.TrimSpace(string(runes[i+1 : nextQuote]))
				output = append(output, []rune(content)...)

				// Добавляем закрывающую кавычку и пробел после неё
				output = append(output, '"')

				if nextQuote+1 < len(runes) {
					j := nextQuote + 1

					// Убираем пробелы до точки, если они есть
					for j < len(runes) && unicode.IsSpace(runes[j]) {
						j++
					}
					if j < len(runes) && runes[j] == '.' || runes[j] == ',' || runes[j] == ':' || runes[j] == ';' || runes[j] == '?' || runes[j] == '!' {
						output = append(output, runes[j])
						i = j + 1 // Перемещаем индекс на точку
					} else {
						// В противном случае добавляем пробел после закрывающей кавычки
						output = append(output, ' ')
						i = nextQuote + 1
					}
				} else {
					// Если кавычки не имеют пары, добавляем одну кавычку и продолжаем
					i = nextQuote + 1
				}
			} else {
				// Если текущий символ не кавычка, просто добавляем его в результат
				output = append(output, runes[i])
				i++
			}
		} else {
			output = append(output, runes[i])
			i++
		}
	}
	// Возвращаем строку из результата
	return string(output)
}

// Функция для корректировки пробелов вокруг одинарных кавычек в строке
func adjustSpacesAroundSingleQuotes(input string) string {
	// Преобразуем строку в срез рун для удобства обработки
	runes := []rune(input)
	var output []rune // Срез для хранения результирующей строки
	i := 0            // Индекс текущего символа

	// Проходим по всем рунам в строке
	for i < len(runes) {
		// Если встречаем одинарные кавычки
		if runes[i] == '\'' {
			nextQuote := -1 // Переменная для хранения индекса следующей кавычки

			// Находим индекс следующей одинарной кавычки
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '\'' {
					nextQuote = j
					break
				}
			}

			// Если нашли пару кавычек
			if nextQuote != -1 {
				// Добавляем пробел перед кавычками, если это необходимо
				if len(output) > 0 && !unicode.IsSpace(output[len(output)-1]) {
					output = append(output, ' ')
				}

				// Добавляем саму кавычку
				output = append(output, '\'')

				// Очищаем пробелы внутри кавычек и добавляем содержимое
				content := strings.TrimSpace(string(runes[i+1 : nextQuote]))
				output = append(output, []rune(content)...)

				// Добавляем закрывающую кавычку
				output = append(output, '\'')

				// Проверяем следующий символ
				if nextQuote+1 < len(runes) {
					j := nextQuote + 1

					// Убираем пробелы до точки, если они есть
					for j < len(runes) && unicode.IsSpace(runes[j]) {
						j++
					}
					if j < len(runes) && runes[j] == '.' || runes[j] == ',' || runes[j] == ':' || runes[j] == ';' || runes[j] == '?' || runes[j] == '!' || runes[j] == '"' {
						i = j // Перемещаем индекс на точку
					} else {
						// В противном случае добавляем пробел после закрывающей кавычки
						output = append(output, ' ')
						i = nextQuote + 1
					}
				} else {
					// Если символов больше нет, добавляем пробел после закрывающей кавычки
					output = append(output, ' ')
					i = nextQuote + 1
				}
			} else {
				// Если кавычки не имеют пары, добавляем одну кавычку и продолжаем
				output = append(output, '\'')
				i++
			}
		} else {
			// Если текущий символ не кавычка, просто добавляем его в результат
			output = append(output, runes[i])
			i++
		}
	}

	// Возвращаем строку из результата
	return strings.TrimSpace(string(output))
}

func removeDoubleSpaces(input string) string {
	// Убираем двойные пробелы
	return regexp.MustCompile(`\s{2,}`).ReplaceAllString(input, " ")
}

// Обработка текста для API и CLI
func processText(content string) string {
	lines := strings.Split(content, "\n")
	if !isASCII(content) {
		return "Error: Text contains non-ASCII characters"
	}
	var processedLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			processedLines = append(processedLines, "")
			continue
		}
		str := " " + line
		str = addSpacesBeforeBrackets(str)
		iterations := 0
		maxIterations := 10
		processedBrackets := make(map[int]bool)
		for iterations < maxIterations {
			iterations++
			startArr := findStartBrackets(str)
			endArr := findEndBrackets(str)
			if len(startArr) == 0 || len(endArr) == 0 {
				break
			}
			startPos := -1
			endPos := -1
			for i := 0; i < len(endArr); i++ {
				if processedBrackets[endArr[i]] {
					continue
				}
				currentEnd := endArr[i]
				for j := len(startArr) - 1; j >= 0; j-- {
					if startArr[j] < currentEnd {
						startPos = startArr[j]
						endPos = currentEnd
						break
					}
				}
				if startPos != -1 {
					break
				}
			}
			if startPos == -1 || endPos == -1 {
				break
			}
			newStr := strings.TrimSpace(str[startPos+1 : endPos])
			if !isValidCommandFormat(newStr) {
				processedBrackets[startPos] = true
				processedBrackets[endPos] = true
				continue
			}
			processedBrackets[startPos] = true
			newStr2, num := getSubStrAndNum(newStr)
			switch strings.TrimSpace(newStr2) {
			case "hex":
				str = hexToDec(str)
			case "bin":
				str = binToDec(str)
			case "up", "cap", "low":
				beforeBracket := str[:startPos]
				afterBracket := str[endPos+1:]
				words := strings.Split(beforeBracket, " ")
				wordCount := 1
				if num > 0 {
					wordCount = num
				}
				lastModifiedIndex := -1
				for i := len(words) - 1; i >= 0 && wordCount > 0; i-- {
					if words[i] != "" {
						switch strings.TrimSpace(newStr2) {
						case "up":
							words[i] = strings.ToUpper(words[i])
						case "cap":
							if len(words[i]) > 0 {
								words[i] = strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
							}
						case "low":
							words[i] = strings.ToLower(words[i])
						}
						lastModifiedIndex = i
						wordCount--
					}
				}
				if lastModifiedIndex == -1 {
					lastModifiedIndex = 0
				}
				beforeWords := strings.Join(words[:lastModifiedIndex], " ")
				afterWords := strings.Join(words[lastModifiedIndex+1:], " ")
				if lastModifiedIndex > 0 && len(beforeWords) > 0 {
					beforeWords += " "
				}
				if lastModifiedIndex < len(words)-1 && len(afterWords) > 0 {
					afterWords = " " + afterWords
				}
				str = beforeWords + words[lastModifiedIndex] + afterWords + afterBracket
			}
		}
		str = transformAtoAn(str)
		str = formatPunctuation(str)
		str = adjustSpacesAroundDoubleQuotes(str)
		str = adjustSpacesAroundSingleQuotes(str)
		str = removeDoubleSpaces(str)
		processedLines = append(processedLines, str)
	}
	result := strings.Join(processedLines, "\n")
	return result
}

// HTTP handler для обработки текста
func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := processText(req.Text)
	resp := map[string]string{"result": result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
