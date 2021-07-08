package models

type Lock struct {
	UUID    string `json:"uuid"`
	TIMEOUT int64  `json:"timeout"`
}
