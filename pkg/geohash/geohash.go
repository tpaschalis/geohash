package geohash

import (
	"strings"
)

type Point struct {
	lat float64
	lon float64
}

type Box struct {
	minLat float64
	maxLat float64
	minLon float64
	maxLon float64
}

type Direction int

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

// Encodes a Point, returning a geohash string
// Uses maximum 12-character precision
func Encode(p Point) (string, error) {
	s, e := EncodeUsingPrecision(p, 12)
	return s, e
}

var base32Charset = "0123456789bcdefghjkmnpqrstuvwxyz"

// Encodes a Point, returning a geohash string
// of length equal to `precision`
func EncodeUsingPrecision(p Point, precision int) (string, error) {
	res := ""
	idx := 0 // Index of base32 charset map
	bit := 0 // each char holds 5 bits
	evenBit := true
	latMin, latMax := -90., 90.
	lonMin, lonMax := -180., 180.

	if err := ValidPoint(p); err != nil {
		return res, err
	}

	if precision > 12 || precision <= 0 {
		return res, InvalidPrecisionError
	}

	for len(res) < precision {
		switch evenBit {
		// Even digits, work on East-West direction
		case true:
			lon := (lonMin + lonMax) / 2
			if p.lon >= lon {
				idx = idx*2 + 1
				lonMin = lon
			} else {
				idx = idx * 2
				lonMax = lon
			}
		// Odd digits, work on North-South direction
		case false:
			lat := (latMax + latMin) / 2
			if p.lat >= lat {
				idx = idx*2 + 1
				latMin = lat
			} else {
				idx = idx * 2
				latMax = lat
			}
		}
		evenBit = !evenBit

		// Completed a five-bit character; move on to the next one
		bit += 1
		if bit == 5 {
			res = res + base32Charset[idx:idx+1]
			bit = 0
			idx = 0
		}
	}

	return res, nil
}

// Decodes a geohash string, return a Box
func Decode(hash string) (Box, error) {
	hash = strings.ToLower(hash)
	//if err := ValidHash(hash); err != nil {
	//	return Box{}, InvalidHashError
	//}

	evenBit := true
	latMin, latMax := -90., 90.
	lonMin, lonMax := -180., 180.

	for i := 0; i < len(hash); i++ {
		char := hash[i : i+1]
		idx := strings.Index(base32Charset, char)

		for n := 4; n >= 0; n-- {
			bitN := idx >> uint(n) & 1
			switch evenBit {
			// Even digits, work on East-West direction
			case true:
				lonMid := (lonMin + lonMax) / 2.
				if bitN == 1 {
					lonMin = lonMid
				} else {
					lonMax = lonMid
				}
			// Odd digits, work on North-South direction
			case false:
				latMid := (latMin + latMax) / 2.
				if bitN == 1 {
					latMin = latMid
				} else {
					latMax = latMid
				}
				evenBit = !evenBit
			}
		}
	}

	return Box{latMin, latMax, lonMin, lonMax}, nil
}

func ValidPoint(p Point) error {
	latMin, latMax := -90., 90.
	lonMin, lonMax := -180., 180.
	if p.lat < latMin || p.lat > latMax || p.lon < lonMin || p.lon > lonMax {
		return InvalidPointError
	}

	return nil
}

func ValidHash(h string) error {
	if len(h) > 12 {
		return ErrorHashExceedsMaxPrecision
	}
	for v := range h {
		s := h[v : v+1]
		if strings.Index(base32Charset) == -1 {
			return
		}
	}
	return nil
}
