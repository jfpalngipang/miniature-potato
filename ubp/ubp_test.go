package ubp

import (
	"testing"

	"math/rand"
	"strconv"
	"strings"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUBP_AuthenticatePartner(t *testing.T) {
	Convey("test getting partner access token", t, func() {
		ubp := UBP{}
		ubp.Init()
		response, err := ubp.AuthenticatePartner()
		So(err, ShouldBeNil)
		So(response.TokenType, ShouldNotBeEmpty)
		So(response.Scope, ShouldNotBeEmpty)
		So(response.AccessToken, ShouldNotBeEmpty)
		So(response.ExpiresIn, ShouldNotBeEmpty)
		So(response.MetaData, ShouldNotBeEmpty)
		So(response.RefreshToken, ShouldNotBeEmpty)
	})
}

func TestInstaPay_TransferFundsFromPartnerAccount(t *testing.T) {
	Convey("test instapay fund transfer", t, func() {
		ubp := UBP{}
		ubp.Init()
		response, err := ubp.AuthenticatePartner()
		So(err, ShouldBeNil)
		So(response.AccessToken, ShouldNotEqual, "")

		companyName := "Palngipang Corp."
		companyPrimaryAddress := "Palngipang Tower"
		companySecondaryAddress := "Fort Bonifacio"
		companyPrimaryCity := "Taguig"
		companyPrimaryState := "Metro Manila"
		companyPrimaryZip := "2634"
		companyPrimaryCountry := "Philippines"

		sender := Sender{
			Name: companyName,
			Address: Address{
				Line1:    companyPrimaryAddress,
				Line2:    companySecondaryAddress,
				City:     companyPrimaryCity,
				Province: companyPrimaryState,
				ZipCode:  companyPrimaryZip,
				Country:  companyPrimaryCountry,
			},
		}

		// NOTE: This should be from the profile service
		beneficiary := Beneficiary{
			AccountNumber: "109453095653", // NOTE:: Should be from the request
			Name:          "Rachelle",
			Address: Address{
				Line1:    "241 A.DEL MUNDO ST BET. 5TH 6TH AVE GRACE",
				Line2:    "PARK CALOOCAN CITY",
				City:     "Caloocan",
				Province: "Manila",
				ZipCode:  "1900",
				Country:  "PH",
			},
		}

		// NOTE:: Should all be from the request other than currency
		remittance := Remittance{
			Amount:        "30.00",
			Currency:      "PHP",
			ReceivingBank: "161408",
			Purpose:       "1001",
			Instructions:  "Just a test case",
		}
		rand.Seed(time.Now().UnixNano())
		refID := strconv.Itoa(rand.Int())

		ubp.FundTransferRequest.Sender = sender
		ubp.FundTransferRequest.SenderRefId = refID[0:8]
		ubp.FundTransferRequest.Beneficiary = beneficiary
		ubp.FundTransferRequest.Remittance = remittance
		fundResponse, err := ubp.TransferFundsFromPartnerAccount(response.AccessToken, "instapay", ApiPathPartnerInstaPayFundTransfer)

		So(err, ShouldBeNil)
		So(fundResponse.SenderRefId, ShouldEqual, ubp.FundTransferRequest.SenderRefId)
		So(fundResponse.CreatedAt, ShouldNotBeEmpty)
		So(fundResponse.TranId, ShouldNotBeEmpty)
		So(strings.ToLower(fundResponse.State), ShouldContainSubstring, strings.ToLower("Credited Beneficiary Account"))
	})
}

func TestPesoNet_TransferFundsFromPartnerAccount(t *testing.T) {
	Convey("test pesonet fund transfer", t, func() {
		ubp := UBP{}
		ubp.Init()
		response, err := ubp.AuthenticatePartner()
		So(err, ShouldBeNil)
		So(response.AccessToken, ShouldNotEqual, "")

		// NOTE: This should be from the profile service
		beneficiary := Beneficiary{
			AccountNumber: "107324511489", // NOTE:: Should be from the request
			Name:          "Rachel",
			Address: Address{
				Line1:    "241 A.DEL MUNDO ST BET. 5TH 6TH AVE GRACE",
				Line2:    "PARK CALOOCAN CITY",
				City:     "Caloocan",
				Province: "Manila",
				ZipCode:  "1900",
				Country:  "PH",
			},
		}

		// NOTE:: Should all be from the request other than currency
		remittance := Remittance{
			Amount:        "30.00",
			Currency:      "PHP",
			ReceivingBank: "161203",
			Purpose:       "1001",
			Instructions:  "Just a test case",
		}

		rand.Seed(time.Now().UnixNano())
		refID := strconv.Itoa(rand.Int())
		ubp.FundTransferRequest.SenderRefId = refID[0:8]
		ubp.FundTransferRequest.Beneficiary = beneficiary
		ubp.FundTransferRequest.Remittance = remittance
		fundResponse, err := ubp.TransferFundsFromPartnerAccount(&ubp.Config, response.AccessToken, ApiPathPartnerPESONetFundTransfer)

		So(err, ShouldBeNil)
		So(fundResponse.SenderRefId, ShouldEqual, ubp.FundTransferRequest.SenderRefId)
		So(fundResponse.CreatedAt, ShouldNotBeEmpty)
		So(fundResponse.UbpTranId, ShouldNotBeEmpty)
		So(strings.ToLower(fundResponse.State), ShouldContainSubstring, "sent for processing")
		//So(strings.ToLower(fundResponse.State), ShouldContainSubstring, strings.ToLower("Credited Beneficiary Account"))
	})
}
