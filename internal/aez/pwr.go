package aez

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/innosat-mats/rac-extract-payload/internal/parquetrow"
)

type pwr uint16

var pwrTemperatures = [...]float64{
	-55, -50, -45, -40, -35, -30, -25, -20, -15, -10,
	-5, 0, 5, 10, 15, 20, 25, 30, 35, 40,
	45, 50, 55, 60, 65, 70, 75, 80, 85, 90,
	95, 100, 105, 110, 115, 120, 125, 130, 135, 140,
	145, 150, 155, 160, 165, 170, 175, 180, 185, 190,
	195, 200, 205, 210, 215, 220, 225, 230, 235, 240,
	245, 250,
} // ⁰C
var pwrResistances = [...]float64{
	1.06208e+05, 7.86360e+04, 5.86500e+04, 4.40600e+04, 3.33320e+04,
	2.53920e+04, 1.94502e+04, 1.50342e+04, 1.16706e+04, 9.13720e+03,
	7.21000e+03, 5.73300e+03, 4.58140e+03, 3.68760e+03, 2.98400e+03,
	2.43080e+03, 2.00000e+03, 1.65952e+03, 1.37270e+03, 1.14206e+03,
	9.60300e+02, 8.10900e+02, 6.83400e+02, 5.79040e+02, 4.94280e+02,
	4.23660e+02, 3.63880e+02, 3.13600e+02, 2.71840e+02, 2.36440e+02,
	2.06800e+02, 1.81482e+02, 1.59284e+02, 1.40204e+02, 1.23778e+02,
	1.09570e+02, 9.74120e+01, 8.68300e+01, 7.74440e+01, 6.92300e+01,
	6.20960e+01, 5.58200e+01, 5.03860e+01, 4.55800e+01, 4.13340e+01,
	3.75600e+01, 3.41800e+01, 3.11640e+01, 2.84540e+01, 2.60240e+01,
	2.38680e+01, 2.19280e+01, 2.02000e+01, 1.86382e+01, 1.71898e+01,
	1.58768e+01, 1.46822e+01, 1.35960e+01, 1.26174e+01, 1.17246e+01,
	1.08974e+01, 1.01410e+01,
} // Ohm

func (data *pwr) voltageADC() float64 {
	return voltageConstant * float64(*data)
}

type pwrt pwr

func (data *pwrt) voltage() float64 {
	pwr := pwr(*data)
	return pwr.voltageADC()
}

func (data *pwrt) resistance() float64 {
	return 3.3*1000/data.voltage() - 1000
}

type pwrp32v pwr

func (data *pwrp32v) voltage() float64 {
	pwr := pwr(*data)
	return 21 * pwr.voltageADC()
}

type pwrp32c pwr

func (data *pwrp32c) current() float64 {
	pwr := pwr(*data)
	return 10.1 / 100 * pwr.voltageADC()
}

type pwrp16v pwr

func (data *pwrp16v) voltage() float64 {
	pwr := pwr(*data)
	return 11 * pwr.voltageADC()
}

type pwrp16c pwr

func (data *pwrp16c) current() float64 {
	pwr := pwr(*data)
	return 10.1 / 5 * pwr.voltageADC()
}

type pwrm16v pwr

func (data *pwrm16v) voltage() float64 {
	pwr := pwr(*data)
	return -10 * pwr.voltageADC()
}

type pwrm16c pwr

func (data *pwrm16c) current() float64 {
	pwr := pwr(*data)
	return 10.1 / 100 * pwr.voltageADC()
}

type pwrp3v3 pwr

func (data *pwrp3v3) voltage() float64 {
	pwr := pwr(*data)
	return 4 * pwr.voltageADC()
}

type pwrp3c3 pwr

func (data *pwrp3c3) current() float64 {
	pwr := pwr(*data)
	return 10.1 / 20 * pwr.voltageADC()
}

func (data *pwrt) temperature() (float64, error) {
	return Interpolate(
		data.resistance(),
		pwrResistances[:],
		pwrTemperatures[:],
	)
}

