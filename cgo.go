package cgobyexample

/*
#include <stdlib.h>
#include <limits.h>
#include "cgo.hpp"

#cgo LDFLAGS: -lstdc++
#cgo CXXFLAGS: -std=c++1z -stdlib=libc++
*/
import "C"
import (
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/gomega"
)

// ConvertPrimitive demonstrates how to convert between Go and C primitives.
func ConvertPrimitive() {
	goIn := C.INT_MAX
	cIn := C.int(goIn)
	cOut := C.integer(cIn)
	goOut := int(cOut)

	spew.Dump(goIn, cIn, cOut, goOut)
	Expect(goOut).To(Equal(goIn))
	Expect(cOut).To(Equal(cIn))
	Expect(goOut).ToNot(Equal(cOut))
}

// ConvertArray demonstrates how to convert a Go slice into a C array and pass it into a C function for processing.
// It also demonstrates how to convert a C array back into a Go slice to read its values in Go.
func ConvertArray() {
	goInSlice := []int{1, 2, 3, 4, 5}

	// Allocate C array.
	cArray := C.calloc(C.size_t(len(goInSlice)), 4) // Set calloc size to 4 because C int is 4 bytes
	// Don't forget to free C allocations at the end to avoid memory leaks.
	defer C.free(cArray)
	cArrayUnsafePtr := unsafe.Pointer(cArray)

	// Synchronously iterate over Go slice and C array, copying values from Go slice to C array.
	for goSliceIndex, cArrayIndex := 0, uintptr(cArrayUnsafePtr); goSliceIndex < len(goInSlice); goSliceIndex, cArrayIndex = goSliceIndex+1, cArrayIndex+4 { // Increment cArrayIndex by 4 because C int is 4 bytes
		// Obtain pointer to C array item.
		cArrayItemUnsafePtr := unsafe.Pointer(cArrayIndex)
		cArrayItemCIntPtr := (*C.int)(cArrayItemUnsafePtr)

		// Write value to C array.
		goInt := goInSlice[goSliceIndex]
		cInt := C.int(goInt)
		*cArrayItemCIntPtr = cInt
	}

	// Pass the C array to a C function for processing.
	C.integer_array(cArrayUnsafePtr, C.int(len(goInSlice)))

	// Allocate Go slice.
	goOutSlice := make([]int, len(goInSlice))

	// Synchronously iterate over Go slice and C array, copying values from C array to Go slice.
	for goSliceIndex, cArrayIndex := 0, uintptr(cArrayUnsafePtr); goSliceIndex < len(goOutSlice); goSliceIndex, cArrayIndex = goSliceIndex+1, cArrayIndex+4 {
		// Obtain pointer to C array item.
		cArrayItemUnsafePtr := unsafe.Pointer(cArrayIndex)
		cArrayItemCIntPtr := (*C.int)(cArrayItemUnsafePtr)

		// Read value from C array.
		cInt := *cArrayItemCIntPtr
		goInt := int(cInt)
		goOutSlice[goSliceIndex] = goInt
	}

	Expect(goInSlice).To(Equal([]int{1, 2, 3, 4, 5}))
	Expect(goOutSlice).To(Equal([]int{5, 4, 3, 2, 1}))
}

// PassStructByValue demonstrates pass by value between Go and C functions.
func PassStructByValue() {
	cInStruct := C.struct_Point{x: 1, y: 2}
	cOutStruct := C.pass_by_value(cInStruct)

	spew.Dump(&cInStruct, &cOutStruct)
	Expect(cOutStruct).To(And(
		Equal(cInStruct),
		BeEquivalentTo(cInStruct),
		BeIdenticalTo(cInStruct),
	))
	Expect(&cOutStruct).To(And(
		Equal(&cInStruct),
		BeEquivalentTo(&cInStruct),
		Not(BeIdenticalTo(&cInStruct)),
	))
}
