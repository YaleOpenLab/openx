package opensolar

// this should contain the future guarantor related functions once we define them concretely

// NewGuarantor returns a new guarantor
func NewGuarantor(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "guarantor")
}

// AddFirstLossGuarantee adds the given entity as a first loss guarantor
func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}
