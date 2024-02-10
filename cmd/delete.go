/*
Copyright © 2024 Gekko Wrld
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
