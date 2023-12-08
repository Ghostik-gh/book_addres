package stdhttp

import (
	"HW-1/gates/psg"
	"HW-1/models/dto"
	"HW-1/pkg"
	"HW-1/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"io"

	"net/http"
)

type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
	ctx context.Context
}

type Record struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}

func NewController(ctx context.Context, addr string, p *psg.Psg) *Controller {
	ctrl := &Controller{}
	ctrl.ctx = ctx

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
		logger.Errorf(c.ctx, "Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(c.ctx, "Error reading request", nil, err.Error())
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		logger.Errorf(c.ctx, "Error JSON", nil, err.Error())
		http.Error(w, "Error JSON", http.StatusBadRequest)
		return
	}

	if record.Name == "" || record.LastName == "" || record.Address == "" || record.Phone == "" {
		err = errors.New("required data is missing")
		logger.Errorf(c.ctx, "Required data is missing", nil, err.Error())
		http.Error(w, "Required data is missing", http.StatusBadRequest)
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		logger.Errorf(c.ctx, "Error: wrong Phone", nil, err.Error())
		http.Error(w, "Error: wrong Phone", http.StatusBadRequest)
		return
	}

	err = c.DB.RecordSave(record)

	if err != nil {
		logger.Errorf(c.ctx, "Error in saving record", nil, err.Error())
		http.Error(w, "Error in saving record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Infof(c.ctx, "Successfully added", nil, "")
}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Errorf(c.ctx, "Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(c.ctx, "Error reading request", nil, err.Error())
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		logger.Errorf(c.ctx, "Error JSON", nil, err.Error())
		http.Error(w, "Error JSON", http.StatusBadRequest)
		return
	}

	if record.Phone != "" {
		record.Phone, err = pkg.PhoneNormalize(record.Phone)
		if err != nil {
			logger.Errorf(c.ctx, "Error: wrong Phone", nil, err.Error())
			http.Error(w, "Error: wrong Phone", http.StatusBadRequest)
			return
		}
	}

	records, err := c.DB.RecordsGet(record)
	if err != nil {
		logger.Errorf(c.ctx, "Error in finding records", nil, err.Error())
		http.Error(w, "Error in finding records", http.StatusInternalServerError)
		return
	}

	recordsJSON, err := json.Marshal(records)
	if err != nil {
		logger.Errorf(c.ctx, "Error JSON", nil, err.Error())
		http.Error(w, "Error JSON", http.StatusInternalServerError)
		return
	}
	w.Write(recordsJSON)
	w.WriteHeader(http.StatusOK)
	logger.Infof(c.ctx, "Success", recordsJSON, "")
}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logger.Errorf(c.ctx, "Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(c.ctx, "Error reading request", nil, err.Error())
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		logger.Errorf(c.ctx, "Error JSON", nil, err.Error())
		http.Error(w, "Error JSON", http.StatusBadRequest)
		return
	}

	if (record.Name == "" && record.LastName == "" && record.MiddleName == "" && record.Address == "") || record.Phone == "" {
		err = errors.New("required data is missing")
		logger.Errorf(c.ctx, "Required data is missing", nil, err.Error())
		http.Error(w, "Required data is missing", http.StatusBadRequest)
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		logger.Errorf(c.ctx, "Error: wrong Phone", nil, err.Error())
		http.Error(w, "Error: wrong Phone", http.StatusBadRequest)
		return
	}

	err = c.DB.RecordUpdate(record)
	if err != nil {
		logger.Errorf(c.ctx, "Error in updating record", nil, err.Error())
		http.Error(w, "Error in updating record", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	logger.Infof(c.ctx, "Success", nil, "")
}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	
	if r.Method != http.MethodPost {
		logger.Errorf(c.ctx, "Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := dto.Record{}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(c.ctx, "Error reading request", nil, err.Error())
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteReq, &record)
	if err != nil {
		logger.Errorf(c.ctx, "Error JSON", nil, err.Error())
		http.Error(w, "Error JSON", http.StatusBadRequest)
		return
	}

	if record.Phone == "" {
		err = errors.New("phone data is missing")
		logger.Errorf(c.ctx, "Phone data is missing", nil, err.Error())
		http.Error(w, "Phone data is missing", http.StatusBadRequest)
		return
	}

	record.Phone, err = pkg.PhoneNormalize(record.Phone)
	if err != nil {
		logger.Errorf(c.ctx, "Error: wrong Phone", nil, err.Error())
		http.Error(w, "Error: wrong Phone", http.StatusBadRequest)
		return
	}

	err = c.DB.RecordDeleteByPhone(record.Phone)
	if err != nil {
		logger.Errorf(c.ctx, "Error in deleting record", nil, err.Error())
		http.Error(w, "Error in deleting record", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	logger.Infof(c.ctx, "Success", nil, "")
}
