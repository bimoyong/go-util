package file

import "os"

// CheckOrMkdirAll function
func CheckOrMkdirAll(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return
		}
	}

	return
}
