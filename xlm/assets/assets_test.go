// +build all

package assets

import (
	"testing"
)

func TestHelpers(t *testing.T) {
	test := AssetID("cool")
	if test != "OXA12a383c1f" {
		t.Fatalf("Asset ID not functioning as expected, quitting!")
	}
}
