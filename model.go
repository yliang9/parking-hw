package main

//var parkingLots  map[string]parkingLot

type car struct {
    Plate string    `json:"plate"` //car license plate number
    Size  int   `json:"size"` //car size: 0 small, 1 medium
}

type ticket struct {
    Plate string `json:"plate"` //car license plate number
    CheckIn string `json:"in"` //time format at time.RFC3339
    LotName string `json:"lotname"`
    Number int `json:"number"` //parking spot number
    SpotType int `json:"spottype"` //0 small, 1 medium
}

type capacity struct {
    Small int `json:"small,omitempty"` //default 50 50 100 Express Daily Value
    Medium int `json:"medium,omitempty"` //default 50 100 150 Express Daily Value
}

type lotSpots struct {
    Small []bool `json:"small"`
    Medium []bool `json:"medium"`
}

type parkingLot struct {
    Name string   `json:"name"`
    Addr string `json:"addr,omitempty"`
    LotType int `json:"lotType"`
    Cap *capacity `json:"cap,omitempty`
    Taken *capacity `json:"cap,omitempty`
    Spots *lotSpots `json:"spots"`
}

const (
    Express = iota
    Daily
    Value
)
