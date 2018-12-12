# Stellar PoC

This sub-repo contains the WIP Stellar Proof of Concept. This has some design
tradeoffs compared to the Ethereum implementation, especially in the concept of
"state" in Ethereum. This PoC takes the three major (actually two, but we implement three)
state variables that are in Ethereum and creates assets on Stellar for each
implementation:

We define roughly three entities in the process of investing in a project:
1. ISSUER - the issuer is the server on which the orders are advertised on
2. INVESTOR - the investor is a person / group of persons who invest in a particular order
3. RECIPIENT - the recipient is the school which receives funding from investors through an order

These three entities interact use the three tokens detailed below:

1. INVTokens - INVTokens are sent from the ISSUER to the INVESTOR one time, as proof of investment
2. DEBTokens - DEBTokens are sent from the ISSUER to the RECIPIENT the time the order is created and then each month, the RECIPIENT pays back the ISSUER each month with a number of DEBTokens like an electricity bill. (This has to be capped at a minimum of the electricity bill in the future, along with an oracle that can attest to the price of electricity at the place)
3. PBTokens - PBTokens are sent from the ISSUER to the RECIPIENT at multiple intervals:
  - At the time of order confirmation, 1 each for each month the school opts in to pay the amount back. eg. if the school opts in for a period of 5 years, the ISSUER issues 60 PBTokens. We also calculate expected_paid_amount as `amount_invested /  payback_amount`
  - Each time the ISSUER receives DEBTokens, the ISSUER confirms the transaction and sends back PBTokens as `paid_amount / expected_paid_amount` eg you pay 210 towards an expected amount of 120. You get paid back 1.75 PBTokens.

Any disparity / failure on the ISSUER's part can be argued with, since the transactions are on chain. The PBToken simplifies this, since one doesn't need to go back in history and calculate how much a given school has paid a person. Percentage paid is simply `PBToken_balance / PBToken_total`.
