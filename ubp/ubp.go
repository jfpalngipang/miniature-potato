package ubp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	disbursement "github.com/jfpalngipang/fund-disbursement"
	"github.com/mitchellh/mapstructure"
)

type UBP struct {
	Config              Config
	FundTransferRequest FundTransferRequest
}

type ApiCall struct {
	Url               string
	Method            string
	Body              io.Reader
	AdditionalHeaders map[string]string
}

const (
	DisbursementMethodInstapay = "instapay"
	DisbursementMethodPesonet  = "pesonet"
	DisbursementMethodUbptoUbp = "ubp"
)

type AuthTokenResponse struct {
	AccessToken  string `mapstructure:"access_token" json:"access_token" yaml:"access_token"`
	ExpiresIn    int    `mapstructure:"expires_in" json:"expires_in" yaml:"expires_in"`
	MetaData     string `mapstructure:"metadata" json:"metadata" yaml:"metadata"`
	RefreshToken string `mapstructure:"refresh_token" json:"refresh_token" yaml:"refresh_token"`
	Scope        string `mapstructure:"scope" json:"scope" yaml:"scope"`
	TokenType    string `mapstructure:"token_type" json:"token_type" yaml:"token_type"`
}

type FundTransferResponse struct {
	TranId      string `json:"tranId"`    // this is the json key for txn id when using instapay fund transfer
	UbpTranId   string `json:"ubpTranId"` // this is the json key for txn id when using pesonet fund transfer
	CreatedAt   string `json:"createdAt"`
	State       string `json:"state"`
	SenderRefId string `json:"senderRefId"`
}

type TransferRequest struct {
	Beneficiary Beneficiary `json:"receiver"`
	Remittance  Remittance  `json:"details"`
}

type Address struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	City     string `json:"city"`
	Province string `json:"province"`
	ZipCode  string `json:"zipCode"`
	Country  string `json:"country"`
}

type Sender struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

var PartnerAddress = &Address{
	Line1:    "Some Tower",
	Line2:    "Some Barangay",
	City:     "Some City",
	Province: "Metro Manila",
	ZipCode:  "4024",
	Country:  "Philippines",
}

var PartnerSender = &Sender{
	Name:    "Palngipang Corp.",
	Address: *PartnerAddress,
}

var ClientBeneficary = &Beneficiary{
	AccountNumber: "107324511489",
	Name:          "Rachelle",
	Address:       *PartnerAddress,
}

var SampleRemittance = &Remittance{
	Amount:        "2000.00",
	Currency:      "PHP",
	ReceivingBank: "161203",
	Purpose:       "5 632",
	Instructions:  "Test Pesonet",
}

type FundTransferRequest struct {
	SenderRefId string      `json:"senderRefId"`
	RequestDate string      `json:"tranRequestDate"`
	Sender      Sender      `json:"sender"`
	Beneficiary Beneficiary `json:"beneficiary"`
	Remittance  Remittance  `json:"remittance"`
}

type Beneficiary struct {
	AccountNumber string  `json:"accountNumber"`
	Name          string  `json:"name"`
	Address       Address `json:"address"`
}

type Remittance struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	ReceivingBank string `json:"receivingBank"`
	Purpose       string `json:"purpose"`
	Instructions  string `json:"instructions"`
}

type Bank struct {
	Code string `json:"code"`
	Bank string `json:"bank"`
}

type GetBanksResponse struct {
	Records      []Bank `json:"records"`
	TotalRecords uint32 `json:"totalRecords"`
}

type InstapayTransaction struct {
	UbpTranID   string `json:"ubpTranId"`
	Type        string `json:"type"`
	CreatedAt   string `json:"createdAt"`
	State       string `json:"state"`
	SenderRefID string `json:"senderRefId"`
}

type InstapayStatusResponse struct {
	Records      []InstapayTransaction `json:"records"`
	TotalRecords uint32                `json:"totalRecords"`
}

