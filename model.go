package main

//This file includes all the data object for parking billing system

//car
type car struct {
	Plate string `json:"plate"` //car license plate number
	Size  int    `json:"size"`  //car size: 0 small, 1 medium
}

//ticket
type ticket struct {
	Plate    string `json:"plate"` //car license plate number
	CheckIn  string `json:"in"`    //time format at time.RFC3339
	LotName  string `json:"lotname"`
	Number   int    `json:"number"`   //parking spot number
	SpotType int    `json:"spottype"` //0 small, 1 medium
}

//capacity all parking spots for a parking lot
// swagger:model
type capacity struct {
	Small  int `json:"small,omitempty"`  //default 50 50 100 Express Daily Value
	Medium int `json:"medium,omitempty"` //default 50 100 150 Express Daily Value
}

//lotSpots parking lot spots status, false, true = empty, taken
// swagger:model
type lotSpots struct {
	Small  []bool `json:"small"`
	Medium []bool `json:"medium"`
}

//parkingLot a parking lot
// swagger:response parkinglot
type parkingLot struct {
	//parking lot name
	Name string `json:"name"`
	//parking lot address
	Addr string `json:"addr,omitempty"`
	//parking lot type, 0/1/2 EXPRESS/DAILY/VALUE
	LotType int `json:"lotType"`
	//total number of the parking spots
	Cap *capacity `json:"cap,omitempty`
	//current taken number of the parking spots
	Taken *capacity `json:"cap,omitempty`
	//current status of the parking spots
	Spots *lotSpots `json:"spots"`
}

const (
	Express = iota
	Daily
	Value
)

const (
	SMALL = iota
	MEDIUM
	LARGE
)
