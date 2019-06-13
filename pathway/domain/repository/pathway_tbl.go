package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"log"
)

const DbFileName = "/go/src/github.com/torsday/gemini_jobcoin_mixer/mutable_data/pathway_ledger.db" // TODO dockerize this
const TableName = "pathway_ledger"
const Col2DepositAddress = "deposit_address"
const Col3OutputAddresses = "unified_output_address_str"
const Col4Debt = "amount_of_debt"
const Col5WhenLastChecked = "when_last_checked" // timestamp
const IndexCol2Depositaddress = "idx_deposit_address"

type PathwayTbl struct{}

// build the table if none exists.
//
// I've chosen to store the output addresses as a single, concatenated string. I've chosen this for simplicity, as it's
// a 1:1 relationship with a pathway. That said, I can envision business reasons not in this mixer where splitting those
// addresses into their own table would make sense. or refactor the persistence altogether (e.g. and addresses table).
func (pwayTbl *PathwayTbl) BuildDbIfNotExists() {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Failed to build database while running BuildDbIfNotExists"))
	}
	defer database.Close()
	//_ = database.Close()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS " +
		TableName + " (id INTEGER PRIMARY KEY, " +
		Col2DepositAddress + " TEXT NOT NULL UNIQUE, " +
		Col3OutputAddresses + " TEXT NOT NULL UNIQUE, " +
		Col4Debt + " TEXT, " +
		Col5WhenLastChecked + " INT)")

	if err != nil {
		log.Fatal(errors.Wrap(err, "query Prepare failed"))
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(errors.Wrap(err, "table creation query execution failed"))
	}

	// add an index for the deposit address to make lookup more efficient.
	statement, err = database.Prepare("CREATE UNIQUE INDEX IF NOT EXISTS " + IndexCol2Depositaddress + " ON " +
		TableName + " (" + Col2DepositAddress + ");")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Index creation preparation failed"))
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Index creation execution failed"))
	}
}

func NewPathwayTbl() *PathwayTbl {
	return &PathwayTbl{}
}

// find pathway by unified output address.
func (pwayTbl *PathwayTbl) FindPathwayByUnifiedOutputAddress(unifiedOutputAddress string) []string {

	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "FindPathwayByUnifiedOutputAddress failed"))
	}
	defer database.Close()

	// TODO confirm this is best practice and works
	pwayRaw, err := database.Query("SELECT * FROM "+TableName+" WHERE "+Col3OutputAddresses+" = ?", unifiedOutputAddress)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "FindPathwayByUnifiedOutputAddress creation failed"))
	}

	simpleRes := getRawPathwaysFromSqlRows(pwayRaw)

	if len(simpleRes) == 0 {
		return []string{}
	}
	return simpleRes[0]
}

// creates a pathway.
//
// defaults to zero debt and null "last checked"
func (pwayTbl *PathwayTbl) CreatePathway(depositAddress string, unifiedOutputAddress string) {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "CreatePathway failed while opening db"))
	}
	defer database.Close()

	statement, err := database.Prepare("INSERT INTO " + TableName + " (" +
		Col2DepositAddress + ", " +
		Col3OutputAddresses + ", " +
		Col4Debt +
		") VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "CreatePathway prep failed"))
	}
	_, err = statement.Exec(depositAddress, unifiedOutputAddress, 0)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "CreatePathway exec failed"))
	}
}

// Update pathway amount given a deposit address.
// TODO handle exceptions
func (pwayTbl *PathwayTbl) UpdatePathwayAmount(depositAddress string, amount string) {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdatePathwayAmount failed while opening db"))
	}
	defer database.Close()
	queryString := "UPDATE " + TableName + " SET " +
		Col4Debt + " = " + amount + " " +
		"WHERE " + Col2DepositAddress + " = '" + depositAddress + "';"
	statement, err := database.Prepare(queryString)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdatePathwayAmount prep failed"))
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdatePathwayAmount exec failed"))
	}
}

// update when a pathway was last checked.
func (pwayTbl *PathwayTbl) UpdateWhenPathwayLastChecked(depositAddress string, whenLastChecked int) {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdateWhenPathwayLastChecked failed while opening db"))
	}
	defer database.Close()
	statement, err := database.Prepare("UPDATE " + TableName + " SET " +
		Col5WhenLastChecked + " =? " +
		"WHERE " + Col2DepositAddress + " =?;")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdateWhenPathwayLastChecked prep failed"))
	}
	_, err = statement.Exec(whenLastChecked, depositAddress)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "UpdateWhenPathwayLastChecked exec failed"))
	}
}

// Get all pathways
func (pwayTbl *PathwayTbl) GetAllPathways() [][]string {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "GetAllPathways failed while opening db"))
	}
	defer database.Close()
	pwayRaw, err := database.Query("SELECT * FROM " + TableName + ";")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "GetAllPathways query failed"))
	}

	simpleRes := getRawPathwaysFromSqlRows(pwayRaw)
	if len(simpleRes) == 0 {
		return [][]string{}
	}
	return simpleRes
}

// Get pathways with debt
func (pwayTbl *PathwayTbl) GetPathwaysWithDebt() [][]string {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "GetPathwaysWithDebt failed while opening db"))
	}
	defer database.Close()
	pwayRaw, _ := database.Query("SELECT * FROM " + TableName + " WHERE " +
		Col4Debt + " IS NOT '0.0' AND " +
		Col4Debt + " IS NOT NULL;")

	simpleRes := getRawPathwaysFromSqlRows(pwayRaw)
	if len(simpleRes) == 0 {
		return [][]string{}
	}
	return simpleRes
}

// Get a pathway's raw data via a search of the deposit address.
func (pwayTbl *PathwayTbl) GetPathwayByDepositAddress(depAdd string) []string {
	database, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "GetPathwayByDepositAddress failed while opening db"))
	}
	defer database.Close()
	queryStr := "SELECT * FROM " + TableName + " WHERE " + Col2DepositAddress + " = ?"
	pwayRaw, _ := database.Query(queryStr, depAdd)
	simpleRes := getRawPathwaysFromSqlRows(pwayRaw) // TODO DRY this
	if len(simpleRes) == 0 {
		return []string{}
	}

	return simpleRes[0]
}

// Get raw string slices from instances of sql rows
//
// importantly, this method brings the amount of debt back to the real number
// TODO consider just using a string rather than persisting as an int
func getRawPathwaysFromSqlRows(rows *sql.Rows) [][]string {
	var rawReturn [][]string

	var id int
	var depositAdd string
	var unifiedOAddStr string
	var debt string
	var whenLastChecked int

	for rows.Next() {
		_ = rows.Scan(&id, &depositAdd, &unifiedOAddStr, &debt, &whenLastChecked)

		// TODO may have to have 'id' catch the first val
		rawReturn = append(rawReturn, []string{depositAdd, unifiedOAddStr, debt, string(whenLastChecked)})
	}

	return rawReturn
}
