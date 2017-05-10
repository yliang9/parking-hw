package main

import (
    "testing"
    "time"
    "fmt"
)

var a App

func TestCalcFee(t *testing.T) {
    total := calcFee(0)
    if total != 100 {
       t.Errorf("Fee was incorrect, got: %d, want: %d.", total, 10)
    }
}

func TestGetDefaultSpace(t *testing.T) {
    //test
    space := getDefaultSpace(0, true)
    if space != 50 {
        t.Errorf("Default available Small spots was incorrect, got: %d, want: %d.", space, 50)
    }
    //TODO add more
}

func TestSameDay(t *testing.T) {
    res1, _ := sameDay("2017-05-07T02:14:37-04:00", "2017-05-08T14:14:37-04:00", 3)
    if res1 {
        t.Errorf("Expected 2 days, but Not")
    }
    res2, _ := sameDay("2017-05-07T20:14:37-04:00", "2017-05-08T02:14:37-04:00", 3)
    if !res2 {
        t.Errorf("Expected 1 day, but Not")
    }
    res3, _ := sameDay("2017-05-07T20:14:37-04:00", "2017-05-08T20:14:37-04:00", 3)
    if !res3 {
        t.Errorf("Expected 2 days, but Not")
    }
    res4, _ := sameDay("2017-05-07T20:14:37-04:00", "2017-05-08T23:14:37-04:00", 3)
    if !res4 {
        t.Errorf("Expected 2 days, but Not")
    }

    //test invalid inputs
    _, err := sameDay("2017-05-07", "2017-05-08T23:14:37-04:00", 3)
    if err == nil {
        t.Errorf("Error " + err.Error())
    }
}

func TestGetFeeForExpress(t *testing.T) {
    t1 , _:= time.Parse(time.RFC3339, "2017-05-07T20:50:37-04:00")
    res1 := getFeeForExpress("2017-05-07T20:14:37-04:00", t1)
    if (res1 != 10) {
        t.Errorf("Expected 10, but got %d", res1)
    }
    fmt.Println(t1)
    res2 := getFeeForExpress("2017-05-06T20:50:37-04:00", t1)
    if (res2 != 240) {
        t.Errorf("Expected 240, but got %d", res1)
    }
}

func TestGetFeeForValue(t *testing.T) {
    t1 , _:= time.Parse(time.RFC3339, "2017-05-07T20:50:37-04:00")
    res1 := getFeeForValue("2017-05-07T20:35:37-04:00", t1, 0)
    if (res1 != 0) {
        t.Errorf("Expected 0, but got %d", res1)
    }
    res1 = getFeeForValue("2017-05-07T20:35:36-04:00", t1, 0)
    if (res1 != 0) {
        t.Errorf("Expected 0, but got %d", res1)
    }
    //charge by min, not second
    res1 = getFeeForValue("2017-05-07T20:34:36-04:00", t1, 1)
    if (res1 != 20) {
        t.Errorf("Expected 0, but got %d", res1)
    }
    res1 = getFeeForValue("2017-05-06T03:00:00-04:00", t1, 0)
    if (res1 != 100) {
        t.Errorf("Expected 0, but got %d", res1)
    }
    t1 , _ = time.Parse(time.RFC3339, "2017-05-08T04:50:00-04:00")
    res1 = getFeeForValue("2017-05-07T19:00:00-04:00", t1, 0)
    if (res1 != 70) {
        t.Errorf("Expected 0, but got %d", res1)
    }

}
