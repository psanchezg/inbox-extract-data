package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
)

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func GetParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func CreateBoltReceipt(params map[string]string, raw string) interfaces.BoltReceipt {
	receipt := interfaces.BoltReceipt{
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
			fmt.Println("error", params)
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

func ParseBodyTravel(body string) (float64, time.Duration, error) {
	rxkms := `\s*(?<Distance>[0-9\.]+) km\s*`
	var distance float64
	var starttime time.Duration
	rxtime := `\s*<span style="word-break: keep-all;">(?<Time>[0-9]+:[0-9]+)<\/span>`
	params2 := GetParams(rxtime, body)
	if params2["Time"] != "" {
		duration := strings.Replace(params2["Time"], ":", "h", 1) + "m0s"
		dur, err := time.ParseDuration(duration)
		if err == nil {
			starttime = dur
		}
	}
	params := GetParams(rxkms, body)
	if params["Distance"] == "" {
		return distance, starttime, fmt.Errorf("no kms found")
	}
	if kms, err := strconv.ParseFloat(params["Distance"], 64); err == nil {
		distance = kms
	}
	return distance, starttime, nil
}

func ParseBodyPlan(body string) (interfaces.BoltPlan, error) {

	plan := interfaces.BoltPlan{
		Duracion: 0,
	}
	rxdetails := `.*<span class="bodyLarge color-contentPrimary".*>(?P<Detalles>.*)<\/span>`
	var compRegEx = regexp.MustCompile(rxdetails)
	matches := compRegEx.FindAllStringSubmatch(body, -1)
	rxdate := `(?P<Fecha>[0-9]{1,2} (January|February|March|April|May) 20[0-9]{2}), (?P<Hora>[0-9]{2}:[0-9]{2}).*`
	rxplan := `(?P<Minutos>[0-9]{1,3}) min al día durante (?P<Duracion>[0-9]+ (mes|días))`
	rxtotal := `(?P<Total>[0-9]+\.[0-9]*)€`
	var inicio time.Time
	var fin time.Time
	if len(matches) == 0 {
		return interfaces.BoltPlan{}, fmt.Errorf("no plan found")
	}
	for _, v := range matches {
		params := GetParams(rxdate, v[1])
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
			params = GetParams(rxtotal, v[1])
			if params["Total"] != "" {
				if total, err := strconv.ParseFloat(params["Total"], 64); err == nil {
					plan.Total = total
				}
			}
			params = GetParams(rxplan, v[1])
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
