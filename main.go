package main

import (
	"encoding/json"
	"fmt"
	config "golang_restapi/config"
	helpers "golang_restapi/helpers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var profiles []Profile = []Profile{}

type User struct {
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
	Email     string `json: "email"`
}

type Profile struct {
	Id          string `json: "id"`
	Department  string `json: "department"`
	Designation string `json: "designation"`
	Employee    User   `json: "employee"`
}

func addItem(q http.ResponseWriter, r *http.Request) {
	var newProfile Profile
	json.NewDecoder(r.Body).Decode(&newProfile)
	q.Header().Set("Content-Type", "application/json")
	profiles = append(profiles, newProfile)
	json.NewEncoder(q).Encode(profiles)

}

func getAllProfiles(q http.ResponseWriter, r *http.Request) {
	q.Header().Set("Content-Type", "application/json")
	json.NewEncoder(q).Encode(profiles)
}

func getProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	if idParam == "" {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be retrieved."))
		return
	}

	profile := helpers.FirstOrDefault(profiles, func(p *Profile) bool { return p.Id == idParam })

	q.Header().Set("Content-Type", "application/json")

	json.NewEncoder(q).Encode(profile)
}

func updateProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	if idParam == "" {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be retrieved."))
		return
	}

	var updatedProfile Profile
	json.NewDecoder(r.Body).Decode(&updatedProfile)

	profile := helpers.FirstOrDefault(profiles, func(p *Profile) bool { return p.Id == idParam })

	profile.Department = updatedProfile.Department
	profile.Designation = updatedProfile.Designation
	profile.Employee = updatedProfile.Employee

	q.Header().Set("Content-Type", "application/json")

	json.NewEncoder(q).Encode(profile)
}

func deleteProfile(q http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	if idParam == "" {
		q.WriteHeader(400)
		q.Write([]byte("ID could not be retrieved."))
		return
	}

	indexOfProfile := helpers.FindIndex(profiles, func(n Profile) bool {
		return n.Id == idParam
	})

	helpers.RemoveElementByIndex(profiles, indexOfProfile)

	q.WriteHeader(200)
}

func main() {

	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var configuration config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("Port is \t", configuration.Server.Port)

	router := mux.NewRouter()
	router.HandleFunc("/profiles", addItem).Methods("POST")
	router.HandleFunc("/profiles", getAllProfiles).Methods("GET")
	router.HandleFunc("/profiles/{id}", getProfile).Methods("GET")
	router.HandleFunc("/profiles/{id}", updateProfile).Methods("PUT")
	router.HandleFunc("/profiles/{id}", deleteProfile).Methods("DELETE")

	uri := fmt.Sprint(":", configuration.Server.Port)
	fmt.Println("URI is \t", uri)
	http.ListenAndServe(uri, router)
}
