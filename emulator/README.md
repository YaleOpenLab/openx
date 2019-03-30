# Emulator

The emulator is meant to offer a CLI interface to investors and recipients who desire it. It encapsulates most of the critical functions of the platform and can be used to test platforms before they make it to production. One can also aim to test the rudimentary workflow on the frontend with the help of the emulator because it has roles and requires user interaction.

Broadly, the emulator has 4 different modes geared at 4 entities right now (2 in the case of the op.zones platform and 4 in the case of the opensolar platform). Some  functions are common to all entities (like viewing which projects are open, displaying balance and similar) but other functions are exclusive to a role (for eg, a recipient can not invest in an order). While this is limited right now, there is no reason this will be expanded in the future and a recipient might be allowed to invest a portion in his own order or a developer might be allowed to invest in a contract he originated.

## Common functions to all entitites

- Ping: pings the platform to see if its still up
- Exchange: exchanges xlm for STABLEUSD as per the rate given by the oracle.
- Ipfs: hashes the string and stores it in ipfs
- Send: sends a speciifc amount of coins to a destination
- Receive: receives an asset or xlm from a remote destination
- Create: creates a local asset that is shareable among peers
- Kyc: Shows users who haven't yet received KYC approval if the user is a KYC inspector

## Investor only functions

- Vote: vote towards a particular stage 2 contract
- Invest: invest in a particular order

## Recipient only functions

- Unlock: unlocks a specific project that was invested in (acts as confirmation that the recipient is willing to accept the investment)
- Payback: pays back towards a specific project that the recipient has already accepted in the past
- Finalize: finalizes a particular project ie moves it from stage 2 to stage 3
- Originate: Originates a particular project ie brings the project from stage 0 to stage 1)
- Calculate: calculates how much the recipient owes the platform during a specific payback period

## Contractor only functions

- Propose: propose a new stage 2 contract that will be voted on by investors and selected by recipients
- MyProposed: displays a list of all the contracts proposed by this specific entity
- AddCollateral: adds a specific amount of collateral towards proposing a contract
- mystage0: view a list of all the stage0 contracts proposed by this entity
- mystage1: view a list of all the stage1 contracts proposed by this entity

## Originator only functions

- Propose: propose a new stage 2 contract that will be voted on by investors and selected by recipients
- PreOriginate: propose a new stage 0 contract that has to be taken to the recipient for upgradation to a stage 1 contract
- AddCollateral: adds a specific amount of collateral towards proposing a contract
- MyProposed: displays a list of all the contracts proposed by this specific entity
- mystage0: view a list of all the stage0 contracts proposed by this entity
- mystage1: view a list of all the stage1 contracts proposed by this entity


In the above list of commands, some are placeholders since they require user action such as uploading pdfs of signed legal contracts, kyc documents and similar. These would also be abstracted in the future if desired by the participants in the platform ecosystem. Also, since a person can create and exchange local assets, this can give rise to a token ecosystem where you (can) have multiple parties interacting with various tokens and continuing to vote on measures like on and off-chain governance.
