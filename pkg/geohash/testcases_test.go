package geohash

var vpTestcases = []ValidPointTestcase{
	{Point{40.23, 23.41}, nil},
	{Point{140.2, 14.51}, ErrInvalidPoint},
	{Point{-140.2, 14.51}, ErrInvalidPoint},
	{Point{20.2, 245.51}, ErrInvalidPoint},
	{Point{20.2, -245.51}, ErrInvalidPoint},
}

var encTestcases = []EncodingTestcase{
	{Point{52.205, 0.119}, 7, "u120fxw", nil},
}

var decTestcases = []DecodingTestcase{
	{"9745rntct4xh", Box{17.482266202569008, 17.48226637020707, -120.62176242470741, -120.62176208943129}, nil},
}
