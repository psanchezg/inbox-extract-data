package outputs

import (
	"os"
)

func FileOutput(lines []string, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := f.WriteString(string(line)); err != nil {
			f.Close()
			return err
		}
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
