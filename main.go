// Api_Encrypt
//
// This API allows encrypt and decrypt documents
//
//	Schemes: [http, https]
//	host: localhost:9076/api
//	Version: 0.0.1
//	Contact: Felipe Morales<fmorales@asicom.cl>
//
// 	Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"api_encrypt/src/functions/documentos"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// swagger:operation GET / Index index
// Check if the service is running.
// ---
// responses:
//     '200':
//         description: Welcome to my Rest API to encrypt and decrypt documents
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Rest API to encrypt and decrypt documents")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", index).Methods("GET")
	router.HandleFunc("/api/encrypt", documentos.EncryptDocument).Methods("POST")
	router.HandleFunc("/api/decrypt", documentos.DecryptDocument).Methods("POST")
	router.HandleFunc("/api/decrypt-hash", documentos.DecryptHash).Methods("POST")
	http.ListenAndServe(":9076", router)
}

//Se deben generar los siguientes directorios en la siguiente ruta
// /home/repositorio-documentos/nombreproyecto/nombredirectoriorelacionadoalarchivo
// ejemplo: /home/repositorio-documentos/essbio/pagos
// reclamos
// pagos
// cron
