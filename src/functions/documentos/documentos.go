package documentos

import (
	"api_encrypt/src/functions/log"
	"api_encrypt/src/models/documentos"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	openssl "github.com/Luzifer/go-openssl/v4"
	openuri "github.com/utahta/go-openuri"
)

// EncryptDocument allows convert files to base64
// swagger:operation POST /api/encrypt EncryptDocument EncryptDocument
// allows convert files to base64
// ---
// parameters:
// - name: files
//   in: formData
//   description: Are the files you want to upload
//   required: true
//   type: file
// - name: servicio
//   in: formData
//   description: Is the number of service of client
//   required: true
//   type: string
// - name: password
//   in: formData
//   description: Is the password for generate hash
//   required: true
//   type: string
// - name: path
//   in: formData
//   description: Is the path where files saved
//   required: true
//   type: string
// responses:
//     '200':
//         description: URL array of uploaded files
//     '400':
//         description: Bad request
//     '404':
//         description: Page not found
//     '500':
//         description: Internal server error
func EncryptDocument(w http.ResponseWriter, r *http.Request) {
	//Permite el Access-Control-Allow-Origin
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// r.ParseMultipartForm(10 << 20) // esto es en MB
	// files := r.MultipartForm.File["files"]
	// servicio := r.FormValue("servicio")
	// password := r.FormValue("password")
	// path := strings.ReplaceAll(r.FormValue("path"), "/", "\\")

	body := documentos.Encryptinput{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		log.CreateLog(err.Error())
		return
	}

	listFiles := []map[string]interface{}{}

	path := fmt.Sprintf("/var/repostorio-documentos/%s/%s", body.Path, body.Servicio) //Certificación y producción
	// path := strings.ReplaceAll(body.Path, "/", "\\")                                   //Desarrollo
	// path = fmt.Sprintf("C:\\home\\repostorio-documentos\\%s\\%s", path, body.Servicio) //Desarrollo
	validatePath(path)

	for _, fh := range body.Files {
		if fh != "" {
			o, err := openuri.Open(fh)
			if err != nil {
				fmt.Println(err)
				log.CreateLog(err.Error())
			}

			namefile := filepath.Base(fh)
			extension := strings.ToLower(filepath.Ext(fh))
			filename := fmt.Sprintf("upload-*%s", extension)

			tempFile, err := ioutil.TempFile(path, filename) //Desarrollo
			if err != nil {
				fmt.Println(err)
				log.CreateLog(err.Error())
			}

			fileBytes, err := ioutil.ReadAll(o)
			if err != nil {
				fmt.Println(err)
				log.CreateLog(err.Error())
			}

			encoded := []byte(base64.StdEncoding.EncodeToString(fileBytes))
			tempFile.Write(encoded)

			hash := generateHash(tempFile.Name(), body.Password)

			// listFiles = append(listFiles, url.QueryEscape(string([]byte(base64.StdEncoding.EncodeToString([]byte(hash))))))
			mapa := map[string]interface{}{
				"hash": url.QueryEscape(string([]byte(base64.StdEncoding.EncodeToString([]byte(hash))))),
				"name": namefile,
			}

			listFiles = append(listFiles, mapa)

			o.Close()
			tempFile.Close()
		}
	}
	json.NewEncoder(w).Encode(listFiles)
}

/*
validatePath check if the path exists and if it does not exist, then the directory (s) are created
*/
func validatePath(directorio string) {
	split := strings.Split(directorio, "/")

	var path string
	for _, p := range split {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			if err != nil {
				path = fmt.Sprintf("%s/%s", path, p)
			}
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err = os.Mkdir(path, 0755)
				if err != nil {
					// Aquí puedes manejar mejor el error, es un ejemplo
					fmt.Println(err)
					log.CreateLog(fmt.Sprintf("validatePath: %s", err.Error()))
				}
			}
		} else {
			path = fmt.Sprintf("%s", p)
		}
	}
}

func generateHash(plaintext string, passphrase string) string {
	o := openssl.New()
	const DefaultPBKDF2Iterations = 10000
	var PBKDF2SHA256 = openssl.NewPBKDF2Generator(sha256.New, DefaultPBKDF2Iterations)

	enc, err := o.EncryptBytes(passphrase, []byte(plaintext), PBKDF2SHA256)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		log.CreateLog(err.Error())
	}

	return string(enc)
}

func DecryptHash(w http.ResponseWriter, r *http.Request) {
	//Permite el Access-Control-Allow-Origin
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	body := documentos.Decrypthashinput{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		log.CreateLog(err.Error())
		return
	}

	if body.Key != "ams-key" {
		w.WriteHeader(400)
		response := map[string]interface{}{
			"error": "No tienes permiso para realizar esta acción",
		}
		//in:body
		json.NewEncoder(w).Encode(response)
		return
	}

	o := openssl.New()

	fmt.Println(body.Password)

	// var BytesToKeyMD5 = openssl.NewBytesToKeyGenerator(md5.New)
	const DefaultPBKDF2Iterations = 10000
	var PBKDF2SHA256 = openssl.NewPBKDF2Generator(sha256.New, DefaultPBKDF2Iterations)

	dec, err := o.DecryptBytes(body.Password, []byte(body.Hash), PBKDF2SHA256)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}

	response := map[string]interface{}{
		"path": string(dec),
	}
	
	json.NewEncoder(w).Encode(response)
}

func generateSHA256Final(plaintext string) string {
	// fmt.Println(plaintext)
	hash := sha256.New()
	hash.Write([]byte(plaintext))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}

func newSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func DecryptDocument(w http.ResponseWriter, r *http.Request) {
	//Permite el Access-Control-Allow-Origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.CreateLog(fmt.Sprintf("r.body -> %s", r.Body))

	body := documentos.Decryptinput{}
	err := json.NewDecoder(r.Body).Decode(&body)
	log.CreateLog(fmt.Sprintf("body -> %s", body))
	if err != nil {
		fmt.Println(err)
		log.CreateLog(fmt.Sprintf("error body -> %s", err.Error()))
		return
	}

	// filename := fmt.Sprintf("C:\\home\\repostorio-documentos\\essbio\\%s", vars["name"])
	// log.CreateLog(body.Path)

	pathString := body.Path
	log.CreateLog(fmt.Sprintf("pathString old %s", pathString))
	log.CreateLog(fmt.Sprintf("pathString %v", strings.Contains(pathString, "home")))
	if strings.Contains(pathString, "home") {
		pathString = strings.Replace(pathString, "home", "var", 1)
	}
	log.CreateLog(fmt.Sprintf("pathString new %s", pathString))

	content, err := ioutil.ReadFile(pathString)
	if err != nil {
		fmt.Println(err)
		log.CreateLog(err.Error())
	}

	fmt.Fprintf(w, "%s", content)
}
