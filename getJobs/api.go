package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
	"net/http"
	"encoding/json"
	"github.com/rs/cors"
	//"github.com/gorilla/mux"
)

type jobs struct {
	Id int
	Url string
	Status string
}



var url string
var id int
var status string

func main() {
	// Connect to the cluster.
	cluster := gocql.NewCluster("192.168.233.128")
	// Use the same timeout as the Java driver.

	cluster.Timeout = 12 * time.Second
	cluster.Consistency = gocql.Quorum
	// Create the session.
	Session, _:= cluster.CreateSession()
	defer Session.Close()



	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

/*

	if err := Session.Query(`INSERT INTO register.urls (id,url,status) VALUES (?, ?, ?)`,
		1,"www.instagram.com", "DONE").Exec(); err != nil {
		log.Fatal(err)
	}

*/



	iter := Session.Query(`SELECT * FROM register.urls`).Iter()
	iter2 := Session.Query(`SELECT * FROM register.urls`).Iter()

	var i=0;
	for iter.Scan(&id, &url, &status) {
		i++;
	}

	jobsList := make([]jobs, i)

	var j=0;
	for iter2.Scan(&id, &url, &status) {

		//fmt.Println("Id:", id, "Url:", url, "Status:", status)

		jobsList[j].Id=id
		jobsList[j].Url=url
		jobsList[j].Status=status
		j++;
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}


		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := jobsList
			json.NewEncoder(w).Encode(u)
		})

		fs := http.FileServer(http.Dir("./api"))
		http.Handle("/monitor/", http.StripPrefix("/monitor/", fs))

		http.Handle("/get_jobs", c.Handler(handler))

		fmt.Println("El servidor se encuentra en ejecución")
		http.ListenAndServe(":8080", nil)


}
