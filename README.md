# What is this project about
This project is a simple program that will draw a number of occurrences of a byte sequence, and save it to out.png.

## how does it work?
For each 2 bytes in file, use first as X, and second as Y. In 256x256 2d array increment value at that place every time that sequence is scaned.
On the end everything is normalized to use values from 0 to 255.
