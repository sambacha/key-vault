package encoder

// IEncoder has the responsibility to consistently encode/ decode objects for backwards compatibility
type IEncoder interface {
	// Encode takes an object and returns encoded bytes or error
	Encode(obj interface{}) ([]byte, error)
	// Decode takes an object and bytes to store decoded data in object, returns error if fails
	Decode(data []byte, v interface{}) error
}
