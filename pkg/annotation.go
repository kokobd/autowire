package pkg

import (
	"errors"
	"strings"
	"unicode"
)

// Annotation ast
type Annotation struct {
	Name string
	Args map[string]string
}

// ParseAnnotation try parse string starts with '@' to Annotation
func ParseAnnotation(str string) (*Annotation, error) {
	if len(str) == 0 || str[0] != '@' {
		return nil, errors.New("annotation must begin with '@'")
	}
	if len(str) == 1 || !unicode.IsLetter(rune(str[1])) {
		return nil, nil
	}
	str = str[1:]
	firstNonLetterIndex := -1
	var firstNonLetterCharacter rune
	for i, ch := range str {
		if !unicode.IsLetter(ch) {
			firstNonLetterIndex = i
			firstNonLetterCharacter = ch
			break
		}
	}
	if firstNonLetterIndex == -1 {
		return &Annotation{
			Name: str,
			Args: nil,
		}, nil
	}
	if unicode.IsSpace(firstNonLetterCharacter) {
		return &Annotation{
			Name: str[0:firstNonLetterIndex],
			Args: nil,
		}, nil
	}
	if firstNonLetterCharacter == '(' {
		name := str[0:firstNonLetterIndex]
		str = str[firstNonLetterIndex:]
		numLeftParens := 0
		for i, ch := range str {
			if ch == '(' {
				numLeftParens++
			}
			if ch == ')' {
				numLeftParens--
			}

			if numLeftParens == 0 {
				str = str[1:i]
				break
			}
		}
		// Here, str is argument array without parens
		argStrList := strings.Split(str, ",")
		if len(argStrList) == 1 && strings.Index(argStrList[0], "=") == -1 {
			value := strings.TrimSpace(argStrList[0])
			var args map[string]string
			if len(value) != 0 {
				args = map[string]string{
					"value": value,
				}
			}
			return &Annotation{
				Name: name,
				Args: args,
			}, nil
		}
		args := make(map[string]string)
		for _, argStr := range argStrList {
			kv := strings.Split(argStr, "=")
			if len(kv) != 2 {
				return nil, errors.New("there must be exactly 1 '=' in named arguments")
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if _, duplicate := args[key]; duplicate {
				return nil, errors.New("named argument can not have duplicate key")
			}
			args[key] = value
		}
		return &Annotation{
			Name: name,
			Args: args,
		}, nil
	} else {
		return nil, errors.New("invalid annotation syntax")
	}
}
