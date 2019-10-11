package ubp

import (
	"testing"

	disbursement "github.com/jfpalngipang/fund-disbursement"
)

func TestDisbursement(m *testing.M) {

}

func Test_TransferFundsSuccess(t *testing.T) {
	Convey("testing transfer transfer funds ", t, func() {
		addr := disbursement.Address{
			Line1:    "Unit 11C15 Fort Victoria",
			Line2:    "23rd street Fort Bonifacio",
			City:     "Taguig",
			Province: "Metro Manila",
			ZipCode:  "1630",
			Country:  "Philippines",
		}
		r := disbursement.Receiver{
			AccountNumber: "100076532781",
			Name:          "Juan Dela Cruz",
			Address:       addr,
		}

		det := disbursement.Details{
			Amount:        "100.00",
			Currncy:       "PHP",
			ReceivingBank: "161312",
			Purpose:       "Fund Transfer",
			Instructions:  "Test Instruction",
		}

		d := disbursement.Disbursement{
			Receiver: r,
			Details:  det,
		}
	})
}
