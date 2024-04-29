package main

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	gmail "google.golang.org/api/gmail/v1"
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
	return paramsMap
}

func createBoltReceipt(params map[string]string, raw string) interfaces.BoltReceipt {
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

func extractMails() {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.MailGoogleComScope)

	msgs, err := inboxer.Query(srv, "from:receipts@bolt.eu after:2024/01/01")
	if err != nil {
		fmt.Println(err)
	}

	var firstMessage time.Time
	totalPagado := 0.0
	totalServicio := 0.0
	totalTiempo := 0
	receipts := []interfaces.BoltReceipt{}
	rx := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .* (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: .*Descuento (?P<Descuento>[0-9\.\-]+)€)?`
	rx2 := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .*(?: (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)?)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: Importe total cobrado (?P<Cobrado>[0-9\.]+)€)?`

	// Range over the messages
	for _, msg := range msgs {
		params := getParams(rx, msg.Snippet)
		if params["Fecha"] == "" {
			params = getParams(rx2, msg.Snippet)
		}
		if params["Fecha"] != "" {
			receipt := createBoltReceipt(params, msg.Snippet)
			if firstMessage.IsZero() || firstMessage.After(receipt.Fecha) {
				firstMessage = receipt.Fecha
			}
			totalPagado += receipt.Total
			totalServicio += receipt.Subtotal
			// fmt.Println("tiempo", receipt.Duracion, totalTiempo)
			totalTiempo += int(receipt.Duracion)
			receipts = append(receipts, receipt)
			// fmt.Println("t", receipt.Subtotal, receipt.Total, receipt.Duracion)
		}
		// fmt.Println("========================================================")
		// time, err := inboxer.ReceivedTime(msg.InternalDate)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("Date: ", time)
		// if firstMessage.IsZero() || firstMessage.After(time) {
		// 	firstMessage = time
		// }
		// md := inboxer.GetPartialMetadata(msg)
		// fmt.Println("From: ", md.From)
		// fmt.Println("Sender: ", md.Sender)
		// fmt.Println("Subject: ", md.Subject)
		// fmt.Println("Delivered To: ", md.DeliveredTo)
		// fmt.Println("To: ", md.To)
		// fmt.Println("CC: ", md.CC)
		// fmt.Println("Mailing List: ", md.MailingList)
		// fmt.Println("Thread-Topic: ", md.ThreadTopic)
		// fmt.Println("Snippet: ", msg.Snippet)
		// body, err := inboxer.GetBody(msg, "text/plain")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(body)
	}

	// fmt.Println("receipts", receipts)
	// Bono 30 días, 20 minutos al día = 30€
	costeBono := 30.0
	fmt.Println("========================================================")
	fmt.Printf("Primer viaje detectado: %v\n", firstMessage.Format("02-01-2006"))
	diff := time.Now().Sub(firstMessage)
	diasUsados := int(diff.Hours() / 24)
	fmt.Printf("Dias restantes del bono: %v\n", 30-diasUsados)
	fmt.Printf("Número de viajes realizados: %v\n", len(receipts))
	minutos := math.Round(float64(totalTiempo) / 60.0)
	fmt.Printf("Tiempo total: %v minutos\n", minutos)
	fmt.Printf("Tiempo adicionales al bono para días usado: %v minutos\n", minutos-(float64(diasUsados)*20))
	fmt.Printf("Coste total del servicio (sin bono): %v €\n", math.Round(totalServicio*100)/100)
	fmt.Printf("Total pagado (sin bono): %v €\n", totalPagado)
	fmt.Printf("Total pagado adicional al bono: %v €\n", math.Round((totalServicio-totalPagado)*100)/100)
	fmt.Printf("Total pagado (incluyendo bono): %v €\n", totalPagado+costeBono)
	fmt.Printf("Coste por minuto real (incluyendo bono): %v €\n", math.Round((totalPagado+costeBono)*100/minutos)/100)
	fmt.Printf("Coste por día (incluyendo bono): %v €\n", math.Round((totalPagado+costeBono)/float64(diasUsados)*100)/100)
	fmt.Println("========================================================")
}

func main() {
	extractMails()
}
