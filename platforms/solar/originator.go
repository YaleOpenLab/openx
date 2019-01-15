package solar

import (
	"fmt"
	database "github.com/OpenFinancing/openfinancing/database"
	utils "github.com/OpenFinancing/openfinancing/utils"
)

func NewOriginator(uname string, pwd string, seedpwd string, Name string, Address string, Description string) (Entity, error) {
	return NewEntity(uname, pwd, seedpwd, Name, Address, Description, "originator")
}

func (contractor *Entity) OriginContract(panelSize string, totalValue int, location string, years int, metadata string, recIndex int) (SolarProject, error) {
	var pc SolarProject
	var err error

	// for this, create a new  contract and store in the contracts db. Wea re sorting
	// by stage, so it shouldn't matter a whole lot
	indexCheck, err := RetrieveAllProjects()
	if err != nil {
		return pc, fmt.Errorf("Projects could not be retrieved!")
	}
	pc.Params.Index = len(indexCheck) + 1
	pc.Params.PanelSize = panelSize
	pc.Params.TotalValue = totalValue
	pc.Params.Location = location
	pc.Params.Years = years
	pc.Params.Metadata = metadata
	pc.Params.DateInitiated = utils.Timestamp()
	iRecipient, err := database.RetrieveRecipient(recIndex)
	if iRecipient.U.Index == 0 { // recipient does not exist
		return pc, fmt.Errorf("Recipient does not exist!")
	}
	pc.Params.ProjectRecipient = iRecipient
	pc.Stage = 0 // 0 since we need to filter this out while retrieving the propsoed contracts
	pc.Originator = *contractor
	// instead of storing in this proposedcontracts slice, store it as a project, but not a contract and retrieve by stage
	err = pc.Save()
	// don't insert the project since the contractor's projects are not final
	return pc, err
}
