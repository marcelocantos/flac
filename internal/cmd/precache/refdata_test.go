package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyVariantREs(t *testing.T) {
	assertVariantREs := func(expected, original string) {
		actual, ok := applyVariantREs(original)
		if assert.True(t, ok) {
			assert.Equal(t, expected, actual)
		}
	}
	assertVariantREs(
		"纔 才 [cai2] /(before an expression of quantity) only/",
		"纔 才 [cai2] /(variant of 才[cai2]) just now/(before an expression of quantity) only/")
}
