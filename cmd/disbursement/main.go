package main

import (
	"fmt"

	"github.com/jfpalngipang/fund-disbursement/http"
	"github.com/jfpalngipang/fund-disbursement/ubp"
)

func main() {
	httpServer := http.NewServer()
	ub := ubp.UBP{}
	ub.Init()
	disbursementService := ubp.NewDisbursementService(&ub)
	httpServer.DisbursementService = disbursementService
	httpServer.Addr = ":8080"
	// httpServer.Host = "127.0.0.1"

	// Open HTTP server.
	err := httpServer.Open()
	if err != nil {
		fmt.Printf("Error opening server: %s\n", err)
	}
	u := httpServer.URL()
	fmt.Printf("Server listening: %s\n", u.String())
}
