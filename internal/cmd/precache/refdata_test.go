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
		`纔 才 [cai2] /(before an expression of quantity) only/`,
		`纔 才 [cai2] /(variant of 才[cai2]) just now/(before an expression of quantity) only/`)
	assertVariant(
		`曏 向 [xiang4] /direction/orientation/to face/to turn toward/to/towards/shortly before/formerly/`,
		`曏 向 [xiang4] /variant of 向[xiang4]/direction/orientation/to face/to turn toward/to/towards/shortly before/formerly/`)
	assertVariant(
		`阯 址 [zhi3] /islet (variant of 沚[zhi3])/`,
		`阯 址 [zhi3] /foundation of a building (variant of 址[zhi3])/islet (variant of 沚[zhi3])/`)

	assertNoVariant(
		`一准 一准 [yi1 zhun3] /also written 一準|一准[yi1 zhun3]/`)
	assertNoVariant(
		`倆錢兒 俩钱儿 [lia3 qian2 r5] /erhua variant of 倆錢兒|俩钱儿[lia3 qian2 r5]/`)
	assertNoVariant(
		`怎麽 怎么 [zen3 me5] /variant of 怎麼|怎么[zen3 me5]/`)
	assertNoVariant(
		`籖 签 [qian1] /Japanese variant of 籤|签[qian1]/`)
	assertNoVariant(
		`閑 闲 [xian2] /(variant of 閒|闲[xian2]) idle/`)
}
