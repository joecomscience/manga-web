package grafana

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request)  {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
}