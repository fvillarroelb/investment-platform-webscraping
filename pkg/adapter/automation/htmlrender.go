package adapter

import (
	"fmt"
	"log"
	"strings"

	entities "fv.io/webscrapping-investment/pkg/domain/entities"
	utils "fv.io/webscrapping-investment/pkg/domain/utils"
	"golang.org/x/net/html"
)

var year string

func recorrerNodosHTML(nodo *html.Node) {

	if nodo != nil {
		// Verificar si el nodo es un párrafo con la clase "period-label"
		if nodo.Type == html.ElementNode && nodo.Data == "p" {
			for _, attr := range nodo.Attr {
				if attr.Key == "class" && attr.Val == "period-label" {
					// Imprimir el contenido del párrafo (Año)
					// Dividir el texto en palabras y obtener año
					palabras := strings.Fields(obtenerTextoNodo(nodo))
					year = strings.Join(palabras[1:], " ")
					//	fmt.Printf("Año: %s\n", year)
				}
			}
		}

		// Verificar si el nodo es un elemento app-investment-movement
		if nodo.Type == html.ElementNode && nodo.Data == "app-investment-movement" {
			// Llamar a la función para obtener los valores de los elementos hijos
			obtenerValoresHijos(nodo)
		}

		// Recorrer los nodos hijos de manera recursiva
		for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
			recorrerNodosHTML(hijo)
		}
	}
}

func obtenerValoresHijos(nodo *html.Node) {
	// Imprimir el contenido de los elementos hijos de app-investment-movement

	for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
		if hijo.Type == html.ElementNode {
			switch hijo.Data {
			case "ion-item":
				move := entities.Movimiento{}
				moveCash := entities.Movimiento{}
				// Aquí puedes procesar los valores dentro de ion-item
				//fmt.Println("Nuevo ion-item encontrado:")
				procesarIonItem(hijo, move, moveCash)

			}
		}
	}
}

func procesarIonItem(nodo *html.Node, move entities.Movimiento, moveCash entities.Movimiento) entities.Movimiento {
	// Imprimir el contenido del ion-item
	//fmt.Printf("Contenido de ion-item: %s\n", renderizarHTML(nodo))

	// Aquí puedes agregar lógica adicional para procesar otros elementos dentro de ion-item
	// Por ejemplo, puedes buscar elementos con las clases "title" y "movement-amount"
	// y extraer sus valores.
	class := "class"

	for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
		if hijo.Type == html.ElementNode {
			switch hijo.Data {
			case "div":
				nodeFind, _, _ := procesarDiv(hijo, class, "movement-amount", false, "", false)
				// Imprimir el contenido del div (valor)
				// Aquí puedes procesar los valores dentro de div
				_, _, attr := procesarDiv(hijo, class, "movement-type ", false, "", true)
				typeFormat, err := utils.GetKeyAfter(attr, " ")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento:", err)
				}
				typeMovementEnum := utils.GetTypeMovement(typeFormat)
				move.Type = typeMovementEnum
				amountText := obtenerTextoNodo(nodeFind)
				amountFormat, _ := utils.GetKeyAfter(amountText, "$")
				move.Currency = utils.GetCurrency(amountText)
				move.Amount = utils.StringToFloatWitoutDot(amountFormat)
				nodeDate, _, _ := procesarDiv(hijo, class, "date", false, "", false)
				move.Date = obtenerTextoNodo(nodeDate) + "/" + year
				move.Platform = entities.RACIONAL

				if typeMovementEnum == entities.PAYMENT || typeMovementEnum == entities.CHARGE || typeMovementEnum == entities.COMMISSION || typeMovementEnum == entities.DIVIDENDS || typeMovementEnum == entities.Intereses {
					nodeDescription, _, _ := procesarDiv(hijo, class, "title", false, "", false)
					move.Description = obtenerTextoNodo(nodeDescription)
					movesCash = append(movesCash, move)
				} else {
					nodeDescription, _, _ := procesarDiv(hijo, class, "title", false, "", false)
					move.Description = obtenerTextoNodo(nodeDescription)
					move.Nemo = utils.GetNemoRacional(obtenerTextoNodo(nodeDescription))
					move.Price = utils.StringToFloatWitoutDot("0")    //TODO se ve en el detalle
					move.Quantity = utils.StringToFloatWitoutDot("0") //TODO se ve en el detalle
					moves = append(moves, move)
				}

			}
		}

	}
	return move
}

