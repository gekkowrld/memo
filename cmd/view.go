/*
Copyright © 2024 Gekko Wrld
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View Your Memo",
	Long:  `View Your Memo`,
	Run: func(cmd *cobra.Command, args []string) {
		argsPassed := len(args)
		if argsPassed > 0 {
			numberView, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Couldn't get the filename, %v", err)
			}

			filename := matchMemoNumber(numberView)
			displayMemo(filename)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func displayMemo(filename string) {
	termSize := CalcTermSize()
	if termSize > 80 {
		termSize = termSize - 10
	}
	binCont, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Couldn't read the file, %v", err)
	}

	strCont := string(binCont)

	re, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(termSize),
	)

	disp, err := re.Render(strCont)
	fmt.Print(disp)
}
