package usecase

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var aferoFS = afero.NewOsFs()

func writeToFile(dir, filename string, data []byte) error {
	filePath := filepath.Join(dir, filename)
	return afero.WriteFile(aferoFS, filePath, data, os.FileMode(0644))
}
