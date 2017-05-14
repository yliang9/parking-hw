package main

import "time"

//calculate the parking fee based on the ticket
func calcFee(t ticket, out time.Time) int {
	//get lot type
	switch t.LotType {
	case Express:
		return getFeeForExpress(t.CheckIn, out)
	case Daily:
		return getFeeForDaily(t.CheckIn, out)
	default:
		return getFeeForValue(t.CheckIn, out, t.SpotType)
	}
}

func sameDay(in time.Time, out time.Time) bool {
	if in.Year() == out.Year() && in.Month() == out.Month() && in.Day() == out.Day() {
		return true
	}
	return false
}

func getDurationHour(d time.Duration) int {
	h := int(d.Hours())
	m := int(d.Minutes())
	s := int(d.Seconds())
	add1 := m % 60
	add2 := s % 60
	if add1 > 0 || add2 > 0 {
		h++
	}
	return h
}

func getHoursForExpressOrDaily(inTime string, outTime time.Time) int {
	in, err := time.Parse(time.RFC3339, inTime)
	//if there is an error, it is our problem, Free parking
	if err != nil {
		return 0
	}
	duration := outTime.Sub(in)
	h := getDurationHour(duration)
	return h
}

func getFeeForExpress(in string, out time.Time) int {
	h := getHoursForExpressOrDaily(in, out)
	return ExpressRate * h
}

func getFeeForDaily(in string, out time.Time) int {
	h := getHoursForExpressOrDaily(in, out)
	d := h / 24
	if h%24 > 0 {
		d++
	}
	return d * DailyRate
}

func dayBefore(in time.Time, out time.Time) bool {
	if in.Year() < out.Year() {
		return true
	}
	if in.Month() < out.Month() {
		return true
	}
	if in.Day() < out.Day() {
		return true
	} else {
		return false
	}
}

func getFeeForMultipleDays(in time.Time, out time.Time, size int) int {
	//a days struct to store how many weekday, weekend passed during (in, out)  only
	type dayCounter struct {
		wd  int
		sat int
		sun int
	}
	days := dayCounter{0, 0, 0}
	//count the days between in and out days exclusively
	var in2 time.Time
	for in2 = in.Add(time.Hour * time.Duration(24)); dayBefore(in2, out); in2 = in2.Add(time.Hour * time.Duration(24)) {
		switch in2.Weekday() {
		case 6:
			days.sat++
		case 0:
			days.sun++
		default:
			days.wd++
		}
	}
	//fmt.Println(days)
	//go back one day
	in2 = in2.Add(-time.Hour * time.Duration(24))
	//there is one special case that we may over charge the user
	//although they are in 2 separate days, but in and out are both work day, and duration <= 5 hours.
	//in this case, the combinaton of 2 separate hours may > 5, user will get mad :(
	dur2 := out.Sub(in2)
	s2 := int(dur2.Seconds())
	if s2 <= 5*60*60 && in.Weekday() < 6 && in.Weekday() > 0 && out.Weekday() < 6 && out.Weekday() > 0 {
		rh := s2 / (60 * 60)
		if s2%(60*60) > 0 {
			rh++
		}
		if size == SMALL {
			return rh*VS + days.wd*5*VS + days.sat*5*VS/2
		} else {
			return rh*VM + days.wd*5*VM + days.sat*5*VM/2
		}
	}
	//get separated hours
	h1 := 24 - in.Hour()
	h2 := out.Hour()
	if out.Second() > 0 {
		h2++
	}
	if h1 > 5 {
		h1 = 5
	}
	if h2 > 5 {
		h2 = 5
	}
	//get fee
	var fee1, fee2 int
	if in.Weekday() == time.Sunday {
		fee1 = 0
	} else {
		fee1 = calcFee2ByHour(h1, in.Weekday() == time.Saturday, size)
	}
	if out.Weekday() == time.Sunday {
		fee2 = 0
	} else {
		fee2 = calcFee2ByHour(h2, out.Weekday() == time.Saturday, size)
	}
	if size == SMALL {
		return fee1 + fee2 + days.wd*5*VS + days.sat*5*VS/2
	} else {
		return fee1 + fee2 + days.wd*5*VM + days.sat*5*VM/2
	}

}

func calcFee2ByHour(h int, sat bool, size int) int {
	res := 0
	if size == SMALL {
		res = h * VS
	} else {
		res = h * VM
	}
	if sat {
		res = res / 2
	}
	return res
}

func getFeeForSameDay(in time.Time, out time.Time, s int, size int) int {
	//sunday?
	if in.Weekday() == time.Sunday {
		return 0
	}
	// > 5 hours
	var h int
	if s > 5*60*60 {
		h = 5
	} else {
		h = s / (60 * 60)
		ss := s % (60 * 60)
		if ss > 0 {
			h++
		}
	}

	return calcFee2ByHour(h, in.Weekday() == time.Saturday, size)
}

func getFeeForValue(inTime string, outTime time.Time, size int) int {
	in, err := time.Parse(time.RFC3339, inTime)
	//if there is an error, it is our problem, Free parking
	if err != nil {
		return 0
	}
	//offset -3 hours
	in2 := in.Add(-time.Hour * time.Duration(VStart))
	out2 := outTime.Add(-time.Hour * time.Duration(VStart))

	duration := out2.Sub(in2)
	s := int(duration.Seconds())

	//free for the first 15 minutes
	if s <= 15*60 {
		return 0
	} else {
		// in the same day?
		sd := sameDay(in2, out2)
		if sd {
			return getFeeForSameDay(in2, out2, s, size)
		} else {
			return getFeeForMultipleDays(in2, out2, size)
		}
	}
}

//use the default value if no input parameter for how many spots for a parking lot
func getDefaultSpace(lotType int, smallSize bool) (num int) {
	switch lotType {
	case Express:
		return 50 //default 50 for either small or medium
	case Daily:
		if smallSize {
			return 50
		} else {
			return 100
		}
	default:
		if smallSize {
			return 100
		} else {
			return 150
		}
	}
}
