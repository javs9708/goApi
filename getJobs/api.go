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
	Id gocql.UUID
	Url string
	Status int
	Retries int
	Priority int
}


var id gocql.UUID
var url string
var status int
var retries int
var priority int

func main() {
	// Connect to the cluster.
	cluster := gocql.NewCluster("192.168.8.128")
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



	iter := Session.Query(`SELECT * FROM queue_buffer.jobs`).Iter()
	iter2 := Session.Query(`SELECT * FROM queue_buffer.jobs`).Iter()

	var i=0;
	for iter.Scan(&id, &url, &status, &retries, &priority) {
		i++;
	}

	jobsList := make([]jobs, i)

	var j=0;
	for iter2.Scan(&id, &url, &status,&retries, &priority) {

		fmt.Println("Id:", id, "Url:", url, "Status:", status, "Retries:", retries, "Priority:", priority)

		jobsList[j].Id=id
		jobsList[j].Url=url
		jobsList[j].Status=status
		jobsList[j].Retries=retries
		jobsList[j].Priority=priority
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
		http.ListenAndServe(":9001", nil)


}
