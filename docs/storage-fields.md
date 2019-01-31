# What is stored on the platform?

This is a comprehensive description of the various things that the database holds at this moment. We shall be focusing on the opensolar platform
since openfinance borrows much stuff from opensolar. We shall also not go into how they are structured on the backend
since that will be covered in a later architecture document / diagram.

## The Solar Project

The Project structure is huge and contains multiple fields relevant to each participant in the system.

-   Index int
    -   An index to keep quick track of how many projects exist
-   PanelSize
    -   Panel Size contains the size and number of solar panels installed at a given location
-   TotalValue  
    -   TotalValue is the total amount that we seek ask from investors
-   Location   
    -   Location is the location of the given installation
-   MoneyRaised
    -   The amount of money that has been raised by investors so far (seed + funding rounds)
-   Years       
    -   The average number of years that it would take to repay the give debt. This would be shown to both recipients and investors so they can have a clear idea on how much they need to pay
-   Metadata    
    -   Any additional identifying metadata such as brand of the panel, how its constructed, etc to provide randomness for generating an AssetID
-   InvestorAssetCode
    -   The code of the asset that the investor receives
-   DebtAssetCode    
    -   The code of the asset that the receiver receives
-   PaybackAssetCode
    -   An optional asset that is assigned for ease of accountability
-   SeedAssetCode    
    -   The asset that is given to seed investors who invest in round 1.5
-   BalLeft
    -   The amount that the recipient needs to payback in order for them to own the given asset that investors have invested in
-   Votes   
    -   The number of votes that that investors have piled on a given proposed (stage 2) contract
-   DateInitiated
    -   The date that the project was confirmed by the recipient (stage 1)
-   DateFunded    
    -   The date that the given project was funded
-   DateLastPaid  
    -   The date that the recipient last paid the platform
-   Originator    
    -   The originator of this particular contract.
-   OriginatorFee
    -   The fee that the originator charges for his service. This is included in the TotalValue field
-   Developer     
    -   The developer behind this particular project.
-   DeveloperFee  
    -   The fee that the originator charges for his service. This is included in the TotalValue field
-   Guarantor     
    -   The guarantor behind this particular project
-   Contractor             
    -   The contractor behind this particular project
-   ContractorFee          
    -   The fee that the contractor charges for his service. This is included in the TotalValue field
-   SecondaryContractor    
    -   The secondary contractor behind this particular project
-   SecondaryContractorFee
    -   The fee that the secondary contractor charges for his service. This is included in the TotalValue field
-   TertiaryContractor     
    -   The tertiary contractor behind this particular project
-   TertiaryContractorFee  
    -   The fee that the tertiary contractor charges for his service. This is included in the TotalValue field
-   ProjectRecipient
    -   The recipient of this particular project
-   ProjectInvestors
    -   The investors who have invested in this particular project
-   SeedInvestors    
    -   The investors who have invested in this particular project in the seed round
-   Stage       
    -   The stage that this particular project is at
-   AuctionType
    -   The type of auction that will be used to determine whether a proposed contract goes towards seeking funding)
-   OriginatorMoUHash       
    -   The ipfs hash of the contract between the originator and the recipient
-   ContractorContractHash  
    -   The ipfs hash of the contract between the contractor and the recipient
-   InvPlatformContractHash
    -   The ipfs hash of the contract between the platform and the investor
-   RecPlatformContractHash
    -   The ipfs hash of the contract between the platform and the recipient
-   Reputation float64
    -   The net reputation of the given user based on feedback from other users invoved in past contracts.
-   Lock    
    -   Once a project is funded, we need the recipient to accept the project and received the debt assets
-   LockPwd
    -   The seed password that is used to lock the contract by the recipient

## The User

The User is a primordial entity on the platform (think dinosaurs). They are the basis for every other entity on the platform and can be thought of as the greatest intersection of all persons on the platform.

-   Index int
    -   index is used for indexing people on the platform
-   EncryptedSeed
    -   The AES256 encrypted seed of the user. Even if the platform is hacked and we lose control of the server, nobody can steal funds because the seeds are encrypted
-   Name string
    -   The Name of the user
-   PublicKey string
    -   The Public Key of the given user
-   LoginUserName string
    -   The username that the user uses to sign in into the platform
-   LoginPassword string
    -   The sha3 hash of the password that the user uses to sign in into the platform. This is different from the seed password, which encrypts the encryptedseed.
-   Address string
    -   The address of the user. Need this for KYC stuff I presume, so added it in here
-   Description string
    -   A short description of the user, can be thought of like a twitter status / about me page
-   Image string
    -   A photo of the person. Optional, but can be added if people like this stuff.
-   FirstSignedUp string
    -   The timestamp of when the user first signed up on the platform
-   Kyc bool
    -   Whether or not the user has passed Kyc. Defaults to false.
-   Inspector bool
    -   Inspector is an authenticated kyc entity that can verify other people on the platform
-   Email string
    -   The email id of the person. Used for sending "forgot password" emails, notifications, etc.
-   Notification bool
    -   A toggle-able option which when set to true sends out notification of activities on the platform to the user.
-   Reputation float64
    -   The reputation of the given user. Reputation increases with good feedback given on the user by other parties in past contracts.

### The Investor

The investor contains all the fields of the user followed by a few specific ones that solely belong to the investor.

-   VotingBalance
    -   the voting balance of the given user. This is 1:1 with the user's stablecoin balance
-   AmountInvested
    -   the amount that the investor has invested so far on the platform
-   InvestedSolarProjects
    -   the set of projects that the investor has invested in
-   InvestedBonds  
    -   the set of bonds that the investor has invested in
-   InvestedCoops         
    -   the set of coops that the investor has invested in

### The Recipient

-   ReceivedSolarProjects
    -   The set of solar projects that the recipient has received
-   DeviceId
    -   The device id of the teller (currently only supports one teller, support for multiple tellers to be added later)
-   DeviceStarts
    -   A log of all times that the device has started in the past
-   DeviceLocation
    -   A log of all the locations the device has been in when the device has (re)started

### The Entity

The entity is a collection of similar entities which perform the same function. IN the future, they can be expanded into different structures but for now are better clumped together in this single structure.

-   Contractor
    -   A contractor is a boolean value which is set to true if the given entity is a contractor
-   Developer
    -   A developer is a boolean value which is set to true if the given entity is a developer
-   Originator
    -   An originator is a boolean value which is set to true if the given entity is an originator
-   Guarantor
    -   A Guarantor is a boolean value which is set to true if the given entity is a guarantor
-   PastContracts
    -   A list of all the contracts that the entity was part of. This is common to all the four sub-entities
-   ProposedContracts
    -   A list of all the contracts that the entity has proposed in the past. This is specific to the Contractor sub-entity.
-   PresentContracts
    -   A list of all the contracts that the entity is currently a part of. This is common to all the four sub-entities
-   PastFeedback
    -   A collection of feedback given on this entity by other entities regarding work done in the past.
-   Collateral
    -   The amount collateral that is put up as guarantee by a Contractor while proposing a project. This is specific to the Contractor sub-entity.
-   CollateralData
    -   The description of the collateral that is put up as guarantee by a Contractor while proposing a project. This is specific to the Contractor sub-entity.

### NOTE

All these fields are editable and shouldn't be taken as sufficient to describe the entities in question. This doc might also be useful in the future under the privacy FAQ section where we can be transparent about what fields we store.
