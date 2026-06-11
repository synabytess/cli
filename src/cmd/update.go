package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/skidify/cli/src/utils"
)

func Update(currentVersion string) bool {
	tagName, err := utils.FetchLatestTag()
	if err != nil {
		utils.PrintError("Cannot fetch latest release info")
		utils.PrintError(err.Error())
		return false
	}
	if currentVersion == tagName {
		utils.PrintSuccess("skidify is up-to-date.")
		return false
	}

	utils.PrintInfo("Latest release: " + tagName)
	var assetURL string = "https://github.com/skidify/cli/releases/download/v" + tagName + "/skidify-" + tagName + "-" + runtime.GOOS + "-"
	var location string = os.TempDir() + "/skidify-" + tagName

	if runtime.GOARCH == "386" && runtime.GOOS == "windows" {
		assetURL += "x32"
	} else if runtime.GOARCH == "arm64" {
		assetURL += "arm64"
	} else if runtime.GOOS == "windows" {
		assetURL += "x64"
	} else {
		assetURL += "amd64"
	}

	if runtime.GOOS == "windows" {
		assetURL += ".zip"
		location += ".zip"
	} else {
		assetURL += ".tar.gz"
		location += ".tar.gz"
	}

	spinner, _ := utils.Spinner.Start("Downloading skidify")

	out, err := os.Create(location)
	if err != nil {
		spinner.Fail("Failed to download skidify")
		utils.Fatal(err)
	}
	defer out.Close()

	resp2, err := http.Get(assetURL)
	if err != nil {
		spinner.Fail("Failed to download skidify")
		utils.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		spinner.Fail("Failed to download skidify")
		utils.Fatal(fmt.Errorf("unexpected HTTP status: %s for %s", resp2.Status, assetURL))
	}

	_, err = io.Copy(out, resp2.Body)
	if err != nil {
		spinner.Fail("Failed to download skidify")
		utils.Fatal(err)
	}
	spinner.Success("Downloaded skidify")

	exe, err := os.Executable()
	if err != nil {
		utils.Fatal(err)
	}
	if exe, err = filepath.EvalSymlinks(exe); err != nil {
		utils.Fatal(err)
	}

	exeOld := exe + ".old"
	utils.CheckExistAndDelete(exeOld)

	if err = os.Rename(exe, exeOld); err != nil {
		permissionError(err)
	}

	switch runtime.GOOS {
	case "windows":
		err = utils.Unzip(location, utils.GetExecutableDir())

	case "linux", "darwin":
		err = exec.Command("tar", "-xzf", location, "-C", utils.GetExecutableDir()).Run()
	}
	if err != nil {
		os.Rename(exeOld, exe)
		permissionError(err)
	}

	utils.CheckExistAndDelete(exeOld)
	utils.PrintSuccess("Successfully updated skidify to v" + tagName)
	return true
}

func permissionError(err error) {
	utils.PrintInfo("If fatal error is \"Permission denied\", please check read/write permission of skidify executable directory.")
	utils.PrintInfo("However, if you used a package manager to install skidify, please upgrade by using the same package manager.")
	utils.Fatal(err)
}
