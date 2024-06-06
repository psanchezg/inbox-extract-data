package utils

import (
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
	"gitlab.com/hartsfield/inboxer"
	"golang.org/x/exp/constraints"
	"google.golang.org/api/gmail/v1"
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if b < a {
		return a
	}
	return b
}

func IterateAndPrintBoltPlans(plans []interfaces.BoltPlan) {
	currentTime := time.Now()

	for _, plan := range plans {
		fmt.Println("========================================================")
		if !plan.Inicio.IsZero() {
			if plan.Purchased {
				fmt.Printf("Plan de %v a %v\n", plan.Inicio.Format("02/01/2006 15:04"), plan.Fin.Format("02/01/2006 15:04"))
				fmt.Printf("Dias del bono: %v\n", plan.Duracion)
				fmt.Printf("Minutos totales del bono: %v\n", plan.Minutos)
				if currentTime.After(plan.Inicio) && currentTime.Before(plan.Fin) {
					fmt.Printf("******* Plan activo **********\n")
				}
			} else {
				fmt.Printf("Periodo de %v a %v\n", plan.Inicio.Format("02/01/2006 15:04"), plan.Fin.Format("02/01/2006 15:04"))
			}
			fmt.Println("||||||||||||||||||||||||||||||||||||||||||||||||")

			diff := time.Since(plan.Uso.PrimerViaje)
			diasUsados := int64(diff.Hours() / 24)
			restantes := fmt.Sprintf("sobre %v dias", plan.Duracion)
			if plan.Purchased && plan.Duracion-diasUsados >= 0 {
				fmt.Printf("Dias restantes del bono: %v\n", plan.Duracion-diasUsados)
				restantes = fmt.Sprintf("sobre %v dias", diasUsados)
			} else {
				// No está activo, calcular al final del bono
				diff = plan.Fin.Sub(plan.Uso.PrimerViaje)
				diasUsados = int64((diff.Hours() - 24) / 24)
			}
			minutos := math.Round(float64(plan.Uso.Tiempo) / 60.0)
			fmt.Printf("Tiempo total: %v minutos\n", minutos)
			fmt.Printf("Distancia total: %v kms\n", math.Round(plan.Uso.Distancia*100)/100)
			if plan.Purchased {
				fmt.Printf("Tiempo adicional usado (fuera bono): %v minutos\n", minutos-(float64(diasUsados)*float64(plan.MinutosDia)))
				fmt.Printf("Coste total del servicio (sin bono): %v €\n", math.Round(plan.Uso.Servicio*100)/100)
				fmt.Printf("Pagado adicional al bono: %v €\n", math.Round(plan.Uso.Pagado*100)/100)
				fmt.Printf("Total incluído en el bono: %v €\n", math.Round((plan.Uso.Servicio-plan.Uso.Pagado)*100)/100)
				fmt.Printf("Total pagado (incluyendo bono): %v €\n", math.Round((plan.Uso.Pagado+plan.Total)*100)/100)
				fmt.Printf("Coste por minuto real (incluyendo bono): %v €\n", math.Round((plan.Uso.Pagado+plan.Total)*100/minutos)/100)
				fmt.Printf("Coste por día (incluyendo bono - %s): %v €\n",
					restantes,
					math.Round((plan.Uso.Pagado+plan.Total)/float64(diasUsados)*100)/100,
				)
				fmt.Printf("Coste por km (incluyendo bono): %v €\n", math.Round((plan.Uso.Pagado+plan.Total)*100/plan.Uso.Distancia)/100)
			} else {
				fmt.Printf("Total pagado: %v €\n", math.Round((plan.Uso.Servicio)*100)/100)
				fmt.Printf("Coste por minuto real: %v €\n", math.Round(plan.Uso.Servicio*100/minutos)/100)
				fmt.Printf("Coste por km: %v €\n", math.Round(plan.Uso.Servicio*100/plan.Uso.Distancia)/100)
			}
		}
	}
}

func WriteFile(msg *gmail.Message) {
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
