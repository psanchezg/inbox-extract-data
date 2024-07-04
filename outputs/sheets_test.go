package outputs

import "testing"

func TestCreateSheet(t *testing.T) {
	ret := []string{"a", "b", "c"}
	SheetsOutput(ret, "")
}
