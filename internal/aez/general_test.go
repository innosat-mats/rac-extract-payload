package aez

import (
	"reflect"
	"testing"
)

func TestRID_String(t *testing.T) {
	tests := []struct {
		name string
		rid  RID
		want string
	}{
		{"RID(0) =''", RID(0), ""},
		{"Stringifies CCD1", CCD1, "CCD1"},
		{"Stringifies CCD2", CCD2, "CCD2"},
		{"Stringifies CCD3", CCD3, "CCD3"},
		{"Stringifies CCD4", CCD4, "CCD4"},
		{"Stringifies CCD5", CCD5, "CCD5"},
		{"Stringifies CCD6", CCD6, "CCD6"},
		{"Stringifies CCD7", CCD7, "CCD7"},
		{"Stringifies PM", PM, "PM"},
		{"Stringifies unknown", RID(42), "Unknown RID: 42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rid.String(); got != tt.want {
				t.Errorf("RID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		rid     RID
		want    []byte
		wantErr bool
	}{
		{
			"Marschals into string representation",
			CCD2,
			[]byte("\"CCD2\""),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rid.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("RID.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RID.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		sid     SID
		want    []byte
		wantErr bool
	}{
		{
			"Marschals into string representation",
			SIDHTR,
			[]byte("\"HTR\""),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sid.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("SID.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SID.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRID_IsCCD(t *testing.T) {
	tests := []struct {
		name string
		rid  RID
		want bool
	}{
		{"CCD1", CCD1, true},
		{"CCD2", CCD2, true},
		{"CCD3", CCD3, true},
		{"CCD4", CCD4, true},
		{"CCD5", CCD5, true},
		{"CCD6", CCD6, true},
		{"CCD7", CCD7, true},
		{"PM", PM, false},
		{"Unkown", RID(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rid.IsCCD(); got != tt.want {
				t.Errorf("RID.IsCCD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRID_CCDNumber(t *testing.T) {
	tests := []struct {
		name string
		rid  RID
		want int32
	}{
		{"CCD1", CCD1, 1},
		{"CCD2", CCD2, 2},
		{"CCD3", CCD3, 3},
		{"CCD4", CCD4, 4},
		{"CCD5", CCD5, 5},
		{"CCD6", CCD6, 6},
		{"CCD7", CCD7, 7},
		{"PM", PM, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rid.CCDNumber(); got != tt.want {
				t.Errorf("RID.CCDNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
