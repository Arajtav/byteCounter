package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "You need to specify input file");
        os.Exit(1);
    }

    if len(os.Args) > 2 {
        fmt.Fprintln(os.Stderr, "Too many arguments");
        os.Exit(1);
    }

    file, err := os.Open(os.Args[1]);
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open file");
        panic(err);
    }

    a, err := file.Stat();
    if err != nil {
        panic(err);
    }
    data := make([]byte, a.Size());
    file.Read(data);    // TODO, DON'T READ WHOLE FILE TO MEMORY, IT WON'T WORK ON LARGE FILES
    file.Close();

    ar := make([][]float32, 256);
    for i := 0; i < 256; i++ {
        ar[i]= make([]float32, 256);
    }

    // for every two bytes in file, use first on as X and second one as Y.
    t := 1;
    for {
        if t > int(a.Size())-1 { break; }
        x := data[t-1];
        y := data[t];
        ar[x][y] += 1.0;
        t += 2;
    }

    // find max value
    mx := float32(0.0);
    for i := 0; i < 256; i++ {
        for j := 0; j < 256; j++ {
            if ar[i][j] > mx { mx = ar[i][j]; }
        }
    }

    // normalize, make every value from 0 to 255
    for i := 0; i < 256; i++ {
        for j := 0; j < 256; j++ {
            ar[i][j] = (ar[i][j]/mx) * 255;
        }
    }

    img := image.NewRGBA(image.Rect(0, 0, 256, 256));
    for i := 0; i < 256; i++ {
        for j := 0; j < 256; j++ {
            v := uint8(ar[i][j]);
            img.Set(i, j, color.RGBA{v, v, v, 255});
        }
    }

    file, err = os.Create("out.png");
    if err != nil {
        panic(err);
    }
    png.Encode(file, img);
    file.Close();
}
