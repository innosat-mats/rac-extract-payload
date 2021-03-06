// Package ccsds implements part of The Consultative Committee for Space Data Systems (CCSDS)
// "Time Code Formats" (CCSDS 301.0-B-2 - CCSDS 301.0-B-4).
package ccsds

import (
	"math"
	"math/big"
	"time"
)

const nanosPerSecond float64 = 1e9

// TAI is the CCSDS recommended epoch for unsegmented time
var TAI = time.Date(1958, time.January, 1, 0, 0, 0, 0, time.UTC)

// UnsegmentedTimeNanoseconds interprets CCSDS Unsegmented time Code (CUC)
// and returns the nanoseconds since epoch.
//
// It assumes the CUC Time follows the recommendation of having seconds as
// units and follows the CCSDS 301.0-B-2 - CCSDS 301.0-B-4 specification
// in the specification.
//
// * coarseTime is time since epoch in seconds
// * fineTime is subseconds is a binary division of a second
func UnsegmentedTimeNanoseconds(coarseTime uint32, fineTime uint16) int64 {
	var nanos int64 = int64(coarseTime) * int64(nanosPerSecond)
	var fine = big.NewFloat(float64(fineTime))
	fine.SetMantExp(fine, -16)
	fineValue, _ := fine.Float64()
	return nanos + int64(math.Round(fineValue*nanosPerSecond))
}

// UnsegmentedTimeDate iterprets CCSDS Unsegmented time Code (CUC)
// and returns the time in UTC.
//
// It assumes the CUC Time follows the recommendation of having seconds as
// units and follows the CCSDS 301.0-B-2 - CCSDS 301.0-B-4 specification
// in the specification.
//
// * coarseTime is time since epoch in seconds
// * fineTime is subseconds is a binary division of a second
// * epoch is the reference time, such as e.g. TAI
func UnsegmentedTimeDate(coarseTime uint32, fineTime uint16, epoch time.Time) time.Time {
	duration := UnsegmentedTimeNanoseconds(coarseTime, fineTime)
	return epoch.Add(time.Duration(duration))
}
