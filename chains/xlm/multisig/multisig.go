package multisig

import (
	"github.com/pkg/errors"
	"log"

	xlm "github.com/Varunram/essentials/crypto/xlm"
	utils "github.com/Varunram/essentials/utils"
	"github.com/stellar/go/keypair"
	build "github.com/stellar/go/txnbuild"
)

// AddSigner is used to add a signer to the pubkey account
func addSigner(seed string, pubkey string, cosignerPubkey string) error {
	memo := "addsigner"
	amount := "1"
	// fun fact: this can be larget than the number of xlm in existence

	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return err
	}

	op1 := build.Payment{
		Destination: pubkey,
		Amount:      amount,
		Asset:       build.NativeAsset{},
	}

	op2 := build.SetOptions{
		Signer: &build.Signer{cosignerPubkey, 1},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op1, &op2},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	_, _, err = xlm.SendTx(mykp, tx)
	if err != nil {
		return errors.Wrap(err, "error while sending tx to horizon")
	}

	return err
}

// constructThresholdTx is used when the number of tx's reaches x-1
func constructThresholdTx(seed string, pubkey string, cosignerPubkey string, y int) error {
	memo := "sealthreshold"
	amount := "1" // some token amount, this can be any number though (even larger than the number of xlm in  xistence)
	x := build.Threshold(y)
	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return err
	}

	op1 := build.Payment{
		Destination: pubkey,
		Amount:      amount,
		Asset:       build.NativeAsset{},
	}

	op2 := build.SetOptions{
		Signer:          &build.Signer{cosignerPubkey, 1},
		MasterWeight:    build.NewThreshold(0),
		LowThreshold:    build.NewThreshold(x),
		MediumThreshold: build.NewThreshold(x),
		HighThreshold:   build.NewThreshold(x),
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op1, &op2},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	_, _, err = xlm.SendTx(mykp, tx)
	if err != nil {
		return errors.Wrap(err, "error while sending tx to horizon")
	}

	return err
}

// Newxofy defines a new x of y multisig contract
func Newxofy(x int, y int, signers ...string) (string, error) {

	if y != len(signers) {
		return "", errors.New("length of multisig tx and number of signers don't match, quitting")
	}

	tempSeed, pubkey, err := xlm.GetKeyPair()
	if err != nil {
		return "", errors.Wrap(err, "error while getting keypair")
	}

	err = xlm.GetXLM(pubkey)
	if err != nil {
		return pubkey, errors.Wrap(err, "error while getting xlm from friendbot")
	}

	for i := 0; i < y-1; i++ {
		err = addSigner(tempSeed, pubkey, signers[i])
		if err != nil {
			return pubkey, errors.Wrap(err, "error whole adding signer to tx")
		}
	}
	// we've reached x-1 = 1 signers, call threshold tx with the x-1'th signer
	err = constructThresholdTx(tempSeed, pubkey, signers[y-1], x)
	if err != nil {
		return pubkey, errors.Wrap(err, "error while constructing threshold tx")
	}

	return pubkey, nil
}

// New1of2 creates a new 1 of 2 multisig
func New1of2(cosigner1Pubkey string, cosigner2Pubkey string) (string, error) {
	return Newxofy(1, 2, cosigner1Pubkey, cosigner2Pubkey)
}

// New2of2 creates a new 2 of 2 multisig
func New2of2(cosigner1Pubkey string, cosigner2Pubkey string) (string, error) {
	return Newxofy(2, 2, cosigner1Pubkey, cosigner2Pubkey)
}

// SendTx broadcasts a multisig tx
func SendTx(txXdr string) (int32, string, error) {
	resp, err := xlm.TestNetClient.SubmitTransactionXDR(txXdr)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not submit tx to horizon")
	}

	log.Printf("Propagated Transaction: %s, sequence: %d\n", resp.Hash, resp.Ledger)
	return resp.Ledger, resp.Hash, nil
}

