package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyVariantREs(t *testing.T) {
	assertVariant := func(expected, original string) {
		t.Helper()
		actual, ok := applyVariantREs(original)
		if assert.True(t, ok) {
			assert.Equal(t, expected, actual)
		}
	}
	assertNoVariant := func(original string) {
		t.Helper()
		actual, ok := applyVariantREs(original)
		assert.False(t, ok, actual)
	}
	assertVariant(
		"纔 才 [cai2] /(before an expression of quantity) only/",
		"纔 才 [cai2] /(variant of 才[cai2]) just now/(before an expression of quantity) only/")
	assertNoVariant(
		"怎麽 怎么 [zen3 me5] /variant of 怎麼|怎么[zen3 me5]/")
}
