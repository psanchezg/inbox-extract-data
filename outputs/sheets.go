package outputs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/psanchezg/inbox-extract-data/utils"
	"google.golang.org/api/sheets/v4"
)

func searchFile(fname string) string {
	// Crear servicio de Google Drive
	ctx := context.Background()
	srv := utils.ConnectToDriveService(ctx)
	// Listar archivos en Google Drive
	r, err := srv.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fileID := ""
	for _, file := range r.Items {
		if file.Title == fname {
			fileID = file.Id
			break
		}
	}
	return fileID
}

// calculateRange calcula el rango de escritura en función del tamaño del array de datos y el nombre de la hoja
func calculateRange(data [][]interface{}, sheetName string) string {
	if len(data) == 0 {
		return ""
	}

	numRows := len(data)
	numCols := len(data[0])

	// Convertir las columnas a la notación de letras
	endCol := utils.ColumnIndexToLetter(numCols)
	endRow := numRows

	// Crear la cadena de rango en formato A1
	rangeStr := fmt.Sprintf("%s!A1:%s%d", sheetName, endCol, endRow)
	return rangeStr
}

func SheetsOutput(values [][]interface{}, path string) {
	ctx := context.Background()
	srv := utils.ConnectToSheetsService(ctx)
	// Create spreadsheet
	aux := strings.Split(path, "|||")
	spreadsheetId := searchFile(utils.DefaultIfEmpty(aux[0], "inbox-extract-data"))
	var err error
	if spreadsheetId != "" {
		fmt.Printf("Found spreadsheet ID: %s\n", spreadsheetId)
	} else {
		spreadsheet := &sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: path,
			},
		}
		spreadsheet, err = srv.Spreadsheets.Create(spreadsheet).Do()
		if err != nil {
			log.Fatalf("Unable to create spreadsheet. %v", err)
		}

		// The ID of the spreadsheet
		spreadsheetId = spreadsheet.SpreadsheetId
		fmt.Printf("New spreadsheet ID: %s\n", spreadsheet.SpreadsheetId)
	}

	// The new values to apply to the spreadsheet
	// values := [][]interface{}{
	// 	// Row 1
	// 	{"sample_A1", "sample_B1"},
	// 	// Row 2
	// 	{"sample_A2", "sample_B2"},
	// 	// Row 3
	// 	{"sample_A3", "sample_A3"},
	// }

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	// Calcular el rango de escritura
	rangeStr := calculateRange(values, utils.DefaultIfEmpty(aux[1], "Hoja 1"))

	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeStr,
		Values: values,
	})
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// values := []*sheets.ValueRange{
	// 	{
	// 		Range:  rangeData,
	// 		Values: [][]interface{}{{"Header 1", "Header 2", "Header 3", "Header 4"}},
	// 	},
	// 	{
	// 		Range:  "Sheet1!A2:D2",
	// 		Values: [][]interface{}{{"Row 2 Col 1", "Row 2 Col 2", "Row 2 Col 3", "Row 2 Col 4"}},
	// 	},
	// }

	// for _, valueRange := range values {
	// 	_, err := srv.Spreadsheets.Values.Update(spreadsheetId, valueRange.Range, &sheets.ValueRange{
	// 		MajorDimension: "ROWS",
	// 		Values:         valueRange.Values,
	// 	}).ValueInputOption("RAW").Do()

	// 	if err != nil {
	// 		log.Fatalf("Unable to set data. %v", err)
	// 	}
	// }

}
