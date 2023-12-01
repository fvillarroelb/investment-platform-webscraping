package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type Movimiento struct {
	Description string
	Amount      float64 //monto total inversion
	Date        string
	Platform    string
	Type        string
	Currency    string
	Nemo        string
	Price       float64 // precio de compra o venta
	Quantity    float64
}

// NewMovimiento es una función que actúa como "constructor" para la estructura Movimiento.
func NewMovimiento(description string, amount float64, date string, platform string, mtype string, currency string, nemo string, price float64, quantity float64) *Movimiento {
	return &Movimiento{
		Description: description,
		Amount:      amount,
		Date:        date,
		Platform:    platform,
		Type:        mtype,
		Currency:    currency,
		Nemo:        nemo,
		Price:       price,
		Quantity:    quantity,
	}
}

// Implementa la interfaz Movimiento para la estructura Movimiento.
func (m Movimiento) quantity() float64 {
	return m.Quantity
}
func (m Movimiento) price() float64 {
	return m.Price
}
func (m Movimiento) nemotecnico() string {
	return m.Nemo
}
func (m Movimiento) platform() string {
	return m.Platform
}
func (m Movimiento) type_() string {
	return m.Type
}
func (m Movimiento) currency() string {
	return m.Currency
}
func (m Movimiento) description() string {
	return m.Description
}
func (m Movimiento) amount() float64 {
	return m.Amount
}
func (m Movimiento) date() string {
	return m.Date
}

// convertToSliceOfSlice convierte un slice de estructuras Movimiento a un slice de slices de strings
func ToInterface(movimientos []Movimiento) []interface{} {
	var movesInterface []interface{}
	for _, move := range movimientos {
		movesInterface = append(movesInterface, move)
	}

	return movesInterface
}

// convertToSliceOfSlice convierte un slice de estructuras Movimiento a un slice de slices de strings
func ConvertToSliceOfSlice(movimientos []Movimiento) [][]string {
	var result [][]string

	for _, movimiento := range movimientos {
		// Crear un slice de strings para cada instancia de Movimiento
		row := []string{
			movimiento.Description,
			strconv.FormatFloat(movimiento.Amount, 'f', -1, 64),
			movimiento.Date,
			movimiento.Platform,
			movimiento.Type,
			movimiento.Currency,
			movimiento.Nemo,
			strconv.FormatFloat(movimiento.Price, 'f', -1, 64),
			strconv.FormatFloat(movimiento.Quantity, 'f', -1, 64),
		}

		// Agregar el slice al resultado
		result = append(result, row)
	}

	return result
}

// convertToSliceOfSlice convierte un slice de estructuras Movimiento a un slice de slices de strings
func ConvertToInterface(moves []Movimiento) []interface{} {
	var movesInterface []interface{}

	for _, move := range moves {
		movesInterface = append(movesInterface, move)
	}

	return movesInterface
}

