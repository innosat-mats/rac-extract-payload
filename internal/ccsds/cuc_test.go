package ccsds

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestUnsegmentedTimeNanoSeconds(t *testing.T) {
	var tests = []struct {
		coarse uint32
		fine   uint16
		want   int64
	}{
		{0, 0, 0000000000},
		{42, 0, 42000000000},
		{0, 0x8000, 500000000},
		{0, 0b1100000000000000, 750000000},
		{0, 0x8000 >> 7, int64(math.Round(math.Pow(2, -8) * math.Pow10(9)))},
		{42, 0x8000 >> 3, 42000000000 + int64(math.Round(math.Pow(2, -4)*math.Pow10(9)))},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("coarse=%d,fine=%d", tt.coarse, tt.fine), func(t *testing.T) {
			got := UnsegmentedTimeNanoseconds(tt.coarse, tt.fine)
			if got != tt.want {
				t.Errorf("UnsegmentedTime(%d, %d) = %d, want %d", tt.coarse, tt.fine, got, tt.want)
			}
		})
	}
}

func TestUnsegmentedTimeDate(t *testing.T) {
	type args struct {
		coarseTime uint32
		fineTime   uint16
		epoch      time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"Returns Epoch/TAI", args{0, 0, TAI}, TAI},
		{
			"Returns expected time",
			args{10, 0b0100000000000000, TAI},
			TAI.Add(time.Second*10 + time.Millisecond*250),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnsegmentedTimeDate(tt.args.coarseTime, tt.args.fineTime, tt.args.epoch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnsegmentedTimeDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
