package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.New("error load env: " + err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return errors.New("error load env: invalid env line " + line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	err = scanner.Err()
	if err != nil {
		return errors.New("error load env: " + err.Error())
	}
	return nil
}
