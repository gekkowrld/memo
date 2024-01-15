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
	"strconv"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit your memo",
	Long: `Edit your memo easily`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
      editMemo, _ := strconv.Atoi(args[0])
      filename := matchMemoNumber(editMemo)
      openEditor(filename)
    }else {
      cmd.Help()
    }
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

