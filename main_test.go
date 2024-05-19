package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/psanchezg/inbox-extract-data/utils"
)

func TestParseBodyTravel(t *testing.T) {
	file, err := os.ReadFile("./test/viaje-viejo.txt")
	if err != nil {
		log.Fatal(err)
	}

	distancia, starttime, err := utils.ParseBodyTravel(string(file))
	if err != nil {
		log.Fatal(err)
	}
	if distancia != 3.82 {
		t.Fatalf(`Parse distancia failed. distancia mustbe "3.82" not "%v"`, distancia)
	}
	if start, err := time.ParseDuration("10h46m0s"); err == nil {
		if starttime != start {
			t.Fatalf(`Parse start time failed. starttime mustbe "10h46m0s" not "%v"`, starttime)
		}
	}

	file, err = os.ReadFile("./test/viaje-nuevo.txt")
	if err != nil {
		log.Fatal(err)
	}

	distancia, starttime, err = utils.ParseBodyTravel(string(file))
	if err != nil {
		log.Fatal(err)
	}
	if distancia != 2.44 {
		t.Fatalf(`Parse distancia failed. distancia mustbe "2.44" not "%v"`, distancia)
	}
	if start, err := time.ParseDuration("15h25m0s"); err == nil {
		if starttime != start {
			t.Fatalf(`Parse start time failed. starttime mustbe "15h25m0s" not "%v"`, starttime)
		}
	}
}

func TestParseBodyPlan(t *testing.T) {
	file, err := os.ReadFile("./test/contratar-plan-20mes.txt")
	if err != nil {
		log.Fatal(err)
	}

	plan, err := utils.ParseBodyPlan(string(file))
	if err != nil {
		log.Fatal(err)
	}
	inicio := plan.Inicio.Format("02/01/2006 15:04")
	if inicio != "13/04/2024 10:43" {
		t.Fatalf(`Inicio plan failed. Inicio mustbe "13/04/2024 10:43" not "%v"`, inicio)
	}
	fin := plan.Fin.Format("02/01/2006 15:04")
	if fin != "13/05/2024 10:43" {
		t.Fatalf(`Fin plan failed. Fin mustbe "13/05/2024 10:43" not "%v"`, fin)
	}
	if plan.MinutosDia != 20 {
		t.Fatalf(`MinutosDia plan failed. MinutosDia mustbe 20 not "%v"`, plan.MinutosDia)
	}
	if plan.Minutos != 20*30 {
		t.Fatalf(`Minutos plan failed. Minutos mustbe 600 not "%v"`, plan.Minutos)
	}
	if plan.Duracion != 30 {
		t.Fatalf(`Duracion plan failed. Duracion mustbe 30 not "%v"`, plan.Duracion)
	}
	if plan.Total != 30 {
		t.Fatalf(`Total plan failed. Duracion mustbe "30"€ not "%v"€`, plan.Total)
	}

	if plan, err := utils.ParseBodyPlan(""); err == nil && plan.Inicio.IsZero() {
		log.Fatal("Must return error when no plan body selected")
	}
}

func TestParseBodyPlan2(t *testing.T) {
	file, err := os.ReadFile("./test/contratar-plan-20mes-2.txt")
	if err != nil {
		log.Fatal(err)
	}

	plan, err := utils.ParseBodyPlan(string(file))
	if err != nil {
		log.Fatal(err)
	}
	inicio := plan.Inicio.Format("02/01/2006 15:04")
	if inicio != "15/05/2024 09:16" {
		t.Fatalf(`Inicio plan failed. Inicio mustbe "15/05/2024 09:16" not "%v"`, inicio)
	}
	fin := plan.Fin.Format("02/01/2006 15:04")
	if fin != "14/06/2024 09:16" {
		t.Fatalf(`Fin plan failed. Fin mustbe "14/06/2024 09:16" not "%v"`, fin)
	}
	if plan.MinutosDia != 20 {
		t.Fatalf(`MinutosDia plan failed. MinutosDia mustbe 20 not "%v"`, plan.MinutosDia)
	}
	if plan.Minutos != 20*30 {
		t.Fatalf(`Minutos plan failed. Minutos mustbe 600 not "%v"`, plan.Minutos)
	}
	if plan.Duracion != 30 {
		t.Fatalf(`Duracion plan failed. Duracion mustbe 30 not "%v"`, plan.Duracion)
	}
	if plan.Total != 30 {
		t.Fatalf(`Total plan failed. Duracion mustbe "30"€ not "%v"€`, plan.Total)
	}

	if plan, err := utils.ParseBodyPlan(""); err == nil && plan.Inicio.IsZero() {
		log.Fatal("Must return error when no plan body selected")
	}
}

func TestParseLines(t *testing.T) {

	rx := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .* (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: .*Descuento (?P<Descuento>[0-9\.\-]+)€)?`
	rx2 := `.*(?P<Fecha>\d{2}\/\d{2}\/\d{4}) .*Total (?P<Total>[0-9\.]+)€ .*Desbloquear (?P<Desbloquear>[0-9\.]+)€ .*(?: (?P<Min>[0-9]+) min(?: (?P<Seg>[0-9]+) s)?)? .*Subtotal (?P<Subtotal>[0-9\.]+)€(?: Importe total cobrado (?P<Cobrado>[0-9\.]+)€)?`

	file, err := os.Open("./test/lines.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		snippets := strings.Split(scanner.Text(), "|||")
		params := utils.GetParams(rx, snippets[0])
		if params["Fecha"] == "" {
			params = utils.GetParams(rx2, snippets[0])
		}
		// Tests
		if params["Fecha"] != snippets[1] {
			t.Fatalf(`Line 1 failed. Fecha mustbe "%v" not "%v"`, params["Fecha"], snippets[1])
		}
		if params["Total"] != snippets[2] {
			t.Fatalf(`Line 1 failed. Total mustbe "%v" not "%v"`, params["Total"], snippets[2])
		}
		if params["Min"] != snippets[3] {
			t.Fatalf(`Line 1 failed. Min mustbe "%v" not "%v"`, params["Min"], snippets[3])
		}
		if params["Seg"] != snippets[4] {
			t.Fatalf(`Line 1 failed. Seg mustbe "%v" not "%v"`, params["Seg"], snippets[4])
		}
		if params["Subtotal"] != snippets[5] {
			t.Fatalf(`Line 1 failed. Subtotal mustbe "%v" not "%v"`, params["Subtotal"], snippets[5])
		}
		if params["Cobrado"] != snippets[6] {
			t.Fatalf(`Line 1 failed. Cobrado mustbe "%v" not "%v"`, params["Cobrado"], snippets[6])
		}
	}
}

func TestParseInputDate(t *testing.T) {
	afterDate := "2024/04/01"
	formattedDate := utils.ParseAndFormatDate(afterDate)
	if formattedDate == afterDate {
		t.Fatalf(`Error en la conversión de fecha de "%v"`, afterDate)
	}
	if formattedDate != "01-04-2024" {
		t.Fatalf(`Error en la conversión de fecha de "%v" a "%v`, afterDate, formattedDate)
	}
}
