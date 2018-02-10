extern "C" {
#include "cgo.hpp"
#include <stdio.h>
}

int integer(int val) { return val; }

void integer_array(void *array, int size) {
  int *intArray = (int *)array;
  int i, j;
  for (i = 0, j = size - 1; i < j; i++, j--) {
    int t = intArray[i];
    intArray[i] = intArray[j];
    intArray[j] = t;
  }
}

Point pass_by_value(Point val) {
  printf("c addr: %p\n", (void *)&val);
  return val;
}
