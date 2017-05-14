package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//use in memory map to store the packing lots, with Factory Pattern, could be easily extended to add other repos, such as a real database
var repo *mapRepo

//APP works as RESTful handler
type App struct {
	Router *mux.Router
}

//Init initialize system
func (a *App) Init() {
	//get repository
	repo = GetMapRepoInstance()
	a.Router = mux.NewRouter()
	a.initRoutes()
}

//Run launch app
func (a *App) Run(addr string) {
	Log.Fatal(http.ListenAndServe(addr, a.Router))
}

//respondWithError wrapper function to prepare for the error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

//respondWithJSON wrapper function for RESTful response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//initRoutes set up routers for all url context
func (a *App) initRoutes() {
	a.Router.HandleFunc("/parkinglot", a.parkinglotPostHandler).Methods("POST")
	a.Router.HandleFunc("/parkinglot", a.parkingLotGetHandler).Methods("GET")
	a.Router.HandleFunc("/parkinglot", a.parkingLotDeleteHandler).Methods("DELETE")
	a.Router.HandleFunc("/checkin", a.checkInHandler).Methods("PUT")
	a.Router.HandleFunc("/checkout", a.checkOutHandler).Methods("PUT")
}

//parkinglotPostHandler handle create parking lot request
//
// swagger:route POST /parkinglot parkinglotPost
//
// Build a parking lot
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       201: parkinglotResponse
//       500: genericError
//
func (a *App) parkinglotPostHandler(w http.ResponseWriter, r *http.Request) {
	//    var plType, small, medium int
	plType, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		Log.Println("use default Express")
		plType = 0 //default Express type
	}
	name := r.FormValue("name")
	addr := r.FormValue("addr")
	small, err := strconv.Atoi(r.FormValue("small"))
	if err != nil {
		Log.Println("No Input small")
		small = 0 //let parkingLotImpl handle the default value
	}
	medium, err := strconv.Atoi(r.FormValue("medium"))
	if err != nil {
		Log.Println("No Input medium")
		medium = 0
	}

	if len(name) <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid Parking Lot Name")
		return
	}
	defer r.Body.Close()

	ret, err := repo.buildParkingLot(plType, name, addr, small, medium)

	fmt.Print(err)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ret)
}

// swagger:route GET /parkinglot parkinglotGet
//
// Retrieve a parking lot information
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: parkinglotResponse
//       500: genericError
//
func (a *App) parkingLotGetHandler(w http.ResponseWriter, r *http.Request) {
	pName := r.FormValue("name")

	defer r.Body.Close()
	p, err := repo.getParkingLot(pName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

// swagger:route DELETE /parkinglot parkinglotDelete
//
// Delete a parking lot. Will fail if the parking lot is not empty.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       201: deleteResponse
//       500: genericError
//
func (a *App) parkingLotDeleteHandler(w http.ResponseWriter, r *http.Request) {
	pName := r.FormValue("name")

	defer r.Body.Close()
	p, err := repo.deleteParkingLot(pName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

// swagger:route PUT /checkin checkin
//
// Car checkin
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       201: ticketIn
//       500: genericError
//
func (a *App) checkInHandler(w http.ResponseWriter, r *http.Request) {
	lot := r.FormValue("lot")
	if len(lot) <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	var c car
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request car payload")
		return
	}

	defer r.Body.Close()
	ret, err := repo.checkIn(&c, lot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, ret)
}

// swagger:route PUT /checkout checkout
//
// Car checkout
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       201: checkoutResponse
//       500: genericError
//
func (a *App) checkOutHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t ticket

	if err := decoder.Decode(&t); err != nil {
		fmt.Println("ERROR:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid Ticket")
		return
	}
	defer r.Body.Close()
	fee, err := repo.checkOut(&t)
	if err != nil {
		fmt.Println("ERROR:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid Ticket")
		return
	}
	respondWithJSON(w, http.StatusCreated, fee)
}
