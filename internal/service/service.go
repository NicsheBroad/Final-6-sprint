package service

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func ConvertText(input string) (string, error) {

	var result string

	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("text is empty")
	}

	isMorse := true

	for _, char := range input {
		if char != '.' && char != '-' && char != ' ' {
			isMorse = false
			break
		}
	}

	if isMorse {
		result = morse.ToText(input)
		if result == "" {
			return "", fmt.Errorf("invalid morse code")
		}
	} else {
		input = strings.ToUpper(input)
		result = morse.ToMorse(input)
		if result == "" {
			return "", fmt.Errorf("invalid text code")
		}
	}

	return result, nil
}
