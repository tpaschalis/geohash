package main

import "errors"

// https://github.com/golang/go/wiki/Errors

var ErrInvalidPoint = errors.New("Invalid Point; Latitude and Longitude should be between [-180, 180] and [-90, 90]")

var ErrInvalidPrecision = errors.New("Invalid Precision as input; Precision should be a between [1, 12]")

var ErrHashExceedsMaxPrecision = errors.New("Invalid Hash; Length exceeds max precision of 12 chars")

var ErrInvalidCharacterInHash = errors.New("Invalid Hash; Contains a character not in the base32 charset")

var ErrInvalidHash = errors.New("")
