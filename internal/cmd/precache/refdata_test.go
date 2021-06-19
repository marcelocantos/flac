package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElideCLTradHanziRE(t *testing.T) {
	assertApply := func(expected, original string) {
		t.Helper()

		actual, ok := applyVariantRE(elideCLTradHanziRE, original)
		if assert.True(t, ok) {
			assert.Equal(t, expected, actual)
		}
	}

	assertApply("/CL:个[ge4]/", "/CL:個|个[ge4]/")
}
