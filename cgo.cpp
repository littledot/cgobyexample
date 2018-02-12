extern "C" {
#include "cgo.hpp"
#include <stdio.h>
}

int integer(int val) { return val; }

void integer_array(int *array, int size) {
  for (int i = 0, j = size - 1; i < j; i++, j--) {
    int t = array[i];
    array[i] = array[j];
    array[j] = t;
  }
}

void integer_grid(int **grid, int size) {
  for (int i = 0; i < size; i++) {
    integer_array(grid[i], size);
  }
}

Point pass_by_value(Point val) {
  printf("c addr: %p\n", (void *)&val);
  return val;
}
