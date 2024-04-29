package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

func TestParseLines(t *testing.T) {

	rx := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .* (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: .*Descuento (?P<Descuento>[0-9\.\-]+)€)?`
	rx2 := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .*(?: (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)?)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: Importe total cobrado (?P<Cobrado>[0-9\.]+)€)?`

	file, err := os.Open("./test/lines.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		snippets := strings.Split(scanner.Text(), "|||")
		params := getParams(rx, snippets[0])
		if params["Fecha"] == "" {
			params = getParams(rx2, snippets[0])
		}
		// Tests
		if params["Fecha"] != snippets[1] {
			t.Fatalf(`Line 1 failed. Fecha mustbe "%v" not "%v"`, params["Fecha"], snippets[1])
		}
		if params["Total"] != snippets[2] {
			t.Fatalf(`Line 1 failed. Total mustbe "%v" not "%v"`, params["Total"], snippets[2])
		}
		if params["Min"] != snippets[3] {
			t.Fatalf(`Line 1 failed. Min mustbe "%v" not "%v"`, params["Min"], snippets[3])
		}
		if params["Seg"] != snippets[4] {
			t.Fatalf(`Line 1 failed. Seg mustbe "%v" not "%v"`, params["Seg"], snippets[4])
		}
		if params["Subtotal"] != snippets[5] {
			t.Fatalf(`Line 1 failed. Subtotal mustbe "%v" not "%v"`, params["Subtotal"], snippets[5])
		}
		if params["Cobrado"] != snippets[6] {
			t.Fatalf(`Line 1 failed. Cobrado mustbe "%v" not "%v"`, params["Cobrado"], snippets[6])
		}
	}
}
