package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyData(t *testing.T) {
	original := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	copied := copyData(original)

	if len(copied) != len(original) {
		t.Errorf("Expected copied map length %d, got %d", len(original), len(copied))
	}

	for key, value := range original {
		if copied[key] != value {
			t.Errorf("Expected value %v for key %s, got %v", value, key, copied[key])
		}
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{"$INT(1:10)", nil},
		{"$INT(10:1)", errors.New("invalid INT range")},
		{"$INT(a:b)", nil},
		{"$INT", nil},
	}

	for _, test := range tests {
		_, err := getInt(test.input)
		if err != nil {
			if err.Error() != test.expected.Error() {
				t.Errorf("getInt(%s) expected %v, got %v", test.input, test.expected, err)
			}
		}

	}
}

func TestRandomString(t *testing.T) {
	input := "$CHAR(10)"
	result, err := randomString(input)
	if err != nil {
		t.Errorf("randomString(%s) unexpected error: %v", input, err)
	}

	if len(result) != 10 {
		t.Errorf("Expected random string length 10, got %d", len(result))
	}
}

func TestParseAndUpdate(t *testing.T) {
	data := map[string]interface{}{
		"key1": "$UUID",
		"key2": "$INT(1:10)",
		"key3": "$CHAR(5)",
	}
	copiedData := copyData(data)
	traverseAndUpdate(data, copiedData)

	if data["key1"] == "$UUID" {
		t.Error("Expected key1 to be updated with UUID, but it wasn't")
	}

	if data["key2"] == "$INT(1:10)" {
		t.Error("Expected key2 to be updated with an integer, but it wasn't")
	}

	if data["key3"] == "$CHAR(5)" {
		t.Error("Expected key3 to be updated with a random string, but it wasn't")
	}
}

func TestWriteToFile(t *testing.T) {
	filename := "test_output.json"
	data := []byte(`{"key": "value"}`)

	err := writeToFile(filename, data)
	if err != nil {
		t.Errorf("writeToFile(%s) unexpected error: %v", filename, err)
	}

	readData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("ioutil.ReadFile(%s) unexpected error: %v", filename, err)
	}

	if string(readData) != string(data) {
		t.Errorf("Expected file content %s, got %s", data, readData)
	}

	os.Remove(filename)
}

func TestMainFunction(t *testing.T) {
	sourceFilename := "test_input.json"
	destinationFilename := "test_output.json"

	sourceData := `{
		"key1": "$UUID",
		"key2": "$INT(1:10)",
		"key3": "$CHAR(5)"
	}`

	err := ioutil.WriteFile(sourceFilename, []byte(sourceData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test input file: %v", err)
	}
	defer os.Remove(sourceFilename)

	// Redirect os.Args to simulate command-line arguments
	os.Args = []string{"cmd", "-s", sourceFilename, "-d", destinationFilename}

	main()

	outputData, err := ioutil.ReadFile(destinationFilename)
	if err != nil {
		t.Fatalf("Failed to read test output file: %v", err)
	}
	defer os.Remove(destinationFilename)

	var outputJSON map[string]interface{}
	if err := json.Unmarshal(outputData, &outputJSON); err != nil {
		t.Fatalf("Failed to unmarshal output JSON: %v", err)
	}

	if outputJSON["key1"] == "$UUID" || outputJSON["key2"] == "$INT(1:10)" || outputJSON["key3"] == "$CHAR(5)" {
		t.Error("Expected keys to be updated with new values, but they weren't")
	}
}
