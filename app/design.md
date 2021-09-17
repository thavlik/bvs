Main Menu
|-> Voting Machine Mode
| |-> Start screen ("Please scan QR code to input your voter information")
|   |-> Information Prompts
|     |-> Candidate selection (either select from list or enter address to write-in)
|       |-> Summary
|         |-> Confirmation
|-> Admin Mode
  |-> User/password authentication
    |-> Manage Elections
    | |-> Search for election by policyID
    | |-> Create Election
    | | |-> Set candidates on ballot
    | |   |-> Set deadline
    | |     |-> Summary
    | |       |-> Confirmation
    | |-> Enumerate Elections
    |   |-> View Election
    |     |-> Manage Ballot
    |     | |-> Enumerate Candidates
    |     | |-> Add Candidate
    |     | |-> Delete Candidate
    |     |-> Voter Roster
    |     | |-> Show Votes page search results for the policyID
    |     |   |-> Void vote (sends .Void token to voter)
    |     |-> Delete Election (deletes private keys from mongo)
    |-> Manage Minters
    | |-> Search for minter by payment address
    | |-> Create Minter
    | | |-> Enter personal information (?)
    | | |-> Assign election access
    | | |-> Summary
    | | |-> Confirmation
    | |-> Enumerate Minters
    |   |-> View Minter
    |     |-> Delete Minter (deletes private keys from mongo)
    |-> Search
      |-> Search for anything by policyID, assetID, address of voter, or address of candidate 

HTTP API:
/elections/get
/elections/list

/elections/delete
/minter/get
/minter/list

/minter/delete
/vote/get
/vote/void