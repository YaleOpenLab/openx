// +build all travis

package xlm

import (
	"testing"
)

// This test would call the remote / local API and check whether balances match
// if not, this quits immediately
func TestAPIs(t *testing.T) {
	var err error
	height, err := GetTransactionHeight("5454e1594d2a6986b094ddf90302d0d838abab258cbec515da75198161091b83")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionHeight("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
	if height != 4335 {
		t.Fatalf("Heights don't match, quitting!")
	}
	_, err = GetTransactionData("5454e1594d2a6986b094ddf90302d0d838abab258cbec515da75198161091b83")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionData("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
}
