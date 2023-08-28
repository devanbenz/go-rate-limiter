package internal

import (
	"fmt"
	"net/http"
	"time"
)

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	currTime := time.Now().Format(time.Kitchen)
	_, err := w.Write([]byte(fmt.Sprintf("Time: %v", currTime)))
	if err != nil {
		timeHandlerError := fmt.Sprintf("%v\n", err)
		http.Error(w, timeHandlerError, http.StatusBadRequest)
		return
	}
}
