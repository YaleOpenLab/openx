package database

// A contract has six Stages (right now an order has 6 stages and later both will be merged)
// TODO: implement 0.5 stages
// seed funding and seeda ssets are also TODOs, thoguh investors can see the assets
// now and can transfer funds if they really want to
// look into state commitments and committing state in the memo field of transactions
// and then having to propagate one transaction for ever major state change
// Stage 0: Originator approaches the recipient to originate an order
// Stage 0.5: Legal contract between the originator and the recipient
// Stage 1: Originator  proposes a contract on behalf of the recipient
// Stage 1.5: The contract, even though not final, is now open to investors' money
// Stage 2: Contractors propose their contracts and i vnestors can vote on them if they want to
// Stage 3: Recipient chooses a particular contract for finalization
// Stage 4: Review the legal contract and finalzie a particular contractor
// Stage 5: Installation of the panels / houses by the developer and contractor
// Stage 6: Power generation and trigerring automatic payments, cover breach, etc.

// A legal contract should ideally be sotred on ipfs and we must keep track of the
// ipfs hash so that we can retrieve it later when required
type Contract struct {
	O Order
}
