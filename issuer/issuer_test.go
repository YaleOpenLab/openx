// +build all travis

package issuer

import (
	"testing"

	xlm "github.com/YaleOpenLab/openx/xlm"
	consts "github.com/YaleOpenLab/openx/consts"
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
	err = InitIssuer(consts.OpenSolarIssuerDir, 1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = FundIssuer(consts.OpenSolarIssuerDir, 1, "blah", platformSeed)
	if err != nil {
		t.Fatal(err)
	}
	err = FundIssuer(consts.OpenSolarIssuerDir, 1, "cool", platformSeed)
	if err == nil {
		t.Fatalf("not able to catch invalid seed error, quitting!")
	}
	_, err = FreezeIssuer(consts.OpenSolarIssuerDir, 1, "blah")
	if err != nil {
		t.Fatal(err)
	}
	err = DeleteIssuer(consts.OpenSolarIssuerDir, 1)
	if err != nil {
		t.Fatal(err)
	}
}
