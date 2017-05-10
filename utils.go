package main

import (
    "log"
    "os"
    "time"
    "fmt"
)

var Log *log.Logger

//check if the two timestamp fall in the same day
//(t.Weekday() != 6 && t.Weekday() != 7)
//the condition is 5 < d2 - d1 < 24
//so if d1 > d2, they should be in d1 < Midnight < d2, need to check if d2 > 3am
//else if d1 < d2, check d1 < 3am < d2

func sameDay(d1 string, d2 string, start int) (bool, error) {
    inTime1, err :=time.Parse(time.RFC3339, d1)
    if err != nil {
        return true, err
    }
    inTime2, err :=time.Parse(time.RFC3339, d2)
    if err != nil {
        return true, err
    }

    h1 := inTime1.Hour()
    h2 := inTime2.Hour()
    fmt.Println(h1)
    fmt.Println(h2)
    //here only consider the hour h1 > h2 means h1 < Midnight < h2 (+ 24 * n)
    if h1 > h2 {
        if h2 > start {
            return false, nil
        }
    }
    //that could have the following cases:
    // h1 < h2 < start
    // h1 < start < h2
    // start < h1 < h2
    if h1 < h2 {
        if h1 < start && start < h2 {
            return false, nil
        }
    }

    return true, nil
}

func getDurationHour(d time.Duration) int {
    h := int(d.Hours())
    m := int(d.Minutes())
    add1 := m % 60
    if add1 > 0 {
        h++
    }
    return h
}

func getFeeForExpress(inTime string, outTime time.Time) int {
    in, err :=time.Parse(time.RFC3339, inTime)
    //if there is an error, it is our problem, Free parking
    if err != nil {
        return 0
    }
    duration := outTime.Sub(in)
    h := getDurationHour(duration)
    return 10 * h
}

func getFeeForDaily() int {
    return 50
}

//if > 5, return 5
//TODO check 3am sharp
func getTwoDaysHours(inTime time.Time, outTime time.Time, start int) (int, int) {
    firstRes, secondRes := 0, 0
    firstHour := inTime.Hour()
    //there are 4 cases
    if firstHour < start {
        firstRes = start - firstHour
        //is it the first 15 minutes of the day?
        if firstRes == 1 {
            fMin := inTime.Minute()
            if fMin >= 45 {
                firstRes = 0
            }
        }
    } else {
        firstRes = 24 + start - firstHour
    }
    if firstRes > 5 {
        firstRes = 5
    }
    secondHour := outTime.Hour()
    if secondHour < start {
        secondRes = 24 + start - secondHour
        fMin := outTime.Minute()
        if fMin > 0 {
            secondRes = secondRes + 1
        }
    } else {
        secondRes = secondHour - start
        fMin := outTime.Minute()
        //first 15 minutes of the day?
        if secondRes == 0 {
            if fMin > 15 {
                secondRes++
            }
        } else {
            if fMin > 0 {
                secondRes++
            }
        }
    }
    if secondRes > 5 {
        secondRes = 5
    }
    return firstRes, secondRes
}

func getFeeForValue(inTime string, outTime time.Time, size int) int {
    in, err :=time.Parse(time.RFC3339, inTime)
    //if there is an error, it is our problem, Free parking
    if err != nil {
        return 0
    }
    //free if less than 15 mins
    duration := outTime.Sub(in)
    //only care about the same day minutes different
    m := int(duration.Minutes()) % (60 * 24)
    h := getDurationHour(duration)
    //multiple days?
    mdays := h / 24
    if m <= 15 {
        return 0
    } else { // 15 mins < duration < 5 hours
        if m <= 5 * 60 {
            if size == 0 {
                return 10 * (h + mdays * 5)
            } else {
                return 20 * (h + mdays * 5)
            }
        } else {
            //> 5 hours
            //same day parking?
            sd, _ := sameDay(inTime, outTime.Format(time.RFC3339), 3)
            if sd {
                if size == 0 {
                    return 10 * (5 + mdays * 5)
                } else {
                    return 20 * (5 + mdays * 5)
                }
            } else {
                first, second := getTwoDaysHours(in, outTime, 3)
                if size == 0 {
                    return 10 * (first + second)
                } else {
                    return 20 * (first + second)
                }
            }
        }
    }

    return 20
}

func calcFee(lotType int) int {
    switch lotType {
    case Express:
        return 100
    case Daily:
        return 50
    default:
        return 10
    }
}

func getDefaultSpace(lotType int, smallSize bool) (num int) {
    switch lotType {
    case Express:
        return 50; //default 50 for either small or medium
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

func setLog(logfile string) (file *os.File)  {
	file, err := os.OpenFile(logfile, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
    return file
}
