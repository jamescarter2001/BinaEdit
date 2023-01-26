#include "stringtable.h"

char* readString(char* stringTable, int offset) {
    char* str = malloc(50);
    strcpy(str, stringTable + offset);

    return str;
}