func procesarDiv(nodo *html.Node, key string, value string, flag bool, atributo string, contiene bool) (*html.Node, bool, string) {
	// Verificar si el nodo tiene un atributo class="movement-amount"
	if flag {
		return nodo, flag, atributo
	}
	for _, attr := range nodo.Attr {
		if attr.Key == key && attr.Val == value {
			//	fmt.Println("nodoElementoEncontrado : ", obtenerTextoNodo(nodo), attr.Val)

			return nodo, true, attr.Val
		}
		if contiene {
			if attr.Key == key && strings.Contains(attr.Val, value) {
				//	fmt.Println("nodoElementoEncontrado : ", obtenerTextoNodo(nodo), attr.Val)

				return nodo, true, attr.Val
			}
		}

	}

	// Si no se encontró un class="movement-amount", continuar procesando los nodos hijos
	for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
		if hijo.Type == html.ElementNode {
			nodoRecursivo, flag, attr := procesarDiv(hijo, key, value, false, "", contiene)
			if flag {
				return nodoRecursivo, flag, attr
			}
			//fmt.Println("nodoRecursivo . ", obtenerTextoNodo(nodoRecursivo))

		}
	}
	return nodo, flag, ""
}

func renderizarHTML(nodo *html.Node) string {
	var resultado strings.Builder
	html.Render(&resultado, nodo)
	return resultado.String()
}

func obtenerTextoNodo(nodo *html.Node) string {
	// Recorrer los nodos hijos para encontrar el texto
	var texto strings.Builder
	for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
		if hijo.Type == html.TextNode {
			texto.WriteString(hijo.Data)
		}
	}
	return texto.String()
}
func extraerNodo(htmlString, etiqueta string) (*html.Node, error) {
	// Parsear el HTML en una estructura de nodos
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, err
	}

	// Función anónima para buscar el nodo con la etiqueta deseada
	var buscarNodo func(*html.Node) *html.Node
	buscarNodo = func(nodo *html.Node) *html.Node {
		if nodo.Type == html.ElementNode && nodo.Data == etiqueta {
			return nodo
		}
		for hijo := nodo.FirstChild; hijo != nil; hijo = hijo.NextSibling {
			if encontrado := buscarNodo(hijo); encontrado != nil {
				return encontrado
			}
		}
		return nil
	}

	// Llamar a la función para buscar el nodo
	nodoDeseado := buscarNodo(doc)
	if nodoDeseado == nil {
		return nil, fmt.Errorf("No se encontró un nodo con la etiqueta %s", etiqueta)
	}

	return nodoDeseado, nil
}

func encontrarXDiv(nodoPadre *html.Node, value int) *html.Node {
	// Inicializar el contador de divs encontrados
	divsEncontrados := 0

	// Iterar a través de los nodos hijos del nodo padre
	for hijo := nodoPadre.FirstChild; hijo != nil; hijo = hijo.NextSibling {
		// Verificar si el nodo actual es un div
		if hijo.Type == html.ElementNode && strings.ToLower(hijo.Data) == "div" {
			// Incrementar el contador de divs encontrados
			divsEncontrados++

			// Verificar si este es el segundo div
			if divsEncontrados == value {
				return hijo
			}
		}

		// Recursivamente buscar en los nodos hijos
		if nodoEncontrado := encontrarXDiv(hijo, value); nodoEncontrado != nil {
			return nodoEncontrado
		}
	}

	// No se encontró un segundo div
	return nil
}
