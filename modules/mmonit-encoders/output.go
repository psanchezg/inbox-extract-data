package mmonitencoders

import (
	"fmt"
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

func ExportData[K interfaces.Any](datas []K) ([]string, [][]interface{}, error) {
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
