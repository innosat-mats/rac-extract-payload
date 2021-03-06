package aez

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"log"
	"path/filepath"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/decodejpeg"
)

func getGrayscaleImage(
	pixels []uint16, width int, height int, shift int, filename string,
) *image.Gray16 {
	nPixels := len(pixels)
	img := image.NewGray16(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	if nPixels != width*height {
		log.Printf(
			"%v: Found %v pixels, but dimension %v x %v says it should be %v\n",
			filename,
			nPixels,
			width,
			height,
			width*height,
		)
	}
	buf := bytes.NewBuffer([]byte{})
	if shift > 0 {
		for idx, pix := range pixels {
			pixels[idx] = pix << shift

		}
	}
	err := binary.Write(buf, binary.BigEndian, pixels)
	if err != nil {
		log.Printf("could not write image data for %v to stream\n", filename)
		return img
	}
	img.Pix = buf.Bytes()
	return img
}

func getGrayscaleImageName(
	originName string,
	imgPackData *CCDImagePackData,
	rid RID,
) string {
	racName := strings.TrimSuffix(filepath.Base(originName), filepath.Ext(originName))
	fileName := fmt.Sprintf("%v_%v_%v.png", racName, imgPackData.Nanoseconds(), rid.CCDNumber())
	return fileName
}

func getImageData(
	buf []byte,
	packData *CCDImagePackData,
	outFileName string,
) []uint16 {
	var imgData []uint16
	var err error
	if packData.JPEGQ != JPEGQUncompressed16bit {
		var height int
		var width int
		imgData, height, width, err = decodejpeg.JpegImageData(buf)
		if err != nil {
			log.Print(err)
			return imgData
		}
		if uint16(height) != packData.NROW || uint16(width) != packData.NCOL+NCOLStartOffset {
			log.Printf(
				"Compressed CCDImage %v has either width %v != %v and/or height %v != %v\n",
				outFileName,
				width,
				packData.NCOL+NCOLStartOffset,
				height,
				packData.NROW,
			)
		}
	} else {
		reader := bytes.NewReader(buf)
		imgData = make([]uint16, reader.Len()/2)
		width, height := int(packData.NCOL+NCOLStartOffset), int(packData.NROW)
		if len(imgData) != width*height {
			log.Printf(
				"Raw CCDImage %v has %v pixels, but dimensions %v x %v\n",
				outFileName,
				len(imgData),
				width,
				height,
			)
		}
		binary.Read(reader, binary.LittleEndian, &imgData)
	}
	return imgData
}
