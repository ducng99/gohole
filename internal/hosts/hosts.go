package hosts

import (
	"bufio"
	"errors"
	"os"

	"github.com/ducng99/gohole/internal/logger"
)

const (
	StartLine = "# Managed by GoHole - Start"
	EndLine   = "# Managed by GoHole - End"
)

func AddDomainsToHosts(domains []string) error {
	path, err := getHostsFilePath()
	if err != nil {
		logger.Printf(logger.LogError, "Failed when getting path to hosts file\n")
		return err
	}

	// Find start line
	file, err := os.Open(path)
	if err != nil {
		logger.Printf(logger.LogError, "Could not open hosts file to read\n")
		return err
	}
	defer file.Close()

	foundStartLine := false
	foundEndLine := false
	textBeforeStartLine := ""
	textAfterEndLine := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == StartLine {
			foundStartLine = true
			break
		}

		textBeforeStartLine += line + lineEnding
	}
	if err := scanner.Err(); err != nil {
		logger.Printf(logger.LogError, "Failed when reading hosts file\n")
		return err
	}

	if !foundStartLine {
		logger.Printf(logger.LogNormal, "Hosts file wasn't modified by gohole before, appending to file\n")
	} else {
		for scanner.Scan() {
			line := scanner.Text()

			if foundEndLine {
				textAfterEndLine += line + lineEnding
			}

			if line == EndLine {
				foundEndLine = true
			}
		}

		if !foundEndLine {
			logger.Printf(logger.LogError, "Hosts file is malformed, cannot find \"%s\" in hosts file. Please fix the file manually, removing all gohole entries.\n", EndLine)
			return errors.New("cannot find end line in hosts file")
		}

		logger.Printf(logger.LogNormal, "Found previously modified entries, replacing with newer\n")
	}

	file.Close()

	return writeToHosts(domains, path, textBeforeStartLine, textAfterEndLine)
}

func writeToHosts(domains []string, filePath string, textBefore, textAfter string) error {
	// Write to a temp file first
	logger.Printf(logger.LogNormal, "Writing %d domains to temp file\n", len(domains))

	file, err := os.CreateTemp("", "gohole_")
	if err != nil {
		logger.Printf(logger.LogError, "Cannot create a temp file to write\n")
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	// Write start
	if _, err = file.WriteString(textBefore + StartLine + lineEnding); err != nil {
		logger.Printf(logger.LogError, "Cannot write to temp file\n")
		return err
	}

	lineDomainsCount := 0

	for _, domain := range domains {
		if lineDomainsCount == 0 {
			if _, err = file.WriteString("0.0.0.0 " + domain); err != nil {
		logger.Printf(logger.LogError, "Cannot write to temp file\n")
				return err
			}

			lineDomainsCount++
		} else if lineDomainsCount == 8 {
			if _, err = file.WriteString(" " + domain + lineEnding); err != nil {
		logger.Printf(logger.LogError, "Cannot write to temp file\n")
				return err
			}

			lineDomainsCount = 0
		} else {
			if _, err = file.WriteString(" " + domain); err != nil {
		logger.Printf(logger.LogError, "Cannot write to temp file\n")
				return err
			}

			lineDomainsCount++
		}
	}

	// Write end line and the rest of the file
	if _, err = file.WriteString(lineEnding + EndLine + lineEnding + textAfter); err != nil {
		logger.Printf(logger.LogError, "Cannot write to temp file\n")
		return err
	}

	file.Close()

	logger.Printf(logger.LogNormal, "Successfully written to temp file\n")
	logger.Printf(logger.LogNormal, "Replacing hosts file...\n")

	// Replace temp file with actual hosts file
	if err := os.Rename(file.Name(), filePath); err != nil {
		logger.Printf(logger.LogError, "Cannot write to hosts file\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Successfully written to hosts file\n")

	return nil
}
