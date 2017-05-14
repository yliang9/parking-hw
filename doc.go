//This file is for swagger doc only
package main

// swagger:parameters parkinglotPost
type ParkingLotPostParam struct {
	// Parking Lot Type 0 -- Express; 1 -- Daily; 2 -- Value
	// in: query
	// minimum: 0
	// maximum: 2
	Type int `json:"type"`
	// Parking Lot Name
	// in: query
	// required: true
	Name string `json:"name"`
	// Parking Lot Address
	// in: query
	Addr string `json:"addr"`
	// Max space for small car
	// in: query
	Small int `json:"small"`
	// Max space for medium car
	// in: query
	Medium int `json:"medium"`
}

// swagger:parameters parkinglotGet parkinglotDelete
type ParkingLotGetParam struct {
	// Parking Lot Name
	// in: query
	// required: true
	Name string `json:"name"`
}

// swagger:parameters checkin
type CheckInParam struct {
	// Parking Lot Name
	// in: query
	// required: true
	Lot string `json:"lot"`
	// Car JSON
	// in: body
	Car car `json:"car"`
}

// A GenericError is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericError
type GenericError struct {
	// in: body
	Body struct {
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// Return a Parking Lot Object in Json
//
// swagger:response parkinglotResponse
type ResponsePost struct {
	// Parking Lot JSON
	// in: body
	Body parkingLot `json:"body"`
}

// Return a Parking ticket in Json
//
// swagger:response ticketIn
type TicketResponse struct {
	// Parking ticket JSON
	// in: body
	Body ticket `json:"body"`
}

// swagger:parameters checkout
type TicketParam struct {
	// Parking ticket JSON
	// in: body
	Body ticket `json:"body"`
}

// swagger:response deleteResponse
type DeleteResponse struct {
	// true -- deleted, false -- failed
	// in: body
	Body bool `json:"body"`
}

// swagger:response checkoutResponse
type CheckoutResponse struct {
	// Parking fee
	// in: body
	Body int `json:"body"`
}
