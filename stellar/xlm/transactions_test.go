package xlm

import (
	"testing"
)

// This test would call the remote / local API and check whether balances match
// if not, this quits immediately
func TestAPIs(t *testing.T) {
	var err error
	height, err := GetTransactionHeight("46c04134b95204b82067f8753dce5bf825365ae58753effbfcc9a7cac2e14f65")
	if err != nil {
		t.Fatal(err)
	}
	if height != 1278685 {
		t.Fatalf("Heights don't match, quitting!")
	}
}
