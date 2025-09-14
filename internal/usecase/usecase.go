package usecase

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson(res http.ResponseWriter, status int, data any) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)

	encoder := json.NewEncoder(res)
	if err := encoder.Encode(data); err != nil {
		log.Printf("coding error  JSON: %v", err)
		return err
	}
	return nil
}

func WriteError(res http.ResponseWriter, status int, msg string) {
	WriteJson(res, status, map[string]string{"error:": msg})
}
