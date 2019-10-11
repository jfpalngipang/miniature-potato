package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	disbursement "github.com/jfpalngipang/fund-disbursement"
	"github.com/pressly/chi"
)

type disbursementHandler struct {
	router chi.Router

	baseUrl             url.URL
	disbursementService disbursement.DisbursementService
}

func newDisbursementHandler() *disbursementHandler {
	h := &disbursementHandler{router: chi.NewRouter()}
	h.router.Get("/instapay/banks", h.handleGetBanksForInstapay)
	h.router.Get("/pesonet/banks", h.handleGetBanksForPesonet)
	h.router.Post("/single/instapay", h.handleSingleDisbursementViaInstapay)
	h.router.Post("/single/pesonet", h.handleSingleDisbursementViaPesonet)
	h.router.Post("/single/ubptoubp", h.handleSingleDisbursementViaUbpToUbp)
	h.router.Get("/status/{method}/{refID}", h.handleGetStatus)
	return h
}

func (h *disbursementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *disbursementHandler) handleGetBanksForInstapay(w http.ResponseWriter, r *http.Request) {
	resp, _ := h.disbursementService.GetBanks("instapay")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *disbursementHandler) handleGetBanksForPesonet(w http.ResponseWriter, r *http.Request) {
	resp, _ := h.disbursementService.GetBanks("pesonet")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *disbursementHandler) handleSingleDisbursementViaInstapay(w http.ResponseWriter, r *http.Request) {
	var fundTransferRequestBody disbursement.Disbursement
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &fundTransferRequestBody)
	if err != nil {
		log.Println("Error in unmarshaling the request body.")
	}
	resp, _ := h.disbursementService.TransferFunds("instapay", &fundTransferRequestBody)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *disbursementHandler) handleSingleDisbursementViaPesonet(w http.ResponseWriter, r *http.Request) {
	var fundTransferRequestBody disbursement.Disbursement
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &fundTransferRequestBody)
	if err != nil {
		log.Println("Error in unmarshaling the request body.")
	}

	resp, _ := h.disbursementService.TransferFunds("pesonet", &fundTransferRequestBody)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *disbursementHandler) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	method := chi.URLParam(r, "method")
	refID := chi.URLParam(r, "refID")
	resp, _ := h.disbursementService.GetStatus(method, refID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
