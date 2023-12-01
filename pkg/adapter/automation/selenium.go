package adapter

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	entities "fv.io/webscrapping-investment/pkg/domain/entities"
	utils "fv.io/webscrapping-investment/pkg/domain/utils"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// Declarar una lista de Movimiento
var moves []entities.Movimiento
var movesCash []entities.Movimiento

// WebDriverManager es una clase para inicializar y gestionar el WebDriver de Selenium.
type WebDriverManager struct {
	driver selenium.WebDriver
}

// NewWebDriverManager crea una nueva instancia de WebDriverManager.
func NewWebDriverManager(port int, flagdebug bool) (*WebDriverManager, error) {
	// Iniciar el servicio ChromeDriver
	// Convert string to int

	_, err := selenium.NewChromeDriverService(
		os.Getenv("SELENIUM_DRIVER"), // Ruta al ejecutable de ChromeDriver
		port,                         // Puerto en el que escuchará ChromeDriver
		//selenium.Output(os.Stderr),   // Salida de ChromeDriver
	)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar el servicio ChromeDriver: %v", err)
	}

	// Configurar opciones del navegador (Chrome en este caso)
	caps := selenium.Capabilities{
		//"browserName": "chrome",
	}

	chromeCaps := chrome.Capabilities{
		Args: []string{
			//"--headless", // Ejecutar en modo headless (sin interfaz gráfica)
			"--headless-new", // comment out this line for testing
			//"--disable-gpu",           // Deshabilitar la aceleración por GPU en modo headless
			//"--no-sandbox",            // Deshabilitar el sandbox para Docker
			//"--disable-dev-shm-usage", // Deshabilitar el uso compartido de memoria en modo headless
		},
	}

	caps.AddChrome(chromeCaps)
	// Convert string to boolean

	selenium.SetDebug(flagdebug)
	// Iniciar el WebDriver
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return nil, fmt.Errorf("Error al iniciar el WebDriver:	", err)
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
		return nil, err
	}
	return &WebDriverManager{
		driver: driver,
	}, nil
}

// Close cierra el WebDriver al finalizar.
func (wm *WebDriverManager) Close() {
	if wm.driver != nil {
		wm.driver.Quit()
	}
}

func (s *WebDriverManager) RacionalGoToMovementPage(movementPage string) {

	err := s.driver.Get(movementPage)
	if err != nil {
		log.Fatal("Error:", err)
	}

	currentWindow, err := s.driver.CurrentWindowHandle()
	if err != nil {
		log.Fatal("Error obteniendo ventana actual:	", err, currentWindow)
	}

	time.Sleep(2 * time.Second)

	if _, err := s.driver.ExecuteScript(fmt.Sprintf("window.open(%q)", os.Getenv("RACIONAL_PAGE_MOVEMENTS")), nil); err != nil {
		log.Fatalf("opening a new window via Javascript returned error: %v", err)
	}

	// tiempo para abrir nueva ventana
	time.Sleep(4 * time.Second)

	// Obtener los identificadores de las ventanas
	windowHandles, err := s.driver.WindowHandles()
	if err != nil {
		fmt.Printf("Error al obtener los identificadores de las ventanas: %v\n", err)
		return
	}

	// Cerrar todas las ventanas excepto la primera
	for i, handle := range windowHandles {
		if i > 0 {
			err := s.driver.SwitchWindow(handle)
			if err != nil {
				fmt.Printf("Error al cambiar a la ventana %s: %v\n", handle, err)
				return
			}

			err = s.driver.Close()
			if err != nil {
				fmt.Printf("Error al cerrar la ventana %s: %v\n", handle, err)
				return
			}
		}
	}
	// Cambiar de nuevo a la primera ventana (opcional)
	err = s.driver.SwitchWindow(windowHandles[0])
	if err != nil {
		fmt.Printf("Error al cambiar a la primera ventana: %v\n", err)
		return
	}
}

