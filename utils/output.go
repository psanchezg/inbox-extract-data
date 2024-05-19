package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/psanchezg/inbox-extract-data/interfaces"
)

func IterateAndPrintBoltPlans(plans []interfaces.BoltPlan) {
	currentTime := time.Now()

	for _, plan := range plans {
		fmt.Println("========================================================")
		if !plan.Inicio.IsZero() {
			fmt.Printf("Plan de %v a %v\n", plan.Inicio.Format("02/01/2006 03:04"), plan.Fin.Format("02/01/2006 03:04"))
			fmt.Printf("Dias del bono: %v\n", plan.Duracion)
			fmt.Printf("Minutos totales del bono: %v\n", plan.Minutos)
			if currentTime.After(plan.Inicio) && currentTime.Before(plan.Fin) {
				fmt.Printf("******* Plan activo **********\n")
			}
			fmt.Println("||||||||||||||||||||||||||||||||||||||||||||||||")
			diff := time.Since(plan.Uso.PrimerViaje)
			diasUsados := int64(diff.Hours() / 24)
			restantes := fmt.Sprintf("sobre %v dias", plan.Duracion)
			if plan.Duracion-diasUsados >= 0 {
				fmt.Printf("Dias restantes del bono: %v\n", plan.Duracion-diasUsados)
				restantes = fmt.Sprintf("sobre %v dias", diasUsados)
			}
			minutos := math.Round(float64(plan.Uso.Tiempo) / 60.0)
			fmt.Printf("Tiempo total: %v minutos\n", minutos)
			fmt.Printf("Distancia total: %v kms\n", math.Round(plan.Uso.Distancia*100)/100)
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
		}
	}
}
