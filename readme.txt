# Text Processing Tool in Go

## Overview

This project is a text completion, editing, and auto-correction tool written in Go. It reads a text file, applies a series of transformations based on specific commands in the text, and outputs the modified text to another file.

The tool can perform the following operations:
- Convert hexadecimal to decimal.
- Convert binary to decimal.
- Change the case of words (uppercase, lowercase, or capitalized).
- Correct punctuation formatting.
- Handle 'a' to 'an' conversions when necessary.

## Features

- **Hexadecimal to Decimal Conversion:** Replace `(hex)` with the decimal equivalent of the preceding hexadecimal number.
- **Binary to Decimal Conversion:** Replace `(bin)` with the decimal equivalent of the preceding binary number.
- **Text Case Transformations:** The tool can apply transformations like:
  - `(up)` - Convert to uppercase.
  - `(low)` - Convert to lowercase.
  - `(cap)` - Capitalize the first letter of each word.
  - `(up, <number>)`, `(low, <number>)`, `(cap, <number>)` - Apply the transformation to a specified number of words.
- **Punctuation Formatting:** Automatically correct spacing around punctuation marks such as `.,!,?,:;`.
- **A vs. An Correction:** Change "a" to "an" before vowels or "h".
- **Handling of Single and Double Quotes:** Remove unnecessary spaces around quotes.
- **Non-ASCII Characters Handling:** The program will check if the text contains any non-ASCII characters and handle them accordingly.

## Usage

The program expects two command-line arguments:
1. **Input file** containing the text to be processed.
2. **Output file** where the modified text will be saved.

### Example

```bash
$ cat sample.txt
it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.

$ go run . sample.txt result.txt

$ cat result.txt
It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.
```

Another example:

```bash
$ cat sample.txt
Simply add 42 (hex) and 10 (bin) and you will see the result is 68.

$ go run . sample.txt result.txt

$ cat result.txt
Simply add 66 and 2 and you will see the result is 68.
```

## Implementation

### Key Functions

- **`readFile(filename string) string`:** Reads the content of a file.
- **`writeFile(filename string, content string)`:** Writes the modified content to a file.
- **`isASCII(text string) bool`:** Checks if the text contains only ASCII characters.
- **`hexToDec(text string) string`:** Converts hexadecimal numbers to decimal.
- **`binToDec(text string) string`:** Converts binary numbers to decimal.
- **`transformAtoAn(text string) string`:** Adjusts the usage of "a" vs. "an" before vowels or "h".
- **`adjustSpacesAroundPunctuation(text string) string`:** Ensures correct spacing around punctuation marks.
- **`formatPunctuation(text string) string`:** Formats punctuation marks (`.,!?:;`) correctly with spaces.
- **`getSubStrAndNum(text string) (string, int)`:** Extracts the command and optional number from a command in the text.

### Example Code Flow

1. **Input Parsing:** The program takes the input file and reads its content.
2. **Text Transformation:** The text is processed line by line. For each line:
   - It checks for commands inside parentheses (e.g., `(hex)`, `(up)`).
   - Applies the relevant transformation based on the command.
   - Corrects punctuation and formatting.
3. **Output:** After processing, the modified text is written to the output file.

## Detailed Steps

1. **Hexadecimal and Binary Conversion:**
   - The program finds commands like `(hex)` or `(bin)` and converts the preceding hexadecimal or binary number into decimal.
   
2. **Text Case Transformation:**
   - The program handles `(up)`, `(low)`, and `(cap)` transformations. It can apply these transformations to a single word or a specified number of words.
   
3. **Punctuation Handling:**
   - The program ensures that punctuation marks such as `.,!?:;` are correctly spaced according to the rules provided.
   - It also handles situations where there are multiple punctuation marks (e.g., `...` or `!?`).

4. **'A' vs. 'An' Correction:**
   - The program automatically changes "a" to "an" when it appears before words starting with a vowel or "h".

## Tests

To test the program, you can create test files and run the following commands:

```bash
$ cat test_input.txt
There is no greater agony than bearing a untold story inside you.

$ go run . test_input.txt test_output.txt

$ cat test_output.txt
There is no greater agony than bearing an untold story inside you.
```

## Dependencies

- **Go Standard Library**: The program uses only the Go standard library, including packages like `fmt`, `os`, `regexp`, `strings`, and `strconv`.

## Conclusion

This project helps you practice string manipulation, file handling, and regular expressions in Go. It demonstrates how to build a text-processing tool that performs various transformations based on commands embedded in the text.