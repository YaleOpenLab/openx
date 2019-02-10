// +build all travis

package issuer

import (
	"testing"

	xlm "github.com/YaleOpenLab/openx/xlm"
)

func TestIssuer(t *testing.T) {
	var err error
	platformSeed, platformPubkey, err := xlm.GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = xlm.GetXLM(platformPubkey)
	if err != nil {
		t.Fatal(err)
	}
	err = InitIssuer(1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = FundIssuer(1, "blah", platformSeed)
	if err != nil {
		t.Fatal(err)
	}
	_, err = FreezeIssuer(1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = DeleteIssuer(1)
	if err != nil {
		t.Fatal(err)
	}
}