// PWR structure 18 octext
type PWR struct {
	PWRT    pwrt    // Temp. sense 0..4095
	PWRP32V pwrp32v // +32V voltage sense 0..4095
	PWRP32C pwrp32c // +32V current sense 0..4095
	PWRP16V pwrp16v // +16V voltage sense 0..4095
	PWRP16C pwrp16c // +16V current sense 0..4095
	PWRM16V pwrm16v // -16V voltage sense 0..4095
	PWRM16C pwrm16c // -16V current sense 0..4095
	PWRP3V3 pwrp3v3 // +3V3 voltage sense 0..4095
	PWRP3C3 pwrp3c3 // +3V3 current sense 0..4095
}

// PWRReport structure in useful units
type PWRReport struct {
	PWRT     float64 // Temp. sense ⁰C
	PWRP32V  float64 // +32V voltage sense voltage
	PWRP32C  float64 // +32V current sense current
	PWRP16V  float64 // +16V voltage sense voltage
	PWRP16C  float64 // +16V current sense current
	PWRM16V  float64 // -16V voltage sense voltage
	PWRM16C  float64 // -16V current sense current
	PWRP3V3  float64 // +3V3 voltage sense voltage
	PWRP3C3  float64 // +3V3 current sense current
	Warnings []error
}

// NewPWR reads a PWR from buffer
func NewPWR(buf io.Reader) (*PWR, error) {
	pwr := PWR{}
	err := binary.Read(buf, binary.LittleEndian, &pwr)
	return &pwr, err
}

//lint:ignore U1000 unused but important for documentation
func pwrVoltageADC(data uint16) float64 {
	return voltageConstant * float64(data)
}

// Report returns a PWRReport with useful units
func (pwr *PWR) Report() PWRReport {
	temp, err := pwr.PWRT.temperature()
	var warnings []error
	if err != nil {
		warning := fmt.Errorf("PWRT: %v", err.Error())
		warnings = append(warnings, warning)
	}
	return PWRReport{
		PWRT:     temp,
		PWRP32V:  pwr.PWRP32V.voltage(),
		PWRP32C:  pwr.PWRP32C.current(),
		PWRP16V:  pwr.PWRP16V.voltage(),
		PWRP16C:  pwr.PWRP16C.current(),
		PWRM16V:  pwr.PWRM16V.voltage(),
		PWRM16C:  pwr.PWRM16C.current(),
		PWRP3V3:  pwr.PWRP3V3.voltage(),
		PWRP3C3:  pwr.PWRP3C3.current(),
		Warnings: warnings,
	}
}

// CSVSpecifications returns the specs used in creating the struct
func (pwr *PWR) CSVSpecifications() []string {
	return []string{"AEZ", Specification}
}

// CSVHeaders returns the field names
func (pwr *PWR) CSVHeaders() []string {
	return csvHeader(pwr.Report())
}

// CSVRow returns the field values
func (pwr *PWR) CSVRow() []string {
	val := reflect.Indirect(reflect.ValueOf(pwr.Report()))
	values := make([]string, val.NumField())
	t := val.Type()
	for i := range values {
		valueField := val.Field(i)
		if t.Field(i).Name == "Warnings" {
			if valueField.Len() == 0 {
				values[i] = ""
			} else {
				var errs = make([]string, valueField.Len())
				for j, l := 0, valueField.Len(); j < l; j++ {
					errs[j] = fmt.Sprintf("%v", valueField.Index(j).Elem())
				}
				values[i] = strings.Join(errs, "|")
			}

		} else {
			values[i] = fmt.Sprintf("%v", valueField.Float())
		}
	}
	return values
}

// SetParquet setsthe parquet representation of the PWR
func (pwr *PWR) SetParquet(row *parquetrow.ParquetRow) {
	report := pwr.Report()
	var warnings []string
	if report.Warnings != nil {
		for err := range report.Warnings {
			warnings = append(warnings, report.Warnings[err].Error())
		}
	}
	row.PWRT = report.PWRT
	row.PWRP32V = report.PWRP32V
	row.PWRP32C = report.PWRP32C
	row.PWRP16V = report.PWRP16V
	row.PWRP16C = report.PWRP16C
	row.PWRM16V = report.PWRM16V
	row.PWRM16C = report.PWRM16C
	row.PWRP3V3 = report.PWRP3V3
	row.PWRP3C3 = report.PWRP3C3
	row.Warnings = warnings
}
