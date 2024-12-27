package services

import "github.com/akinolaemmanuel49/gocommerce/internal/models"

func ifNotEmpty(newVal, oldVal string) string {
	if newVal != "" {
		return newVal
	}
	return oldVal
}

func mergeAddress(newAddr, oldAddr models.Address) models.Address {
	return models.Address{
		Street:  ifNotEmpty(newAddr.Street, oldAddr.Street),
		City:    ifNotEmpty(newAddr.City, oldAddr.City),
		State:   ifNotEmpty(newAddr.State, oldAddr.State),
		Zip:     ifNotEmpty(newAddr.Zip, oldAddr.Zip),
		Country: ifNotEmpty(newAddr.Country, oldAddr.Country),
	}
}
