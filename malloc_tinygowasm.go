//go:build tinygo.wasm && !gc.custom && !custommalloc

package wasilibs

import (
	"unsafe"
)

/*
#include <stdint.h>
void* malloc(size_t size);
void free(void* ptr);
*/
import "C"

// TinyGo currently only includes a subset of malloc functions by default, so we
// reimplement the remaining here.

//export posix_memalign
func posix_memalign(memptr *uint32, align uint32, size uint32) int32 {
	if align < 4 {
		return 22 /* EINVAL */
	}

	// Ignore alignment and hope for best, TinyGo by default does not
	// provide a way to allocate aligned memory.
	mem := uint32(uintptr(C.malloc(C.size_t(size))))
	if mem == 0 {
		// TODO(anuraaga): Needs to read errno to be precise
		return 1
	}

	*memptr = mem

	return 0
}

//export calloc
func calloc(num uint32, size uint32) uint32

//export __libc_calloc
func __libc_calloc(num uint32, size uint32) uint32 {
	return calloc(num, size)
}

//export __libc_malloc
func __libc_malloc(size uint32) uint32 {
	return uint32(uintptr(C.malloc(C.size_t(size))))
}

//export __libc_free
func __libc_free(ptr uint32) {
	C.free(unsafe.Pointer(uintptr(ptr)))
}
