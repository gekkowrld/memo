/*
Copyright © 2024 Gekko Wrld
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your environment",
	Long:  `Configure how memo works and what to use `,
	Run: func(cmd *cobra.Command, args []string) {
		editFlag := cmd.Flag("edit").Changed
		viewFlag := cmd.Flag("view").Changed
		defaultFlag := getKeyValue("EditConfig")

		if editFlag {
			editConfig()
		} else if viewFlag {
			viewConfig()
		} else {
			if defaultFlag.(bool) {
				editConfig()
			} else {
				viewConfig()
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("edit", "e", false, "Edit the config file")
	configCmd.PersistentFlags().BoolP("view", "v", false, "View the configuration file")
}

type Config struct {
	MemoDir      string `toml:"memodir"`
	Editor       string `toml:"editor"`
	ListFGColour string `toml:"listfgcolour"`
	ListBGColour string `toml:"listbgcolour"`
	DisplayWidth int    `toml:"displaywidth"`
	EditConfig   bool   `toml:"editconfig"`
	Git          bool   `toml:"git"`
  StaticFiles  string `toml:"staticfiles"`
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
	case "configDir":
		return filepath.Dir(config_location)
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

func editConfig()  {
  // Open the default editor instead of doing it myself
  configFilename := getKeyValue("config_location").(string)
  err := openEditor(configFilename)
  if err != nil {
    fmt.Println("Something went wrong while editing the config file")
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

func viewConfig() {
	cellSize := CalcTermSize()/2 - 5

	const (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")
	)
	re := lipgloss.NewRenderer(os.Stdout)
	var (
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		CellStyle = re.NewStyle().Padding(1, 2).Width(cellSize)
		OddRowStyle = CellStyle.Copy().Foreground(gray)
		EvenRowStyle = CellStyle.Copy().Foreground(lightGray)
	)
	memoDir := getKeyValue("MemoDir").(string)
	editor := getKeyValue("Editor").(string)
	listfg := getKeyValue("ListFGColour").(string)
	listbg := getKeyValue("ListBGColour").(string)
	editconf := strconv.FormatBool(getKeyValue("EditConfig").(bool))
	configLoc := getKeyValue("config_location").(string)
  staticFiles := getKeyValue("StaticFiles").(string)

	if listfg == "" {
		listfg = "NO Colour!"
	}
	if listbg == "" {
		listbg = "NO Colour!"
	}

	rows := [][]string{
		{"Memo Directory", memoDir},
		{"Config File Location", configLoc},
		{"Editor", editor},
		{"Foreground Colour", listfg},
		{"Background Colour", listbg},
		{"Config default to Edit", editconf},
    {"Static files directory", staticFiles},
	}

	di := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("Key", "Value").
		Rows(rows...)

	fmt.Println(di)
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
