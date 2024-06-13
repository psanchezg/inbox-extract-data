package bolt

import (
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/psanchezg/inbox-extract-data/utils"
	"google.golang.org/api/gmail/v1"
)

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	if paramsMap["Seg"] != "" && paramsMap["Min"] == "" {
		paramsMap["Min"] = "0"
	}
	return paramsMap
}

func createBoltReceipt(params map[string]string, raw string, msg *gmail.Message) BoltReceipt {
	receipt := BoltReceipt{
		//Fecha: params["Fecha"],
		Snippet: raw,
	}
	// Convert date
	l := "02/01/2006"
	tt, err := time.Parse(l, params["Fecha"])
	if err == nil {
		receipt.Fecha = tt
	} else {
		fmt.Println("error", params["Fecha"], err, raw)
	}
	// Convert values to float
	if total, err := strconv.ParseFloat(params["Total"], 64); err == nil {
		receipt.Total = total
	}
	if params["Cobrado"] != "" {
		if total, err := strconv.ParseFloat(params["Cobrado"], 64); err == nil {
			receipt.Total = total
		}
	}
	if subtotal, err := strconv.ParseFloat(params["Subtotal"], 64); err == nil {
		receipt.Subtotal = subtotal
	}
	if descuento, err := strconv.ParseFloat(params["Descuento"], 64); err == nil {
		receipt.Descuento = descuento
	}
	if desbloquear, err := strconv.ParseFloat(params["Desbloquear"], 64); err == nil {
		receipt.Desbloquear = desbloquear
	}
	// Convert time to seconds

	if params["Seg"] == "" && params["Min"] == "" {
		if params["Cobrado"] == "0.00" || params["Subtotal"] == "0.00" {
			// El máximo
			var segundos int32 = 20 * 60
			receipt.Duracion = int32(segundos)
		} else {
			fmt.Println("error", params, raw)
			utils.WriteFile(msg)
			receipt.Duracion = 0
		}
	} else {
		var segundos int32 = 0
		if min, err := strconv.ParseFloat(params["Min"], 64); err == nil {
			segundos = int32(math.Round(min * 60))
		}
		if seg, err := strconv.ParseFloat(params["Seg"], 64); err == nil {
			segundos = segundos + int32(seg)
		}
		receipt.Duracion = int32(segundos)
	}

	return receipt
}

func parseBodyTravel(body string) (float64, time.Duration, error) {
	rxkms := `\s*(?<Distance>[0-9\.]+) km\s*`
	var distance float64
	var starttime time.Duration
	rxtime := `\s*<span style="word-break: keep-all;">(?<Time>[0-9]+:[0-9]+)<\/span>`
	params2 := getParams(rxtime, body)
	if params2["Time"] != "" {
		duration := strings.Replace(params2["Time"], ":", "h", 1) + "m0s"
		dur, err := time.ParseDuration(duration)
		if err == nil {
			starttime = dur
		}
	}
	params := getParams(rxkms, body)
	if params["Distance"] == "" {
		return distance, starttime, fmt.Errorf("no kms found")
	}
	if kms, err := strconv.ParseFloat(params["Distance"], 64); err == nil {
		distance = kms
	}
	return distance, starttime, nil
}

func parseBodyPlan(body string) (BoltPlan, error) {

	plan := BoltPlan{
		Duracion:  0,
		Purchased: true,
		Uso: BoltUsePlan{
			Tiempo:    0,
			Distancia: 0,
			Pagado:    0,
			Servicio:  0,
		},
	}
	rxdetails := `.*<span class="bodyLarge color-contentPrimary".*>(?P<Detalles>.*)<\/span>`
	var compRegEx = regexp.MustCompile(rxdetails)
	matches := compRegEx.FindAllStringSubmatch(body, -1)
	rxdate := `(?P<Fecha>[0-9]{1,2} (January|February|March|April|May|June|July|Agoust|September|October|November|December) 20[0-9]{2}), (?P<Hora>[0-9]{2}:[0-9]{2}).*`
	rxplan := `(?P<Minutos>[0-9]{1,3}) min al día durante (?P<Duracion>[0-9]+ (mes|días))`
	rxtotal := `(?P<Total>[0-9]+\.[0-9]*)€`
	var inicio time.Time
	var fin time.Time
	if len(matches) == 0 {
		return BoltPlan{}, fmt.Errorf("no plan found")
	}
	for _, v := range matches {
		params := getParams(rxdate, v[1])
		if params["Fecha"] != "" {
			tt, err := time.Parse("02 January 2006", params["Fecha"])
			if err == nil {
				if inicio.IsZero() {
					inicio = tt
				} else {
					fin = tt
				}
			}
			if params["Hora"] != "" {
				duration := strings.Replace(params["Hora"], ":", "h", 1) + "m0s"
				dur, err := time.ParseDuration(duration)
				if err == nil {
					if fin.IsZero() {
						inicio = inicio.Add(dur)
					} else {
						fin = fin.Add(dur)
					}
				}
			}
		} else {
			params = getParams(rxtotal, v[1])
			if params["Total"] != "" {
				if total, err := strconv.ParseFloat(params["Total"], 64); err == nil {
					plan.Total = total
				}
			}
			params = getParams(rxplan, v[1])
			if params["Duracion"] == "1 mes" {
				plan.Duracion = 30
			}
			if params["Minutos"] != "" {
				// Convert values to float
				if minutos, err := strconv.ParseInt(params["Minutos"], 10, 16); err == nil {
					plan.MinutosDia = minutos
					if plan.Duracion > 0 {
						plan.Minutos = plan.Duracion * plan.MinutosDia
					}
				}
			}
		}
	}
	plan.Inicio = inicio
	plan.Fin = fin
	return plan, nil
}

