package cmd

import (
	"github.com/skidify/cli/src/utils"
)

// ShowConfigDirectory shows config directory in user's default file manager application
func ShowConfigDirectory() {
	configDir := utils.GetskidifyFolder()
	err := utils.ShowDirectory(configDir)
	if err != nil {
		utils.PrintError("Error opening config directory:")
		utils.Fatal(err)
	}
}
