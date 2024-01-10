/*
Copyright © 2024 Gekko Wrld

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
	"path/filepath"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the memos already created",
	Long:  `List the memos that you have already crated in a list form`,
	Run: func(cmd *cobra.Command, args []string) {
		List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List() {

	memoDir, err := strconv.Unquote(strconv.Quote(getKeyValue("MemoDir").(string)))
	if err != nil {
		log.Fatalf("Error converting MemoDir to string: %v", err)
	}

	dirExists := DirectoryExists(memoDir) 

	if dirExists {

		files, err := os.ReadDir(memoDir)
		if err != nil {
			log.Fatalf("Error reading %s with error %v", memoDir, err)
		}

		// Sort the files before further processing
		sort.Slice(files, func(i, j int) bool {
			numI, _ := strconv.Atoi(strings.SplitN(files[i].Name(), "-", 2)[0])
			numJ, _ := strconv.Atoi(strings.SplitN(files[j].Name(), "-", 2)[0])
			return numI < numJ
		})
		var memoList string
		for _, file := range files {
			if !file.IsDir() {
				// Get the memo number
				numberStr := strings.SplitN(file.Name(), "-", 2)[0]
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					log.Printf("Error getting memo number from file %s with error code %v", file.Name(), err)
					continue
				}
				content, err := os.ReadFile(filepath.Join(memoDir, file.Name()))
				if err != nil {
			    log.Fatal("Error reading content from %s: %v", file.Name(), err)
			    continue
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
				memoInfo := fmt.Sprintf("Memo %d: %s", number, strings.TrimSpace(firstNonSpaceLine))
				memoList += "\n" + memoInfo
			}
		}

		terminalWidth, _, err := terminalSize(int(syscall.Stdin))
		if err != nil {
			// Default to 80 if unable to determine terminal width
			terminalWidth = 80
		}

		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(4).
			Width(terminalWidth)

		fmt.Println(style.Render(memoList))

	}
}

func terminalSize(fd int) (int, int, error) {
	var dimensions [4]uint16

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0)
	if errno != 0 {
		return 0, 0, errno
	}

	return int(dimensions[1]), int(dimensions[0]), nil
}
