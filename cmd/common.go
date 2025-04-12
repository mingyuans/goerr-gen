package cmd

import (
	"os"
)

func writeToFile(src []byte, outputName string) error {
	return os.WriteFile(outputName, src, 0600)
}
