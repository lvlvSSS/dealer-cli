package files

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func GetAllSubFiles(parent string) ([]string, error) {
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return nil, err
	}
	parent = parentAbs

	result := make([]string, 0, 16)
	stat, err := os.Stat(parent)
	if err != nil || !stat.IsDir() {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("extract file[%s] fails, because it is not directory ", parent))
	}
	files, err := os.ReadDir(parent)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			subParent := filepath.Join(parent, file.Name())
			subFiles, err := GetAllSubFiles(subParent)
			if err != nil {
				continue
			} else {
				result = append(result, subFiles...)
			}
		} else {
			result = append(result, filepath.Join(parent, file.Name()))
		}
	}
	return result, nil
}
