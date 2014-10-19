package main

const (
	b62_chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var b62_base = uint64(len(b62_chars))

func b62_Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	arr := []uint8{}

	for num > 0 {
		rem := num % b62_base
		num = num / b62_base
		arr = append(arr, b62_chars[rem])
	}

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return string(arr)
}
