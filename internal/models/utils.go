package models

func IfNotNil[T any](newVal *T, oldVal T) T {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func StringPtr(s string) *string {
	return &s
}

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
