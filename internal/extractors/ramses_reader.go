package extractors

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/ramses"
)

// Packet is complete Ramses packet
type Packet struct {
	Header             ramses.Ramses
	OhbseCcsdsTMPacket ramses.OhbseCcsdsTMPacket
	Payload            []byte
}

// StreamBatch tells origin of batch
type StreamBatch struct {
	Buf    io.Reader
	Origin common.OriginDescription
}

//DecodeRamses reads Ramses packages from buffer
func DecodeRamses(recordChannel chan<- common.DataRecord, streamBatch ...StreamBatch) {
	defer close(recordChannel)
	var err error
	for _, stream := range streamBatch {
		for {
			header := ramses.Ramses{}
			err = header.Read(stream.Buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				recordChannel <- common.DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}
				break
			}

			if !header.Valid() {
				err := fmt.Errorf("Not a valid RAC-file")
				recordChannel <- common.DataRecord{Origin: stream.Origin, Error: err, Buffer: []byte{}}
				break
			}

			ccsdsTMPacketHeader := ramses.OhbseCcsdsTMPacket{}
			err = ccsdsTMPacketHeader.Read(stream.Buf)
			if err != nil {
				recordChannel <- common.DataRecord{Origin: stream.Origin, RamsesHeader: header, Error: err, Buffer: []byte{}}
				break
			}
			header.Length -= uint16(binary.Size(ccsdsTMPacketHeader))

			payload := make([]byte, header.Length)
			_, err = stream.Buf.Read(payload)
			if err != nil {
				recordChannel <- common.DataRecord{Origin: stream.Origin, RamsesHeader: header, Error: err, Buffer: []byte{}}
				break
			}
			recordChannel <- common.DataRecord{
				Origin:                   stream.Origin,
				RamsesHeader:             header,
				OhbseCcsdsTMPacketHeader: ccsdsTMPacketHeader,
				Error:                    nil,
				Buffer:                   payload,
			}
		}
	}
}
