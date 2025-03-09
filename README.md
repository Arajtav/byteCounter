For each 2 bytes in file, program uses first one as X, and second as Y, in 256x256 table value at position is incremented.
After scanning the whole thing, everything is normalized and saved to a png.

## Usage
There are 2 flags you can use:
- `-b` - offset every value in the table by minimal value, before normalization.
- `-i` - ignore most frequent byte sequence for normalization (replaces it with magenta pixel).

## notes
Each time the program is run, it creates temporary `out_` directory, with `N.png` images for each processed input file.
Stdout has a map of which output files correspond to which input files.
All logs are written to stderr.
