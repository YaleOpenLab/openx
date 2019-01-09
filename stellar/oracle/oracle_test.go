// +build all travis

package oracle

// this test function actually does nothing, since the oracle itself is a placeholder
// until we arrive at consensus on how it should be structured and stuff
import (
	"log"
	"testing"
)

func TestOracle(t *testing.T) {
	var err error
	bill := MonthlyBill()
	if err != nil {
		t.Fatal(err)
	}
	if bill != "120.000000" {
		log.Println(bill)
		t.Fatalf("Oracle does not output constant value")
	}
	billF := MonthlyBillInFloat()
	if err != nil {
		t.Fatal(err)
	}
	if billF != 120.0 {
		t.Fatalf("Oracle does not output constant value")
	}
	exchangeFloat := ExchangeXLMforUSD("1")
	if err != nil {
		t.Fatal(err)
	}
	if exchangeFloat != 99990 {
		log.Println("EXCHHFLO: ", exchangeFloat)
		t.Fatalf("Exchange value does not match")
	}
}
