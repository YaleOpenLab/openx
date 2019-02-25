// +build all

package xlm

import (
  "testing"
  utils "github.com/YaleOpenLab/openx/utils"
)

// go test -run=XXX -tags="all" -bench=.
func BenchmarkGetBlockHash(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _= GetBlockHash(utils.ItoS(i))
  }
}

func BenchmarkGetLedgerData(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _= GetLedgerData(utils.ItoS(i))
  }
}

func BenchmarkGetLatestBlockHash(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _= GetLatestBlockHash()
  }
}

func BenchmarkGetNativeBalance(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetNativeBalance("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
  }
}

func BenchmarkGetAccountData(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetAccountData("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
  }
}

func BenchmarkGetAssetBalance(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetAssetBalance("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34", "YOL77fa301ef")
  }
}

func BenchmarkGetAllBalances(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetAllBalances("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
  }
}

func BenchmarkHasStableCoin(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _ = HasStableCoin("GD5TMZITWOGORE4AVRSESLHAFXAF4YTHJOOJ2CE5RG7RA5WT73QZQURK")
  }
}

// test with an accoutn that does nto have stabelcoin in its account
func BenchmarkHasStableCoin2(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _ = HasStableCoin("GC6Z2KKU4EDTIHAYTJC3Y3AER4ZS5GDSX7S5IKJRRTHLRMJIMCPKQY34")
  }
}

func BenchmarkGetTransactionHeight(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetTransactionHeight("46c04134b95204b82067f8753dce5bf825365ae58753effbfcc9a7cac2e14f65")
  }
}

func BenchmarkGetTransactionData(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetTransactionData("46c04134b95204b82067f8753dce5bf825365ae58753effbfcc9a7cac2e14f65")
  }
}

// test with a fake tx
func BenchmarkGetTransactionData2(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _ = GetTransactionData("blah")
  }
}

func BenchmarkGetKeyPair(b *testing.B) {
  b.ResetTimer()
  for i := 1 ; i < b.N ; i ++ {
    _, _, _ = GetKeyPair()
  }
}
