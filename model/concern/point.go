package model

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"
)

type Point struct {
	x float64
	y float64
}

// CombineLatLng into float array
func (c *Point) FloatArray() []float64 {
	return []float64{c.x, c.y}
}

// CombineLatLng into string
func (c *Point) String() string {

	latStr := strconv.FormatFloat(c.x, 'f', -1, 64)
	lngStr := strconv.FormatFloat(c.y, 'f', -1, 64)

	return "(" + latStr + "," + lngStr + ")"
}

// Scan scan value into Point, implements sql.Scanner interface
func (p *Point) Scan(value interface{}) error {
	byteArrValue, ok := value.([]byte)
	if ok {
		value = string(byteArrValue)
	}

	latLng, ok := value.([]float64)
	if ok {
		p.x = latLng[0]
		p.y = latLng[1]

		return nil
	}
	latLngStr, ok := value.(string)
	if ok {
		trimPrefix := strings.TrimPrefix(latLngStr, "(")
		trimSuffixAfterPrefix := strings.TrimSuffix(trimPrefix, ")")
		dataCoordinate := strings.Split(trimSuffixAfterPrefix, ",")

		if len(dataCoordinate) >= 2 {
			p.x, _ = strconv.ParseFloat(dataCoordinate[0], 64)
			p.y, _ = strconv.ParseFloat(dataCoordinate[1], 64)

			return nil
		}
	}
	return errors.New("failed to parse point: not supported format (float array or string)")
}

func (p *Point) ScanFloatArray(value []float64) error {
	if len(value) >= 2 {
		p.x = value[0]
		p.y = value[1]

		return nil
	}
	return errors.New("failed to parse point: not enough coordinate")
}

func (p *Point) ScanString(value string) error {
	trimPrefix := strings.TrimPrefix(value, "(")
	trimSuffixAfterPrefix := strings.TrimSuffix(trimPrefix, ")")
	dataCoordinate := strings.Split(trimSuffixAfterPrefix, ",")

	if len(dataCoordinate) >= 2 {
		p.x, _ = strconv.ParseFloat(dataCoordinate[0], 64)
		p.y, _ = strconv.ParseFloat(dataCoordinate[1], 64)

		return nil
	}
	return errors.New("failed to parse point: not enough coordinate")
}

// Value return json value, implement driver.Valuer interface
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}
