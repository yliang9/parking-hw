package main

import (
	"testing"
	"time"
)

type testCase struct {
	in      string
	out     string
	size    int
	lotType int
	expect  int
}

func TestCalcFee(t *testing.T) {
	testCases := []testCase{
		//Express ----
		//< 1hour, $5
		{"2017-05-07T20:14:37-04:00", "2017-05-07T20:15:38-04:00", 0, 0, 5},
		//24 hours 5 * 24 = 120
		{"2017-05-06T20:14:37-04:00", "2017-05-07T20:14:37-04:00", 0, 0, 120},

		//Daily ---
		//1 second
		{"2017-05-07T20:14:37-04:00", "2017-05-07T20:14:38-04:00", 0, 1, 50},
		//24 hours exactly
		{"2017-05-07T20:14:37-04:00", "2017-05-08T20:14:37-04:00", 0, 1, 50},
		//12 days
		{"2017-05-07T20:14:37-04:00", "2017-05-18T20:14:38-04:00", 0, 1, 600},

		//sat, sun

		//Value ---
		//test >5 hours in the same day, Sat.
		{"2017-05-06T20:14:37-04:00", "2017-05-07T02:29:37-04:00", 1, 2, 50},
		//sat 4 hours
		{"2017-05-06T23:14:37-04:00", "2017-05-07T02:29:37-04:00", 1, 2, 40},
		//less than 15 Minutes
		{"2017-05-06T02:50:37-04:00", "2017-05-06T03:05:37-04:00", 1, 2, 0},
		//in 2 workdays, 15 minutes + 1 sec, treated as 1 hour
		{"2017-05-10T02:50:37-04:00", "2017-05-10T03:05:38-04:00", 1, 2, 20},
		//15 minutes + 1 sec in the same day
		{"2017-05-09T02:44:59-04:00", "2017-05-09T03:00:00-04:00", 1, 2, 20},
		//4 days, but 1 Sat, 1 Sun, Fri full day, Thur > 5 hours 2 * 5 * 20 + 5 * 20 / 2 = 250
		{"2017-05-04T10:14:37-04:00", "2017-05-07T20:29:36-04:00", 1, 2, 250},
		//test case provided by doc
		{"2017-05-04T19:00:37-04:00", "2017-05-05T04:50:37-04:00", 0, 2, 70},
		//one month
		//21 full workday, 5 Sat. > 5 on 05-04, 2 hours on 06-05
		//21 * 5 * 10 + 5 * 5 * 10 / 2 + 5 * 10 + 2 * 10 = 1245
		{"2017-05-04T19:00:37-04:00", "2017-06-05T04:50:37-04:00", 0, 2, 1245},
		//21 full workday, 5 Sat. but less than 5 hours on start and end day
		//21 * 5 * 10 + 5 * 5 * 10 / 2 + 4 * 10 = 1215
		{"2017-05-04T23:59:37-04:00", "2017-06-05T03:50:37-04:00", 0, 2, 1215},
	}

	for _, tc := range testCases {
		ti := ticket{}
		ti.CheckIn = tc.in
		ti.LotType = tc.lotType
		ti.SpotType = tc.size
		out, _ := time.Parse(time.RFC3339, tc.out)
		fee := calcFee(ti, out)
		if fee != tc.expect {
			t.Errorf("Expected %d, but got %d.", tc.expect, fee)
		}
	}
}

func TestGetFeeForValue(t *testing.T) {
	testCases := []testCase{
		// this test case doesn't consider 3 am offset, is based on 12 am for the day
		//sat, sun
		//it will be 4 hours Sat 4 * 20 / 2 = 40
		{"2017-05-06T20:14:37-04:00", "2017-05-07T20:14:38-04:00", 1, 2, 40},
		//full saturday
		//it will be max 5 hours Sat, 5 * 20 / 2 = 50
		{"2017-05-06T00:00:00-04:00", "2017-05-07T00:00:00-04:00", 1, 2, 50},
		//full saturday + some hours in sunday
		//it will be max 5 hours Sat, 5 * 20 / 2 = 50
		{"2017-05-06T00:00:00-04:00", "2017-05-07T20:00:00-04:00", 1, 2, 50},
		//some hours in friday, some in saturday
		//it will be 1 hour Fri, 4 hours sat. should be 20 * 1 + 4 * 20 / 2 = 60
		{"2017-05-05T23:59:59-04:00", "2017-05-06T03:59:59-04:00", 1, 2, 60},
		// weekdays cross time, but total 4 hours (4 * 20)
		{"2017-05-17T23:59:59-04:00", "2017-05-18T03:59:59-04:00", 1, 2, 80},
		// 2 sat, + 8 workdays + 4 hours workday 2 * 20 * 5 / 2 + 8 * 20 * 5 + 4 * 20 = 980
		{"2017-05-05T23:59:59-04:00", "2017-05-18T03:59:59-04:00", 1, 2, 980},
		//1 hour Fri, + 4 hours Sat 20 + 4 * 20 / 2 = 60
		{"2017-05-05T23:59:59-04:00", "2017-05-06T03:59:59-04:00", 1, 2, 60},
	}
	for _, tc := range testCases {
		ti := ticket{}
		ti.CheckIn = tc.in
		ti.LotType = tc.lotType
		ti.SpotType = tc.size
		out, _ := time.Parse(time.RFC3339, tc.out)
		in, _ := time.Parse(time.RFC3339, tc.in)
		fee := getFeeForMultipleDays(in, out, tc.size)
		//fee := calcFee(ti, out)
		if fee != tc.expect {
			t.Errorf("Expected %d, but got %d.", tc.expect, fee)
		}
	}
}

func TestGetDefaultSpace(t *testing.T) {
	//test
	space := getDefaultSpace(0, true)
	if space != 50 {
		t.Errorf("Default available Small spots was incorrect, got: %d, want: %d.", space, 50)
	}
	space = getDefaultSpace(1, true)
	if space != 50 {
		t.Errorf("Default available Small spots was incorrect, got: %d, want: %d.", space, 50)
	}

}

func TestGetFeeForExpress(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2017-05-07T20:50:37-04:00")
	res1 := getFeeForExpress("2017-05-07T20:14:37-04:00", t1)
	if res1 != 5 {
		t.Errorf("Expected 10, but got %d", res1)
	}
	res2 := getFeeForExpress("2017-05-06T20:50:37-04:00", t1)
	if res2 != 120 {
		t.Errorf("Expected 240, but got %d", res2)
	}
}
