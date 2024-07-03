package mmonitencoders

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetOnlyValidData(t *testing.T) {
	rxbody := `Date:\s+(?P<Date>\d{1,2} (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{2}:\d{2}:\d{2} \+0000)\nHost:\s+(?P<Encoder>[a-z0-9\.]+)\nService:\s+(?P<Channel>[a-z0-9_]+)\sAction:\s+Alert\nDescription:\s+(?P<Action>(stop|start){1}) action done`
	body := `Date:        01 Jun 19:29:34 +0000
Host:        encoder10
Service:     anoticias_canal3_video_hls
Action:      Alert
Description: stop action done

Date:        01 Jun 19:29:34 +0000
Host:        encoder10
Service:     anoticias_canal3
Action:      Alert
Description: stop action done


Your faithful employee,
M/Monit`
	params := getParams(rxbody, body)
	if len(params) != 1 {
		t.Fatal("must be capture only the writh message")
	}

}

func TestParseBodies(t *testing.T) {

	rxbody := `Date:\s+(?P<Date>\d{1,2} (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{2}:\d{2}:\d{2} \+0000)\nHost:\s+(?P<Encoder>[a-z0-9\.]+)\nService:\s+(?P<Channel>[a-z0-9_]+)\sAction:\s+Alert\nDescription:\s+(?P<Action>(stop|start){1}) action done`

	file, err := os.Open("../../test/mmonit-bodies.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|||")
		parts[0] = strings.ReplaceAll(parts[0], "\\n", "\n")
		dates := strings.Split(parts[1], ",")
		params := getParams(rxbody, parts[0])
		// t.Fatal(parts[0])
		for i, v := range params {
			if v["Date"] != dates[i] {
				t.Fatalf(`Line 1 failed. Fecha mustbe "%v" not "%v"`, dates[i], v["Date"])
			}
		}
		// Tests
		// if params["Fecha"] == "" {
		// 	params = utils.GetParams(rx2, snippets[0])
		// }
		// // Tests
		// if params["Fecha"] != snippets[1] {
		// 	t.Fatalf(`Line 1 failed. Fecha mustbe "%v" not "%v"`, snippets[1], params["Fecha"])
		// }
		// if params["Total"] != snippets[2] {
		// 	t.Fatalf(`Line 1 failed. Total mustbe "%v" not "%v"`, snippets[2], params["Total"])
		// }
		// if params["Seg"] != snippets[4] {
		// 	t.Fatalf(`Line 1 failed. Seg mustbe "%v" not "%v"`, snippets[4], params["Seg"])
		// }
		// if params["Min"] != snippets[3] {
		// 	t.Fatalf(`Line 1 failed. Min mustbe "%v" not "%v"`, snippets[3], params["Min"])
		// }
		// if params["Subtotal"] != snippets[5] {
		// 	t.Fatalf(`Line 1 failed. Subtotal mustbe "%v" not "%v"`, snippets[5], params["Subtotal"])
		// }
		// if params["Cobrado"] != snippets[6] {
		// 	t.Fatalf(`Line 1 failed. Cobrado mustbe "%v" not "%v"`, snippets[6], params["Cobrado"])
		// }
	}
}
