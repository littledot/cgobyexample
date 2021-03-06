#ifndef __CGO_HPP__
#define __CGO_HPP__

typedef struct Point {
  int x;
  int y;
} Point;

int integer(int val);
void integer_array(int* array, int size);
void integer_grid(int** grid, int size);
unsigned long unsigned_long(unsigned long a);

Point pass_by_value(Point val);

#endif
