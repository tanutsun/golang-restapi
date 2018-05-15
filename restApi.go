package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	db "gopkg.in/gorethink/gorethink.v4"

	config "restApi/config" //get database connection

	"encoding/json"
)

var session *db.Session = config.GetSession()

type Person struct {
	Id    string `gorethink:"id,omitempty"`
	Name  string `gorethink:"name"`
	Place string `gorethink:"place"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/get", Get)
	router.HandleFunc("/insert", Insert)
	router.HandleFunc("/update", Update).Methods("POST")
	router.HandleFunc("/delete/{id}", Delete)

	log.Fatal(http.ListenAndServe(":8081", router))
}

func Index(w http.ResponseWriter, r *http.Request) {

	
	createTable(w)

	fmt.Fprintf(w, "Rest Api Golang")
}

func Get(w http.ResponseWriter, r *http.Request) {

	fetchAllRecords(w)
	recordCount(w)
}

func Insert(w http.ResponseWriter, r *http.Request) {

	id := insertRecord(w)

	printStr(w, id)
	recordCount(w)
}


func Update(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	if id != "" {
		// Update a record
		updateRecord(w, id)
	}
	recordCount(w)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	if vars["id"] != "" {
		// Delete a record
		deleteRecord(w, vars["id"])
	}
	recordCount(w)
}



func createTable(w http.ResponseWriter) interface{} {
	result, err := db.DB("test").TableCreate("people").RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}

	printStr(w, "*** Create table result: ***")
	printObj(w, result)
	printStr(w, "\n")
	return result
}

func insertRecord(w http.ResponseWriter) string {
	var data = map[string]interface{}{
		"Name":  "David Davidson",
		"Place": "Somewhere",
	}

	result, err := db.Table("people").Insert(data).RunWrite(session)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	printStr(w, "*** Insert result: ***")
	printObj(w, result)
	printStr(w, "\n")

	return result.GeneratedKeys[0]
}

func updateRecord(w http.ResponseWriter, id string) {
	var data = map[string]interface{}{
		"Name":  "Steve Stevenson",
		"Place": "Anywhere",
	}

	result, err := db.Table("people").Get(id).Update(data).RunWrite(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	printStr(w, "*** Update result: ***")
	printObj(w, result)
	printStr(w, "\n")
}

func fetchOneRecord(w http.ResponseWriter) {
	cursor, err := db.Table("people").Run(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	var person interface{}
	cursor.One(&person)
	cursor.Close()

	printStr(w, "*** Fetch one record: ***")
	printObj(w, person)
	printStr(w, "\n")
}

func recordCount(w http.ResponseWriter) {
	cursor, err := db.Table("people").Count().Run(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	var cnt int
	cursor.One(&cnt)
	cursor.Close()

	printStr(w, "*** Count: ***")
	printObj(w, cnt)
	printStr(w, "\n")
}

func fetchAllRecords(w http.ResponseWriter) {
	rows, err := db.Table("people").Run(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read records into persons slice
	var persons []Person
	err2 := rows.All(&persons)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	printStr(w, "*** Fetch all rows: ***")
	for _, p := range persons {
		printObj(w, p)
	}
	printStr(w, "\n")
}

func deleteRecord(w http.ResponseWriter, id string) {
	result, err := db.Table("people").Get(id).Delete().Run(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	printStr(w, "*** Delete result: ***")
	printObj(w, result)
	printStr(w, "\n")
}

func printStr(w http.ResponseWriter, v string) {
	//response to api
	fmt.Fprintf(w, v)
	//response to console
	fmt.Println(v)
	
}

func printObj(w http.ResponseWriter, v interface{}) {
	vBytes, _ := json.Marshal(v)
	//response to api
	json.NewEncoder(w).Encode(v)
	//response to console
	fmt.Println(string(vBytes))
	
}