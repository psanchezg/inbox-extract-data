package bolt

import (
	"fmt"
	"math"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
	"github.com/psanchezg/inbox-extract-data/utils"
)

// func parseTimeWithTimezone(date string) (time.Time, error) {
// 	if strings.HasSuffix(date, "Z") {
// 		date = date[:len(date)-1] + "+00:00"
// 	}
// 	layout := "2006-01-02T15:04:05-07:00"
// 	tt, err := time.Parse(layout, date)
// 	if err != nil {
// 		return time.Time{}, err
// 	}
// 	loc := tt.Location()
// 	zone, offset := tt.Zone()
// 	// Offset will be 0 if timezone is not recognized (or UTC, but that's ok).
// 	// Read carefully https://pkg.go.dev/time#Parse
// 	// In this case we'll try to load location from zone name.
// 	// Timezones that are recognized: local, UTC, GMT, GMT-1, GMT-2, ..., GMT+1, GMT+2, ...
// 	if offset == 0 {
// 		// Make sure you have timezone database available in your system for
// 		// time.LoadLocation to work. Read https://pkg.go.dev/time#LoadLocation
// 		// about where Go looks for timezone database.
// 		// Perhaps the simplest solution is to `import _ "time/tzdata"`, but
// 		// note that it increases binary size by few hundred kilobytes.
// 		// See https://golang.org/doc/go1.15#time/tzdata
// 		loc, err = time.LoadLocation(zone)
// 		if err != nil {
// 			return time.Time{}, err // or `return tt, nil` if you more prefer
// 			// the original Go semantics of returning time with named zone
// 			// but zero offset when timezone is not recognized.
// 		}
// 	}
// 	return time.ParseInLocation(layout, date, loc)
// }

