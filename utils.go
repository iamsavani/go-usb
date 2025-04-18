package gadget

import (
	"bytes"
	"os"
)

// Convert boolean to "1" or "0" string representation.
func boolToIntStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// GetUdcs returns a list of UDCs (USB Device Controllers).
func GetUdcs() []string {
	var udcs []string

	files, err := os.ReadDir(UdcPathGlob)
	if err != nil {
		return nil
	}

	for _, file := range files {
		udcs = append(udcs, file.Name())
	}

	return udcs
}

func WriteIfDifferent(path string, content []byte, perm os.FileMode) error {
	if _, err := os.Stat(path); err == nil {
		oldContent, err := os.ReadFile(path)
		if err == nil {
			if bytes.Equal(oldContent, content) {
				return nil
			}

			if len(oldContent) == len(content)+1 &&
				bytes.Equal(oldContent[:len(content)], content) &&
				oldContent[len(content)] == '\n' {
				return nil
			}
		}
	}

	return os.WriteFile(path, content, perm)
}
