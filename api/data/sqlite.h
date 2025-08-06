#ifndef SQLITE_COLLATION_H
#define SQLITE_COLLATION_H

#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Function declarations
int setup_accent_insensitive_like(sqlite3 *db);
void normalize_string_wrapper(const char* input, int input_len, char* output, int* output_len);
int contains_pattern_wrapper(const char* text, int text_len, const char* pattern, int pattern_len);

#endif // SQLITE_COLLATION_H
