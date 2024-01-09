package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"errors"
	"path/filepath"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new memo",
	Long:  `Add a new memo!`,
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
	openEditor(fileName, title)
}

func createFileName(title string) string {
	// Get the user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Define the memo directory
	memoDir := userHomeDir + "/memo"

	// Check if the memoDir exists else create it
	if _, err := os.Stat(memoDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(memoDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
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

	// Format the new file name
	newFileName := fmt.Sprintf("%d-%s-%s.md", nextNumber, formattedDate, strings.ReplaceAll(strings.ToLower(title), " ", "_"))

	return filepath.Join("/", memoDir, newFileName)
}

func openEditor(fileName string, title string) {
	// Assume that the editor in this case is vim
	editor := "vim"

	// Format the title string
	title = "# " + title

	// Write some content to the file before opening
	err := os.WriteFile(fileName, []byte(title), 0644)
	if err != nil {
		log.Fatal (err)
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
}

