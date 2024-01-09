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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"log"
	"errors"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the memos already created",
	Long: `List the memos that you have already crated in a list form`,
	Run: func(cmd *cobra.Command, args []string) {
		List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List(){
	// Should read from memoDir but assume it for now
	
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	memoDir := filepath.Join("/", userHomeDir, "/memo")

	dirExists := true

	if _, err := os.Stat(memoDir); errors.Is(err, os.ErrNotExist) {
		dirExists = false
	}

	if dirExists {

		files, err := os.ReadDir(memoDir)
		if err != nil {
			log.Fatalf("Error reading %s with error %v", memoDir, err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}
