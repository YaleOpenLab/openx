package xlm

// println.go pretty prints the balances retrieved by calling the stellar testnet
// API
import (
	"fmt"

	"github.com/stellar/go/protocols/horizon"
)

/*
	type Balance struct {
	    Balance            string `json:"balance"`
	    Limit              string `json:"limit,omitempty"`
	    BuyingLiabilities  string `json:"buying_liabilities"`
	    SellingLiabilities string `json:"selling_liabilities"`
	    base.Asset
	}
*/
func PrettyPrintBalances(balances []horizon.Balance) {
	fmt.Println("   LIST OF ALL YOUR BALANCES: ")
	for _, balance := range balances {
		if balance.Asset.Code == "" {
			fmt.Printf("    ASSET CODE: XLM, ASSET BALANCE: %s\n", balance.Balance)
			continue
		}
		fmt.Printf("    ASSET CODE: %s, ASSET BALANCE: %s\n", balance.Asset.Code, balance.Balance)
	}
}
