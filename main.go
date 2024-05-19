package main

import (
	"context"
	"encoding/base64"
	"fmt"
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

// func writeFile(msg *gmail.Message) {
// 	time, err := inboxer.ReceivedTime(msg.InternalDate)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	f, err := os.Create(fmt.Sprintf("./dump/%s-%s.txt", time.Format("2006-02-01"), msg.Id))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	decoded, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
// 	if err != nil {
// 		fmt.Println(err)
// 		return

// 	}
// 	if _, err := f.WriteString(string(decoded)); err != nil {
// 		fmt.Println(err)
// 		f.Close()
// 		return
// 	}
// 	err = f.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

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
	receipts := []interfaces.BoltReceipt{}
	otherreceipts := []interfaces.BoltReceipt{}
	planes := []interfaces.BoltPlan{}
	plan := interfaces.BoltPlan{
		Total:      0.0,
		Minutos:    0,
		MinutosDia: 0,
		Duracion:   0,
		Purchased:  false,
		Uso: interfaces.BoltUsePlan{
			Tiempo:    0,
			Distancia: 0,
			Pagado:    0,
			Servicio:  0,
		},
	}
	parsedAfterDate, err := time.Parse("2006/01/02", afterDate)
	if err == nil {
		plan.Inicio = parsedAfterDate
		plan.Fin = time.Now()
		diferencia := plan.Fin.Sub(plan.Inicio)
		plan.Duracion = int64(diferencia.Hours() / 24)
		planes = append(planes, plan)
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
			if !receipt.Fecha.IsZero() && (firstMessage.IsZero() || firstMessage.After(receipt.Fecha)) {
				firstMessage = receipt.Fecha
			}
			currentPlanIdx := utils.GetCurrentPlanIdxForDate(planes, receipt.Fecha)
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
				// fmt.Println("Viaje fuera de plan: ", receipt.Fecha.Format("02/01/2006"))
				otherreceipts = append(otherreceipts, receipt)
			}
			// fmt.Println("t", receipt.Subtotal, receipt.Total, receipt.Duracion)
		} else {
			// Test parse plan
			if errbody == nil {
				if detectplan, err := utils.ParseBodyPlan(string(decoded)); err == nil {
					plan = detectplan
					fmt.Printf("Plan encontrado el %v\n", plan.Inicio.Format("02/01/2006 15:04"))
					planes = append(planes, plan)
				}
			}
		}
		// writeFile(msg)
	}

	// fmt.Println("receipts", receipts)
	// Bono 30 días, 20 minutos al día = 30€
	fmt.Printf("Primer viaje detectado: %v\n", firstMessage.Format("02-01-2006 15:04"))
	fmt.Printf("Total número de viajes realizados: %v\n", len(receipts))
	utils.IterateAndPrintBoltPlans(planes)
	fmt.Println("========================================================")
}

func main() {
	extractMails()
}
