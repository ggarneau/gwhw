package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusCode int

func JSON(w http.ResponseWriter, j interface{}, c_optional ...StatusCode) {
	w.Header().Set("Content-Type", "application/json")
	c := StatusCode(200) // Default
	if len(c_optional) > 0 {
		c = c_optional[0]
	}
	w.WriteHeader(int(c))
	jsonResp, err := json.Marshal(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", jsonResp)
}
