package cgobyexample_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCgoByExample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CgoByExample Suite")
}