func (s *WebDriverManager) RacionalGoToLoginPage(loginPage string) {
	err := s.driver.Get(loginPage)
	if err != nil {
		log.Fatal("Error:", err)
	}

	//Obtiene html
	html, err := s.driver.PageSource()
	if err != nil {
		log.Fatal("Error:", err, html)
	}

	// select the login form
	formElement, err := s.driver.FindElement(selenium.ByCSSSelector, "form")
	if err != nil {
		log.Fatal("Error al seleccionar form	:", err)
	}
	//fmt.Print(formElement)
	// fill in the login form fields
	user, err := formElement.FindElement(selenium.ByCSSSelector, "#main > app-login > ion-content > div > div.login-container > div.login-form > form > app-simple-credentials > form > ion-item > ion-input > input")
	if err != nil {
		log.Fatal("Error buscar user:", err, user)
	}
	user.SendKeys(os.Getenv("RACIONAL_USER"))
	password, err := formElement.FindElement(selenium.ByCSSSelector, "#main > app-login > ion-content > div > div.login-container > div.login-form > form > app-simple-credentials > form > app-password-input > form > ion-item > ion-input > input")
	if err != nil {
		log.Fatal("Error buscar password:", err, password)
	}
	password.SendKeys(os.Getenv("RACIONAL_PASSWORD"))

	//GET  login button
	loginButton, err := s.driver.FindElement(selenium.ByCSSSelector, "#main > app-login > ion-content > div > div.login-container > div.login-actions > button")
	if err != nil {
		log.Fatal("Error buscar boton login:", err, loginButton)
	}
	//Wait up to 5 seconds for the page to load
	s.driver.SetPageLoadTimeout(10 * time.Second)
	//Click login
	loginButton.Click()

	// wait up to 10 seconds for page charge
	err = s.driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		totalInversiones, _ := driver.FindElement(selenium.ByCSSSelector, "#main > app-tabs > app-racional-tabs > ion-tabs > div > ion-router-outlet > app-home > ion-content > div > div.desktop-home-center > div.investment-summary > div > div.investment-total-header > div.main-value > span.investment-amount")
		if totalInversiones != nil {
			return totalInversiones.IsDisplayed()
		}
		return false, nil
	}, 10*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func (s *WebDriverManager) RacionalMovementPageScrap(scroll bool) ([]entities.Movimiento, []entities.Movimiento) {

	if scroll {
		//Scrolleamos hasta el año configurado
		// Convertir la cadena a entero
		num, err := strconv.Atoi(os.Getenv("RACIONAL_YEAR_TO"))
		if err != nil {
			fmt.Println("Error al convertir la cadena a entero:", err)
		}

		lastYear := num - 1
		condition := fmt.Sprintf("%s %d", "Año", lastYear)
		//fmt.Println("condition year : ", condition)

		var count int
		var flag bool
		for {
			//body > app-root > ion-split-pane > app-side-menu > div > div.middle-column > app-filter-movements > ion-content > app-investment-movement-list > div:nth-child(2) > div
			count = count + 1
			movesNodes, err := s.driver.FindElements(selenium.ByCSSSelector, "body > app-root > ion-split-pane > app-side-menu > div > div.middle-column > app-filter-movements > ion-content > app-investment-movement-list > div:nth-child(2)  > div")
			if err != nil {
				log.Fatal("Error get element Ver todos", err)
			}
			//scrolleamos
			scriptMove := "arguments[0].scrollIntoView();"
			_, err = s.driver.ExecuteScript(scriptMove, []interface{}{movesNodes[count]})
			if err != nil {
				fmt.Println("Error al hacer scroll hasta el elemento:", err)

			}
			//esperamos que cargue
			time.Sleep(1 * time.Second)

			//verificamos año
			for _, movimientoElement := range movesNodes {
				yearElement, _ := movimientoElement.FindElement(selenium.ByCSSSelector, "div > p")
				yearText, _ := yearElement.Text()

				if yearText == condition {
					flag = true
					break
				}
			}
			if flag {
				break
			}

		}
	}

	//obtengo lista de elementos  movimientos
	/*
		movimientosNode, err := s.driver.FindElements(selenium.ByCSSSelector, "body > app-root > ion-split-pane > app-side-menu > div > div.middle-column > app-filter-movements > ion-content > app-investment-movement-list > div:nth-child(2) > div")
		if err != nil {
			log.Fatal("Error buscar movimientos:", err)
		}
		fmt.Println("Count Movimientos Node ", len(movimientosNode))
	*/
	htmlPage, _ := s.driver.PageSource()
	//utils.MakeFile("codigo.html", []byte(htmlPage))

	// Llamada a la función para extraer el nodo

	nodoList, err := extraerNodo(htmlPage, entities.AppInvestmentMovemenetList)
	if err != nil {
		fmt.Println("Error extraer nodo:", err)
	}
	nodo := encontrarXDiv(nodoList, 2)
	//imprimir html en file: utils.MakeFile("nodo.html", []byte(renderizarHTML(nodo)))

	// Llamar a la función para recorrer y obtener los valores
	recorrerNodosHTML(nodo)

	return moves, movesCash
}

