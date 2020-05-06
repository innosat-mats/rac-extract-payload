package extractors

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/aez"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

func Test_instrumentHK(t *testing.T) {
	type args struct {
		sid aez.SID
		buf io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    common.Exportable
		wantErr bool
	}{
		{
			"STAT",
			args{sid: aez.SIDSTAT, buf: bytes.NewReader(make([]byte, 100))},
			aez.STAT{},
			false,
		},
		{
			"HTR",
			args{sid: aez.SIDHTR, buf: bytes.NewReader(make([]byte, 100))},
			aez.HTR{},
			false,
		},
		{
			"PWR",
			args{sid: aez.SIDPWR, buf: bytes.NewReader(make([]byte, 100))},
			aez.PWR{},
			false,
		},
		{
			"CPRUA",
			args{sid: aez.SIDCPRUA, buf: bytes.NewReader(make([]byte, 100))},
			aez.CPRU{},
			false,
		},
		{
			"CPRUB",
			args{sid: aez.SIDCPRUB, buf: bytes.NewReader(make([]byte, 100))},
			aez.CPRU{},
			false,
		},
		{"Unknown", args{sid: aez.SID(0)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instrumentHK(tt.args.sid, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("instrumentHK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instrumentHK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instrumentTransparentData(t *testing.T) {
	type args struct {
		rid aez.RID
		buf io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    common.Exportable
		wantErr bool
	}{
		{
			"CCD1",
			args{rid: aez.CCD1, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD2",
			args{rid: aez.CCD2, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD3",
			args{rid: aez.CCD3, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD4",
			args{rid: aez.CCD4, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD5",
			args{rid: aez.CCD5, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD6",
			args{rid: aez.CCD6, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"CCD7",
			args{rid: aez.CCD7, buf: bytes.NewReader(make([]byte, 100))},
			aez.CCDImage{BadColumns: []uint16{}},
			false,
		},
		{
			"PM",
			args{rid: aez.PM, buf: bytes.NewReader(make([]byte, 100))},
			nil,
			true,
		},
		{"Unknown", args{rid: aez.RID(0)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instrumentTransparentData(tt.args.rid, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("instrumentTransparentData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instrumentTransparentData() = %v, want %v", got, tt.want)
			}
		})
	}
}