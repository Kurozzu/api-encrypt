package log

import (
	"log"
	"os"
	"time"
)

/*
CreateLog permite crear un log para llevar un control de lo que sucede en la ejecución del web service
*/
func CreateLog(mensaje string) {
	//Si el archivo no existe, entonces los genera en la ruta especificada. Si existe, entonces lo abre
	//linux
	// f, err := os.OpenFile("/var/log/api_encriptacion/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//windows, ruta por defecto del proyecto. Puedeser cambiada por una ruta opt en windows.
	f, err := os.OpenFile("C:\\var\\log\\api_encriptacion\\error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// f, err := os.OpenFile("src/logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//Agrega la hora al mensaje enviado
	error := time.Now().String() + " " + mensaje + "\n"
	//Escribe en el archivo
	_, err = f.Write([]byte(error))
	if err != nil {
		log.Fatal(err)
	}
	//Cierra el archivo
	f.Close()
}
