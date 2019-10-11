package disbursement

type Disbursement struct {
	Receiver Receiver `json:"receiver"`
	Details  Details  `json:"transfer_details"`
	// IdemKey  string   `json:"idem_key"`
}

type Receiver struct {
	AccountNumber string  `json:"accountNumber"`
	Name          string  `json:"name"`
	Address       Address `json:"address"`
}

type Details struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	ReceivingBank string `json:"receivingBank"`
	Purpose       string `json:"purpose"`
	Instructions  string `json:"instructions"`
}

type Address struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	City     string `json:"city"`
	Province string `json:"province"`
	ZipCode  string `json:"zipCode"`
	Country  string `json:"country"`
}

type DisbursementService interface {
	TransferFunds(method string, d *Disbursement) (interface{}, error)
	GetBanks(method string) (interface{}, error)
	GetStatus(method string, referenceID string) (interface{}, error)
}
