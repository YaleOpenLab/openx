Assets in Stellar are used to keep track of the 'state' of a user's account. They can be thought of as tokens or digital 'receipts'

The entities in the platform are described in the main README file, but this explains how the PaybackAssets, InvestorAssets and DebtAssets work.

1. InvestorAsset - An InvestorAsset is issued by the issuer for every USD that the investor
// has invested in the contract. This peg needs to be ensured maybe in protocol
// with stablecoins on Stellar or we need to provide an easy onboarding scheme
// for users into the crypto world using other means. The investor receives
// InvestorAssets as proof of investment but profit return mechanism is not taken into
// account here, since that needs clear definition on how much investors get each
// period for investing in the project.

2. DebtAsset - for each InvestorAsset (and indirectly, USD invested in the project),
// we issue a DebtAsset to the recipient of the assets so that they can pay us back.
// DebtAssets are also lunked with PaybackAssets and they should be immutable as well,
// so that the issuer can not change the amount of debt at any point in the future.
// MW: Mention that DebtAssets are not equal to InvestorAssets since there must be an interest %
// that needs to be paid to investors, which is also part of the DebtAsset

3. PaybackAsset - each PaybackAsset denotes a month of appropriate payback. A month's worth
// of payback is decided by the recipient, who decides the payback period of the
// given assets at the time of creation. PaybackAssets are non-fungible, it means
// that one project's payback asset is not worth the same as the other project's PaybackAsset.
// the other two assets are fungible - each InvestorAsset is worth +1USD and each DebtAsset
// is worth -1 USD and can be transferred to other peers willing to take profit / debt
// on behalf of the above entities. Since PaybackAsset is not fungible, the flag
// authorization_required needs to be set and a party without a trustline with
// the issuer can not trade in this asset (and ideally, the issuer will not accept
// trustlines in this new asset)
// PaybackAssets in general are not always an arbitrary decision of the recipient,
// rather its set by an agreement of utility or rent payment, tied to the information from
//  an IoT device (i.e a powermeter in the case of solar)
