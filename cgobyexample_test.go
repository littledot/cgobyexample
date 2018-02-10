package cgobyexample_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	. "github.com/littledot/cgo-by-example"
)

var _ = Describe("CgoByExample", func() {
	It("Should Pass", func() {
		ConvertPrimitive()
		ConvertArray()
		PassStructByValue()
	})
})
