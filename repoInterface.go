package main

type repository interface {
    buildParkingLot(plType int, name string, addr string, small int, medium int) (parkingLot, error)
    getParkingLot(name string) (parkingLot, error)
    deleteParkingLot(name string) (bool, error)
    checkIn(mycar * car, lot string) (ticket, error)
    checkOut(t *ticket) (int, error) //use int for money for now, since the int type meet the design doc requirement
}
