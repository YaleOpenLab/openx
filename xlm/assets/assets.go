package assets

import (
	//"log"

	utils "github.com/YaleOpenLab/openx/utils"
	xlm "github.com/YaleOpenLab/openx/xlm"
	"github.com/stellar/go/network"
	build "github.com/stellar/go/txnbuild"
)

// Package assets contains asset related functions like calculating AssetID and
// sets up DebtAssets, PaybackAssets and InvestorAssets for a specific project that it has been
// passed
// The entities in the system are described in the README file and this part
// will explain how the PaybackAssets, InvestorAssets and DebtAssets work.

// 1. InvestorAsset - An InvestorAsset is issued by the issuer for every USD that the investor
// has invested in the contract. This peg needs to be ensured maybe in protocol
// with stablecoins on Stellar or we need to provide an easy onboarding scheme
// for users into the crypto world using other means. The investor receives
// InvestorAssets as proof of investment but profit return mechanism is not taken into
// account here, since that needs clear definition on how much investors get each
// period for investing in the project.

// 2. DebtAsset - for each InvestorAsset (and indirectly, USD invested in the project),
// we issue a DebtAsset to the recipient of the assets so that they can pay us back.
// DebtAssets are also lunked with PaybackAssets and they should be immutable as well,
// so that the issuer can not change the amount of debt at any point in the future.
// MW: Mention that DebtAssets are not equal to InvestorAssets since there must be an interest %
// that needs to be paid to investors, which is also part of the DebtAsset

// 3. PaybackAsset - each PaybackAsset denotes a month of appropriate payback. A month's worth
// of payback is decided by the recipient, who decides the payback period of the
// given assets at the time of creation. PaybackAssets are non-fungible, it means
// that one project's payback asset is not worth the same as the other project's PaybackAsset.
// the other two assets are fungible - each InvestorAsset is worth +1USD and each DebtAsset
// is worth -1 USD and can be transferred to other peers willing to take profit / debt
// on behalf of the above entities. SInce PaybackAsset is not fungible, the flag
// authorization_required needs to be set and a party without a trustline with
// the issuer can not trade in this asset (and ideally, the issuer will not accept
// trustlines in this new asset)
// PaybackAssets in general are not always an arbitrary decision of the recipient,
// rather its set by an agreement of utility or rent payment, tied to the information from
//  an IoT device (i.e a powermeter in the case of solar).

// The hard part is ensuring that the assets are pegged to the USD in a stable way.
// we could ensure the peg ourselves by accepting USD off chain, but that's not provable
// on chain and the investor has to trust the issuer with that. Also, in this case,
// anonymous investors wouldn't be able to invest, which is something that would be
// nice to have

// AssetID assigns a unique assetID to each asset. We assume that there won't be more
// than 68719476736 (16^9) assets that are created at any point, so we're good.
// the total AssetID must be less than 12 characters in length, so we take the first
// three for a human readable identifier and then the last 9 are random hex characaters
// passed through SHA3
func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "OXA" + x[64:73] // max length of an asset in stellar is 12 (OXA: OpenX Asset)
}

// CreateAsset creates a new asset belonging to the public key referenced above
func CreateAsset(assetName string, PublicKey string) build.Asset {
	// need to set a couple flags here
	return build.CreditAsset{assetName, PublicKey}
}

// TrustAsset trusts a specific asset issued by a particular public key and signs
// a transaction with a preset limit on how much it is willing to trsut that issuer's
// asset for
func TrustAsset(assetCode string, assetIssuer string, limit string, seed string) (string, error) {
	// TRUST is FROM Seed TO assetIssuer
	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return "", err
	}

	op := build.AllowTrust{
		Trustor:   mykp.Address(),
		Type:      build.CreditAsset{assetCode, assetIssuer},
		Authorize: true,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	_, txHash, err := xlm.SendTx(mykp, tx)
	if err != nil {
		return txHash, err
	}

	op2 := build.ChangeTrust{
		Line:  build.CreditAsset{assetCode, assetIssuer},
		Limit: limit,
	}

	tx = build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op2},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	_, txHash, err = xlm.SendTx(mykp, tx)
	return txHash, err
}

// SendAssetFromIssuer transfers _amount_ number of assets from the caller to the destination
// and returns an error if the destination doesn't have a trustline with the issuer
// This method is called by the issuer of the asset
func SendAssetFromIssuer(assetCode string, destination string, amount string,
	seed string, issuerPubkey string) (int32, string, error) {

	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, issuerPubkey},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	return xlm.SendTx(mykp, tx)
}

// SendAssetToIssuer sends a specific asset back to the issuer
func SendAssetToIssuer(assetCode string, destination string, amount string,
	seed string) (int32, string, error) {

	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, destination},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
	}

	return xlm.SendTx(mykp, tx)
}

// SendAsset sends the asset to a destination which already has an established trustline with the issuer
func SendAsset(assetCode string, issuerPubkey string, destination string, amount string,
	seed string, memo string) (int32, string, error) {
	passphrase := network.TestNetworkPassphrase
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return -1, "", err
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.CreditAsset{assetCode, issuerPubkey},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	return xlm.SendTx(mykp, tx)
}
