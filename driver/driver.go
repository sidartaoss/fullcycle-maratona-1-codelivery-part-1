package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Driver struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Drivers struct {
	Drivers []Driver
}

func loadDrivers() []byte {
	jsonFile, err := os.Open("drivers.json")
	if err != nil {
		panic(err.Error())
	}

	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err.Error())
	}
	return data
}

func ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers := loadDrivers()
	w.Write(drivers)
}

func GetDriverById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadDrivers()

	var drivers Drivers
	if err := json.Unmarshal(data, &drivers); err != nil {
		panic(err.Error())
	}

	for _, v := range drivers.Drivers {
		if v.Uuid == vars["id"] {
			driver, err := json.Marshal(v)
			if err != nil {
				panic(err.Error())
			}
			w.Write(driver)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/drivers", ListDrivers)
	r.HandleFunc("/drivers/{id}", GetDriverById)
	http.ListenAndServe(":8081", r)
}
