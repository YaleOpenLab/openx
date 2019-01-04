package assets

import (
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
)

// AssetID assigns a unique assetID to each asset. We assume that there won't be more
// than 68719476736 (16^9) assets that are created at any point, so we're good.
// the total AssetID must be less than 12 characters in length, so we take the first
// three for a human readable identifier and then the last 9 are random hex characaters
// passed through SHA3
func AssetID(inputString string) string {
	// so the assetID right now is a hash of the asset name, concatenated investor public keys and nonces
	x := utils.SHA3hash(inputString)
	return "YOL" + x[64:73] // max length of an asset in stellar is 12
	// log.Fatal(fmt.Errorf("All good"))
	// return nil
}

// CalculatePayback is a TODO function that should simply some up the PBToken
// balance and then return them to the frontend UI for a nice display
func CalculatePayback(project database.Project, amount string) string {
	// the idea is that we should be able ot pass an assetId to this function
	// and it must calculate how much time we have left for payback. For this example
	// until twe do the db stuff, lets pass a few params (although this could be done
	// separately as well).
	// TODO: this functon needs to be the payback function
	amountF := utils.StoF(amount)
	amountPB := (amountF / float64(project.Params.TotalValue)) * float64(project.Params.Years*12)
	amountPBString := utils.FtoS(amountPB)
	return amountPBString
}
