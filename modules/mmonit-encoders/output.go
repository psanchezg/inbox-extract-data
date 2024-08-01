package mmonitencoders

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
	"github.com/psanchezg/inbox-extract-data/utils"
)

func GetAggregateStats[K interfaces.Any](datas []K) MmonitUsageStats {
	ret := MmonitUsageStats{
		Total:   0.0,
		Minutes: 0,
	}
	for _, data := range datas {
		ret.Total += float64(data[Props["total"].Key].(float64))
		ret.Minutes += int64(data[Props["duration"].Key].(float64))
	}
	return ret
}

func ExportData[K interfaces.Any](datas map[string][]K) ([]string, [][]interface{}, error) {
	// Human export
	var lines []string
	var values [][]interface{}
	for key := range datas {
		data := datas[key]
		// Export data
		var partialLines []string
		var partialValues [][]interface{}
		var err error
		if partialLines, partialValues, err = exportData(data); err != nil {
			fmt.Println(err)
		}
		lines = append(lines, partialLines...)
		if len(values) == 0 {
			values = append(values, partialValues...)
		} else {
			values = append(values, partialValues[1:]...)
		}
	}
	if len(values) == 0 {
		return lines, values, nil
	}
	headers := values[0]
	vals := values[1:]
	sort.Slice(vals, func(a, b int) bool {
		timea, err1 := utils.ParseDateWithFormat(vals[a][2].(string), "02/01/2006 15:04:05")
		timeb, err2 := utils.ParseDateWithFormat(vals[b][2].(string), "02/01/2006 15:04:05")
		if err1 != nil || err2 != nil {
			// Manejar errores de parsing (aquí simplemente se considera que la fecha inválida es menor)
			fmt.Println("Error parsing date:", err1, err2)
			return err1 != nil
		}
		return timea.Before(timeb)
	})
	return lines, utils.InsertAtBeginning(vals, headers), nil
}

func exportData[K interfaces.Any](datas []K) ([]string, [][]interface{}, error) {
	props := Props
	// Lines return
	ret := []string{}
	values := [][]any{}
	// Values header
	values = append(values, []any{
		props["client"].Column[0],
		props["channel"].Column[0],
		props["from"].Column[0],
		props["to"].Column[0],
		props["duration"].Column[0],
	})

	// currentTime := time.Now()
	for _, data := range datas {
		var aux any
		var ok bool
		// Channel
		if aux, ok = data[props["channel"].Key]; !ok {
			return ret, values, fmt.Errorf("error getting channel")
		}
		client := strings.Split(aux.(string), "_")[0]
		ret = append(ret, ("========================================================\n"))
		ret = append(ret, fmt.Sprintf(props["client"].Text[0], client, props["client"].Unit))
		ret = append(ret, fmt.Sprintf(props["channel"].Text[0], aux, props["channel"].Unit))
		// Values
		vals := []any{
			client,
			aux,
		}
		// From date
		if aux, ok = data[props["from"].Key]; !ok {
			return ret, values, fmt.Errorf("error getting start date")
		}
		from, err := utils.ParseDateWithFormat(aux.(string))
		if err != nil {
			return ret, values, err
		}
		if !from.IsZero() {
			_, offset := from.Zone()
			from = from.Add(time.Duration(offset) * time.Second)
			// To date
			var aux any
			var ok bool
			if aux, ok = data[props["to"].Key]; !ok {
				return ret, values, fmt.Errorf("error getting end date")
			}
			to, err := utils.ParseDateWithFormat(aux.(string))
			if err != nil {
				return ret, values, err
			}
			_, offset = to.Zone()
			to = to.Add(time.Duration(offset) * time.Second)
			ret = append(ret, fmt.Sprintf(props["from"].Text[0], from.Format("02/01/2006 15:04:05"), to.Format("02/01/2006 15:04:05")))
			vals = append(vals, from.Format("02/01/2006 15:04:05"))
			vals = append(vals, to.Format("02/01/2006 15:04:05"))
		}
		// Duration
		if aux, ok = data[props["duration"].Key]; !ok {
			return ret, values, fmt.Errorf("error getting duration")
		}
		ret = append(ret, fmt.Sprintf(props["duration"].Text[0], aux, props["duration"].Unit))
		vals = append(vals, aux)
		values = append(values, vals)
	}

	return ret, values, nil
}