func (s *WebDriverManager) Renta4GoToLoginPage(loginPage string) {
	err := s.driver.Get(loginPage)
	if err != nil {
		log.Fatal("Error:", err)
	}

	//Obtiene html
	html, err := s.driver.PageSource()
	if err != nil {
		log.Fatal("Error:", err, html)
	}

	// select the login form
	formElement, err := s.driver.FindElement(selenium.ByCSSSelector, "form")
	if err != nil {
		log.Fatal("Error al seleccionar form	:", err)
	}
	//fmt.Print(formElement)
	// fill in the login form fields
	user, err := formElement.FindElement(selenium.ByXPATH, "//*[@id='rut']")
	if err != nil {
		log.Fatal("Error buscar user:", err, user)
	}
	user.SendKeys(os.Getenv("RENTA4_USER"))
	password, err := formElement.FindElement(selenium.ByXPATH, "//*[@id='password']")
	if err != nil {
		log.Fatal("Error buscar password:", err, password)
	}
	password.SendKeys(os.Getenv("RENTA4_PASSWORD"))

	//GET  login button
	loginButton, err := s.driver.FindElement(selenium.ByXPATH, "//*[@id='formContent']/div[2]/div/form/input")
	if err != nil {
		log.Fatal("Error buscar boton login:", err, loginButton)
	}
	//Wait up to 10 seconds for the page to load
	s.driver.SetPageLoadTimeout(10 * time.Second)
	//Click login
	loginButton.Click()
}

func (s *WebDriverManager) Renta4GoToMovementPage(movementsPage string) {

	err := s.driver.Get(movementsPage)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Llenar fechas
	dateInit, err := s.driver.FindElement(selenium.ByCSSSelector, "#fechaIni")
	if err != nil {
		log.Fatal("Error buscar fecha inicio:", err, dateInit)
	}

	// Cambiar el valor del atributo usando JavaScript
	script := "arguments[0].setAttribute('value', '" + os.Getenv("RENTA4_FECHA_INICIO") + "');"
	_, err = s.driver.ExecuteScript(script, []interface{}{dateInit})
	if err != nil {
		fmt.Println("Error al ejecutar el script:", err)
		return
	}

	dateEnd, err := s.driver.FindElement(selenium.ByCSSSelector, "#fechaFin")
	if err != nil {
		log.Fatal("Error buscar password:", err, dateEnd)
	}
	// Cambiar el valor del atributo usando JavaScript
	fechaFin := os.Getenv("RENTA4_FECHA_FIN")

	if fechaFin == "" {
		// Obtener la fecha actual
		fechaActual := time.Now()
		// Formatear la fecha en el formato deseado "dd-mm-yyyy"
		formatoDeseado := "02-01-2006"
		fechaFin = fechaActual.Format(formatoDeseado)
	}
	scriptDateEnd := "arguments[0].setAttribute('value', '" + fechaFin + "');"
	_, err = s.driver.ExecuteScript(scriptDateEnd, []interface{}{dateEnd})
	if err != nil {
		fmt.Println("Error al ejecutar el script:", err)
		return
	}

	//GET  search button
	//Wait up to 30 - por la cantidad de movimientos
	s.driver.SetPageLoadTimeout(30 * time.Second)
	searchButton, err := s.driver.FindElement(selenium.ByCSSSelector, "#botonCont")
	if err != nil {
		log.Fatal("Error buscar button search:", err, searchButton)
	}
	//Click search button
	err = searchButton.Click()
	if err != nil {
		log.Fatal("Error click button search:", err)
	}

}

