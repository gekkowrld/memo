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
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a memo",
	Long:  `Delete a memo from the collection of your memos`,
	Run: func(cmd *cobra.Command, args []string) {
		argsPassed := len(args)
		if argsPassed > 0 {

			fistArgPassed, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("You passed a non int")
			}
			filename := matchMemoNumber(fistArgPassed)
			err = removeFile(filename, fistArgPassed)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func removeFile(filename string, memoNumber int) error {
	// Confirm that the value given exists in the first place
	ex := FileExists(filename)

	var e error
	if ex {
		e = os.Remove(filename)
		fmt.Printf("Deleted %s", filename)
	} else {
		fmt.Printf("[%d]: Can't delete memo, couldn't match any file\n", memoNumber)
		os.Exit(1)
	}

	return e
}

