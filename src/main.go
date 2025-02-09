package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flagBn := flag.Bool("bn", false, "If set to true, program will make contrast higher by subtracting minimal value from everything")
	flagDebug := flag.Bool("debug", false, "Prints debug information")
	flagIgf := flag.Bool("igf", false, "Makes program ignore most occurring byte sequence when normalizing data. It will also replace pixel at that place to magenta")
	flag.Parse() // TODO: currently it does not work when flags are after the input file

	var args = flag.Args()
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Multiple input files are not suppoted yet.")
		os.Exit(1)
	}
	fn := os.Stdin.Name()
	if len(args) > 0 {
		fn = args[0]
	}
	var ar [256][256]float64

	count := 0
	file, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode().IsRegular() {
		file.Close()
		data, err := os.ReadFile(fn)
		if err != nil {
			panic(err)
		}
		if *flagDebug {
			fmt.Printf("input file size: %d\n", len(data))
		}
		for {
			if count+2 > len(data) {
				break
			}
			ar[data[count]][data[count+1]] += 1.0
			count += 2
		}
	} else {
		if *flagDebug {
			fmt.Println("irregular file")
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		data := make([]byte, 4096)
		for {
			if len(signals) > 0 {
				fmt.Printf("Program interrupted, %d bytes read\n", count)
				break
			}
			n, err := file.Read(data)
			if err != nil {
				panic(err)
			}
			if n == 0 {
				break
			}
			count2 := 0
			for {
				if count2+2 > n {
					break
				}
				ar[data[count2]][data[count2+1]] += 1.0
				count2 += 2
			}
			count += count2
			if n != 4096 {
				break
			}
		}
		file.Close()
	}

	// find max value and min value
	mx := float64(0.0)
	mx2 := float64(0.0)
	var mxv struct{ X, Y uint8 }
	mi := math.MaxFloat64
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			if ar[i][j] > mx {
				mx2 = mx
				mx = ar[i][j]
				mxv.X = uint8(i)
				mxv.Y = uint8(j)
			}
			if ar[i][j] < mi {
				mi = ar[i][j]
			}
		}
	}

	if *flagIgf {
		ar[mxv.X][mxv.Y] = 0.0
		mx = mx2
	}

	if *flagDebug {
		fmt.Printf("max: %f\nmin: %f\nmxv: %02x %02x\n", mx, mi, mxv.X, mxv.Y)
	}

	// normalization so everything will be from 0 to 255
	mx /= 255 // one operation instead multiplying by 255 in loop
	if *flagBn {
		mx -= mi / 255
		if mx == 0.0 {
			mx = 127
		}

		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				ar[i][j] = ((ar[i][j] - mi) / mx)
			}
		}
	} else {
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				ar[i][j] = (ar[i][j] / mx)
			}
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			v := uint8(ar[i][j])
			img.Set(i, j, color.RGBA{v, v, v, 255})
		}
	}

	if *flagIgf {
		img.Set(int(mxv.X), int(mxv.Y), color.RGBA{255, 0, 255, 255})
	}

	file, err = os.Create("out.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create output file")
		panic(err)
	}
	png.Encode(file, img)
	file.Close()
}
