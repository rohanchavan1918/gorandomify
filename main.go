package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed    Color = "\u001b[31m"
	ColorGreen  Color = "\u001b[32m"
	ColorYellow Color = "\u001b[33m"
	ColorBlue   Color = "\u001b[34m"
	ColorReset  Color = "\u001b[0m"
)

type Updater interface {
	Update(value string) (interface{}, error)
}

type UUIDUpdater struct{}
type IntUpdater struct{}
type CharUpdater struct{}

func (u UUIDUpdater) Update(value string) (interface{}, error) {
	return uuid.New().String(), nil
}

func (u IntUpdater) Update(value string) (interface{}, error) {
	return getInt(value)
}

func (u CharUpdater) Update(value string) (interface{}, error) {
	return randomString(value)
}

var updaters = map[string]Updater{
	"uuid": UUIDUpdater{},
	"int":  IntUpdater{},
	"char": CharUpdater{},
}

func main() {
	sourcePath := flag.String("t", "", "Source of template file")
	destinationPath := flag.String("o", "", "Destination path")
	flag.Parse()

	if *sourcePath == "" {
		colorize(ColorRed, "template file not passed")
		return
	}

	plan, err := ioutil.ReadFile(*sourcePath)
	if err != nil {
		colorize(ColorRed, err.Error())
		return
	}

	var originalData map[string]interface{}
	if err := json.Unmarshal(plan, &originalData); err != nil {
		colorize(ColorRed, err.Error())
		return
	}

	copiedData := copyData(originalData)
	traverseAndUpdate(originalData, copiedData)

	updatedJSON, err := json.MarshalIndent(copiedData, "", "  ")
	if err != nil {
		colorize(ColorRed, "Error: "+err.Error())
		return
	}

	if *destinationPath != "" {
		if err := writeToFile(*destinationPath, updatedJSON); err != nil {
			colorize(ColorRed, "Error: "+err.Error())
			return
		}

		colorize(ColorGreen, "JSON generated successfully: "+*destinationPath)
	} else {
		fmt.Println(string(updatedJSON))
	}

}

func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func copyData(original map[string]interface{}) map[string]interface{} {
	copied := make(map[string]interface{})
	for key, value := range original {
		copied[key] = value
	}
	return copied
}

func traverseAndUpdate(data, copiedData map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			traverseAndUpdate(v, copiedData[key].(map[string]interface{}))
		case string:
			parseAndUpdate(key, v, data, copiedData)
		}
	}
}

func writeToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func parseAndUpdate(key, value string, data, copiedData map[string]interface{}) {
	if updater := getUpdater(value); updater != nil {
		if newVal, err := updater.Update(value); err == nil {
			data[key] = newVal
			copiedData[key] = newVal
		} else {
			colorize(ColorRed, "Error updating key: "+key+", value: "+value+", err: "+err.Error())
		}
	}
}

func getUpdater(value string) Updater {
	switch {
	case strings.HasPrefix(value, "$UUID"):
		return updaters["uuid"]
	case strings.HasPrefix(value, "$INT"):
		return updaters["int"]
	case strings.HasPrefix(value, "$CHAR"):
		return updaters["char"]
	}
	return nil
}

func getInt(value string) (int, error) {
	matches := regexp.MustCompile(`^\$INT\((\d+):(\d+)\)$`).FindStringSubmatch(value)
	if len(matches) == 0 {
		return rand.Intn(10000), nil
	}

	lower, err1 := strconv.Atoi(matches[1])
	upper, err2 := strconv.Atoi(matches[2])
	if err1 != nil || err2 != nil || lower > upper {
		return 0, fmt.Errorf("invalid INT range")
	}

	return rand.Intn(upper-lower+1) + lower, nil
}

func randomString(value string) (string, error) {
	length := 10
	if parts := strings.Split(value, "("); len(parts) == 2 {
		if l, err := strconv.Atoi(strings.TrimSuffix(parts[1], ")")); err == nil {
			length = l
		}
	}
	return getRandomStrNlen(length), nil
}

func getRandomStrNlen(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