func ConvertToInterfaceMatriz(movimientos []Movimiento) [][]interface{} {
	resultado := make([][]interface{}, len(movimientos)+1) // +1 para incluir una fila de encabezados

	// Encabezados
	resultado[0] = []interface{}{"Plataforma", "Fecha", "TipoMovimiento", "Descripcion", "Nemotecnico", "Moneda", "Monto", "Precio", "Cantidad Acciones"}

	// Datos
	for i, movimiento := range movimientos {
		resultado[i+1] = []interface{}{
			movimiento.Platform,
			movimiento.Date,
			movimiento.Type,
			movimiento.Description,
			movimiento.Nemo,
			movimiento.Currency,
			movimiento.Amount,
			movimiento.Price,
			movimiento.Quantity,
		}

	}

	return resultado
}
func GroupAndSumInterface(moves [][]string) []Movimiento {

	// Mapa para almacenar la suma de amounts y prices por múltiples criterios
	sumaPorCriterios := make(map[string]map[string]float64)

	// Iterar sobre los movimientos y agrupar por día
	for i, movimiento := range moves {
		// Ignorar la primera fila (encabezados)
		//if i == 0 {
		//	continue
		//}
		// Obtener criterios clave para agrupar
		nemo := movimiento[6]       // El índice 6 corresponde a la columna de "Nemo"
		plataforma := movimiento[3] // El índice 3 corresponde a la columna de "Plataforma"
		fecha := movimiento[2]      // El índice 2 corresponde a la columna de "Date"
		tipo := movimiento[4]       // El índice 4 corresponde a la columna de "Type"
		currency := movimiento[5]

		// Crear clave de agrupación única
		clave := fmt.Sprintf("%s|%s|%s|%s|%s", nemo, plataforma, fecha, tipo, currency)

		// Inicializar el mapa interno para el grupo si aún no existe
		if _, ok := sumaPorCriterios[clave]; !ok {
			sumaPorCriterios[clave] = make(map[string]float64)
		}

		// Convertir el Amount y el Price a enteros

		amountFormateado := strings.ReplaceAll(movimiento[1], ".", "")
		amountFormateado = strings.ReplaceAll(movimiento[1], ",", ".")
		var amount float64
		var err error
		if amountFormateado == "" {
			amount = 0
		} else {
			amount, err = strconv.ParseFloat(amountFormateado, 64)
			if err != nil {
				fmt.Printf("Error al convertir Amount a entero en el movimiento %d: %v\n", i, err)
				continue
			}
		}

		// Reemplazar la coma por el punto y eliminar caracteres no numéricos
		priceFormateado := strings.ReplaceAll(movimiento[7], ".", "")
		priceFormateado = strings.ReplaceAll(movimiento[7], ",", ".")
		var price float64
		if priceFormateado == "" {
			price = 0
		} else {
			// Convertir a un número (float64)
			price, err = strconv.ParseFloat(priceFormateado, 64)
			if err != nil {
				fmt.Println("Error al convertir a float:", err)
				continue
			}
		}

		// Reemplazar la coma por el punto y eliminar caracteres no numéricos
		quantityFormateado := strings.ReplaceAll(movimiento[8], ".", "")
		quantityFormateado = strings.ReplaceAll(movimiento[8], ",", ".")
		var quantity float64
		if quantityFormateado == "" {
			quantity = 0
		} else {
			// Convertir a un número (float64)
			quantity, err = strconv.ParseFloat(quantityFormateado, 64)
			if err != nil {
				fmt.Println("Error al convertir a float quantity:", err)
				continue
			}
		}

		// Sumar el amount y el price al acumulador del grupo
		sumaPorCriterios[clave]["amount"] += amount
		sumaPorCriterios[clave]["price"] += price
		sumaPorCriterios[clave]["count"]++ // Contador para calcular el promedio
		sumaPorCriterios[clave]["quantity"] += quantity
	}
	// Crear un slice de tipo Movimiento para almacenar los resultados
	var groupMoves []Movimiento

	// Iterar sobre los resultados y agregarlos al slice
	for criterios, suma := range sumaPorCriterios {
		parts := strings.Split(criterios, "|")
		amountSum := suma["amount"]
		priceSum := suma["price"]
		count := suma["count"]
		quantity := suma["quantity"]

		// Calcular el precio promedio
		var avgPrice float64
		if count > 0 {
			avgPrice = float64(priceSum) / float64(count)
		}

		groupMoves = append(groupMoves, Movimiento{
			Description: "GroupedXDay",
			Amount:      amountSum,
			Date:        parts[2], // La fecha está en la cuarta posición en la clave
			Platform:    parts[1], // La plataforma está en la segunda posición en la clave
			Type:        parts[3], // El tipo está en la quinta posición en la clave
			Currency:    parts[4], // Puedes ajustar esto según tus necesidades
			Nemo:        parts[0], // El Nemo está en la primera posición en la clave
			Price:       avgPrice, // Precio promedio con dos decimales
			Quantity:    quantity, // Puedes ajustar esto según tus necesidades
		})
	}

	return groupMoves

}
