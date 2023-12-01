package domain

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	entities "fv.io/webscrapping-investment/pkg/domain/entities"
)

// stringToBool converts a string to a boolean
func StringToBool(s string) (bool, error) {
	switch s {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid boolean string: %s", s)
	}
}

func GetPathAndFiles(s string) (string, []fs.DirEntry, error) {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	// lee archivos actuales
	files, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range files {
		fmt.Println(e.Name())
	}
	return dir, files, nil
}

func GetCurrency(montoString string) string {
	if strings.Contains(montoString, "US$") {
		return "USD"
	} else if strings.Contains(montoString, "$") {
		return "CLP"
	} else {
		return "Desconocido"
	}
}

func GetKeyAfter(input string, clave string) (string, error) {

	if input == "" {
		return "", nil
	}

	// Dividir el string usando el espacio como delimitador
	parts := strings.Split(input, clave)

	// Verificar si hay al menos dos partes
	if len(parts) >= 2 {
		// El segundo elemento del array es lo que sigue después de "movement-type"
		valor := parts[1]
		return valor, nil
	} else {
		fmt.Println("Error al hacer el split de ", input)
		return "", nil
	}

}

func GetNemoRacional(description string) string {

	descriptionUpper := strings.ToUpper(description)
	palabras := strings.Fields(descriptionUpper)
	count := len(palabras)
	switch count {
	case 1:
		primeraPalabra := palabras[0]
		return primeraPalabra
	case 2:
		segundaPalabra := palabras[1]
		return segundaPalabra
	case 3:
		terceraPalabra := palabras[2]
		return terceraPalabra
	default:
		fmt.Println("GetNemoRacional-Nemo-Desconocido	:", description)
		return "Desconocido"
	}

}

func GetTypeMovement(typeMovement string) string {
	// Eliminar los puntos de la cadena
	cadenaSinPuntos := strings.ReplaceAll(typeMovement, ".", "")
	cadenaSinPuntos = strings.ReplaceAll(cadenaSinPuntos, ",", ".")
	typeFormat := strings.TrimSpace(cadenaSinPuntos)
	typeFormat = strings.ToUpper(typeFormat)
	if strings.Contains(typeFormat, entities.Venta) || strings.Contains(typeFormat, entities.SELL) {
		return entities.SELL
	} else if strings.Contains(typeFormat, entities.Compra) || strings.Contains(typeFormat, entities.BUY) {
		return entities.BUY
	} else if strings.Contains(typeFormat, entities.Gastos) {
		return entities.COMMISSION
	} else if strings.Contains(typeFormat, entities.Abono) || strings.Contains(typeFormat, entities.Deposito) {
		return entities.PAYMENT
	} else if strings.Contains(typeFormat, entities.Cargo) || strings.Contains(typeFormat, entities.Retiro) {
		return entities.CHARGE
	} else if strings.Contains(typeFormat, entities.Dividendo) {
		return entities.DIVIDENDS
	} else if strings.Contains(typeFormat, entities.Intereses) {
		return entities.Intereses
	} else {
		fmt.Println("TypeMovement Desconocido :", typeMovement)
		return "Desconocido " + typeMovement
	}
}

func StringToIntWitoutDot(cadena string) (int, error) {

	// Eliminar los puntos de la cadena
	cadenaSinPuntos := strings.ReplaceAll(cadena, ".", "")

	// Convertir la cadena sin puntos a un número entero
	numeroEntero, err := strconv.Atoi(cadenaSinPuntos)
	if err != nil {
		fmt.Printf("Error al convertir la cadena a número entero: %v\n", err)
		return 0, err
	}
	return numeroEntero, nil
}

func StringToFloatWitoutDot(cadena string) float64 {

	// Eliminar los puntos de la cadena
	cadenaSinPuntos := strings.ReplaceAll(cadena, ".", "")
	cadenaSinPuntos = strings.ReplaceAll(cadenaSinPuntos, ",", ".")
	// Convertir string a float64
	valorFloat, err := strconv.ParseFloat(cadenaSinPuntos, 64)
	if err != nil {
		fmt.Println("Error al convertir el string a float64:", err)
		return 0
	}

	return valorFloat
}

func MakeFile(nombreArchivo string, datos []byte) error {
	// Abrir o crear el archivo en modo escritura
	archivo, err := os.Create("./output/" + nombreArchivo)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo: %v", err)
	}
	defer archivo.Close() // Asegúrate de cerrar el archivo al finalizar la función

	// Escribir datos en el archivo
	_, err = archivo.Write(datos)
	if err != nil {
		return fmt.Errorf("Error al escribir en el archivo: %v", err)
	}

	fmt.Printf("Se ha creado el archivo '%s' y se han escrito datos en él.\n", nombreArchivo)
	return nil
}
