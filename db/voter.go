package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// VoterHistory is the struct that represents a single VoterHistory item
type VoterHistory struct{
	PollId int
	VoteId int
	VoteDate time.Time
}

// Voter is the struct that represents a single Voter item
type Voter struct{
	VoterId int
	Name string
	Email string
	VoteHistory []VoterHistory
}

type VoterList struct {
	Voters map[int]Voter //A map of VoterIDs as keys and Voter structs as values
}

//constructor for VoterList struct
func NewVoterList() (*VoterList, error) {

	//Now that we know the file exists, at at the minimum we have
	//a valid empty DB, lets create the ToDo struct
	voterList := &VoterList{
		Voters: make(map[int]Voter),
	}

	// We should be all set here, the ToDo struct is ready to go
	// so we can support the public database operations
	return voterList, nil
}

//Add receivers to any structs you want, but at the minimum you should add the API behavior to the
//VoterList struct as its managing the collection of voters.  Also dont forget in the constructor
//that you need to make the map before you can use it - make map[int]Voter

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// AddItem accepts a ToDoItem and adds it to the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must not already exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if so, return an error
//
// Postconditions:
//
//	 (1) The item will be added to the DB
//		(2) The DB file will be saved with the item added
//		(3) If there is an error, it will be returned
func (t *VoterList) AddVoter(voter Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	_, ok := t.Voters[voter.VoterId]
	if ok {
		return errors.New("item already exists")
	}

	//Now that we know the item doesn't exist, lets add it to our map
	t.Voters[voter.VoterId] = voter

	//If everything is ok, return nil for the error
	return nil
}

// DeleteItem accepts an item id and removes it from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be removed from the DB
//		(2) The DB file will be saved with the item removed
//		(3) If there is an error, it will be returned
func (t *VoterList) DeleteVoter(id int) error {

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(t.Voters, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /todo endpoint
func (t *VoterList) DeleteAll() error {
	//To delete everything, we can just create a new map
	//and assign it to our existing map.  The garbage collector
	//will clean up the old map for us
	t.Voters = make(map[int]Voter)

	return nil
}

// UpdateItem accepts a ToDoItem and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (t *VoterList) UpdateVoter(voter Voter) error {

	// Check if item exists before trying to update it
	// this is a good practice, return an error if the
	// item does not exist
	_, ok := t.Voters[voter.VoterId]
	if !ok {
		return errors.New("item does not exist")
	}

	//Now that we know the item exists, lets update it
	t.Voters[voter.VoterId] = voter

	return nil
}

// GetItem accepts an item id and returns the item from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be returned, if it exists
//		(2) If there is an error, it will be returned
//			along with an empty ToDoItem
//		(3) The database file will not be modified
func (t *VoterList) GetVoter(id int) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	item, ok := t.Voters[id]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return item, nil
}

// GetAllItems returns all items from the DB.  If successful it
// returns a slice of all of the items to the caller
// Preconditions:   (1) The database file must exist and be a valid
//
// Postconditions:
//
//	 (1) All items will be returned, if any exist
//		(2) If there is an error, it will be returned
//			along with an empty slice
//		(3) The database file will not be modified
func (t *VoterList) GetAllVoters() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voterList []Voter

	//Now lets iterate over our map and add each item to our slice
	for _, voter := range t.Voters {
		voterList = append(voterList, voter)
	}

	//Now that we have all of our items in a slice, return it
	return voterList, nil
}

// GetVoterPolls retrieves the voting history for a specific voter.
// It takes voter ID as input and returns their voting history as a slice of VoterHistory.
func (t *VoterList) GetVoterPolls(voterID int) ([]VoterHistory, error) {
	voter, err := t.GetVoter(voterID)
	if err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

// GetVoterPoll retrieves a specific voting record for a voter.
// It takes voter ID and poll ID as input and returns the corresponding VoterHistory if found.
func (t *VoterList) GetVoterPoll(voterID, pollID int) (VoterHistory, error) {
	voter, err := t.GetVoter(voterID)
	if err != nil {
		return VoterHistory{}, err
	}

	for _, history := range voter.VoteHistory {
		if history.PollId == pollID {
			return history, nil
		}
	}

	return VoterHistory{}, errors.New("poll not found for this voter")
}

// AddVoterPoll adds a new voting record for a voter.
// It takes voter ID, poll ID, and vote date as input and adds the record to the corresponding voter.
func (t *VoterList) AddVoterPoll(voterID, pollID int, voteDate time.Time) error {
	voter, err := t.GetVoter(voterID)
	if err != nil {
		return err
	}

	newVoterHistory := VoterHistory{
		PollId:   pollID,
		VoteId:   len(voter.VoteHistory) + 1, // Assuming vote ID increments linearly
		VoteDate: voteDate,
	}

	voter.VoteHistory = append(voter.VoteHistory, newVoterHistory)

	err = t.UpdateVoter(voter)
	if err != nil {
		return err
	}

	return nil
}

// UpdateVoterPoll updates a voting record for a voter.
// It takes voter ID, poll ID, and new vote date as input and updates the corresponding record.
func (t *VoterList) UpdateVoterPoll(voterID, pollID int, newVoteDate time.Time) error {
	voter, err := t.GetVoter(voterID)
	if err != nil {
		return err
	}

	for i, history := range voter.VoteHistory {
		if history.PollId == pollID {
			voter.VoteHistory[i].VoteDate = newVoteDate
			err := t.UpdateVoter(voter)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("poll not found for this voter")
}

// DeleteVoterPoll deletes a voting record for a voter.
// It takes voter ID and poll ID as input and removes the corresponding record.
func (t *VoterList) DeleteVoterPoll(voterID, pollID int) error {
	voter, err := t.GetVoter(voterID)
	if err != nil {
		return err
	}

	for i, history := range voter.VoteHistory {
		if history.PollId == pollID {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			err := t.UpdateVoter(voter)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("poll not found for this voter")
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (t *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (t *VoterList) PrintAllVoters(voterList []Voter) {
	for _, voter := range voterList {
		t.PrintVoter(voter)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (t *Voter) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
