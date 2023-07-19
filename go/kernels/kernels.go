package kernels

type Kernel [3][3]uint32

var None = Kernel{
	{0, 0, 0},
	{0, 1, 0},
	{0, 0, 0},
}

var Vertical = Kernel{
	{0, 1, 0},
	{0, 2, 0},
	{0, 1, 0},
}

var Horizontal = Kernel{
	{0, 0, 0},
	{1, 2, 1},
	{0, 0, 0},
}

var Soft = Kernel{
	{1, 2, 1},
	{2, 12, 2},
	{1, 2, 1},
}
