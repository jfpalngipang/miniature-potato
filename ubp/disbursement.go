package ubp

import disbursement "github.com/jfpalngipang/fund-disbursement"

type DisbursementService struct {
	ubp *UBP
}

func NewDisbursementService(ubp *UBP) *DisbursementService {
	return &DisbursementService{ubp: ubp}
}

func (s *DisbursementService) TransferFunds(method string, transferRequest *disbursement.Disbursement) (interface{}, error) {
	authResponse, _ := s.ubp.AuthenticatePartner()
	fundTransferResponse, _ := s.ubp.TransferFundsFromPartnerAccount(authResponse.AccessToken, method, transferRequest)
	return fundTransferResponse, nil
}

func (s *DisbursementService) GetBanks(method string) (interface{}, error) {
	banksResponse, _ := s.ubp.GetBanksForTransfer(method)
	return banksResponse, nil
}

func (s *DisbursementService) GetStatus(method string, referenceID string) (interface{}, error) {
	var statusResponse interface{}
	if method == "instapay" {
		statusResponse, _ = s.ubp.GetInstapayTransferStatus(referenceID)
	} else if method == "pesonet" {
		statusResponse, _ = s.ubp.GetPesonetTransferStatus(referenceID)
	}

	return statusResponse, nil
}