func (s *WebDriverManager) Renta4MovementPageScrap(flagTest bool) ([]entities.Movimiento, []entities.Movimiento) {

	var contador int
	for {
		contador = contador + 1
		fmt.Println("Page ", contador)
		//solo para pruebas
		if flagTest {
			if contador == 5 {
				break
			}
		}

		// Esperar a que aparezca un elemento específico
		if err := waitForElement(s.driver, contador); err != nil {
			log.Fatalf("Error al esperar a que aparezca el elemento después de hacer clic en el botón: %v", err)
		}

		//Obtenemos paginacion y movemos mouse
		paginations, err := s.driver.FindElements(selenium.ByCSSSelector, "body > div.nk-body > div.nk-wrap > div:nth-child(1) > div > div > div:nth-child(4) > div > ul > li")
		if err != nil {
			log.Printf("Error al buscar el elemento 'Siguiente': %v", err)
			break
		}

		var countPagination int
		for _, pagina := range paginations {
			countPagination = countPagination + 1

			dateElement, err := pagina.FindElement(selenium.ByCSSSelector, "a")
			if err != nil {
				log.Printf("Error obteniendo elemento pagina : a	': %v", err)
				break
			}
			paginaText, err := dateElement.Text()
			if err != nil {
				log.Printf("Error obteniendo texto de pagina : a': %v", err)
				break
			}
			//fmt.Println(paginaText)
			if paginaText == entities.SIGUIENTE1 {
				//fmt.Println("Posicion Boton Siguiente	:", countPagination)
				location, _ := pagina.Location()
				err = pagina.MoveTo(location.X, location.Y)
				if err != nil {
					log.Printf("Error al mover mouse : a 'Siguiente': %v", err)
					break
				}

			}
		}

		//Capturamos data
		movimientosNode, err := s.driver.FindElements(selenium.ByCSSSelector, "body > div > div.nk-wrap > div:nth-child(1) > div > div > div:nth-child(3) > div > div > div > div > table > tbody > tr")
		if err != nil {
			log.Fatal("Error buscar movimientos:", err)
		}

		for _, movimientoElement := range movimientosNode {
			var descriptionElement, amountElement, nemoElement, priceElement, quantityElement selenium.WebElement
			var description, amount, date, nemo, price, quantity string
			// select the name and price nodes
			dateElement, err := movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(1)")
			if err != nil {
				log.Fatal("Error obteniendo valor de elemento:", err)
			}
			date, _ = dateElement.Text()
			date = strings.ReplaceAll(date, "-", "/")
			typeElement, err := movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
			if err != nil {
				log.Fatal("Error obteniendo valor de elemento:", err)
			}
			typeText, _ := typeElement.Text()
			typeMovementEnum := utils.GetTypeMovement(typeText)

			switch typeMovementEnum {
			case entities.BUY:
				descriptionElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 2:", err)
				}
				amountElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(8)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 8:", err)
				}
				nemoElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(3)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 3:", err)
				}
				priceElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(5)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 5:", err)
				}
				quantityElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(4)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 4:", err)
				}
				description, _ = descriptionElement.Text()
				amount, _ = amountElement.Text()
				nemo, _ = nemoElement.Text()
				price, _ = priceElement.Text()
				quantity, _ = quantityElement.Text()
			case entities.SELL:
				descriptionElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 2:", err)
				}
				amountElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(7)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 7:", err)
				}
				nemoElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(3)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 3:", err)
				}
				priceElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(5)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 5:", err)
				}
				quantityElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(4)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 4:", err)
				}
				description, _ = descriptionElement.Text()
				amount, _ = amountElement.Text()
				nemo, _ = nemoElement.Text()
				price, _ = priceElement.Text()
				quantity, _ = quantityElement.Text()
			case entities.COMMISSION:
				descriptionElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 2:", err)
				}
				amountElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(7)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 7:", err)
				}
				description, _ = descriptionElement.Text()
				amount, _ = amountElement.Text()
			case entities.CHARGE:
				descriptionElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 2:", err)
				}
				amountElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(7)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 8:", err)
				}
				description, _ = descriptionElement.Text()
				amount, _ = amountElement.Text()
			case entities.PAYMENT:
				descriptionElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(2)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 2:", err)
				}
				amountElement, err = movimientoElement.FindElement(selenium.ByCSSSelector, "td:nth-child(6)")
				if err != nil {
					log.Fatal("Error obteniendo valor de elemento 8:", err)
				}
				description, _ = descriptionElement.Text()
				amount, _ = amountElement.Text()
			default:
				fmt.Println("Tipo  movimiento renta4 desconocido 	:", typeMovementEnum, "-", date, "-", typeText)
				continue
			}

			if typeMovementEnum == entities.PAYMENT || typeMovementEnum == entities.CHARGE || typeMovementEnum == entities.COMMISSION {
				// movimientos dinero
				moveCash := entities.Movimiento{}
				moveCash.Description = description
				moveCash.Amount = utils.StringToFloatWitoutDot(amount)
				moveCash.Date = date
				moveCash.Platform = entities.RENTA4
				moveCash.Type = typeMovementEnum
				moveCash.Currency = "CLP" // TODO pendiente us movement  en renta4 - no tengo

				movesCash = append(movesCash, moveCash)

			} else {
				// movimientos acciones
				move := entities.Movimiento{}
				move.Description = description + " " + nemo
				move.Amount = utils.StringToFloatWitoutDot(amount)
				move.Date = date
				move.Platform = entities.RENTA4
				move.Type = typeMovementEnum
				move.Currency = "CLP" // TODO pendiente us movement  en renta4 - no tengo
				move.Nemo = nemo
				move.Price = utils.StringToFloatWitoutDot(price)
				move.Quantity = utils.StringToFloatWitoutDot(quantity)
				moves = append(moves, move)
			}

		}

		//Hacemos Click en Siguiente Pagina
		err = paginations[countPagination-1].Click()
		if err != nil {
			log.Printf("Error al hacer click en pagina : a 'Siguiente': %v", err)
			break
		}
	}

	return moves, movesCash
}

// waitForPageLoad espera a que la página cargue completamente
func waitForElement(wd selenium.WebDriver, contador int) error {
	timeout := 20 * time.Second

	// Utilizar WebDriverWait para esperar a que el estado de la página sea "complete"
	wait := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		page, err := wd.FindElement(selenium.ByCSSSelector, "body > div > div.nk-wrap > div:nth-child(1) > div > div > div:nth-child(4) > div > ul > li.page-item.active")
		if err != nil {
			fmt.Println("Error encontrando next-page-active	:", err)
			return false, nil // Continuar esperando hasta que aparezca el elemento
		}
		pageNumber, _ := page.Text()
		// Convertir a cadena
		contadorNumerico := strconv.Itoa(contador)
		//	fmt.Println("Comparacion Page contadorvspagenumber ", contador, "-", pageNumber)
		if pageNumber != contadorNumerico {
			return false, nil
		}
		return true, nil
	}, timeout)

	return wait
}
