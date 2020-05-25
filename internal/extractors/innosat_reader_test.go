package extractors

import (
	"reflect"
	"testing"

	"github.com/innosat-mats/rac-extract-payload/internal/innosat"
)

func TestDecodeSource(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    SourcePackage
		wantErr bool
	}{
		{
			"empty byte array",
			args{[]byte{}},
			SourcePackage{},
			true,
		},
		{
			"too short byte array",
			args{make([]byte, 2)},
			SourcePackage{},
			true,
		},
		{
			"OK package",
			args{[]byte{0x08, 0x64, 0x88, 0x97, 0x00, 0x41, 0x10, 0x80, 0x19, 0x00, 0x00, 0x12, 0x19,
				0xda, 0x7e, 0x00, 0x17, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x31,
				0x01, 0x2e, 0x01, 0x32, 0x01, 0x30, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2d, 0x01, 0x30,
				0x01, 0x2e, 0x01, 0x30, 0x01, 0x2b, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x32, 0x01, 0x32, 0x01, 0x33,
				0x01, 0x2e, 0x01, 0x31, 0x01, 0x2d, 0x01, 0x2e, 0x01, 0x20, 0xf7}},
			SourcePackage{
				Header:             &innosat.SourcePacketHeader{PacketID: 2148, PacketSequenceControl: 34967, PacketLength: 65},
				Payload:            &innosat.TMHeader{PUS: 16, ServiceType: 128, ServiceSubType: 25, CUCTimeSeconds: 4633, CUCTimeFraction: 55934},
				ApplicationPayload: []byte{0, 23, 1, 47, 1, 49, 1, 47, 1, 47, 1, 47, 1, 49, 1, 46, 1, 50, 1, 48, 1, 47, 1, 47, 1, 47, 1, 45, 1, 48, 1, 46, 1, 48, 1, 43, 1, 47, 1, 49, 1, 50, 1, 50, 1, 51, 1, 46, 1, 49, 1, 45, 1, 46, 1},
			},
			false,
		},
		{
			"Bad Checksum",
			args{[]byte{0x08, 0x64, 0x88, 0x97, 0x00, 0x41, 0x10, 0x80, 0x19, 0x00, 0x00, 0x12, 0x19,
				0xda, 0x7e, 0x00, 0x17, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x31,
				0x01, 0x2e, 0x01, 0x32, 0x01, 0x30, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2f, 0x01, 0x2d, 0x01, 0x30,
				0x01, 0x2e, 0x01, 0x30, 0x01, 0x2b, 0x01, 0x2f, 0x01, 0x31, 0x01, 0x32, 0x01, 0x32, 0x01, 0x33,
				0x01, 0x2e, 0x01, 0x31, 0x01, 0x2d, 0x01, 0x2e, 0x01, 0x00, 0x00}},
			SourcePackage{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeSource(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
