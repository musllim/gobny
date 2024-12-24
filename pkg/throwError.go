package pkg

import (
	"fmt"
	"net/http"
)

func ThrowError(w http.ResponseWriter, message string, status int, err error) {
	w.WriteHeader(status)
	fmt.Println(message, err.Error())
	w.Write([]byte(message))
}
