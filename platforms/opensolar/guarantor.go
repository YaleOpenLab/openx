package opensolar

// this should contain the future guarantor related functions once we define them concretely
func NewGuarantor(uname string, pwd string, seedpwd string, Name string,
	Address string, Description string) (Entity, error) {
	return newEntity(uname, pwd, seedpwd, Name, Address, Description, "guarantor")
}
