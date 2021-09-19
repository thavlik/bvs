Public blockchain explorer-type website
|-> Home page (layout like https://www.flexpool.io/)
| |-> Big search bar
| |-> Statistics, graphs, leaderboards, etc.
|-> Search
|-> Votes
| |-> List votes
| |-> View vote
|   |-> Link to election
|   |-> Link to voter
|   |-> Link to minter
|   |-> Link to asset search on external cardano blockchain
|-> Elections
| |-> List elections
| |-> View election
|   |-> Graphs and stats specific to election 
|   |-> Link to vote search on external cardano blockchain
|-> Minters
| |-> List minters
| |-> View minter
|   |-> Link to search on external cardano blockchain
Voting/admin app
|-> Stats dashboard
| |-> Leaderboard
| |-> Votes/time graphs
| |-> Robust query & visualization tools
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
    | | |-> Summary
    | | |-> Confirmation
    | |-> Enumerate Minters
    |   |-> View Minter
    |     |-> Delete Minter (deletes private keys from mongo)
    |-> Search
      |-> Search for anything by policyID, assetID, address of voter, or address of candidate 

HTTP API:

POST /search

GET /elections/get
GET /elections/list
POST /elections/create
POST /elections/delete
GET /election/candidates
PUT /election/candidates

GET /minter/get
GET /minter/list
POST /minter/create
POST /minter/delete

GET /vote/list
GET /vote/get
POST /vote/void
POST /vote/mint
POST /vote/cast
