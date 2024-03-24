package main

import (
    "os"
    "bufio"

    "fmt"

    "flag"

    "image"
    "image/color"
    "image/png"
)

func main() {
    flagBn      := flag.Bool("bn",      false, "If set to true, program will make contrast higher by subtracting minimal value from everything");
    flagDebug   := flag.Bool("debug",   false, "Prints debug information");
    flagIgf     := flag.Bool("igf",     false, "Makes program ignore most occurring byte sequence when normalizing data. It will also replace pixel at that place to magenta");
    flag.Parse();

    if len(flag.Args()) < 1 {
        fmt.Fprintln(os.Stderr, "You need to specify input file");
        os.Exit(1);
    }

    if len(flag.Args()) > 1 {
        fmt.Fprintln(os.Stderr, "Too many arguments");
        os.Exit(1);
    }

    file, err := os.Open(flag.Args()[0]);
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open file");
        panic(err);
    }
    reader := bufio.NewReader(file);

    if *flagDebug { fmt.Printf("reader size: %d\n", reader.Size()); }

    var ar [256][256]float64;

    // for every two bytes in file, use first on as X and second one as Y.
    for {
        x, err := reader.ReadByte(); if err != nil { break; }
        y, err := reader.ReadByte(); if err != nil { break; }
        ar[x][y] += 1.0;
    }
    file.Close();

    // find max value and min value
    mx := float64(0.0);
    var mxv struct{X, Y uint8};
    mi := float64(1.79769313486231570814527423731704356798070e+308);    // float64 max
    for i := 0; i < 256; i++ {
        for j := 0; j < 256; j++ {
            if ar[i][j] > mx { mx = ar[i][j]; mxv.X = uint8(i); mxv.Y = uint8(j);}
            if ar[i][j] < mi { mi = ar[i][j]; }
        }
    }

    if *flagIgf {
        ar[mxv.X][mxv.Y] = 0.0;
        mx = float64(0.0);
        for i := 0; i < 256; i++ {
            for j := 0; j < 256; j++ {
                if ar[i][j] > mx { mx = ar[i][j]; }
            }
        }
    }

    if *flagDebug { fmt.Printf("max: %f\nmin: %f\nmxv: %02x %02x\n", mx, mi, mxv.X, mxv.Y); }

    // normalization so everything will be from 0 to 255
    mx /= 255;  // one operation instead multiplying by 255 in loop
    if *flagBn {
        mx -= mi/255;
        if mx == 0.0 { mx = 127; }

        for i := 0; i < 256; i++ {
            for j := 0; j < 256; j++ {
                ar[i][j] = ((ar[i][j]-mi)/mx);
            }
        }
    } else {
        for i := 0; i < 256; i++ {
            for j := 0; j < 256; j++ {
                ar[i][j] = (ar[i][j]/mx);
            }
        }
    }

    img := image.NewRGBA(image.Rect(0, 0, 256, 256));
    for i := 0; i < 256; i++ {
        for j := 0; j < 256; j++ {
            v := uint8(ar[i][j]);
            img.Set(i, j, color.RGBA{v, v, v, 255});
        }
    }

    if *flagIgf { img.Set(int(mxv.X), int(mxv.Y), color.RGBA{255, 0, 255, 255}); }

    file, err = os.Create("out.png");
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to create output file");
        panic(err);
    }
    png.Encode(file, img);
    file.Close();
}
