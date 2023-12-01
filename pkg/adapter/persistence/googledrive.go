package adapter

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// GoogleDriveClient representa el cliente para interactuar con Google Sheets
type GoogleDriveClient struct {
	service *sheets.Service
}

// NewGoogleDriveClient crea un nuevo cliente para interactuar con Google Sheets
func NewGoogleDriveClient(credentialsPath string) (*GoogleDriveClient, error) {
	// Lee el archivo JSON de credenciales
	credBytes, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("No se pudo leer el archivo de credenciales: %v", err)
	}

	// Configura la estructura de credenciales
	config, err := google.JWTConfigFromJSON(credBytes, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("Error al configurar el cliente OAuth2: %v", err)
	}

	// Crea un cliente autenticado
	client := config.Client(context.Background())

	// Crea el servicio de Google Sheets
	sheetsService, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el servicio de Google Sheets: %v", err)
	}

	return &GoogleDriveClient{
		service: sheetsService,
	}, nil
}

// AddSheet agrega una nueva pestaña a una hoja de cálculo de Google Sheets
func (g *GoogleDriveClient) AddSheet(spreadsheetID, tabName string) error {
	_, err := g.service.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: tabName,
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return fmt.Errorf("Error al añadir la pestaña: %v", err)
	}

	return nil
}

// AddData añade datos a una pestaña específica de una hoja de cálculo
func (g *GoogleDriveClient) AddData(spreadsheetID, tabName string, values [][]interface{}) error {
	// Rango donde se insertarán los datos (por ejemplo, "Sheet1!A1")
	rangeToWrite := fmt.Sprintf("%s!A1", tabName)

	// Crear el objeto de solicitud para escribir los datos
	request := sheets.ValueRange{
		Values: values,
	}

	// Hacer la solicitud de escritura
	_, err := g.service.Spreadsheets.Values.Update(spreadsheetID, rangeToWrite, &request).
		ValueInputOption("RAW").Do()

	if err != nil {
		return fmt.Errorf("Error al añadir datos a la pestaña: %v", err)
	}

	return nil
}
