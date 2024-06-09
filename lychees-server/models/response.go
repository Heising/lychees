package models

type ResponseERR struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
type TokenInfo struct {
	Token      *string `json:"token,omitempty"`
	ExpireUnix int64   `json:"expireUnix"`
}
