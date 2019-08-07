// +build all

package xlm

import (
	utils "github.com/Varunram/essentials/utils"
	"testing"
)

// go test -run=XXX -tags="all" -bench=.
func BenchmarkGetBlockHash(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetBlockHash(utils.ItoS(i))
	}
}

func BenchmarkGetLedgerData(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetLedgerData(utils.ItoS(i))
	}
}

func BenchmarkGetLatestBlockHash(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetLatestBlockHash()
	}
}

func BenchmarkGetNativeBalance(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetNativeBalance("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	}
}

func BenchmarkGetAccountData(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetAccountData("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	}
}

func BenchmarkGetAssetBalance(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetAssetBalance("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH", "OXA6fd8ca6bc")
	}
}

func BenchmarkGetAllBalances(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetAllBalances("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	}
}

func BenchmarkHasStableCoin(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = HasStableCoin("GD5TMZITWOGORE4AVRSESLHAFXAF4YTHJOOJ2CE5RG7RA5WT73QZQURK")
	}
}

// test with an accoutn that does not have stabelcoin in its account
func BenchmarkHasStableCoin2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = HasStableCoin("GAMTX6MDG65OFU42WGZH2W73AODEKILBW3IWZYPY5SMIVSDSMYXGTTUH")
	}
}

func BenchmarkGetTransactionHeight(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetTransactionHeight("7c4a995b2cb881618fc3b799d0492d24c38af67f764c5e0c66984a291204a6ad")
	}
}

func BenchmarkGetTransactionData(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetTransactionData("7c4a995b2cb881618fc3b799d0492d24c38af67f764c5e0c66984a291204a6ad")
	}
}

// test with a fake tx
func BenchmarkGetTransactionData2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = GetTransactionData("blah")
	}
}

func BenchmarkGetKeyPair(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _, _ = GetKeyPair()
	}
}
