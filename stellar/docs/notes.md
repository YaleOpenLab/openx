# General Notes

## Changes suggested in call 1:
1. See bank anchors and if possible, use their tethered assets for INVTokens instead of doing one ourselves (docs/anchors.md)
2. Split the PBToken idea into its own model and develop a model without the payback token in question
3. Add new entities - developers, contractors,  originators and provisions for them to do their assigned roles and duties.
4. Assume two models overall -
 - All banks are anchors on stellar - we can trade in tethered assets and will not face a problem technically
 - Banks involved with entities are not anchors and we need to come up with solutions to accommodate them in our model
5. Separate the role of an issuer and a platform. An issuer would be someone like neighbourly or swytch whereas the platform is something that we'd be developing upon

Observations:
 - Might be difficult to get banks as anchors since there seem to be only 4 popular ones around

## TODOs copied from the Ethereum Smart Contract

### TODO SOLAR ENGINEERING DOCUMENTS
The developers must present the blueprints of the proposed work before working on it. These must be stored in IPFS and hashed to the contract.
These documents will have to get updated at the end of the install and are saved as the blueprint of the 'installed system'

### TODO VERIFICATION OF INSTALMENT & SENSORS
The signature of a 3rd party verifier should be considered as part of onboarding a system and its data.
The verifier confirms system was built according to plan, is compliant with regulation, and has the appropriate working sensors and associated public keys.
This digital signature will allow the sensor data and oracle to commit payment transactions and REC minting.

### TODO CONTRACTOR PAYMENTS
Allow contractor to collect payout on system, once its confirmed

### TODO Consider payouts to be based on instalments throughout the buildout process and with energy generation data
Eg. the contractors receives an upfront payment before the process begins but once the deployment is confirmed, and receives another payment once the witness sensors shows the first generation data. Final payment to contractor should occur once the system shows steady data and behavior after a selected period of time. It also requires a 3rd party verification of system engineering and metering sensors.

### TODO Renewable Energy certification

### TODO Define the payment cycle
Consider making payments every two weeks or every month

### TODO LEGAL OWNERSHIP
Come up with a step so that when ownership is fully transferred, there is an automatic report that can change a registry that has legal validity
functions that require no gas, but will check the state of the system

Returns true if consumer has completely paid off any outstanding balance on the panel at ssAddress within their consumerBuffer period

### TODO BREACH SCENARIOS
Add all the situations and scenarios that need to be considered if the payments are not made or the wallet accounts have insufficient funds
Consider sending email notifications, bringing in the guarantors, activating hardware etc.