func getCurrentPlanIdxForDate(plans []BoltPlan, searchDate time.Time) int {
	defaultPlanIdx := -1
	for i, p := range plans {
		if searchDate.After(p.Inicio) && searchDate.Before(p.Fin) {
			if p.Purchased {
				return i
			} else {
				defaultPlanIdx = i
			}
		}
	}
	if defaultPlanIdx > -1 {
		return defaultPlanIdx
	}
	return -1
}

func parseAndFormatDate(afterDate string, format string) string {
	parsedAfterDate, err := time.Parse("2006/01/02", afterDate)
	if err == nil {
		if format == "" {
			format = "02-01-2006"
		}
		return parsedAfterDate.Format(format)
	}
	return afterDate
}

func ProcessRawData(msgs []*gmail.Message) ([]BoltPlan, error) {
	var firstMessage time.Time
	var lastMessage time.Time
	receipts := []BoltReceipt{}
	// otherreceipts := []BoltReceipt{}
	planes := []BoltPlan{}
	plan := BoltPlan{
		Total:      0.0,
		Minutos:    0,
		MinutosDia: 0,
		Duracion:   0,
		Purchased:  false,
		Fin:        time.Now(),
		Uso: BoltUsePlan{
			Tiempo:    0,
			Distancia: 0,
			Pagado:    0,
			Servicio:  0,
		},
	}
	planes = append(planes, plan)
	rx := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .* (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: .*Descuento (?P<Descuento>[0-9\.\-]+)€)?`
	rx2 := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ Duración (?P<Duracion>[0-9\.]+)€(?: (?P<Min>[0-9]+) min)?(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: Importe total cobrado (?P<Cobrado>[0-9\.]+)€)?`

	// Range over the messages
	for _, msg := range msgs {
		decoded, errbody := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		params := getParams(rx, msg.Snippet)
		if params["Fecha"] == "" {
			params = getParams(rx2, msg.Snippet)
		}
		if params["Fecha"] != "" {
			receipt := createBoltReceipt(params, msg.Snippet, msg)
			if errbody == nil {
				if distancia, starttime, err := parseBodyTravel(string(decoded)); err == nil {
					receipt.Distancia = distancia
					if !receipt.Fecha.IsZero() {
						receipt.Fecha = receipt.Fecha.Add(starttime)
					}
				}
			}
			if !receipt.Fecha.IsZero() {
				if firstMessage.IsZero() || firstMessage.After(receipt.Fecha) {
					firstMessage = receipt.Fecha.Add(time.Duration(-30) * time.Minute)
					planes[0].Inicio = firstMessage
					diferencia := planes[0].Fin.Sub(planes[0].Inicio)
					planes[0].Duracion = int64(diferencia.Hours() / 24)
				}
				lastMessage = receipt.Fecha
			}
			currentPlanIdx := getCurrentPlanIdxForDate(planes, receipt.Fecha)
			if currentPlanIdx > -1 {
				receipts = append(receipts, receipt)
				planes[currentPlanIdx].Uso.Pagado += receipt.Total
				planes[currentPlanIdx].Uso.Servicio += receipt.Subtotal
				planes[currentPlanIdx].Uso.Tiempo += int64(receipt.Duracion)
				planes[currentPlanIdx].Uso.Distancia += receipt.Distancia
				first := planes[currentPlanIdx].Uso.PrimerViaje
				if !receipt.Fecha.IsZero() && (first.IsZero() || first.After(receipt.Fecha)) {
					planes[currentPlanIdx].Uso.PrimerViaje = receipt.Fecha
				}
			} else {
				fmt.Println("Viaje fuera de plan: ", receipt.Fecha.Format("02/01/2006"))
				// TODO
				// otherreceipts = append(otherreceipts, receipt)
			}
			// fmt.Println("t", receipt.Subtotal, receipt.Total, receipt.Duracion)
		} else {
			// Test parse plan
			if errbody == nil {
				if detectplan, err := parseBodyPlan(string(decoded)); err == nil {
					plan = detectplan
					fmt.Printf("Plan encontrado el %v\n", plan.Inicio.Format("02/01/2006 15:04"))
					planes = append(planes, plan)
				}
			}
		}
		// writeFile(msg)
	}
	planes[0].Fin = lastMessage
	fmt.Printf("Primer viaje detectado: %v\n", firstMessage.Format("02-01-2006 15:04"))
	fmt.Printf("Total número de viajes realizados: %v\n", len(receipts))
	return planes, nil
}
