# OpenSolar

The OpenSolar Project aims to use blockchain and IoT based smart contracts to help finance community solar projects on the stellar blockchain. This idea was [presented at the MIT Media Lab's Fall 2018 Members Week Demo](https://github.com/YaleOpenLab/smartPropertyMVP) using the Ethereum blockchain. Stellar has some design tradeoffs compared to Ethereum, especially with regard to the concept of "state" in Ethereum.

Stellar doesn't support a Turing complete VM like Ehtereum nor a stack based system like Bitcoin. Instead, it provides a standard set of operations which are decided on the protocol level and offers SDKs to build upon the platform. While limited in function, stellar possesses powerful fundamentals which could be used to partially replicate the notion of st ate in Ethereum.

One of the core features of stellar is the creation of [Assets](https://www.stellar.org/developers/guides/concepts/assets.html). An Asset is anything that can be issued by a party to track ownership or provide proofs of investment / debt, etc. In this prototype, we use these Assets to track multiple things:

1. Investment in a particular asset: When a person invests say 2000 USD in a project, he needs some sort of irreversible proof that he has invested in that particular project. This would ideally be just stored in a variable if in Ethereum, but here, we issue an asset called "Investor Asset" once an investor invests in a particular project. This investor asset is issued 1:1 with the amount invested, meaning an investor who invests 2000 USD like above would get 2000 Investor Assets in return. The Investor Asset is pseudo-unique (more on this to follow) and can be tracked with the help of the 12 character identifier that stellar provides us with.

NOTE: Stellar limits the character limit of each asset to 12 characters, so the identifier is not unique. In the case collisions arise, the project ID, Debt Assets and time of creation can be used to identify the asset. There seems to be no workaround for this limit, so we are forced to go ahead with this scheme.

2. When a project has reached its target goal in USD, we need to assure investors that their amount was invested in this project. For this, we issue "Recipient Assets" which denote that a particular group of investors have invested in the given project. Like the investor asset, this is 1:1 with the amount invested.

3. Once a project has been installed and can generate electricity, the recipient starts to pay back towards the project. After confirmation of each payment, we issue a payback asset, which is proportional to the monthly payment bill to provide ease of accounting and quick look back on whether the recipient is not defaulting on its payments.

For more notions of state and ownership, we can continue to issue relevant assets which would track ownership and history.

Apart from ownership, the assets above serve other functions  that are useful:

1. Investor Assets are tradable: Since the investor asset is a proof of investment in a particular project, we can trade them for other investor assets like traditional property markets or we could use them to trade with USD / take a loan against this asset similar to a secondary mortgage market.

2. Parties that are willing to donate to a particular recipient can choose to payback their electricity bill on their behalf or choose to buy some of their Recipient Assets in order to hedge some risk on behalf of them. This is useful for big charity organizations which invest in multiple projects, who can keep track of their charity donations in an easy way and provide publicly auditable proof of their donation towards a charity.

While dealing with real world entities, we are inevitably faced with dealing with legal contracts. The platform should not worry about what's in the contract as long as it has been agreed to and vetted by both parties, so we store the contract in ipfs and commit the resulting hash in two split stellar transactions' memo fields. The memo field of stellar can only hold 28 characters, so we split the 46 character ipfs hash into two parts and pad the second hash with characters to denote that it is an ipfs hash and not some garbage value. This ipfs hash has to be checked on all parties' ends to ensure that this is the same contract that they agreed to earlier

Each solar system deployment is defined as a "Project" and a given project right now has stages till 7, with decimal values denoting smaller increments regarding the project's stage:

 - Stage 0: A pre origin contract, which is given by an originator to a recipient
 - Stage 0.5: A MOU between the originator and the recipient regarding the originator's participation in the project.
 - Stage 1: An originated contract, which is the evolution of a pre origin contract at the behest of the recipient, who may require a legal contract to be signed by the originator for doing so
 - Stage 2: A proposed contract, which is a contract given by a contractor who is willing to lend his services for the particular project at a fee that will be paid by the recipient and in turn, the investor as well
 - Stage 3: A finalised contract, which is open for investor funds. A contract can be promoted from stage 2 to 3 only by a recipient, who chooses the best contract that fits in with his needs and requirements
 - Stage 4: A contract that has raised the money required for investment from investors and is ready to be installed. The control of the contract now passes over from the recipient to the developer, who is responsible to take this project to stage 5.
 - Stage 5: the stage where the developer should install all the required devices and the contractor must do his job of supervising.
 - Stage 6: Power Generation s tarts and the recipient starts to payback the investor in monthly instalments similar to an electricity bill
 - Stage 7: The recipient has finished paying back the investors and can now own the solar panels installed in his space.

There are various entities defined (and more on the way) that emulate different functions that would be performed by entities in the real world:
 - ISSUER - the issuer is the server on which the orders are advertised on
 - RECIPIENT - the recipient is the school which receives funding from investors through an order
 - ORIGINATOR - the person who proposes a pre-origin contract to a recipient and requires approval from the recipient to originate the contract
 - CONTRACTOR - takes in a originated order and proposes a contract which fits in with the requirements described in the originated contract
 - INVESTOR - the investor is a person / group of persons who invest in a particular finalised order
 - DEVELOPER - the person who is responsible for installing the hardware and making sure the solar panels work and produce electricity
 - GUARANTOR - the person who is liable to make the recipient pay or pays itself in the case of a delayed payment from the recipient's side

A rough path taken by a specific project would be:
1. Originator(s) approaches recipient with an idea for using a space owned by the recipient to start a new project.
2. Recipient reviews all received proposals and chooses a specific one. An MOU is commonly agreed upon by both the originator and the recipient.
3. Recipient promotes the project from stage 0 to stage 1, hence opening the project for proposals from contractors. If seed funding is desired, he promotes the contract to stage 1.5, asking investors to invest in the seed stage.
4. Contractors view all originated projects and put out proposals for working with specific projects. A given project may have multiple contractors performing different roles in the system.
5. At this stage, investors can see all proposed projects and van vote on them, giving reviewing recipients a chance at which proposal would get the most funding.
6. Recipient views all the proposed contracts by various contractors and chooses a specific one to open to investor funding. It promotes the project from stage 2 to stage 3.
7. Investors view all contracts open to funding and choose to invest in them. They get Investor Assets in return.
8. Once the funding goal is reached, the platform hands out Recipient Assets and Payback Assets to the Recipient and the Contractor is expected to start working as per previously proposed contract. The project's stage is upgraded from stage 3 to stage 4.
9. Once installation is complete, the contractor promotes the project from stage 4 to stage 5.
10. Power generation is tested. If the installation works fine, the project's stage is upgraded from stage 5 to stage 6 and the recipient pays back an amount based on electricity consumption with inputs from the energy oracle.
11. Once the paid amount matches that agreed upon in the contract, ownership of the installed system is granted to the recipient. The Project's stage is upgraded from stage 6 to stage 7.

## Running the code in this repo

In order to be able to run this, you need to have the latest version of go installed. [Here](https://medium.com/@patdhlk/how-to-install-go-1-9-1-on-ubuntu-16-04-ee64c073cd79) is a quick tutorial on how to get go installed on a Linux / macOS machine.

Once you have go installed, you need to get the packages in this repo. This can be done using `go get -v ./...`. The `stellar/go` package might print out errors due to its problems with `go get`, in which case you need to get the package separately and then run `dep ensure -v` on the project, proceeded by `go get`ing the other packages as normal.

Then you need to build the stellar package `go build` and then run the executable like `./stellar`.

You will be faced with a CLI interface with which you can interact with

### Installing IPFS

ipfs is used by some parts of the program to store legal contracts, files that the user might want to store permanently. Download a release from https://github.com/ipfs/go-ipfs/releases and run install.sh. In case you face an issue with migration between various ipfs versions, you might need to run [fs-repo-migrations](https://github.com/ipfs/fs-repo-migrations/blob/master/run.md) to migrate to a newer version. If you don't have anything valuable, you can delete the directory and run `ipfs init` again (this will delete the data stored in ipfs prior to deletion)

You need to keep your peer key (`ipfs.key` usually) in a safe place for future reference. Start ipfs using `ipfs daemon` and you can test if it worked by creating a test file `test.txt` and run `ipfs add test.txt` to see if it succeeds. The resultant hash can be decrypted using `curl "http://127.0.0.1:8080/ipfs/hash"` where 8080 is the endpoint of the ipfs server or by doing `cat /ipfs/hash` directly. You can also refer to [this helpful tutorial](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) in order to get easily started with ipfs.

### Running tests

Running tests is mostly simple with `go test` but the tests have flags since some require running other software in the background (`ipfs`). There are two kinds of flags right now - `travis` and `all`. If you need the coverage stats as well, you need to
```
go get golang.org/x/tools/cmd/cover
```
if you already don't have the package. Then running `go test --tags="all" -coverprofile=test.txt ./...` should run all the tests and provide coverage data. Running with the tag `travis` will omit the tests in `ipfs/` which requires [a local `go-ipfs` node running](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) as described above.
