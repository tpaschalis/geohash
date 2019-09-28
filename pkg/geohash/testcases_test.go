package geohash

var vpTestcases = []ValidPointTestcase{
	{Point{40.23, 23.41}, nil},
	{Point{140.2, 14.51}, InvalidPointError},
	{Point{-140.2, 14.51}, InvalidPointError},
	{Point{20.2, 245.51}, InvalidPointError},
	{Point{20.2, -245.51}, InvalidPointError},
}

var encTestcases = []EncodingTestcase{
	{Point{52.205, 0.119}, 7, "u120fxw", nil},
}

var decTestcases = []DecodingTestcase{
	{"zjtptsj2vn", Box{77.2779929638, 77.2779983282, 142.280373573, 142.280384302}, nil},
}
