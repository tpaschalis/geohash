package geohash

import "errors"

var InvalidPointError = errors.New("Invalid Point; Latitude and Longitude should be between [-180, 180] and [-90, 90]")

var InvalidPrecisionError = errors.New("Invalid Precision as input; Precision should be a between [1, 12]")

var ErrorHashExceedsMaxPrecision = errors.New("Invalid Hash; Length exceeds max precision of 12 chars")
