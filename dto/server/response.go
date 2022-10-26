package server

type Response struct {
	Ok    bool        `json:"ok"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}
