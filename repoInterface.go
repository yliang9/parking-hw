//go:generate swagger generate spec
package main

//repository is the interface that defines all the basic parking lot methods
//repository can be implemented by different classes as polymorphism
//in this project, I use map object as persistence, it could be easily switch to support Database
type repository interface {
    //buildParkingLot create a parking lot, return the created parking lot or error if failed
    //plType: parking lot type, EXPRESS/DAILY/VALUE, optional, default EXPRESS
    //name: parking lot name, not nil
    //addr: parking lot address, optional
    //small: number of parking spot for small car
    //medium: number of parking spot for medium car
    buildParkingLot(plType int, name string, addr string, small int, medium int) (parkingLot, error)
    //getParkingLot find parking lot by name, return error is not found
    //name: parking lot name
    getParkingLot(name string) (parkingLot, error)
    //deleteParkingLot delete a parking lot by name
    //name: parking lot name
    deleteParkingLot(name string) (bool, error)
    //checkIn car check in, return a ticket which include time and other parking information
    //mycar: car information includes plate number, size
    //lot: parking lot name
    checkIn(mycar * car, lot string) (ticket, error)
    //checkOut car check out, return a parking fee in int format
    //t: parking ticket received when checkIn
    checkOut(t *ticket) (int, error) //use int for money for now, since the int type meet the design doc requirement
}