func ExportData[K interfaces.Any](datas []K) ([]string, [][]interface{}, error) {
	props := Props
	// Lines return
	ret := []string{}
	values := [][]any{}
	// Values header
	values = append(values, []any{
		props["from"].Column[0],
		props["to"].Column[0],
		props["purchased"].Column[0],
		props["duration"].Column[0],
		props["duration2"].Column[0],
		props["duration"].Column[1],
		props["usage_time"].Column[0],
		props["usage_distance"].Column[0],
		props["duration3"].Column[0],
		props["usage_service"].Column[1],
		props["usage_paid"].Column[1],
		props["usage_paid"].Column[0],
	})

	currentTime := time.Now()
	for _, data := range datas {
		// Start date
		var aux any
		var ok bool
		if aux, ok = data[props["from"].Key]; !ok {
			return ret, values, fmt.Errorf("error getting start date")
		}
		from, err := utils.ParseDateWithFormat(aux.(string))
		if err != nil {
			return ret, values, err
		}
		ret = append(ret, ("========================================================\n"))
		if !from.IsZero() {
			// End date
			var aux any
			var ok bool
			if aux, ok = data[props["to"].Key]; !ok {
				return ret, values, fmt.Errorf("error getting end date")
			}
			to, err := utils.ParseDateWithFormat(aux.(string))
			if err != nil {
				return ret, values, err
			}
			ret = append(ret, fmt.Sprintf(props["from"].Text[0], from.Format("02/01/2006 15:04"), to.Format("02/01/2006 15:04")))
			aux2, ok := data[props["purchased"].Key]
			if !ok {
				return ret, values, fmt.Errorf("error getting purchased")
			}
			purchased := aux2.(bool)
			// Values
			vals := []any{
				from.Format("02/01/2006 15:04:05"),
				to.Format("02/01/2006 15:04:05"),
				purchased,
			}
			if purchased {
				ret = append(ret, fmt.Sprintf(props["duration"].Text[0], data[props["duration"].Key], props["duration"].Unit))
				ret = append(ret, fmt.Sprintf(props["duration2"].Text[0], data[props["duration2"].Key], props["duration2"].Unit))
				vals = append(vals, data[props["duration"].Key])
				vals = append(vals, data[props["duration2"].Key])
				if currentTime.After(from) && currentTime.Before(to) {
					ret = append(ret, fmt.Sprintf(props["to"].Text[0]))
				}
			} else {
				vals = append(vals, "")
				vals = append(vals, "")
			}
			ret = append(ret, "||||||||||||||||||||||||||||||||||||||||||||||||\n")

			usage := data[props["usage"].Key].(map[string]interface{})
			firstTravel, err := utils.ParseDateWithFormat(usage[props["usage_firsttravel"].Key].(string))
			if err != nil {
				return ret, values, err
			}
			diff := time.Since(firstTravel)
			diasUsados := int64(diff.Hours() / 24)
			diasBono := int64(data[props["duration"].Key].(float64))
			restantes := fmt.Sprintf("sobre %v dias", data[props["duration"].Key])
			if purchased {
				if diasBono-diasUsados >= 0 {
					ret = append(ret, fmt.Sprintf(props["duration"].Text[1], diasBono-diasUsados))
					restantes = fmt.Sprintf("sobre %v dias", diasUsados)
					vals = append(vals, diasBono-diasUsados)
				} else {
					// Caducado
					diff = to.Sub(firstTravel)
					diasUsados = int64((diff.Hours() - 24) / 24)
					vals = append(vals, 0)
				}
			} else {
				// No está activo, calcular al final del bono
				diff = to.Sub(firstTravel)
				diasUsados = int64((diff.Hours() - 24) / 24)
				vals = append(vals, "")
			}
			minutosUsados := math.Round(float64(usage[props["usage_time"].Key].(float64)) / 60.0)
			ret = append(ret, fmt.Sprintf(props["usage_time"].Text[0], minutosUsados, props["usage_time"].Unit))
			vals = append(vals, minutosUsados)
			distance := utils.Round(usage[props["usage_distance"].Key])
			ret = append(ret, fmt.Sprintf(props["usage_distance"].Text[0], distance, props["usage_distance"].Unit))
			vals = append(vals, distance)
			costeServicio := utils.Round(usage[props["usage_service"].Key])
			if purchased {
				// Computed
				minutosDia := utils.Round(data[props["duration3"].Key])
				tiempoAdicional := minutosUsados - (float64(diasUsados) * minutosDia)
				ret = append(ret, fmt.Sprintf(props["duration3"].Text[0], tiempoAdicional, props["duration3"].Unit))
				vals = append(vals, tiempoAdicional)
				ret = append(ret, fmt.Sprintf(props["usage_service"].Text[1], costeServicio, props["usage_service"].Unit))
				vals = append(vals, costeServicio)
				pagadoAdicional := utils.Round(usage[props["usage_paid"].Key])
				ret = append(ret, fmt.Sprintf(props["usage_paid"].Text[1], pagadoAdicional, props["usage_paid"].Unit))
				vals = append(vals, pagadoAdicional)
				// Computed
				importeCubiertoBono := utils.Round(costeServicio - pagadoAdicional)
				ret = append(ret, fmt.Sprintf("Total incluído en el bono: %v €\n", importeCubiertoBono))
				// Computed
				paid := utils.Round(usage[props["usage_paid"].Key])
				totalBono := utils.Round(data[props["total"].Key])
				totalConBono := utils.Round(paid + totalBono)
				ret = append(ret, fmt.Sprintf(props["total"].Text[0], totalConBono, props["usage_paid"].Unit))
				vals = append(vals, totalConBono)
				// Computed
				costeMinutoConBono := utils.Round(totalConBono / minutosUsados)
				ret = append(ret, fmt.Sprintf("Coste por minuto real (incluyendo bono): %v €\n", costeMinutoConBono))
				// Computed
				costeDia := utils.Round((paid + totalBono) / float64(diasUsados))
				ret = append(ret, fmt.Sprintf("Coste por día (incluyendo bono - %s): %v €\n", restantes, costeDia))
				// Computed
				costePorKm := utils.Round(totalConBono / distance)
				ret = append(ret, fmt.Sprintf("Coste por km (incluyendo bono): %v €\n", costePorKm))
			} else {
				costeServicio := math.Round((usage[props["usage_service"].Key].(float64))*100) / 100
				ret = append(ret, fmt.Sprintf(props["usage_service"].Text[0], costeServicio, props["usage_service"].Unit))
				vals = append(vals, "")
				vals = append(vals, "")
				vals = append(vals, "")
				vals = append(vals, costeServicio)
				// Computed
				costeMinuto := utils.Round(costeServicio / minutosUsados)
				ret = append(ret, fmt.Sprintf("Coste por minuto real: %v €\n", costeMinuto))
				// Computed
				costePorKm := utils.Round(costeServicio / distance)
				ret = append(ret, fmt.Sprintf("Coste por km: %v €\n", costePorKm))
			}
			values = append(values, vals)
		}
	}

	return ret, values, nil
}
