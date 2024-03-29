package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

func matchMemoNumber(memoNumber int) string {
	memoDir := getKeyValue("MemoDir").(string)

	isDirectoryThere := DirectoryExists(memoDir)

	if !isDirectoryThere {
		log.Fatalf("You can't delete any memo, you have none :)")
	}

	files, err := os.ReadDir(memoDir)
	if err != nil {
		log.Fatalf("Couldn't read the contents of %s", memoDir)
	}

	var matchedFile string

	for _, file := range files {
		matches := regexp.MustCompile(`^(\d+)-\d{4}-\d{2}-\d{2}-(.+)\.md$`).FindStringSubmatch(file.Name())

		if len(matches) >= 3 {
			currentMemoNumber, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			// Check if the memo number matches the provided memoNumber
			if currentMemoNumber == memoNumber {
				// Append to MemoDir
				matchedFile = filepath.Join(memoDir, file.Name())
			}
		}
	}

	return matchedFile
}

func openEditor(fileName string, oTitle ...string) error {
	editor, err := strconv.Unquote(strconv.Quote(getKeyValue("Editor").(string)))
	if err != nil {
		log.Fatalf("Error converting Editor to string: %v", err)
	}

	var title string
	if len(oTitle) == 0 {
		title = ""
	} else {
		title = oTitle[0]
	}
	// If title is provided, write it to the file
	if title != "" {
		title = "# " + title + "\n\n"
		err := os.WriteFile(fileName, []byte(title), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run the editor with the specified file
	cmd := exec.Command(editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// Display any errors that occur during execution
	err = cmd.Run()
	if err != nil {
		log.Fatalf("%s exited with error, couldn't open %s: %v", editor, fileName, err)
	}

	return err
}

func getFileTitle(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading content from", filename)
	}

	lines := strings.Split(string(content), "\n")

	var firstNonSpaceLine string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			firstNonSpaceLine = trimmedLine
			break
		}
	}
	if firstNonSpaceLine == "" {
		firstNonSpaceLine = "No title for this file"
	}

	firstNonSpaceLine = strings.TrimPrefix(firstNonSpaceLine, "#")
	return firstNonSpaceLine
}
