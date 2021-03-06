package main

//implementation of the map repository
//use Singleton pattern and Mutex to provent concurrent write to the parking lot
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//a map object to hold the parking lots
type mapRepo struct {
	parkingLots map[string]parkingLot
	mu          sync.RWMutex
}

var instance *mapRepo
var once sync.Once

//a singleton
func GetMapRepoInstance() *mapRepo {
	once.Do(func() {
		instance = &mapRepo{
			parkingLots: make(map[string]parkingLot),
		}
	})
	return instance
}

//implement repoInterface
func (r *mapRepo) buildParkingLot(plType int, name string, addr string, small int, medium int) (parkingLot, error) {
	//if empty, set default values
	p := parkingLot{}
	p.Name = name
	p.Addr = addr
	p.LotType = plType
	if small == 0 {
		small = getDefaultSpace(plType, true)
	}
	if medium == 0 {
		medium = getDefaultSpace(plType, false)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p0, exist := r.parkingLots[name]

	if exist {
		return p0, errors.New(fmt.Sprintf("Parking Lot '%s' Exists", name))
	}

	cap := capacity{}
	cap.Small = small
	cap.Medium = medium
	p.Cap = &cap

	count := capacity{}
	count.Small = 0
	count.Medium = 0
	p.Taken = &count

	spots := lotSpots{}
	spots.Small = make([]bool, small)
	spots.Medium = make([]bool, medium)
	p.Spots = &spots

	r.parkingLots[name] = p
	return p, nil
}

//implements repoInterface
func (r *mapRepo) getParkingLot(name string) (parkingLot, error) {
	var lot parkingLot
	if len(name) <= 0 {
		return lot, errors.New("Invalid input")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	lot, exist := r.parkingLots[name]
	if !exist {
		return lot, errors.New(fmt.Sprintf("Parking Lot '%s' does not exists", name))
	}
	return lot, nil
}

//assume that we cannot delete a parking lot if not empty
//implements repoInterface
func (r *mapRepo) deleteParkingLot(name string) (bool, error) {
	p, err := r.getParkingLot(name)
	if err != nil {
		return false, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.Taken.Small > 0 || p.Taken.Medium > 0 {
		return false, errors.New(fmt.Sprintf("Parking Lot '%s' not Empty", name))
	}
	delete(r.parkingLots, name)
	return true, nil
}

//implements repoInterface
func (r *mapRepo) checkIn(mycar *car, lotName string) (ticket, error) {
	t := ticket{}
	t.Plate = mycar.Plate
	lot, err := r.getParkingLot(lotName)
	if err != nil {
		return t, err
	}
	//check Size and Available Spots
	switch mycar.Size {
	case LARGE:
		return t, errors.New("Your car is too big to park here")
	case MEDIUM:
		if lot.Cap.Medium == lot.Taken.Medium {
			return t, errors.New("Lot Full")
		}
	default:
		//assume that small car can park in medium lot
		if lot.Cap.Small == lot.Taken.Small && lot.Cap.Medium == lot.Taken.Medium {
			return t, errors.New("Lot Full")
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	//find the first Spot
	if mycar.Size == MEDIUM {
		for i, val := range lot.Spots.Medium {
			if !val {
				lot.Spots.Medium[i] = true
				t.Number = i
				t.SpotType = MEDIUM
				lot.Taken.Medium++
				break
			}
		}
	} else if mycar.Size == SMALL {
		//check if there is space for small
		if lot.Taken.Small < lot.Cap.Small {
			for i, val := range lot.Spots.Small {
				if !val {
					lot.Spots.Small[i] = true
					t.Number = i
					t.SpotType = SMALL
					lot.Taken.Small++
					break
				}
			}
		} else {
			for i, val := range lot.Spots.Medium {
				if !val {
					lot.Spots.Medium[i] = true
					t.Number = i
					t.SpotType = MEDIUM
					lot.Taken.Medium++
					break
				}
			}
		}
	}

	r.parkingLots[lot.Name] = lot
	t.CheckIn = time.Now().Format(time.RFC3339)
	t.LotName = lotName
	t.LotType = lot.LotType
	return t, nil
}

//implements repoInterface
func (r *mapRepo) checkOut(t *ticket) (int, error) {
	lot, err := r.getParkingLot(t.LotName)
	if err != nil {
		return 0, err
	}
	out := time.Now()
	//clear parking spot
	if t.SpotType == SMALL {
		if lot.Taken.Small == 0 {
			return 0, errors.New("lot is alreay empty, invalid ticket")
		}
		lot.Taken.Small--
		lot.Spots.Small[t.Number] = false
	} else {
		if lot.Taken.Medium == 0 {
			return 0, errors.New("lot is alreay empty, invalid ticket")
		}
		lot.Taken.Medium--
		lot.Spots.Medium[t.Number] = false
	}
	fee := calcFee(*t, out)
	return fee, nil
}
