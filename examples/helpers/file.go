package helpers

import (
	"fmt"
	"os"
)

// SaveToFile writes content to a file at the specified path.
func SaveToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing content to file: %v", err)
	}

	fmt.Println("File successfully written to", filePath)
	return nil
}
