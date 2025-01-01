package hosts

import (
	"bufio"
	"io"
	"strings"
)

func ParseFromReader(reader io.Reader) ([]string, error) {
	domains := make([]string, 0, 300000)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineDomains := handleLine(scanner.Text())
		domains = append(domains, lineDomains...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

func handleLine(line string) []string {
	domains := make([]string, 0, 10)

	if strings.HasPrefix(line, "#") {
		// This line is a comment
		return domains
	}

	fields := strings.Fields(line)

	// Skip first field as it contains IP
	for i := 1; i < len(fields); i++ {
		if strings.HasPrefix(fields[i], "#") {
			break
		}

		domains = append(domains, fields[i])
	}

	return domains
}
