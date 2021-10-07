package documentos

/*
Encryptinput algo
*/
type Encryptinput struct {
	Files    []string `json:"files"`
	Servicio string   `json:"servicio"`
	Password string   `json:"password"`
	Path     string   `json:"path"`
}

type Decryptinput struct {
	Path string `json:"path"`
}

type Decrypthashinput struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
	Key      string `json:"key"`
}
