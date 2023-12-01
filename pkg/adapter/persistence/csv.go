package adapter

import (
	"encoding/csv"
	"fmt"
	"os"
)

// CSVCreator es una estructura que encapsula la lógica para crear archivos CSV.
type CSVCreator struct {
	FileName string
}

// NewCSVCreator crea una nueva instancia de CSVCreator.
func NewCSVCreator(fileName string) *CSVCreator {
	return &CSVCreator{FileName: fileName}
}

// CreateCSV crea un archivo CSV con datos proporcionados.
func (c *CSVCreator) CreateCSV(moves [][]string) error {

	// Crear el archivo CSV
	file, err := os.Create(c.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Crear un escritor CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir el encabezado
	header := []string{"Plataforma", "Tipo", "Descripcion", "Nemotecnico", "Moneda", "Monto", "Fecha", "Precio", "Cantidad"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Escribir los movimientos al archivo CSV
	for _, move := range moves {
		record := []string{
			move[3],
			move[4],
			move[0],
			move[6],
			move[5],
			move[1],
			move[2],
			move[7],
			move[8],
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	fmt.Println("Archivo CSV creado exitosamente:", c.FileName)

	return nil
}
