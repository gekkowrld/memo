package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
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


func TerminalSize(fd int) (int, int, error) {
	var dimensions [4]uint16

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0)
	if errno != 0 {
		return 0, 0, errno
	}

	return int(dimensions[1]), int(dimensions[0]), nil
}

func CalcTermSize() int {
		terminalWidth, _, err := TerminalSize(int(syscall.Stdin))
		if err != nil {
			// Default to 80 if unable to determine terminal width
			terminalWidth = 80
		}

	return terminalWidth
}
