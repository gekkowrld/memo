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
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your environment",
	Long:  `Configure how memo works and what to use `,
	Run: func(cmd *cobra.Command, args []string) {
		editFlag := cmd.Flag("edit").Changed

		if editFlag {
			editConfig()
		} else {
			editConfig()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("edit", "e", false, "Edit the config file")
}

type Config struct {
	MemoDir      string `toml:"memodir"`
	Editor       string `toml:"editor"`
	ListFGColour string `toml:"listfgcolour"`
	ListBGColour string `toml:"listbgcolour"`
	DisplayWidth int    `toml:"displaywidth"`
	// A specialkey "config_dir" is where this config file lives
	// it will be useless (redundant even) to add it in the file
}

func getKeyValue(key string) any {
	// Check if the environment variable is set
	config_file_env := os.Getenv("GMEMOCONF")

	var config_location string

	if config_file_env == "" {
		config_location = getDefaultConfigFile()
	} else {
		config_location = config_file_env
	}

	// For keys that are not found in the config file but required
	// 	in other locations
	switch key {
	case "configFile", "config_location", "configLocation":
		return config_location
	case "programLocation":
		return "$GOPATH/bin/memo"
	case "programName":
		return "memo"
	}

	// For all the keys that can be found in the config files
	// 	or a typo?

	var conf Config
	if _, err := toml.DecodeFile(config_location, &conf); err != nil {
		log.Fatal(err)
	}
	value := reflect.ValueOf(conf)
	field := value.FieldByName(key)

	// Check if the field is valid
	if field.IsValid() {
		return field.Interface()
	}

	return nil
}

func editConfig() {
	var conf Config
	var saveFile bool
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	config_location, _ := ConvertToString(getKeyValue("configFile"))
	if _, err := toml.DecodeFile(config_location, &conf); err != nil {
		log.Fatal(err)
	}
	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("memo").
			Description("You are now editing your config file")),

		huh.NewGroup(
			huh.NewInput().
				Title("MemoDir").
				Value(&conf.MemoDir),
			huh.NewInput().
				Title("Editor").
				Value(&conf.Editor),
		),
		huh.NewGroup(huh.NewConfirm().
			Title("Would you like to save your configs?").
			Value(&saveFile).
			Affirmative("Yes!").
			Negative("Nah!")),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if saveFile {
		saveFileSleep := func() {
			time.Sleep(2 * time.Second)
		}
		if err = saveConfigToFile(config_location, conf); err != nil {
			log.Fatal(err)
		}
		displayText := "Saving your config to " + config_location
		err = spinner.New().Title(displayText).Accessible(accessible).Action(saveFileSleep).Run()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Succesfully saved your file in ", config_location)
	} else {
		fmt.Println("Didn't save your config, please try again if you wish to change the config")
	}

}

func saveConfigToFile(filename string, conf Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode and write the config to the file
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(conf); err != nil {
		return err
	}

	return nil
}

func getDefaultMemoDir() string {

	// Get the user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Define the memo directory
	memoDir := filepath.Join(userHomeDir, "memo")

	// Check if the memoDir exists else create it
	if _, err := os.Stat(memoDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(memoDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return memoDir
}

func getDefaultConfigFile() string {

	// Get the user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Define the  config dir
	configDir := filepath.Join(userHomeDir, ".config/memo")

	// Check if the configDir exists else create it
	if !DirectoryExists(configDir) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Define the config File
	configFile := filepath.Join(configDir, "config.toml")
	// Check if the config file exists else create it
	file, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Close the file after trying whatever that was
	file.Close()

	return configFile

}
