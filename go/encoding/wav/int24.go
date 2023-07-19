/**
 * Little endian 24 bit audio integer
 *
 */
package wav

const mask = int32(0x7fffff) ^ -1

type int24 [3]byte

func (i int24) ToInt32() int32 {
	var v int32 = 0

	v |= int32(i[2]) & 0x00007f
	v <<= 8
	v |= int32(i[1]) & 0x0000ff
	v <<= 8
	v |= int32(i[0]) & 0x0000ff

	if (i[2] & 0x80) == 0x80 {
		v |= mask
	}

	return v
}

/* Returns the value scaled to the interval [-1.0,+1.0]
 */
func (i int24) ToFloat() float32 {
	var v int32 = 0

	v |= int32(i[2]) & 0x00007f
	v <<= 8
	v |= int32(i[1]) & 0x0000ff
	v <<= 8
	v |= int32(i[0]) & 0x0000ff

	if (i[2] & 0x80) == 0x80 {
		v |= mask
	}

	return float32((2*v)+1) / 16777216.0
}
