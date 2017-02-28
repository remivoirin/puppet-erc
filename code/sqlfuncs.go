package main

import (
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const dblocation = "./sql/roles.db"

// Create DB, the role table and the default record.
func Dbinitialize() (bool, string) {
	if _, err := os.Stat(dblocation); err == nil {
		return false, "DB file already exists"
	}

	os.Create(dblocation)
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return false, "Could not create and open " + dblocation
	}

	createstmt := `CREATE TABLE roles(
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
				host_regex VARCHAR(200) NOT NULL ,
				role VARCHAR(200) NOT NULL ,
				comment VARCHAR(2000) NULL 
			);`

	_, err = db.Exec(createstmt)
	if err != nil {
		return false, "Could not create table"
	}

	insertstmt := `INSERT INTO roles(host_regex, role, comment)
			VALUES(?,?,?)
			;`
	stmt, err := db.Prepare(insertstmt)
	if err != nil {
		return false, "Could not prepare SQL insert statement"
	}
	_, err = stmt.Exec("(.*)", "default", "Default role matching everything")
	if err != nil {
		return false, "Could not insert default role"
	}

	db.Close()

	return true, "OK"
}

// Insert a record
func Dbinsert(host_regex string, role string, comment string) (bool, string) {
	if _, err := os.Stat(dblocation); err != nil {
		return false, "DB file not found"
	}

	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return false, "Could not open " + dblocation
	}

	defer db.Close()

	insertstmt := `INSERT INTO roles(host_regex, role, comment)
			VALUES(?,?,?)
			;`
	stmt, err := db.Prepare(insertstmt)
	if err != nil {
		return false, "Could not prepare SQL insert statement"
	}

	_, err = stmt.Exec(host_regex, role, comment)
	if err != nil {
		return false, "Could not insert data"
	}

	return true, "OK"
}

// Return all entries
func Dblist () []Fullentry {
	if _, err := os.Stat(dblocation); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	sql_readall := `
		SELECT id, host_regex, role, comment FROM roles
		ORDER BY id DESC
		`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var entries []Fullentry
	for rows.Next() {
		entry := Fullentry{}
		err2 := rows.Scan(&entry.Id, &entry.Host_regex, &entry.Role, &entry.Comment)
		if err2 != nil {
			panic(err2)
		}
		entries = append(entries, entry)
	}

	return entries
}

// Delete a record
func Dbdeletebyid(myid int) (bool, string) {
	if _, err := os.Stat(dblocation); err != nil {
		return false, "DB file not found"
	}

	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return false, "Could not open " + dblocation
	}

	defer db.Close()

	deletestmt := `DELETE FROM roles
			WHERE id=?
			;`
	stmt, err := db.Prepare(deletestmt)
	if err != nil {
		return false, "Could not prepare SQL delete statement"
	}

	_, err = stmt.Exec(myid)
	if err != nil {
		return false, "Could not delete data"
	}

	return true, "OK"
}
