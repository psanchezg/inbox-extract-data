package mmonitencoders

import (
	"fmt"
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

func ExportDataAsStrings[K interfaces.Any](datas []K) ([]string, error) {
	ret := []string{}
	props := Props

	// currentTime := time.Now()
	for _, data := range datas {
		var aux any
		var ok bool
		// Channel
		if aux, ok = data[props["channel"].Key]; !ok {
			return ret, fmt.Errorf("error getting channel")
		}
		ret = append(ret, ("========================================================\n"))
		ret = append(ret, fmt.Sprintf(props["channel"].Text[0], aux, props["channel"].Unit))
		// From date
		if aux, ok = data[props["from"].Key]; !ok {
			return ret, fmt.Errorf("error getting start date")
		}
		from, err := utils.ParseDateWithFormat(aux.(string))
		if err != nil {
			return ret, err
		}
		if !from.IsZero() {
			_, offset := from.Zone()
			from = from.Add(time.Duration(offset) * time.Second)
			// To date
			var aux any
			var ok bool
			if aux, ok = data[props["to"].Key]; !ok {
				return ret, fmt.Errorf("error getting end date")
			}
			to, err := utils.ParseDateWithFormat(aux.(string))
			if err != nil {
				return ret, err
			}
			_, offset = to.Zone()
			to = to.Add(time.Duration(offset) * time.Second)
			ret = append(ret, fmt.Sprintf(props["from"].Text[0], from.Format("02/01/2006 15:04:05"), to.Format("02/01/2006 15:04:05")))
		}
		// Duration
		if aux, ok = data[props["duration"].Key]; !ok {
			return ret, fmt.Errorf("error getting duration")
		}
		ret = append(ret, fmt.Sprintf(props["duration"].Text[0], aux, props["duration"].Unit))
	}

	return ret, nil
}
