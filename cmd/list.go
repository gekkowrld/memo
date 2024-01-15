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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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

	nothingMessage := "You currently have no memo.\nRun `memo new` to get started or `memo help` to get help"
	var memoList string
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
		for _, file := range files {
			if !file.IsDir() {
				// Get the memo number
				numberStr := strings.SplitN(file.Name(), "-", 2)[0]
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					log.Printf("Error getting memo number from file %s with error code %v", file.Name(), err)
					continue
				}
        firstNonSpaceLine := getFileTitle(filepath.Join(memoDir, file.Name()))
				memoInfo := fmt.Sprintf("Memo %d: %s", number, strings.TrimSpace(firstNonSpaceLine))
				memoList += "\n" + memoInfo
			}
		}

		if memoList == "" {
			memoList = nothingMessage
		}
	} else {
		memoList = nothingMessage
	}

	terminalWidth := CalcTermSize()
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
