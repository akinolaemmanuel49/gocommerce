package models

// IfNotNil returns a pointer to a new value of generic type
// if the new value is not nil else it returns the old value
func IfNotNil[T any](newVal *T, oldVal T) T {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// MergeAddress merges old and new addressess
func MergeAddress(newAddr *UpdateAddress, oldAddr Address) Address {
	if newAddr == nil {
		return oldAddr
	}

	return Address{
		Street:  IfNotNil(newAddr.Street, oldAddr.Street),
		City:    IfNotNil(newAddr.City, oldAddr.City),
		State:   IfNotNil(newAddr.State, oldAddr.State),
		Zip:     IfNotNil(newAddr.Zip, oldAddr.Zip),
		Country: IfNotNil(newAddr.Country, oldAddr.Country),
	}
}
