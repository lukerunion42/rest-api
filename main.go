package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Eaddress  string `json:"emailaddress,omitempty"`
	Pnumber   string `json:"phonenumber,omitempty"`
}

var people []Person

// This function will import a CSV file to manipulate
//  must have file in Working Directory
func LoadPeopleCsv() []Person {
	var people []Person

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//Open CSV File
	csvFile, err := os.Open(pwd + "/testdata.csv")
	if err != nil {
		log.Fatal(err)
	}

	//This will give you the working Directory in commmand line
	//println(pwd + "/data.csv")

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		people = append(people, Person{
			ID:        line[0],
			Firstname: line[1],
			Lastname:  line[2],
			Eaddress:  line[3],
			Pnumber:   line[4],
		})

	}
	return people
}

//This function will write the data to a CSV File
func WriteToCsvFile(people []Person) {
	//open file
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//Open CSV File
	csvFile, err := os.Create(pwd + "/testdata2.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(bufio.NewWriter(csvFile))
	defer writer.Flush()

	for _, value := range people {
		var s []string
		s = append(s, value.ID)
		s = append(s, value.Firstname)
		s = append(s, value.Lastname)
		s = append(s, value.Eaddress)
		s = append(s, value.Pnumber)
		err := writer.Write(s)
		if err != nil {
			log.Fatal("Cannot Write to File", err)
		}

	}

}

// Show all people
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	people := LoadPeopleCsv()
	json.NewEncoder(w).Encode(people)
}

// Show one person
func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	people := LoadPeopleCsv()
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// Create a person
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	people := LoadPeopleCsv()
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
	WriteToCsvFile(people)
}

// Delete a person
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	people := LoadPeopleCsv()
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
	WriteToCsvFile(people)
}

// Update a Person
func UpdatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	people := LoadPeopleCsv()
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			var person Person
			_ = json.NewDecoder(r.Body).Decode(&person)
			if person.ID == item.ID {
				people = append(people[:index], people[index+1:]...)

				//This appends item to the end
				people = append(people, person)
			}

			//This will sort the slice
			for i := range people {
				for j := range people {
					if people[i].ID < people[j].ID {
						var sub = people[i]
						people[i] = people[j]
						people[j] = sub
					}
				}
			}

		}
		json.NewEncoder(w).Encode(people)
	}
	WriteToCsvFile(people)
}

// main function
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", UpdatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}