// Tx2of2 constructs a tx where the source account pubkey1 is the 2of2 account, we need 2 signers for this tx
func Tx2of2(pubkey1 string, destination string, signer1 string, signer2 string, amountx float64, memo string) error {
	sourceAccount, err := xlm.ReturnSourceAccountPubkey(pubkey1)
	if err != nil {
		return errors.Wrap(err, "could not load account details, quitting")
	}

	amount, err := utils.ToString(amountx)
	if err != nil {
		return errors.Wrap(err, "could not convert to float, quitting")
	}

	op := build.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       build.NativeAsset{},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	_, kp1, err := xlm.ReturnSourceAccount(signer1)
	if err != nil {
		return err
	}

	_, kp2, err := xlm.ReturnSourceAccount(signer2)
	if err != nil {
		return err
	}

	txe, err := tx.BuildSignEncode(kp1.(*keypair.Full), kp2.(*keypair.Full))
	if err != nil {
		return errors.Wrap(err, "second party couldn't sign tx")
	}

	_, _, err = SendTx(txe)
	return err
}

// AuthImmutable2of2 sets the auth immutable flag on a multisig account
func AuthImmutable2of2(pubkey1 string, signer1 string, signer2 string) error {
	sourceAccount, err := xlm.ReturnSourceAccountPubkey(pubkey1)
	if err != nil {
		return errors.Wrap(err, "could not load account details, quitting")
	}

	op := build.SetOptions{
		SetFlags: []build.AccountFlag{build.AuthImmutable},
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
	}

	_, kp1, err := xlm.ReturnSourceAccount(signer1)
	if err != nil {
		return err
	}

	_, kp2, err := xlm.ReturnSourceAccount(signer2)
	if err != nil {
		return err
	}

	txe, err := tx.BuildSignEncode(kp1.(*keypair.Full), kp2.(*keypair.Full))
	if err != nil {
		return errors.Wrap(err, "second party couldn't sign tx")
	}

	_, _, err = SendTx(txe)
	return err
}

// TrustAssetTx trusts a specific asset
func TrustAssetTx(assetCode string, assetIssuer string, limit string, pubkey string, signer1 string, signer2 string) error {
	sourceAccount, err := xlm.ReturnSourceAccountPubkey(pubkey)
	if err != nil {
		return errors.Wrap(err, "could not load account details, quitting")
	}

	op := build.ChangeTrust{
		Line:  build.CreditAsset{assetCode, assetIssuer},
		Limit: limit,
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
	}

	_, kp1, err := xlm.ReturnSourceAccount(signer1)
	if err != nil {
		return err
	}

	_, kp2, err := xlm.ReturnSourceAccount(signer2)
	if err != nil {
		return err
	}

	txe, err := tx.BuildSignEncode(kp1.(*keypair.Full), kp2.(*keypair.Full))
	if err != nil {
		return errors.Wrap(err, "second party couldn't sign tx")
	}

	_, _, err = SendTx(txe)
	return err
}

// Convert2of2 converts the account with pubkey myPubkey to a 2of2 multisig account
func Convert2of2(myPubkey string, seed string, cosignerPubkey string) error {
	// account should exist before calling this route
	memo := "testsign"
	amount := "1"

	sourceAccount, mykp, err := xlm.ReturnSourceAccount(seed)
	if err != nil {
		return err
	}

	op1 := build.Payment{
		Destination: myPubkey,
		Amount:      amount,
		Asset:       build.NativeAsset{},
	}

	op2 := build.SetOptions{
		Signer:          &build.Signer{cosignerPubkey, 1},
		MasterWeight:    build.NewThreshold(1),
		LowThreshold:    build.NewThreshold(2),
		MediumThreshold: build.NewThreshold(2),
		HighThreshold:   build.NewThreshold(2),
	}

	tx := build.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []build.Operation{&op1, &op2},
		Timebounds:    build.NewInfiniteTimeout(),
		Network:       xlm.Passphrase,
		Memo:          build.Memo(build.MemoText(memo)),
	}

	log.Println("CHECK THIS: ", mykp, tx)
	_, _, err = xlm.SendTx(mykp, tx)
	if err != nil {
		return errors.Wrap(err, "error while sending tx to horizon")
	}

	return nil
}