type PesonetTransaction struct {
	UbpTranID   string `json:"ubpTranId"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	State       string `json:"state"`
	SenderRefID string `json:"senderRefId"`
}

type PesonetStatusResponse struct {
	Records      []PesonetTransaction `json:"records"`
	TotalRecords uint32               `json:"totalRecords"`
}

func (u *UBP) Init() {
	// u.Config.LoadConfiguration("/Users/jfpalngipang/fund-disbursement/ubp/config.dev.json")
	u.Config.LoadConfiguration("/app/ubp/config.dev.json")
}

// AuthenticatePartner to get token from UBP API
func (u *UBP) AuthenticatePartner() (AuthTokenResponse, error) {
	values := url.Values{}
	values.Add("client_id", u.Config.ClientId)
	values.Add("username", u.Config.Username)
	values.Add("password", u.Config.Password)
	values.Add("grant_type", "password")
	values.Add("scope", u.Config.Scope)

	apiCall := ApiCall{
		Method: http.MethodPost,
		Url:    u.Config.BaseUrl + u.Config.PartnerAuthPath,
		Body:   strings.NewReader(values.Encode()),
		AdditionalHeaders: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	response, err := ApiRequest(&apiCall, &u.Config)
	if err != nil {
		fmt.Printf("Partner Auth API call error: %s\n", err)
		return AuthTokenResponse{}, err
	}

	var authResponse AuthTokenResponse

	err = mapstructure.Decode(response, &authResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return AuthTokenResponse{}, err
	}

	return authResponse, nil
}

func (u *UBP) TransferFundsFromPartnerAccount(token string, method string, transferRequest *disbursement.Disbursement) (FundTransferResponse, error) {
	reqDate := time.Now().Format("2006-01-02T15:04:05.000")
	u.FundTransferRequest.RequestDate = reqDate[0:23] // max 23 chars oly
	u.FundTransferRequest.Sender = *PartnerSender

	addr := Address{
		Line1:    transferRequest.Receiver.Address.Line1,
		Line2:    transferRequest.Receiver.Address.Line2,
		City:     transferRequest.Receiver.Address.City,
		Province: transferRequest.Receiver.Address.Province,
		ZipCode:  transferRequest.Receiver.Address.ZipCode,
		Country:  transferRequest.Receiver.Address.Country,
	}

	var ben = &Beneficiary{
		AccountNumber: transferRequest.Receiver.AccountNumber,
		Name:          transferRequest.Receiver.Name,
		Address:       addr,
	}

	var r = &Remittance{
		Amount:        transferRequest.Details.Amount,
		Currency:      transferRequest.Details.Currency,
		ReceivingBank: transferRequest.Details.ReceivingBank,
		Purpose:       "1001",
		Instructions:  "Fund Transfer via Instapay",
	}

	u.FundTransferRequest.Beneficiary = *ben
	u.FundTransferRequest.Remittance = *r
	refID := u.GenerateRefId()
	u.FundTransferRequest.SenderRefId = refID

	b, err := json.Marshal(u.FundTransferRequest)
	if err != nil {
		fmt.Printf("Unmarshal Error: %s\n", err)
		return FundTransferResponse{}, err
	}

	var apiPath string
	if method == DisbursementMethodInstapay {
		apiPath = u.Config.InstapayPath
	} else if method == DisbursementMethodPesonet {
		apiPath = u.Config.PesonetPath
	}

	apiCall := ApiCall{
		Method: http.MethodPost,
		Url:    u.Config.BaseUrl + apiPath,
		Body:   bytes.NewBuffer(b),
		AdditionalHeaders: map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": "Bearer " + token,
		},
	}

	response, err := ApiRequest(&apiCall, &u.Config)
	if err != nil {
		fmt.Printf("Error in Instapay Request: %s\n", err)
		return FundTransferResponse{}, err
	}

	var fundTransferResponse FundTransferResponse

	err = mapstructure.Decode(response, &fundTransferResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return FundTransferResponse{}, err
	}

	return fundTransferResponse, nil
}

func (u *UBP) GetBanksForTransfer(method string) (GetBanksResponse, error) {
	var apiPath string
	if method == DisbursementMethodInstapay {
		apiPath = u.Config.InstapayGetBanksPath
	} else if method == DisbursementMethodPesonet {
		apiPath = u.Config.PesonetGetBanksPath
	}

	apiCall := ApiCall{
		Method: http.MethodGet,
		Url:    u.Config.BaseUrl + apiPath,
		Body:   nil,
		AdditionalHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	response, err := ApiRequest(&apiCall, &u.Config)
	if err != nil {
		fmt.Printf("Error in retrieving bank list: %s\n", err)
		return GetBanksResponse{}, err
	}

	var getBanksResponse GetBanksResponse

	err = mapstructure.Decode(response, &getBanksResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return GetBanksResponse{}, err
	}

	return getBanksResponse, nil
}

func (u *UBP) GetPesonetTransferStatus(referenceId string) (PesonetStatusResponse, error) {
	apiPath := strings.Replace(u.Config.GetTransferStatusPath, "{referenceId}", referenceId, -1)
	apiPath = strings.Replace(apiPath, "{method}", "pesonet", -1)
	apiCall := ApiCall{
		Method: http.MethodGet,
		Url:    u.Config.BaseUrl + apiPath,
		Body:   nil,
		AdditionalHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	response, err := ApiRequest(&apiCall, &u.Config)
	if err != nil {
		fmt.Printf("Error in retrieving bank list: %s\n", err)
		return PesonetStatusResponse{}, err
	}

	var pesonetStatusResponse PesonetStatusResponse

	err = mapstructure.Decode(response, &pesonetStatusResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return PesonetStatusResponse{}, err
	}

	return pesonetStatusResponse, nil
}

func (u *UBP) GetInstapayTransferStatus(referenceId string) (InstapayStatusResponse, error) {
	apiPath := strings.Replace(u.Config.GetTransferStatusPath, "{referenceId}", referenceId, -1)
	apiPath = strings.Replace(apiPath, "{method}", "instapay", -1)
	apiCall := ApiCall{
		Method: http.MethodGet,
		Url:    u.Config.BaseUrl + apiPath,
		Body:   nil,
		AdditionalHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	response, err := ApiRequest(&apiCall, &u.Config)
	if err != nil {
		fmt.Printf("Error in retrieving bank list: %s\n", err)
		return InstapayStatusResponse{}, err
	}

	var instapayStatusResponse InstapayStatusResponse

	err = mapstructure.Decode(response, &instapayStatusResponse)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return InstapayStatusResponse{}, err
	}

	return instapayStatusResponse, nil
}

// ApiRequest for forming and sending request to UBP API
func ApiRequest(api *ApiCall, conf *Config) (interface{}, error) {
	req, _ := http.NewRequest(api.Method, api.Url, api.Body)

	for k, v := range api.AdditionalHeaders {
		req.Header.Add(k, v)
	}

	conf.SetHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: %s\n", err)
		return nil, err
	}

	fmt.Println("RESPONSE:")
	fmt.Println("For api call: " + req.URL.String())
	fmt.Println("Status: ", res.Status)
	fmt.Println()

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error in reading response body: %s\n", err)
		return nil, err
	}

	fmt.Println("Response Body: ", string(resBody))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		fmt.Println("response error: ", string(resBody))
		return nil, errors.New(string(resBody))
	}

	var response interface{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		fmt.Printf("Error in unmarmshaling response body: %s\n", err)
		return nil, err
	}

	return response, nil
}

func (c *Config) SetHeaders(req *http.Request) {
	req.Header.Add("x-ibm-client-id", c.ClientId)
	req.Header.Add("x-ibm-client-secret", c.ClientSecret)
	req.Header.Add("x-partner-id", c.PartnerId)
	req.Header.Add("x-client-id", c.ClientId)
	req.Header.Add("x-client-secret", c.ClientSecret)
}

func (u *UBP) GenerateRefId() string {
	rand.Seed(time.Now().UnixNano())
	sRef := strconv.Itoa(rand.Int())
	return sRef[0:8]
}
