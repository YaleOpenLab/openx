package solar

import (
	"encoding/json"
	"fmt"

	database "github.com/OpenFinancing/openfinancing/database"
	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/boltdb/bolt"
)

// A legal contract should ideally be stored on ipfs and we must keep track of the
// ipfs hash so that we can retrieve it later when required
// A Project is what is stored in the database and what is used by other packages
// Project imports SolarParams since having everythin inside one struct is tedious
// and SolarParams already has lots of keys. Also, this doesn't affect the way its
// actually stored in the database, so its a nice way to do it.
// SolarParams is also what's needed by the assets and other stuff whereas the other fields
// are needed in other parts, another nice distinction
type Project struct {
	Params SolarParams // Params is the former Order struct imported into the new Project structure

	// List of entities other than the contractor
	Originator    Entity  // a specific contract must hold the person who originated it
	OriginatorFee float64 // fee paid to the originator from the total fee of the project
	Developer     Entity  // the developer who would be responsible for isntallign the solar panels and the IoT hubs
	DeveloperFee  Entity  // the fee charged by the developer
	Guarantor     Entity  // the person guaranteeing the specific project in question

	// List of contractor entities
	Contractor             Entity  // the person with the proposed contract
	ContractorFee          float64 // fee paid to the contractor from the total fee of the project
	SecondaryContractor    Entity  // this is the secondary contractor involved in the project
	SecondaryContractorFee float64 // the fee to be paid towards the secondary contractor
	TertiaryContractor     Entity  // tertiary contractor if any can be added to the system
	TertiaryContractorFee  float64 // the fee to be paid towards the tertiary contractor

	ProjectRecipient database.Recipient  // The recipient of the project in question
	ProjectInvestors []database.Investor // The various investors who are invested in the project
	SeedInvestors    []database.Investor // investors who took part before the contract was at stage 3

	Stage       float64 // the stage at which the contract is at, float due to potential support of 0.5 state changes in the future
	AuctionType string  // the type of the auction in question. Default is blind auction unless explicitly mentioned

	// Various ipfs hashes that we need to store
	OriginatorMoUHash       string // the contract between the originator and the recipient at stage LegalContractStage
	ContractorContractHash  string // the contract between the contractor and the platform at stage ProposeProject
	InvPlatformContractHash string // the contract between the investor and the platform at stage FundedProject
	RecPlatformContractHash string // the contract between the recipient and the platform at stage FundedProject

	Reputation float64 // the positive reputation associated with a given project
	// MWTODO: get feedback on Reputation weighting
	Lock    bool   // lock investment in order to wait for receipient's confirmation
	LockPwd string // the recipient's seedpwd. Will be set to null as soon as we use it.
}

// so a project's rough workflow is like
// origincontract (0) -> approval by recipient (1) -> OpenForMoneyStage (1.5) -> ...
// NewOriginProject returns a new project passed a project and originator to assign to
// Save or Insert inserts a specific Project into the database
func (a *Project) Save() error {
	db, err := database.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return err
		}
		return b.Put([]byte(utils.ItoB(a.Params.Index)), encoded)
	})
	return err
}

// RetrieveProject retrieves the project with the specified index from the database
func RetrieveProject(key int) (Project, error) {
	var inv Project
	db, err := database.OpenDB()
	if err != nil {
		return inv, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return fmt.Errorf("Retrieved project nil")
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

// RetrieveAllProjects retrieves all projects from the database
func RetrieveAllProjects() ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return err
			}
			// append only contracts which are open for funding and below
			arr = append(arr, rProject)
		}
		return nil
	})
	return arr, err
}

func RetrieveProjectsAtStage(stage float64) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return err
			}
			// append only contracts which are open for funding and below
			if rProject.Stage == stage {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveContractorProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			if rProject.Stage == stage && rProject.Contractor.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// RetrieveOriginProjectsIO is used when we want to display the list of originated
// contracts to the originator
func RetrieveOriginatorProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return err
			}
			if rProject.Stage == stage && rProject.Originator.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func RetrieveRecipientProjects(stage float64, index int) ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		// this is Update to cover the case where the  bucket doesn't exists and we're
		// trying to retrieve a list of keys
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return nil
			}
			if rProject.Stage == stage && rProject.ProjectRecipient.U.Index == index {
				// return contracts which have been originated and are not final yet
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

// TODO: Consider that for this authorization to happen, there could be a
// verification requirement (eg. that the project is relatively feasible),
// and that it may need several approvals for it (eg. Recipient can be two
// figures here — the school entity (more visible) and the department of
// education (more admin) who is the actual issuer) along with a validation
// requirement
func VerifyBeforeAuthorizing(projIndex int) bool {
	// here we verify some information related to the originator
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return false
	}
	// print out the originator's name here. In the future, this would involve
	// the kyc operator to check the originator's credentials
	fmt.Printf("ORIGINATOR'S NAME IS: %s and PROJECT's METADATA IS: %s", project.Originator.U.Name, project.Params.Metadata)
	return true
}

