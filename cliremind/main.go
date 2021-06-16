package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)
var home = os.Getenv("HOME")
var path = home + "/.clireminder"

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	db := getDb()

	var arg string
	if len(os.Args) < 2 {
		arg = "0"
	} else {
		arg = os.Args[1]
	}


	switch arg {
	case "hello":
		printHello()
	case "add":
		insertReminder(db, os.Args)
	case "del":
		delReminder(db, os.Args[2])
	case "get":
		getReminders(db)
	default:
		showUsage()
	}

}

func printHello() {
	fmt.Println("Hello to you")
}

func showUsage() {
	fmt.Println("A command line reminder tool")
	fmt.Println("usage: \n" +
		"\tcliremind add \"feed the fish\" \"due date\"\n" +
	    "\tcliremind del id" +
		"\tcliremind hello",
	)
}

func getDb() *sql.DB {

	var db *sql.DB
	_, err := os.Stat(path)
	if err != nil {
		db, _ = createDb()
		createTable(db)
	} else {
		db, err = sql.Open("sqlite3", path + "/datastore.db")
		if err != nil {
			fmt.Println("database not found")
		}
	}
	return db
}

func createDb() (*sql.DB, error) {
	log.Println("Creating datastore.db...")
	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf("error creating directory with error: %v\n", err.Error())
	}

	file, err := os.Create(path + "/datastore.db")
	if err != nil {
		fmt.Println("oh shit")
		log.Fatal(err.Error())
	}
	file.Close()

	return sql.Open("sqlite3", path + "/datastore.db")
}

func createTable(db *sql.DB) {
	createReminderTableSQL := `CREATE TABLE reminders (
		"reminder_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"reminder" TEXT,
		"duedate" DATE
	  );`

	statement, err := db.Prepare(createReminderTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
		if err != nil {

		}
}

func insertReminder(db *sql.DB, args []string) {
	if len(args) < 4 {
		fmt.Println("usage:  cliremind add description date")
		os.Exit(0)
	}

	log.Println("Inserting record ...")
	insert := `INSERT INTO reminders(reminder, duedate) VALUES (?, ?)`
	statement, err := db.Prepare(insert) // Prepare statement.
	if err != nil {
		log.Fatalln(err.Error())
	}
	res, err := statement.Exec(args[2], args[3])
	if err != nil {
		fmt.Printf("there was an error adding the reminder: %s\n", err.Error())
		log.Fatalln(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("unable to get the last inserted id: %s\n", err.Error())
		log.Fatalln(err.Error())
	}
	fmt.Printf("added your reminder as id: %v\n", id)
}

func getReminders(db *sql.DB) {
	row, err := db.Query("select reminder_id, reminder, duedate from reminders")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var reminder string
		var duedate time.Time
		var reminderId int
		row.Scan(&reminderId, &reminder, &duedate)

		fmt.Printf( "\n%d \t%s  %s",reminderId,reminder, formatDate(duedate))
	}
	fmt.Println("\n")
}

func delReminder(db *sql.DB, id string) {
	insert := `DELETE FROM reminders where reminder_id = ?`
	statement, err := db.Prepare(insert) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id)
	if err != nil {
		fmt.Printf("unable to delete that id: %s", err.Error())
	}
}

func formatDate(d time.Time) string {

	ro, _ := time.Parse("0001-01-01 00:00:00", "0000-01-01 00:00:00")

	if ro.Equal(d) {
		return "(daily reminder)"
	}
	diff := d.Sub(time.Now())

	if diff.Hours() < 0 {
		return "(past due)"
	}
	return fmt.Sprintf("(due in %v hours)", int(diff.Hours()))
}