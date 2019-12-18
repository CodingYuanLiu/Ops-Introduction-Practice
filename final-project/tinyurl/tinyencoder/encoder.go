package tinyencoder

func Decode(s string) (val int) {
	val = 0
	for _, c := range s {
		val = val << 6
		if c == '-' {
			val |= 0x3f
		} else if c == '_' {
			val |= 0x3e
		} else if c >= 'a' && c <= 'z' {
			val |= int(c) - 0x3d
		} else if c >= 'A' && c <= 'Z' {
			val |= int(c) - 0x37
		} else {
			val |= int(c) - 0x30
		}
	}
	return val
}

func Encode(val int) (s string) {
	var table = [64]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b',
		'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
		'x', 'y', 'z', '_', '-'}
	data := make([]byte, 6)
	for i := 0; i < 6; i++ {
		index := val & 0x3f
		val = val >> 6
		data[5 - i] = table[index]
	}
	return string(data)
}