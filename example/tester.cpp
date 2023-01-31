#include <iostream>
#include <cstring>

#include "funcs.h"

int main(int argc, char const * const * argv) {
    if (argc < 2) {
        return EXIT_FAILURE;
    }

    if (strcmp(argv[1], "test") == 0) {
        test();
    } else {
        std::cerr << "Tundmatu funktsioon\n";
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}