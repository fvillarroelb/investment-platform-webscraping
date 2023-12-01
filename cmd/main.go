package main

import (
	"fmt"
	"log"
	"os"

	automation "fv.io/webscrapping-investment/pkg/adapter/automation"
	adapter "fv.io/webscrapping-investment/pkg/adapter/persistence"
	entities "fv.io/webscrapping-investment/pkg/domain/entities"
	utils "fv.io/webscrapping-investment/pkg/domain/utils"
	"github.com/joho/godotenv"
)

// define a custom data type for the scraped data

// variables
var repo *adapter.Repository

func main() {
	//cargar config.json
	configFilename := "config.json"

	// Carga la configuración desde el archivo
	appConfig, err := entities.LoadConfig(configFilename)
	if err != nil {
		fmt.Printf("Error cargando la configuración: %v\n", err)
		return
	}

	//Leer configuraciones
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error Obtener variables de entorno	:", err)
	}
	// Utiliza algo del paquete persistence para resolver el error
	flagDb, _ := utils.StringToBool(os.Getenv("DATABASE_FLAG"))

	if flagDb {
		//0 QUERY INVESTMENT
		// Crear una instancia de Database
		db, err := adapter.NewDatabase()
		if err != nil {
			fmt.Println("Error al crear la conexión a MongoDB:", err)
			return
		}
		defer db.Close()
		//Inicio de collections
		collectionInvestment := db.GetCollection(appConfig.DataBase.CollectionInvestment)
		collectionInvestmentGroup := db.GetCollection(appConfig.DataBase.CollectionInvestmentGroup)

		// Crear una instancia de Repository
		repo = adapter.NewRepository(db, collectionInvestment, collectionInvestmentGroup)

		//repo.DropCollectionInversiones()
		//repo.DropCollectionInversionesGroup()

		err = repo.CreateUniqueIndexGroup()
		err = repo.CreateUniqueIndex()

		_, err = repo.QueryCountInvestment()
		if err != nil {
			fmt.Println("Error al contar movimientos actuales:", err)
			return
		}
	}

	//1 SELENIUM config
	seleniumService, err := automation.NewWebDriverManager(appConfig.Selenium.SeleniumPort, appConfig.Selenium.SeleniumDebug)
	if err != nil {
		log.Fatal("Error Inicio Selenium:", err)
	}
	defer seleniumService.Close()

	// FLAG PARA IR A RACIONAL
	flagRacional, _ := utils.StringToBool(os.Getenv("RACIONAL_FLAG"))
	var moves, movesCash []entities.Movimiento
	if flagRacional {
		//1 SECCION LOGIN
		seleniumService.RacionalGoToLoginPage(appConfig.Racional.LoginPageURL)

		// 2 SECCION MOVIMIENTOS
		seleniumService.RacionalGoToMovementPage(appConfig.Racional.MovementsPageURL)
		// 3 SCRAP
		moves, movesCash = seleniumService.RacionalMovementPageScrap(appConfig.Racional.FlagScroll)
	}

	// FLAG PARA IR A RENTA4
	flagRenta4, _ := utils.StringToBool(os.Getenv("RENTA4_FLAG"))
	if flagRenta4 {
		//1 SECCION LOGIN
		seleniumService.Renta4GoToLoginPage(appConfig.Renta4.LoginPageURL)

		// 2 SECCION MOVIMIENTOS
		seleniumService.Renta4GoToMovementPage(appConfig.Renta4.MovementsPageURL)

		// 3 SCRAP
		moves, movesCash = seleniumService.Renta4MovementPageScrap(appConfig.Renta4.FlagTest)
	}

	//AGRUPAMOS DATA
	flagGroupData, _ := utils.StringToBool(os.Getenv("DATA_AGRUPADA_XDAY"))
	var movesGroup []entities.Movimiento

	if flagGroupData {
		//asignamos
		movesGroup = entities.GroupAndSumInterface(entities.ConvertToSliceOfSlice(moves))
	}

	//4 GUARDAR EN BD
	if flagDb {
		repo.InsertManyInversiones(entities.ToInterface(moves))
		if flagGroupData {
			repo.InsertManyInversionesAgrupadas(entities.ConvertToInterface(movesGroup))
		}

	}

	//5 GENERAR CSV
	csvService := adapter.NewCSVCreator("./output/" + os.Getenv("FILE_CSV_NAME") + ".csv")
	csvService.CreateCSV(entities.ConvertToSliceOfSlice(moves))
	csvServiceCash := adapter.NewCSVCreator("./output/" + os.Getenv("FILE_CSV_NAME") + "-cashflow.csv")
	csvServiceCash.CreateCSV(entities.ConvertToSliceOfSlice(movesCash))
	csvServiceGroup := adapter.NewCSVCreator("./output/" + os.Getenv("FILE_CSV_NAME") + "-group.csv")
	csvServiceGroup.CreateCSV(entities.ConvertToSliceOfSlice(movesGroup))

	//6 GOOGLESHEETS
	flagGoogleSheet, _ := utils.StringToBool(os.Getenv("GOOGLESHEETS_FLAG"))
	if flagGoogleSheet {
		// Crear un nuevo cliente de Google Sheets
		googleClient, err := adapter.NewGoogleDriveClient(os.Getenv("GOOGLESHEETS_CREDENTIALS_PATH"))
		if err != nil {
			log.Fatalf("Error al crear el cliente de Google Sheets: %v", err)
		}

		// Añadir la nueva pestaña a la hoja de cálculo
		err = googleClient.AddSheet(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsInvestment)
		if err != nil {
			fmt.Println(err)
		}
		// Añadir datos a la pestaña existente
		err = googleClient.AddData(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsInvestment, entities.ConvertToInterfaceMatriz(moves))
		if err != nil {
			log.Fatalf("Error al añadir datos a la pestaña: %v", err)
		}

		// Añadir la nueva pestaña a la hoja de cálculo
		err = googleClient.AddSheet(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsCashFlow)
		if err != nil {
			fmt.Println(err)
		}
		// Añadir datos a la pestaña existente
		err = googleClient.AddData(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsCashFlow, entities.ConvertToInterfaceMatriz(movesCash))
		if err != nil {
			log.Fatalf("Error al añadir datos a la pestaña: %v", err)
		}

		// Añadir la nueva pestaña a la hoja de cálculo
		err = googleClient.AddSheet(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsInvestmentGroup)
		if err != nil {
			fmt.Println(err)
		}
		// Añadir datos a la pestaña existente
		err = googleClient.AddData(os.Getenv("GOOGLESHEETS_SPREADSHEETS_ID"), appConfig.GoogleSheets.SpreadsheetsInvestmentGroup, entities.ConvertToInterfaceMatriz(movesGroup))
		if err != nil {
			log.Fatalf("Error al añadir datos a la pestaña: %v", err)
		}

		fmt.Println("Datos añadidos exitosamente.")

	}

}
