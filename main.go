package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
	"github.com/psanchezg/inbox-extract-data/utils"
	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	gmail "google.golang.org/api/gmail/v1"
)

var (
	// destinationDir is the path to the directory where the attachments will be saved.
	afterDate = os.Getenv("AFTER_DATE")
)

func writeFile(msg *gmail.Message) {
	time, err := inboxer.ReceivedTime(msg.InternalDate)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(fmt.Sprintf("./dump/%s-%s.txt", time.Format("2006-02-01"), msg.Id))
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
	if err != nil {
		fmt.Println(err)
		return

	}
	if _, err := f.WriteString(string(decoded)); err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func extractMails() {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailReadonlyScope)

	if afterDate == "" {
		afterDate = "2024/01/01"
	}
	msgs, err := inboxer.Query(srv, fmt.Sprintf("from:receipts@bolt.eu after:%s", afterDate))
	if err != nil {
		fmt.Println(err)
	}

	var firstMessage time.Time
	totalPagado := 0.0
	totalServicio := 0.0
	totalTiempo := 0
	totalDistancia := 0.0
	receipts := []interfaces.BoltReceipt{}
	plan := interfaces.BoltPlan{
		Total:      30.0,
		Minutos:    20 * 30,
		MinutosDia: 20,
		Duracion:   30,
	}
	rx := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .* (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: .*Descuento (?P<Descuento>[0-9\.\-]+)€)?`
	rx2 := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .*(?: (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)?)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: Importe total cobrado (?P<Cobrado>[0-9\.]+)€)?`

	sort.Slice(msgs, func(a, b int) bool {
		timea, err := inboxer.ReceivedTime(msgs[a].InternalDate)
		if err != nil {
			return true
		}
		timeb, err := inboxer.ReceivedTime(msgs[b].InternalDate)
		if err != nil {
			return false
		}
		return timea.Before(timeb)
	})

	// Range over the messages
	for _, msg := range msgs {
		decoded, errbody := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		params := utils.GetParams(rx, msg.Snippet)
		if params["Fecha"] == "" {
			params = utils.GetParams(rx2, msg.Snippet)
		}
		if params["Fecha"] != "" {
			receipt := utils.CreateBoltReceipt(params, msg.Snippet)
			if errbody == nil {
				if distancia, starttime, err := utils.ParseBodyTravel(string(decoded)); err == nil {
					receipt.Distancia = distancia
					if !receipt.Fecha.IsZero() {
						receipt.Fecha = receipt.Fecha.Add(starttime)
					}
				}
			}
			if !plan.Inicio.IsZero() && plan.Inicio.Before(receipt.Fecha) && plan.Fin.After(receipt.Fecha) {
				receipts = append(receipts, receipt)
				if firstMessage.IsZero() || firstMessage.After(receipt.Fecha) {
					firstMessage = receipt.Fecha
				}
				totalPagado += receipt.Total
				totalServicio += receipt.Subtotal
				// fmt.Println("tiempo", receipt.Duracion, totalTiempo)
				totalTiempo += int(receipt.Duracion)
				totalDistancia += receipt.Distancia
			} else {
				fmt.Println("Viaje fuera de fecha: ", receipt.Fecha.Format("02/01/2006"))
			}
			// fmt.Println("t", receipt.Subtotal, receipt.Total, receipt.Duracion)
		} else {
			// Test parse plan
			fmt.Println("Test parse plan")
			if errbody == nil {
				if detectplan, err := utils.ParseBodyPlan(string(decoded)); err == nil {
					plan = detectplan
					fmt.Printf("Plan encontrado el %v\n", plan.Inicio.Format("02/01/2006 03:04"))
				}
			}
		}
		// writeFile(msg)

		// for _, v := range msg.Payload.Body.Data {
		// 	fmt.Println(v)
		// }
		// body, err := inboxer.GetBody(msg, "text/plain")
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(body)
		// }
		// body, err = inboxer.GetBody(msg, "text/html")
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(body)
		// }

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
	fmt.Println(plan)
	fmt.Println("========================================================")
	if !plan.Inicio.IsZero() {
		fmt.Printf("Plan activo de %v a %v\n", plan.Inicio.Format("02/01/2006 03:04"), plan.Fin.Format("02/01/2006 03:04"))
		fmt.Printf("Dias del bono: %v\n", plan.Duracion)
		fmt.Printf("Minutos totales del bono: %v\n", plan.Minutos)
		fmt.Println("========================================================")
	}
	fmt.Printf("Primer viaje detectado: %v\n", firstMessage.Format("02-01-2006 03:04"))
	diff := time.Now().Sub(firstMessage)
	diasUsados := int64(diff.Hours() / 24)
	fmt.Printf("Dias restantes del bono: %v\n", plan.Duracion-diasUsados)
	fmt.Printf("Número de viajes realizados: %v\n", len(receipts))
	minutos := math.Round(float64(totalTiempo) / 60.0)
	fmt.Printf("Tiempo total: %v minutos\n", minutos)
	fmt.Printf("Distancia total: %v kms\n", math.Round(totalDistancia*100)/100)
	fmt.Printf("Tiempo adicional usado (fuera bono): %v minutos\n", minutos-(float64(diasUsados)*float64(plan.MinutosDia)))
	fmt.Printf("Coste total del servicio (sin bono): %v €\n", math.Round(totalServicio*100)/100)
	fmt.Printf("Pagado adicional al bono: %v €\n", math.Round(totalPagado*100)/100)
	fmt.Printf("Total incluído en el bono: %v €\n", math.Round((totalServicio-totalPagado)*100)/100)
	fmt.Printf("Total pagado (incluyendo bono): %v €\n", math.Round((totalPagado+plan.Total)*100)/100)
	fmt.Printf("Coste por minuto real (incluyendo bono): %v €\n", math.Round((totalPagado+plan.Total)*100/minutos)/100)
	fmt.Printf("Coste por día (incluyendo bono): %v €\n", math.Round((totalPagado+plan.Total)/float64(diasUsados)*100)/100)
	fmt.Println("========================================================")
}

func main() {
	extractMails()
}
