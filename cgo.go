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
	"reflect"
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
	cArrayUnsafePtr := C.calloc(C.size_t(len(goInSlice)), C.sizeof_int) // Set calloc size to 4 because C int is 4 bytes
	// Don't forget to free C allocations at the end to avoid memory leaks.
	defer C.free(cArrayUnsafePtr)

	// Synchronously iterate over Go slice and C array, converting values from Go slice to C array.
	for goSliceIndex, cArrayIndex := 0, uintptr(cArrayUnsafePtr); goSliceIndex < len(goInSlice); goSliceIndex, cArrayIndex = goSliceIndex+1, cArrayIndex+C.sizeof_int { // Increment cArrayIndex by 4 because C int is 4 bytes
		// Obtain pointer to C array item.
		cArrayItemUnsafePtr := unsafe.Pointer(cArrayIndex)
		cArrayItemCIntPtr := (*C.int)(cArrayItemUnsafePtr)

		// Write value to C array.
		goInt := goInSlice[goSliceIndex]
		cInt := C.int(goInt)
		*cArrayItemCIntPtr = cInt
	}

	// Pass the C array to a C function for processing.
	C.integer_array((*C.int)(cArrayUnsafePtr), C.int(len(goInSlice)))

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

// ConvertArray2 demonstrates how to convert between a Go slice and C array using pointers to Go slices' internal array and reflect.SliceHeaders.
func ConvertArray2() {
	goInSlice := []int{1, 2, 3, 4, 5}

	// Convert Go items to C.
	goCIntSlice := make([]C.int, len(goInSlice))
	for i, goInt := range goInSlice {
		cInt := C.int(goInt)
		goCIntSlice[i] = cInt
	}

	// Obtain pointer to Go slice's internal array by referencing the 1st item in the slice.
	goCIntArray := (*C.int)(&goCIntSlice[0])

	C.integer_array(goCIntArray, C.int(len(goInSlice)))

	goCIntArrayUnsafePtr := unsafe.Pointer(goCIntArray)
	goCIntArrayRawPtr := uintptr(goCIntArrayUnsafePtr)
	goCIntSliceHeader := reflect.SliceHeader{
		Data: goCIntArrayRawPtr,
		Len:  len(goInSlice),
		Cap:  len(goInSlice),
	}
	goCIntSliceUnsafePtr := unsafe.Pointer(&goCIntSliceHeader)
	goCIntSlicePtr := (*[]C.int)(goCIntSliceUnsafePtr)
	goOutSlice := make([]int, goCIntSliceHeader.Len)
	for i, item := range *goCIntSlicePtr {
		goOutSlice[i] = int(item)
	}

	Expect(goInSlice).To(Equal([]int{1, 2, 3, 4, 5}))
	Expect(goOutSlice).To(Equal([]int{5, 4, 3, 2, 1}))
}

// BadConvertGrid demonstrates that the technique used to convert 1D Go slices to C arrays does not work for 2D Go slices and C arrays.
func BadConvertGrid() {
	goInGrid := [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}

	// Create a slice of pointers to model C's pointer to pointer (int**).
	goCIntPtrSlice := make([]*C.int, len(goInGrid))
	for i, goIntSlice := range goInGrid {
		// Convert Go items to C.
		goCIntSlice := make([]C.int, len(goIntSlice))
		for j, goInt := range goIntSlice {
			cInt := C.int(goInt)
			goCIntSlice[j] = cInt
		}
		// Obtain pointer to internal array and store it in the master slice.
		goCIntPtrSlice[i] = &goCIntSlice[0]
	}

	// Obtain pointer to internal array of the master slice.
	goCIntPtrArray := (**C.int)(&goCIntPtrSlice[0])

	// panic: cgo argument has Go pointer to Go pointer
	C.integer_grid(goCIntPtrArray, C.int(len(goInGrid)))
}

// GoodConvertGrid demonstrates how to properly convert between 2D Go slices and C arrays by wrapping Go slices around C arrays directly.
// "Turning C arrays into Go slices" https://github.com/golang/go/wiki/cgo
func GoodConvertGrid() {
	goInGrid := [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}

	// Create a slice of pointers to model C's pointer to pointer (int**).
	goCIntPtrSlice := make([]*C.int, len(goInGrid))
	for i, goIntSlice := range goInGrid {
		// Allocate C array.
		cArrayUnsafePtr := C.calloc(C.size_t(len(goIntSlice)), C.sizeof_int)
		defer C.free(cArrayUnsafePtr)

		// Store pointer to C array in Go slice.
		cArrayPtr := (*C.int)(cArrayUnsafePtr)
		goCIntPtrSlice[i] = cArrayPtr

		// "Turning C arrays into Go slices" https://github.com/golang/go/wiki/cgo
		// Create a Go slice wrapped around a C array.
		goCIntSlice := (*[1 << 30]C.int)(cArrayUnsafePtr)[:len(goIntSlice):len(goIntSlice)]
		for j, goInt := range goIntSlice {
			// Write values to the Go slice wrapper, which writes to the internal C array.
			cInt := C.int(goInt)
			goCIntSlice[j] = cInt
		}
	}

	// Obtain pointer to internal array of the master slice.
	goCIntPtrArray := (**C.int)(&goCIntPtrSlice[0])

	C.integer_grid(goCIntPtrArray, C.int(len(goInGrid)))

	// Create a Go slice wrapped around goCIntPtrArray
	goCIntPtrArrayUnsafePtr := unsafe.Pointer(goCIntPtrArray)
	goCIntPtrArraySlice := (*[1 << 30]*C.int)(goCIntPtrArrayUnsafePtr)[:len(goInGrid):len(goInGrid)]

	goOutGrid := make([][]int, len(goInGrid))
	// Read values from the Go slice wrapper, which reads from the internal C array.
	for i, goCIntArray := range goCIntPtrArraySlice {
		goIntSlice := make([]int, len(goInGrid))

		// Create a Go slice wrapped around goCIntArray
		goCIntArrayUnsafePtr := unsafe.Pointer(goCIntArray)
		goCIntArraySlice := (*[1 << 30]C.int)(goCIntArrayUnsafePtr)[:len(goInGrid):len(goInGrid)]

		// Read values from the Go slice wrapper, which reads from the internal C array.
		for j, item := range goCIntArraySlice {
			goIntSlice[j] = int(item)
		}
		goOutGrid[i] = goIntSlice
	}

	Expect(goInGrid).To(Equal([][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}))
	Expect(goOutGrid).To(Equal([][]int{{3, 2, 1}, {3, 2, 1}, {3, 2, 1}}))
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
