package cmd

import (
	"os"
	"strconv"
	"fmt"
)

// FileExists checks if a file exists.
func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists checks if a directory exists.
func DirectoryExists(dirName string) bool {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ConvertToString(value any) (string, error) {
    val, err := strconv.Unquote(strconv.Quote(value.(string)))
    if err != nil {
        return "", fmt.Errorf("error converting value to string: %v", err)
    }
    return val, nil
}

