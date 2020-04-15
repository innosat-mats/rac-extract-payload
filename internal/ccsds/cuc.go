// Package implements part of The Consultative Committee for Space Data Systems (CCSDS)
// "Time Code Formats" (CCSDS 301.0-B-2 / CCSDS 301.0-B-4).
package ccsds

import (
	"math"
	"time"
)

var milliesPerSecond float64 = math.Pow10(9)

// TAI is the CCSDS recommended epoch for unsegmented time
var TAI = time.Date(1958, time.January, 1, 0, 0, 0, 0, time.UTC)

// UnsegmentedTimeNanoseconds interprets CCSDS Unsegmented time Code (CUC)
// and returns the nanoseconds since epoch.
//
// It assumes the CUC Time follows the recommendation of having seconds as
// units and follows the CCSDS 301.0-B-2 / CCSDS 301.0-B-4 specification
// in the specification.
//
// * coarseTime is time since epoch in seconds
// * fineTime is subseconds is a binary division of a second
func UnsegmentedTimeNanoseconds(coarseTime uint32, fineTime uint16) int64 {
	var nanos int64 = int64(coarseTime) * int64(milliesPerSecond)
	var fine float64 = 0
	var flag uint16 = 1
	for pos := 0; pos < 16; pos++ {
		if (fineTime & flag) == flag {
			fine += math.Pow(2, -float64(pos))
		}
		flag = flag << 1
	}
	nanos += int64(math.Round(fine * milliesPerSecond))
	return nanos
}

// UnsegmentedTimeDate iterprets CCSDS Unsegmented time Code (CUC)
// and returns the time in UTC.
//
// It assumes the CUC Time follows the recommendation of having seconds as
// units and follows the CCSDS 301.0-B-2 / CCSDS 301.0-B-4 specification
// in the specification.
//
// It sets epoch as TAI (1958-01-01)
//
// * coarseTime is time since epoch in seconds
// * fineTime is subseconds is a binary division of a second
func UnsegmentedTimeDate(coarseTime uint32, fineTime uint16) time.Time {
	duration := UnsegmentedTimeNanoseconds(coarseTime, fineTime)
	return TAI.Add(time.Duration(duration))
}