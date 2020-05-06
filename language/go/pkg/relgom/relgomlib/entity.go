package relgomlib

import (
	"strings"

	"github.com/pkg/errors"
)

type EntityTypeStaticMetadata struct {
	PKMask       []uint64
	RequiredMask []uint64
}

// UpdateMaskForFieldButPanicIfAlreadySet checks the mask and panics if the field was already set.
func UpdateMaskForFieldButPanicIfAlreadySet(entityMask *uint64, fieldMask uint64) {
	if *entityMask&fieldMask != 0 {
		panic(errors.New("field already set"))
	}
	*entityMask |= fieldMask
}

// PanicIfRequiredFieldsNotSet checks the mask and panics if any required fields were not set.
// It takes a fieldsCommaList instead of a slice to ensure efficiency of successful scenarios.
func PanicIfRequiredFieldsNotSet(entityMasks []uint64, requiredFieldsMasks []uint64, fieldsCommalist string) {
	var fields []string
	var missingFields []string
	for i, entityMask := range entityMasks {
		if entityMask&requiredFieldsMasks[i] != requiredFieldsMasks[i] {
			gap := ^entityMask & requiredFieldsMasks[i]
			if fields == nil {
				fields = strings.Split(fieldsCommalist, ",")
			}
			for i, field := range fields {
				if gap&(uint64(1)<<uint(i)) != 0 {
					missingFields = append(missingFields, field)
				}
			}
		}
	}
	if len(missingFields) != 0 {
		quantifier := "field"
		if len(missingFields) != 1 {
			quantifier += "s"
		}
		panic(errors.Errorf("required field%s %s not set", quantifier, strings.Join(missingFields, ", ")))
	}
}
