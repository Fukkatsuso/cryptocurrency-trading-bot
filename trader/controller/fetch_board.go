package controller

import (
	"fmt"
	"net/http"
)

func FetchBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Fetch Board")
}
