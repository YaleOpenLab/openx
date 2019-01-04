package main

// test.go runs the PoC stellar implementation calling various functions
import (
	"fmt"
	"log"
	"os"

	assets "github.com/YaleOpenLab/smartPropertyMVP/stellar/assets"
	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	ipfs "github.com/YaleOpenLab/smartPropertyMVP/stellar/ipfs"
	rpc "github.com/YaleOpenLab/smartPropertyMVP/stellar/rpc"
	stablecoin "github.com/YaleOpenLab/smartPropertyMVP/stellar/stablecoin"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	// TODO: define default values for each and then use them if not passed
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	// InvAmount int    `short:"i" description:"Desired investment" required:"true"`
	InvAmount int `short:"i" description:"Desired investment"`
	// RecYears  int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years." required:"true"`
	RecYears int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years."`
	Port     string `short:"p" description:"The port on which the server runs on"`
}

func main() {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	/*
		fmt.Println("WHICH PLATFORM WOULD YOU IKE TO ENTER INTO?")
		fmt.Println("1. Platform of Contracts")
		fmt.Println("1. Platform of Platforms")
		platformArg, err := utils.ScanForInt()
		if err != nil {
			log.Fatal(err)
		}
		switch platformArg {
		case 1:
			fmt.Println("WELCOME TO THE PLATFORM OF CONTRACTS IDEA")
			// the platform of contracts idea is the idea of having an open platform with
			// various stakeholders and optimizing their game theoretic objectives. Specific
			// conttracts coud emulate the function of platforms themselves but in this case,
			// we handle all aspects of game theory within what we have and ecnourage people
			// to tak part in the system.
			// TODO: think of how to implement this in a nice way
		case 2:
			fmt.Println("WELCOME TO THE PLATFORM OF PLATFORMS IDEA")
			// the platform of platforms idea is similar to what we want in the opensolar
			// project where we can implement various partners as entities in the system
			// TODO: emulate different partners in the system that we have now.
		}
	*/
	// the memo field in stellar is particularly interesting for us wrt preserving state
	// as shown in the preliminary pdf example. We need to test out whether we can commit
	// more state in the memo field and find a way for third party clients to be able
	// to replicate this state so that this is verifiable and not just an entry of state
	// in the system.
	// TODO: think of a reasonable way to hash the current state of the system with
	// its hash in the memo field
	// Open the database

	database.CreateHomeDir()
	platformPublicKey, platformSeed, err := StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s", platformSeed, platformPublicKey)
	// TODO: how much do we pay the investor? how does it work
	// Do we sell the REC created from the solar panels only to the investor? If so,
	// isn't that enough to propel investment in the solar contract itself?
	// TODO: need a server to run a public stellar node to test out stuff
	// change the API mapping
	// move current number of years metric to a separate package since that is
	// more suitable for a model like affordable housing.
	// look into what kind of data we get from the pi and checkout pi specific code
	// to see if we can get something from there.

	// TODO: so the idea would be to split the current PoC into two parts, one
	// focused on the platform of platform ideas and one on the platform of contracts
	// ideas. The current implementation that we have is focused more on the platform
	// of platforms idea, with ideas to integrate various partners at differnet stages
	// to use their input in some parts, but we could do it in an automated way with
	// assets and tokens for everything as well, which is similar to the platform
	// of contracts idea.
	// TODO: how do we emulate various partners? need to get input / have some stuff
	// that we presume they do and then fill in the rest.
	// also need to implement stages in contracts based on finalization
	// need to integrate ipfs into the workflow so that we can store copies of the
	// specific contract and then reference it when needed using the ipfs hash.
	// need to start working on the base contract that connects investors and
	// recipients and the investors and the platform. Need to transition automatically
	// and also cover breach scenarios in case the recipient doesn't pay for a specific
	// period of time
	// TODO: upgrade the RPC and tests to fit in with recent changes
	// TODO: also write a Makefile so that its easy for people to get started with stuff
	fmt.Println("------------STELLAR HOUSE INVESTMENT CLI INTERFACE------------")

	// init stablecoin stuff
	err = stablecoin.InitStableCoin()
	if err != nil {
		log.Fatal(err)
	}

	// start a goroutine to listen for stablecoin payments and issuance
	go stablecoin.ListenForPayments()
	// don't have an error catching thing because if this fails, the platform should
	// not initialize
	// insert an investor with the relevant details
	// add dummy investor and recipient data for the demo
	// uname: john, password: password
	// need to ask for user role as well here, to know whether the user is an investor
	// or recipient so that we can show both sides
	// After this, ask what the user wants to do - there are roughly three options:
	// 1. Create a new investor account
	// 2. Create a new recipient account
	// 3. Login (Are you an investor / recipient)
	fmt.Println("------WHAT DO YOU WANT TO DO?------")
	fmt.Println("1. CREATE A NEW INVESTOR ACCOUNT")
	fmt.Println("2. CREATE A NEW RECIPIENT ACCOUNT")
	fmt.Println("3: ALREADY HAVE AN ACCOUNT")
	opt, err := utils.ScanForInt()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	switch opt {
	case 1:
		log.Println("You have chosen to create a new investor account, welcome")
		err := NewInvestorPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		log.Println("You have chosen to create a new recipient account, welcome")
		err := NewRecipientPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	default:
	}

	investor, recipient, contractor, isRecipient, isContractor, err := LoginPrompt()
	if err != nil {
		log.Println(err)
		return
	}
	// check if the user is a recipient here
	if isRecipient {
		// we already have the recipient, so no need to make a call to the database
		// have a for loop here with various options
		for {
			fmt.Println("------------RECIPIENT INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Projects (STAGE 3)")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Payback towards an Project (STAGE 6)")
			fmt.Println("  4. Exchange XLM for USD")
			fmt.Println("  5. Finalize a specific Project (STAGE 2->3)")
			fmt.Println("  6. View all Pre Origin Projects (STAGE 0)")
			fmt.Println("  7. View all Projects (ALL STAGES)")
			fmt.Println("  8. View all Origin Projects (STAGE 1)")
			fmt.Println("  default: Exit")
			optI, err := utils.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				fmt.Println("------------LIST OF ALL AVAILABLE PROJECTS------------")
				allProjects, err := database.RetrieveStage3Projects()
				if err != nil {
					log.Println("Error retrieving all projects from the database")
				}
				PrintProjects(allProjects)
				break
			case 2:
				PrintRecipient(recipient)
				break
			case 3:
				PrintPBProjects(recipient.ReceivedProjects)
				fmt.Println("WHICH ORDER DO YOU WANT TO PAY BACK TOWARDS? (ENTER ORDER NUMBER)")
				projectNumber, err := utils.ScanForInt()
				if err != nil {
					log.Println("INPUT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				// check if we can get the project using the project number that we have here
				rtContract, err := database.RetrieveProject(projectNumber)
				if err != nil {
					log.Println("Couldn't retrieve project, try again!")
					continue
				}
				// so we can retrieve the project using the project Index, nice
				PrintPBProject(rtContract.Params)
				fmt.Println("HOW MUCH DO YOU WANT TO PAYBACK?")
				paybackAmount, err := utils.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Printf(" DO YOU WANT TO CONFIRM THAT YOU WANT TO PAYBACK %s TOWARDS THIS PROJECT? (PRESS N IF YOU DON'T WANT TO)\n", paybackAmount)
				confirmOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THE PAYBACK ORDER")
					break
				}
				fmt.Printf("PAYING BACK %s TOWARDS PROJECT NUMBER: %d\n", paybackAmount, rtContract.Params.Index) // use the rtContract.Params here instead of using projectNumber from long ago
				// now we need to call back the payback function to payback the asset
				// Here, we will simply payback the DEBTokens that was sent to us earlier
				if rtContract.Params.DEBAssetCode == "" {
					log.Fatal("Project not found")
				}

				err = recipient.Payback(rtContract, rtContract.Params.DEBAssetCode, platformPublicKey, paybackAmount)
				// TODO: right now, the payback asset directly sends back, change
				if err != nil {
					log.Println("PAYBACK TX FAILED, PLEASE TRY AGAIN!")
					break
				}
				// now send back the PBToken from the platform to the issuer
				// this function is optional and can be deleted in case we don't need PBAssets
				err = assets.SendPBAsset(rtContract.Params, recipient.U.PublicKey, paybackAmount, platformSeed, platformPublicKey)
				if err != nil {
					log.Println("PBAsset sending back FAILED, PLEASE TRY AGAIN!")
					break
				}
				rtContract, err = database.RetrieveProject(projectNumber)
				if err != nil {
					log.Println("Couldn't retrieve updated project, check again!")
					continue
				}
				// we should update the local slice to keep track of the changes here
				recipient.UpdateProjectSlice(rtContract.Params)
				// so we can retrieve the project using the project Index, nice
				PrintParams(rtContract.Params)
				// print the project in a nice way
				break
			case 4:
				log.Println("Enter the amount you want to convert into STABLEUSD")
				convAmount, err := utils.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					break
				}
				hash, err := assets.TrustAsset(stablecoin.StableUSD, consts.StablecoinTrustLimit, recipient.U.PublicKey, recipient.U.Seed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for trusting stableUSD: ", hash)
				// now send coins across and see if our tracker detects it
				_, hash, err = xlm.SendXLM(stablecoin.Issuer.PublicKey, convAmount, recipient.U.Seed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipient.U.PublicKey)
				break
			case 5:
				// we shouild  finalize the contract that we want
				// can be imagined as some sort of voting mechanism to choose the winning
				// contract.
				// now we display a list of options for the recipient to choose which parameter
				// he would like to decide the winning contract
				// 1. Price
				// 2. Completion time
				// 3. Select Manually
				fmt.Println("CHOOSE THE METRIC BY WHICH YOU WANT TO SELECT THE WINNING BID: ")
				allContracts, err := database.RetrieveProposedProjectsIR(recipient.U.Index)
				// retrieve all contracts towards the project
				if err != nil {
					log.Fatal(err)
				}
				PrintProjects(allContracts)
				fmt.Println("1. PRICE")
				fmt.Println("2. COMPLETION TIME (IN YEARS)")
				fmt.Println("3. SELECT MANUALLY")
				fmt.Println("ENTER YOUR CHOICE AS A NUMBER (1 / 2 / 3)")
				opt, err := utils.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// the length of the above two slices must be the same
				for _, elem := range allContracts {
					log.Println("======================================================================================")
					PrintProject(elem)
				}
				switch opt {
				case 1:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY LEAST PRICE")
					// here we assume that the timeout period for the auction is up and that
					// price is the winning metric of a specific bid, like in traditional contract
					bestContract, err := database.SelectContractByPrice(allContracts)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					// we need the contractor who proposed this contract
					database.FinalizeProject(bestContract)
					PrintProject(bestContract)
					// now at this point, we need to mark this specific contract as completed.
					// do we set a flag? db entry? how do we do that
				case 2:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY NUMBER OF YEARS")
					bestContract, err := database.SelectContractByTime(allContracts)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					database.FinalizeProject(bestContract)
					PrintProject(bestContract)
				case 3:
					for i, contract := range allContracts {
						log.Println("BEST BID CHOICE NUMBER: ", i)
						PrintProject(contract)
					}
					fmt.Println("ENTER YOUR OPTION AS A NUMBER")
					opt, err := utils.ScanForInt()
					if err != nil {
						log.Println(err)
						break
					}
					log.Println("BEST CONTRACT IS: ")
					// we need the contractor who proposed this contract
					database.FinalizeProject(allContracts[opt])
					PrintProject(allContracts[opt])
				default:
					break
				}
				// now we need to finalize this person, potentially
				// move funds from the investor money to this person and so on.
				// another question is that whether we raise money before and then we have a
				// blind auction or whether we take in their feedback and then present this to
				// investors. Investors would ideally want to know more about what they are
				// investing in, so I guess the second option is better for now.
			case 6:
				fmt.Println("LIST OF ALL PRE ORIGIN PROJECTS BY ORIGINATORS (STAGE 0)")
				allMyProjects, err := database.RetrievePreOriginProjects()
				if err != nil {
					log.Println(err)
					continue
				}
				PrintProjects(allMyProjects)
				fmt.Println("ENTER THE PROJECT INDEX")
				contractIndex, err := utils.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				// we need to upgrade the contract's whose index is contractIndex to stage 1
				// we already have a contract, so just upgrade that part

				for _, elem := range allMyProjects {
					if elem.Params.Index == contractIndex {
						// increase this contract's stage
						log.Println("UPGRADING PROJECT INDEX", elem.Params.Index)
						err = elem.SetOriginContractStage()
						if err != nil {
							log.Println(err)
							break
						}
						break
					}
				}
			case 7:
				fmt.Println("PRINTING ALL PROJECTS: ")
				allContracts, err := database.RetrieveAllProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allContracts)
			case 8:
				fmt.Println("PRINTING ALL ORIGINATED PROJECTS: ")
				x, err := database.RetrieveOriginProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(x)
			default: // this default is for the larger switch case
				// check whether he wants to go back to the display all screen again
				fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
				exitOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if exitOpt == "Y" || exitOpt == "y" {
					fmt.Println("YOU HAVE DECIDED TO EXIT")
					log.Fatal("")
				}
				break
			}
		}
		PrintRecipient(recipient)
		return
	} else if isContractor {
		log.Println("WELCOME BACK!!")
		fmt.Println("----------CONTRACTOR INTERFACE-------------")
		fmt.Println("AVAILABLE ROLES: ")
		fmt.Println("  1. CONTRACTOR")
		fmt.Println("  2. ORIGINATOR")
		var optI int
		if contractor.Contractor {
			optI = 1
		} else if contractor.Originator {
			optI = 2
		}
		switch optI {
		case 1:
			fmt.Println("-------------WELCOME BACK CONTRACTOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. VIEW ALL ORIGINATED (STAGE 1) PROJECTS")
				fmt.Println("  2. VIEW PROFILE")
				fmt.Println("  3. CREATE A PROPOSED (STAGE 2) PROJECT")
				fmt.Println("  4. VIEW ALL MY PROPOSED (STAGE 2) PROJECTS")
				optI, err := utils.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					// TODO: add voting scheme here
					fmt.Println("LIST OF ALL ORIGINATED PROJECTS: ")
					originatedProjects, err := database.RetrieveOriginProjects()
					if err != nil {
						log.Println(err)
						break
					}
					PrintProjects(originatedProjects)
				case 2:
					PrintEntity(contractor)
				case 3:
					fmt.Println("YOU HAVE CHOSEN TO CREATE A NEW PROPOSED PROJECT")
					err = ProposeContractPrompt(&contractor)
					if err != nil {
						log.Println(err)
						continue
					}
				case 4:
					fmt.Println("LIST OF ALL PROPOSED CONTRACTS BY ME: ")
					allMyProjects, err := database.RetrieveProposedProjectsIC(contractor.U.Index)
					if err != nil {
						log.Println(err)
						continue
					}
					PrintProjects(allMyProjects)
				}
			}
		case 2:
			fmt.Println("-------------WELCOME BACK ORIGINATOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. PROPOSE A PRE-ORIGIN (STAGE 0) PROJECT TO A RECIPIENT")
				fmt.Println("  2. VIEW PROFILE")
				fmt.Println("  3. VIEW ALL MY PRE-ORIGINATED (STAGE 0) PROJECTS")
				fmt.Println("  4. VIEW ALL MY ORIGINATED (STAGE 1) PROJECTS")
				optI, err = utils.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					err := OriginContractPrompt(&contractor)
					if err != nil {
						fmt.Println(err)
						fmt.Println("RETURNING BACK TO THE MAIN LOOP")
						continue
					}
				case 2:
					PrintEntity(contractor)
				case 3:
					allMyProjects, err := database.RetrievePreOriginProjectsI(contractor.U.Index)
					if err != nil {
						fmt.Println(err)
						fmt.Println("RETURNING BACK TO THE MAIN LOOP")
						continue
					}
					PrintProjects(allMyProjects)
				case 4:
					allMyProjects, err := database.RetrieveOriginProjectsIO(contractor.U.Index)
					if err != nil {
						fmt.Println(err)
						fmt.Println("RETURNING BACK TO THE MAIN LOOP")
						continue
					}
					PrintProjects(allMyProjects)
				default:
					// check whether he wants to go back to the display all screen again
					fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
					exitOpt, err := utils.ScanForString()
					if err != nil {
						log.Println(err)
						break
					}
					if exitOpt == "Y" || exitOpt == "y" {
						fmt.Println("YOU HAVE DECIDED TO EXIT")
						log.Fatal("")
					}
				}
			}
		}
		PrintEntity(contractor)
	} else {
		// User is an investor
		for {
			// Main investor loop
			fmt.Println("------------INVESTOR INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Projects (STAGE 3)")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Invest in an Project (STAGE 3)")
			fmt.Println("  4. Display All Balances")
			fmt.Println("  5. Exchange XLM for USD")
			fmt.Println("  6. Display all Origin (STAGE 1) Projects")
			fmt.Println("  7. Vote towards a specific proposed project (STAGE 2)")
			fmt.Println("  8. Get ipfs hash of a contract")
			fmt.Println("  default: Exit")
			optI, err := utils.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				fmt.Println("------------LIST OF ALL AVAILABLE ORDERS------------")
				allProjects, err := database.RetrieveStage3Projects()
				if err != nil {
					log.Println("Error retrieving all projects from the database")
				}
				PrintProjects(allProjects)
				break
			case 2:
				PrintInvestor(investor)
				break
			case 3:
				fmt.Println("----WHICH PROJECT DO YOU WANT TO INVEST IN? (ENTER ORDER NUMBER WITHOUT SPACES)----")
				oNumber, err := utils.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// now the user has decided to invest in the asset with index uInput
				// we need to retrieve the project and ask for confirmation
				uContract, err := database.RetrieveProject(oNumber)
				if err != nil {
					log.Println("Couldn't retrieve project, try again!")
					continue
				}
				PrintProject(uContract)
				fmt.Println(" HOW MUCH DO YOU WANT TO INVEST?")
				investmentAmount, err := utils.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println(" DO YOU WANT TO CONFIRM THIS ORDER? (PRESS N IF YOU DON'T WANT TO)")
				confirmOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
					break
				}
				// when I am creating an account, I will have a PublicKey and Seed, so
				// don't need them here
				// check whether the investor has XLM already
				balance, err := xlm.GetNativeBalance(platformPublicKey)
				if err != nil {
					log.Fatal(err)
				}
				// balance is in string, convert to int
				balanceI := utils.StoF(balance)
				log.Println("Platform's balance is: ", balanceI)
				if balanceI < 21 { // 1 to account for fees
					// get coins if balance is this low
					log.Println("Refilling platform balance")
					err := xlm.GetXLM(platformPublicKey)
					// TODO: in future, need to refill platform sufficiently well and interact
					// with a cold wallet that we have previously set
					if err != nil {
						log.Fatal(err)
					}
				}

				balance, err = xlm.GetNativeBalance(platformPublicKey)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Platform balance updated is: ", balance)
				fmt.Printf("Platform seed is: %s and platform's publicKey is %s", platformSeed, platformPublicKey)
				log.Println("Investor's publickey is: ", investor.U.PublicKey)
				balance, err = xlm.GetNativeBalance(investor.U.PublicKey)
				if balance == "" || err != nil {
					// means we need to setup an account first
					// Generating a keypair on stellar doesn't mean that you can send funds to it
					// you need to call the CreateAccount method in project to be able to send funds
					// to it
					log.Println("Investor balance empty, refilling: ", investor.U.PublicKey)
					_, _, err = xlm.SendXLMCreateAccount(investor.U.PublicKey, consts.DonateBalance, platformSeed)
					if err != nil {
						log.Println("Investor Account doesn't have funds")
						log.Fatal(err)
					}
				}
				// balance is in string, convert to float
				balance, err = xlm.GetNativeBalance(investor.U.PublicKey)
				if err != nil {
					log.Fatal(err)
				}
				balanceI = utils.StoF(balance)
				log.Println("Investor balance is: ", balanceI)
				if balanceI < 3 { // to setup trustlines
					_, _, err = xlm.SendXLM(investor.U.PublicKey, consts.DonateBalance, platformSeed)
					if err != nil {
						log.Println("Investor Account doesn't have funds")
						log.Fatal(err)
					}
				}

				recipient := uContract.Params.ProjectRecipient
				// from here on, reference recipient
				balance, err = xlm.GetNativeBalance(recipient.U.PublicKey)
				if balance == "" || err != nil {
					// means we need to setup an account first
					// Generating a keypair on stellar doesn't mean that you can send funds to it
					// you need to call the CreateAccount method in project to be able to send funds
					// to it
					_, _, err = xlm.SendXLMCreateAccount(recipient.U.PublicKey, consts.DonateBalance, platformSeed)
					if err != nil {
						log.Println("Recipient Account doesn't have funds")
						log.Fatal(err)
					}
				}
				balance, err = xlm.GetNativeBalance(recipient.U.PublicKey)
				if err != nil {
					log.Fatal(err)
				}
				// balance is in string, convert to float
				balanceI = utils.StoF(balance)
				log.Println("Recipient balance is: ", balanceI)
				if balanceI < 3 { // to setup trustlines
					_, _, err = xlm.SendXLM(recipient.U.PublicKey, consts.DonateBalance, platformSeed)
					if err != nil {
						log.Println("Recipient Account doesn't have funds")
						log.Fatal(err)
					}
				}
				log.Println("The investor's public key and private key are: ", investor.U.PublicKey, " ", investor.U.Seed)
				log.Println("The recipient's public key and private key are: ", recipient.U.PublicKey, " ", recipient.U.Seed)

				log.Println(&investor, &recipient, investmentAmount, uContract.Params)
				// so now we have three entities setup, so we create the assets and invest in them
				cProject, err := assets.InvestInProject(platformPublicKey, platformSeed, &investor, &recipient, investmentAmount, uContract) // assume payback period is 5
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println("YOUR PROJECT INVESTMENT HAS BEEN CONFIRMED: ")
				PrintParams(cProject)
				fmt.Println("PLEASE CHECK A BLOCKCHAIN EXPLORER TO CONFIRM BALANCES: ")
				fmt.Println("https://testnet.steexp.com/account/" + investor.U.PublicKey + "#balances")
				break
			case 4:
				balances, err := xlm.GetAllBalances(investor.U.PublicKey)
				if err != nil {
					log.Fatal(err)
				}
				// need to pr etty print this, experiment out with stuff
				PrintBalances(balances)
				break
			case 5:
				// this should be expanded in the future to make use of the inbuilt DEX
				// on stellar (checkout stellarterm)
				log.Println("Enter the amount you want to convert into STABLEUSD")
				// this would also mean that you need to check whether we have the balance
				// here and then proceed further
				convAmount, err := utils.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					return
				}
				// maybe don't trust asset again when you've trusted it already? check if that's
				// possible and save on the tx fee for a single transaction. But I guess its
				// difficult to retrieve trustlines, so we'll go ahead with it
				hash, err := assets.TrustAsset(stablecoin.StableUSD, consts.StablecoinTrustLimit, investor.U.PublicKey, investor.U.Seed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for trusting stableUSD: ", hash)
				// now send coins across and see if our tracker detects it
				_, hash, err = xlm.SendXLM(stablecoin.Issuer.PublicKey, convAmount, investor.U.Seed)
				if err != nil {
					log.Fatal(err)
				}

				log.Println("tx hash for sent xlm: ", hash, "pubkey: ", investor.U.PublicKey)
				rpc.StartServer("8080") // run this in project to check whether the go routine is running
				break
			case 6:
				// TODO: add voting scheme here
				fmt.Println("LIST OF ALL ORIGINATED PROJECTS: ")
				originatedProjects, err := database.RetrieveOriginProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(originatedProjects)
			case 7:
				// this is the case where an investor can vote on a particular proposed project
				// ie stage 2 projects for the recipient to have an understanding about
				// which contracts are popular and can receive more investor money
				fmt.Println("LIST OF ALL PROPOSED ORDERS: ")
				allProposedProjects, err := database.RetrieveProposedProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allProposedProjects)
				fmt.Println("WHICH CONTRACT DO YOU WANT TO VOTE TOWARDS?")
				vote, err := utils.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("You have voted for contract number: ", vote)
				err = investor.VoteTowardsProposedProject(allProposedProjects, vote)
				if err != nil {
					log.Println(err)
					break
				}
			case 8:
				fmt.Println("WELCOME TO THE IPFS HASHING INTERFACE")
				fmt.Println("ENTER THE STRING THAT YOU WOULD LIKE THE IPFS HASH FOR")
				// the UI should ideally have a menu that asks the user for a file and then
				// produces the hash of it. In this case, we shall use a sample ipfs file
				// and then hash it.
				// this uses the platform's ipfs key though, not the user's. If the user
				// wants to serve his own ipfs files, he is better off running a client on
				// his own
				hashString, err := utils.ScanForString()
				if err != nil {
					fmt.Println("Couldn't read user input, going back to the main loop")
				}
				hash, err := ipfs.AddStringToIpfs(hashString)
				if err != nil {
					fmt.Println("Couldn't hash user input, exiting to main menu", err)
					break
				}
				hashCheck, err := ipfs.GetStringFromIpfs(hash)
				if err != nil || hashCheck != hashString {
					fmt.Println("Hashed strings and retrieved strings don't match, don't use this hash!")
					break
				}
				// don't print this hash unless we can decrypt it and be sure that it behaves as expected
				log.Println("THE HASH OF THE PROVIDED STRING IS: ", hash)
				// try to retrieve the string back from ipfs and check if it works correctly
			default:
				// check whether he wants to go back to the display all screen again
				fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
				exitOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if exitOpt == "Y" || exitOpt == "y" {
					fmt.Println("YOU HAVE DECIDED TO EXIT")
					log.Fatal("")
				}
			} // end of switch
		}
		log.Fatal("")
		rpc.StartServer(opts.Port) // this must be towards the end
	}
}
