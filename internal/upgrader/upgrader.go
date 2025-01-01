package upgrader

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ducng99/gohole/globals"
	"github.com/ducng99/gohole/internal/logger"
	"golang.org/x/mod/semver"
)

type GHReleaseResponse struct {
	TagName string `json:"tag_name"`
	Assets []GHReleaseAsset `json:"assets"`
}

type GHReleaseAsset struct {
	Name string `json:"name"`
	Url string `json:"browser_download_url"`
}

func RunTemp() error {
	currentFilePath, err := os.Executable()
	if err != nil {
		logger.Printf(logger.LogError, "Cannot find path to current executable\n")
		return err
	}

	currentFilePath, err = filepath.EvalSymlinks(currentFilePath)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot find path to current executable\n")
		return err
	}

	extension := filepath.Ext(currentFilePath)

	tempFile, err := os.CreateTemp("", "gohole_*"+extension)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot create a temporary gohole for self upgrade\n")
		return err
	}
	defer tempFile.Close()

	currentFile, err := os.Open(currentFilePath)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot read current executable\n")
		return err
	}
	defer currentFile.Close()

	_, err = io.Copy(tempFile, currentFile)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot write to temporary executable file\n")
		return err
	}

	tempFile.Close()

	tempCommand := exec.Command(tempFile.Name(), "upgrade", "--file-path", currentFilePath)
	tempCommand.Stdout = os.Stdout
	tempCommand.Stderr = os.Stderr

	if err := tempCommand.Start(); err != nil {
		logger.Printf(logger.LogError, "Cannot start temporary executable\n")
		return err
	}

	if err := tempCommand.Process.Release(); err != nil {
		logger.Printf(logger.LogError, "Cannot start temporary executable\n")
		return err
	}

	return nil
}

func CheckAndUpgrade(originalFilePath string) error {
	logger.Printf(logger.LogNormal, "Current version is %s\n", globals.Version)

	resp, err := http.Get("https://api.github.com/repos/ducng99/gohole/releases/latest")
	if err != nil {
		logger.Printf(logger.LogError, "Cannot connect to GitHub.\n")
		return err
	}

	if resp.StatusCode >= 400 {
		logger.Printf(logger.LogError, "Cannot get latest release metadata\n")
		return nil
	}

	decoder := json.NewDecoder(resp.Body)

	var latestVersionData GHReleaseResponse

	if err := decoder.Decode(&latestVersionData); err != nil {
		logger.Printf(logger.LogError, "Cannot parse latest release metadata\n")
		return err
	}

	if semver.Compare(globals.Version, latestVersionData.TagName) != -1 {
		logger.Printf(logger.LogNormal, "This is already the latest version\n")
		return nil
	}

	logger.Printf(logger.LogNormal, "New version %s found\n", latestVersionData.TagName)

	var downloadUrl string

	for _, asset := range latestVersionData.Assets {
		if asset.Name == getNewFileName() {
			downloadUrl = asset.Url
			break
		}
	}

	if downloadUrl == "" {
		logger.Printf(logger.LogWarn, "Cannot find a download link for %s, the release might be broken\n", latestVersionData.TagName)
		return nil
	}

	downloadResponse, err := http.Get(downloadUrl)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot download new version\n")
		return err
	}

	if downloadResponse.StatusCode >= 400 {
		logger.Printf(logger.LogError, "Cannot download new version\n")
		return nil
	}

	if err := pkill(filepath.Base(originalFilePath)); err != nil {
		logger.Printf(logger.LogError, "Cannot stop existing gohole process\n")
		return err
	}

	originalFile, err := os.Create(originalFilePath)
	if err != nil {
		logger.Printf(logger.LogError, "Cannot open original file to upgrade\n")
		return err
	}
	defer originalFile.Close()

	if _, err := io.Copy(originalFile, downloadResponse.Body); err != nil {
		logger.Printf(logger.LogError, "Cannot write to original file to upgrade\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Successfully upgraded to %s\n", latestVersionData.TagName)
	logger.Printf(logger.LogNormal, "If you are using gohole DNS server, you have to manually start it again\n")

	return nil
}