// RecipientAuthorize is used to originate a specific project and take it to stage 1
func RecipientAuthorize(projIndex int, recpIndex int) error {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	if project.Stage != 0 { // project stage not at zero, shouldn't be called here
		return fmt.Errorf("Project stage not zero")
	}
	if !VerifyBeforeAuthorizing(projIndex) {
		// not verified, quit Here
		return fmt.Errorf("Originator not verified")
	}
	recipient, err := database.RetrieveRecipient(recpIndex)
	if err != nil {
		return err
	}
	if project.ProjectRecipient.U.Name != recipient.U.Name {
		return fmt.Errorf("You can't authorize a project which is not assigned to you!")
	}
	// set the project as both originated and ready for investors' money
	err = project.SetOriginProject()
	if err != nil {
		return err
	}
	// once a specific project has been originated, set the reputation for the originator
	// depending on the value proposed by the recipient. This reputation will increase
	// automatically once the project has been invested in users because the total value
	// would be higher, so we needn't have a separate mechanism that deals with this.
	err = RepOriginatedProject(project.Originator.U.Index, project.Params.Index)
	if err != nil {
		return err
	}
	/*
		err = project.SetOpenForMoneyStage()
		if err != nil {
			return err
		}
	*/
	return nil
}

func VoteTowardsProposedProject(invIndex int, votes int, projectIndex int) error {
	inv, err := database.RetrieveInvestor(invIndex)
	if err != nil {
		return err
	}
	if votes > inv.VotingBalance {
		return fmt.Errorf("Can't vote with an amount greater than available balance")
	}
	project, err := RetrieveProject(projectIndex)
	if err != nil {
		return err
	}
	if project.Stage != 2 {
		return fmt.Errorf("You can't vote for a project with stage less than 2")
	}
	// we have the specific contract and need to upgrade the number of votes on this one
	project.Params.Votes += votes
	err = project.Save()
	if err != nil {
		return err
	}
	err = inv.DeductVotingBalance(votes)
	if err != nil {
		return err
	}
	fmt.Println("CAST VOTE TOWARDS PROJECT SUCCESSFULLY")
	return nil
}

// stage is set automatically to 1 by the call to SetOriginProject
// this function is used exclusively for testing
func newOriginProject(project SolarParams, originator Entity) (Project, error) {
	// need variadic params to store optional stuff
	var proposedProject Project
	proposedProject.Params = project
	proposedProject.Originator = originator
	proposedProject.Stage = 1
	err := proposedProject.Save()
	return proposedProject, err
}

// A function to find a project within an array of projects, given the key or index
func findInKey(key int, arr []Project) (Project, error) {
	var dummy Project
	for _, elem := range arr {
		if elem.Params.Index == key {
			return elem, nil
		}
	}
	return dummy, fmt.Errorf("Not found")
}

func (project *Project) updateRecipient(a database.Recipient) error {
	pos := -1
	for i, mem := range a.ReceivedSolarProjects {
		if mem == project.Params.DebtAssetCode {
			// rewrite the thing in memory that we have
			pos = i
			break
		}
	}
	if pos != -1 {
		// rewrite the thing in memory
		a.ReceivedSolarProjects[pos] = project.Params.DebtAssetCode
		err := a.Save()
		return err
	}
	return nil
}

func SaveOriginatorMoU(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	a.OriginatorMoUHash = hash
	return a.Save()
}

func SaveContractHash(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	a.ContractorContractHash = hash
	return a.Save()
}

func SaveInvPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	a.InvPlatformContractHash = hash
	return a.Save()
}

func SaveRecPlatformContract(projIndex int, hash string) error {
	a, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	a.RecPlatformContractHash = hash
	return a.Save()
}

func RetrieveLockedProjects() ([]Project, error) {
	var arr []Project
	db, err := database.OpenDB()
	if err != nil {
		return arr, err
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(database.ProjectsBucket)
		for i := 1; ; i++ {
			var rProject Project
			x := b.Get(utils.ItoB(i))
			if x == nil {
				break
			}
			err := json.Unmarshal(x, &rProject)
			if err != nil {
				return err
			}
			// append only contracts which are open for funding and below
			if rProject.Lock {
				arr = append(arr, rProject)
			}
		}
		return nil
	})
	return arr, err
}

func UnlockProject(username string, pwhash string, projIndex int, seedpwd string) error {
	fmt.Println("UNLOCKING PROJECT")
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	recipient, err := database.ValidateRecipient(username, pwhash)
	if err != nil {
		return err
	}

	if recipient.U.LoginPassword != project.ProjectRecipient.U.LoginPassword {
		return fmt.Errorf("Seeds don't match, quitting!")
	}

	if !project.Lock {
		return fmt.Errorf("Project not locked")
	}

	project.LockPwd = seedpwd
	project.Lock = false
	fmt.Println("Project unlocked and lock password set to seed password! SEEDPWD: ", seedpwd)
	err = project.Save()
	if err != nil {
		return err
	}
	return nil
}
