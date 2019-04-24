// +build all travis

package xlm

import (
	"testing"
)

// This test would call the remote / local API and check whether balances match
// if not, this quits immediately
func TestAPIs(t *testing.T) {
	var err error
	height, err := GetTransactionHeight("bea5f00c6327a2d76dbe427c242c5087230191a9c83778b68f3d1fda5a7534a8")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionHeight("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
	if height != 2452 {
		t.Fatalf("Heights don't match, quitting!")
	}
	_, err = GetTransactionData("bea5f00c6327a2d76dbe427c242c5087230191a9c83778b68f3d1fda5a7534a8")
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetTransactionData("blah")
	if err == nil {
		t.Fatalf("Shouldn't work, invalid tx hash")
	}
}
