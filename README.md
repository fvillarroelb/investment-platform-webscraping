#  Webscraping-Investment

Este proyecto busca simplificar la vida de inversionistas pequeños para analizar datos de operaciones  y poder obtener estadisticas y orden de lo ejecutado .

## Índice
- [Características](#características)
- [Instalación](#instalación)
- [Uso](#uso)
- [Contribución](#contribución)
- [Licencia](#licencia)

## Características
Este programa incorpora las siguientes funcionalidades:
-   Recuperar movimientos de Racional.cl .
-   Recuperar movimientos de Renta4.cl .
-   Generar un CSV con toda la informacion guardada .
-   Puedes persistir la informacion en una BD No-Relacional .
-   Puedes persistir la informacion en pestañas de GoogleSheets .
-   Se generan 3 pestañas/tablas :
     - Lista de todos los movimientos .
     - Lista de movimientos de dinero entrantes,salientes y generados de los portales de inversion . 
     - Lista de los movimientos agrupados por dia y nemotecnico.

## Instalación
Proporciona instrucciones claras y concisas sobre cómo instalar tu proyecto. Puedes incluir comandos específicos, requisitos previos y pasos adicionales si es necesario.
```sh
# Debes configurar  las variables de entorno necesarias para la aplicacion en un archivo .env
# Debes configurar  las variables de entorno necesarias para la aplicacion en un archivo .env
SELENIUM_DRIVER="./recursos/chromedriver/chromedriver.exe" #ruta al driver de chrome bajo la carpeta /recursos

FILE_CSV_NAME=mis-inversiones                              #nombre del archivo csv que se generara
DATA_AGRUPADA_XDAY=true                                    #Si queremos generar un archivo agrupado por dia

DATABASE_FLAG=false                                         #flag para grabar en BD. Si no necesitas grabar los movimientos en una BD colocar "false"
RACIONAL_FLAG_SCROLL=true
DATABASE_URI="uri-string-mongodb-con-certificado"          #uri de conexion a mongo atlas con url de certificado bajo carpeta ./recursos
DATABASE_NAME=mis-finanzas                                 #nombre de la base de datos

RACIONAL_FLAG=true                                         #flag para ir a racional a buscar movimientos. Si no tienes cuenta en racional colocar "false"
RACIONAL_USER=micorreo@gmail.com                           #usuario racional
RACIONAL_PASSWORD=mi-password                              #password racional
RACIONAL_YEAR_TO=2023                                      #año a buscar movimientos en racional

RENTA4_FLAG=true                                          #flag para ir a renta4 a buscar movimientos. Si no tienes cuenta en renta4 colocar "false"
RENTA4_USER=mi-rut                                         #usuario renta4
RENTA4_PASSWORD=mi-password                                #password renta4 
RENTA4_FECHA_INICIO=01-01-2023                             #fecha inicial para buscar movimientos
RENTA4_FECHA_FIN=                                          #fecha final para buscar movimientos. Si la fecha fin no se setea se busca hata el dia de hoy.


GOOGLESHEETS_FLAG=true                                      #flag para grabar en googlesheets

```
Si quieres ver un ejemplo completo del archivo .env puedes ir a la siguiente ruta ./example/.env-example-full

```bash
# Ejemplo de comandos de instalación
go build .


```


GoogleSheets   
 
Si deseas guardar tus datos en un archivo de GoogleSheets (Aka Excel de google) puedes configurar en la consola una cuenta de servicio para que acceda a tu archivo y almacene tus datos de inversiones.
```
Habilitar la API de Google Sheets:  

    -   Ve a la Consola de Desarrolladores de Google (https://console.cloud.google.com/).
    -   Crea un nuevo proyecto o selecciona uno existente.
    -   Habilita la API de Google Sheets para tu proyecto.
    -   Crea credenciales para la API, selecciona el tipo de cuenta de servicio y descarga el archivo JSON con tus credenciales.
```
MongoDB   

Se puede generar un cluster free-tier de mongo DB para almacenar tus datos de inversiones.
```  
https://cloud.mongodb.com
```
## Uso
Para usar este programa puedes ejecutar el compilado  que se genera  de la siguiente forma :

```bash
#Debes tener en tu carpeta ./recursos el driver de chromiun 
# Ejecutar
webscrapping-investment.exe
```
## Contribución
¡Gracias por considerar contribuir al proyecto! A continuación, se presentan algunas pautas para contribuir:

Forkea el proyecto y clónalo localmente.
Crea una nueva rama para tu contribución.
Realiza tus cambios y asegúrate de seguir las pautas de estilo.
Haz pruebas para asegurarte de que todo funciona como se espera.
Envía una solicitud de extracción a la rama principal del proyecto.
## Licencia
Este proyecto está bajo la Licencia MIT. Consulta el archivo LICENSE para obtener más detalles.