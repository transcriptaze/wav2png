package wav2png

var none = [][]uint32{
	{0, 0, 0},
	{0, 1, 0},
	{0, 0, 0},
}

var vertical = [][]uint32{
	{0, 1, 0},
	{0, 2, 0},
	{0, 1, 0},
}

var horizontal = [][]uint32{
	{0, 0, 0},
	{1, 2, 1},
	{0, 0, 0},
}

var soft = [][]uint32{
	{1, 2, 1},
	{2, 12, 2},
	{1, 2, 1},
}
