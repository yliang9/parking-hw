package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "encoding/json"
    "fmt"
    "strconv"
)

//use in memory map to store the packing lots, with Factory Pattern, could be easily extended to add other repos, such as a real database
var repo *mapRepo

type App struct {
    Router *mux.Router
}

func (a *App) Init() {
    //get repository
    repo = GetMapRepoInstance()
    a.Router = mux.NewRouter()
    a.initRoutes()
}

func (a *App) Run(addr string) {
    Log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) initRoutes() {
    a.Router.HandleFunc("/parkinglot", a.parkinglotPostHandler).Methods("POST")
    a.Router.HandleFunc("/parkinglot", a.parkingLotGetHandler).Methods("GET")
    a.Router.HandleFunc("/checkin", a.checkInHandler).Methods("PUT")
    a.Router.HandleFunc("/checkout", a.checkOutHandler).Methods("PUT")
}

func (a *App) parkinglotPostHandler(w http.ResponseWriter, r *http.Request) {
//    var plType, small, medium int
    plType, err := strconv.Atoi(r.FormValue("lotType")[:0])
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

func (a *App) parkingLotGetHandler(w http.ResponseWriter, r *http.Request) {
    pName := r.FormValue("name");

    defer r.Body.Close()
    p, err := repo.getParkingLot(pName)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) checkInHandler(w http.ResponseWriter, r *http.Request) {
    lot := r.FormValue("lot")
    if (len(lot) <= 0) {
        respondWithError(w, http.StatusBadRequest, "Invalid input")
        return
    }
    var c car
    fmt.Println(r.Body)
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

func (a *App) checkOutHandler(w http.ResponseWriter, r *http.Request) {
    plate := r.FormValue("plate")
    decoder := json.NewDecoder(r.Body)

    var t ticket

    if (len(plate) <=0) {
        respondWithError(w, http.StatusBadRequest, "Invalid input")
        return
    }

    if err := decoder.Decode(&t); err != nil {
        fmt.Println("ERROR:", err)
        respondWithError(w, http.StatusBadRequest, "Invalid Ticket")
        return
    }
    defer r.Body.Close()

    fee := calcFee(0)
    respondWithJSON(w, http.StatusCreated, fee)
}
