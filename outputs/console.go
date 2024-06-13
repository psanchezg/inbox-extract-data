package outputs

import (
	"fmt"
)

func ConsoleOutput(lines []string) {
	for _, line := range lines {
		fmt.Print(line)
	}
}
