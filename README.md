# What is this project about
This project is a simple program that will draw a number of occurrences of 2 bytes sequences in file, and save it to out.png.

## How does it work?
For each 2 bytes in file, program uses first one as X, and second as Y. Program increments value in 2d array at that place every time that sequence is scanned.
On the end everything is normalized to use values from 0 to 255.

## Usage
There are 3 flags you can use:
- `bn` - By default program uses ratio of value to the max value to determine brightness of the pixel. That means, if every sequence occurred 1000 times and one 1001 times, every pixel of output will look the same. With this flag minimal value (1000 in this example) will be subtracted from every value (so it will work like if every sequence occurred 0 times, and one 1 time).
- `debug` - Prints debug information.
- `igf` - Makes program ignore sequence that occurred most times (for example without this flag, in binary files that has a lot of null bytes, output will have only one bright pixel).
