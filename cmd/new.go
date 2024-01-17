package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new memo",
	Long:  `Add something memorable to your collection of memos`,
	Run: func(cmd *cobra.Command, args []string) {
		title()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func title() {
	var title string
	huh.NewInput().
		Title("Memo Title: ").
		Value(&title).Run()

	// Pass the title to createFileName and OpenEditor functions
	fileName := createFileName(title)
	exitStatus := openEditor(fileName, title)
	if exitStatus == nil {
		commitMsg := fmt.Sprintf("[New]: %s", title)
		commit(commitMsg, fileName)
	}
}

func createFileName(title string) string {

	// Convert from an interface (or 'any') to string
	memoDir, err := strconv.Unquote(strconv.Quote(getKeyValue("MemoDir").(string)))
	if err != nil {
		log.Fatalf("Error converting MemoDir to string: %v", err)
	}

	isMemoDirPresent := DirectoryExists(memoDir)

	if !isMemoDirPresent {
		os.MkdirAll(memoDir, 0700)
	}
	// Read existing files in the memo directory
	files, err := os.ReadDir(memoDir)
	if err != nil {
		log.Fatalf("Error reading %s with error %v", memoDir, err)
	}

	// Extract and store numerical parts in a slice
	var numbers []int
	re := regexp.MustCompile(`^(\d+)-`)

	for _, file := range files {
		match := re.FindStringSubmatch(file.Name())
		if len(match) > 1 {
			num, _ := strconv.Atoi(match[1])
			numbers = append(numbers, num)
		}
	}

	// Sort numerical parts in ascending order
	sort.Ints(numbers)

	// Find the maximum number
	maxNumber := 0
	if len(numbers) > 0 {
		maxNumber = numbers[len(numbers)-1]
	}

	// Increment the maximum number for the next file
	nextNumber := maxNumber + 1

	// Format the current date
	formattedDate := time.Now().Format("2006-01-02")

	sanitizedTitle := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(title, "_")
	// Format the new file name
	newFileName := fmt.Sprintf("%d-%s-%s.md", nextNumber, formattedDate, strings.ToLower(sanitizedTitle))

	return filepath.Join("/", memoDir, newFileName)
}
