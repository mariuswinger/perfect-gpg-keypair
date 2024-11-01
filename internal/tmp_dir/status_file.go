package tmpdir

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type StatusFile struct {
	Path string
}

func (status_file StatusFile) ReadCreatedKeyId() (string, error) {
	file, err := os.Open(status_file.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	re := regexp.MustCompile(`KEY_CREATED \w (\w*)`)
	var id string
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		if len(match) > 1 {
			id = match[1]
		}
	}
	if id == "" {
		return "", fmt.Errorf("not present in status file")
	}
	if len(id) != 40 {
		return "", fmt.Errorf("invalid key format")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return id, nil
}
