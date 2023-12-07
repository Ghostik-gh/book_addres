package stdhttp

import (
	"HW-1/gates/psg"
	"HW-1/models/dto"
	"HW-1/pkg"
	"encoding/json"
	"errors"
	"io"
	"log"

	"net/http"
)

type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
}

type Record struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}

func NewController(addr string, p *psg.Psg) *Controller {
	ctrl := &Controller{}

	mux := http.NewServeMux()

	mux.Handle("/create", http.HandlerFunc(ctrl.RecordAdd))
	mux.Handle("/get", http.HandlerFunc(ctrl.RecordsGet))
	mux.Handle("/update", http.HandlerFunc(ctrl.RecordUpdate))
	mux.Handle("/delete", http.HandlerFunc(ctrl.RecordDeleteByPhone))

	ctrl.Srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	ctrl.DB = p

	return ctrl
}

func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request", nil, err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		log.Println("Error JSON", nil, err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	if record.Name == "" || record.LastName == "" || record.Address == "" || record.Phone == "" {
		err = errors.New("required data is missing")
		log.Println("Required data is missing", nil, err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		log.Println("Error: wrong Phone", nil, err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = c.DB.RecordSave(record)

	if err != nil {
		log.Println("Error in saving record", nil, err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	log.Println("Successfully added", nil, "")

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request", nil, err.Error())
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		log.Println("Error JSON", nil, err.Error())
		return
	}

	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			log.Println("Error: wrong Phone", nil, err.Error())
			return
		}
	}

	records, err := c.DB.RecordsGet(record)
	if err != nil {
		log.Println("Error in finding records", nil, err.Error())
		return
	}

	recordsJSON, err := json.Marshal(records)
	if err != nil {
		log.Println("Error JSON", nil, err.Error())
		return
	}
	w.Write(recordsJSON)
	w.WriteHeader(http.StatusOK)
	log.Println("Success", recordsJSON, "")
}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request", nil, err.Error())
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		log.Println("Error JSON", nil, err.Error())
		return
	}

	if (record.Name == "" && record.LastName == "" && record.MiddleName == "" && record.Address == "") || record.Phone == "" {
		err = errors.New("required data is missing")
		log.Println("Required data is missing", nil, err.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		log.Println("Error: wrong Phone", nil, err.Error())
		return
	}

	err = c.DB.RecordUpdate(record)
	if err != nil {
		log.Println("Error in updating record", nil, err.Error())
		return
	}
	log.Println("Success", nil, "")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request", nil, err.Error())
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		log.Println("Error JSON", nil, err.Error())
		return
	}

	if record.Phone == "" {
		err = errors.New("phone data is missing")
		log.Println("Phone data is missing", nil, err.Error())
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		log.Println("Error: wrong Phone", nil, err.Error())
		return
	}

	err = c.DB.RecordDeleteByPhone(record.Phone)
	if err != nil {
		log.Println("Error in deleting record", nil, err.Error())
		return
	}
	log.Println("Success", nil, "")
	w.WriteHeader(http.StatusOK)
}
