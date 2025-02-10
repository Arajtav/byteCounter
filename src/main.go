package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"os/signal"
	"syscall"
)

var done chan int
var signals chan os.Signal

func byteCounter(betterNormalization bool, ignoreMostFrequent bool, inputFileName string, outputFileName string) {
	var ar [256][256]float64
	count := 0
	file, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Mode().IsRegular() {
		file.Close()
		data, err := os.ReadFile(inputFileName)
		if err != nil {
			panic(err)
		}
		for {
			if count+2 > len(data) {
				break
			}
			ar[data[count]][data[count+1]] += 1.0
			count += 2
		}
	} else {
		fmt.Println("irregular file")
		data := make([]byte, 4096)
		for {
			if len(signals) > 0 {
				fmt.Printf("Program interrupted, %d bytes read\n", count)
				break
			}
			n, err := file.Read(data)
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				println("EOF")
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
				mx = ar[i][j]
				mxv.X = uint8(i)
				mxv.Y = uint8(j)
			} else if ar[i][j] > mx2 {
				mx2 = ar[i][j]
			} else if ar[i][j] < mi {
				mi = ar[i][j]
			}
		}
	}

	fmt.Printf("max: %f\nmin: %f\nmxv: %02x %02x\n", mx, mi, mxv.X, mxv.Y)
	if ignoreMostFrequent {
		mx = mx2
	}
	mx /= 255
	if betterNormalization {
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

	if ignoreMostFrequent {
		img.Set(int(mxv.X), int(mxv.Y), color.RGBA{255, 0, 255, 255})
	}

	file, err = os.Create(outputFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create output file")
		panic(err)
	}
	png.Encode(file, img)
	file.Close()
	done <- count
}

func main() {
	flagBetterNormalization := flag.Bool("b", false, "If set to true, program will make contrast higher by subtracting minimal value from everything")
	flagIgnoreMostFrequent := flag.Bool("i", false, "Makes program ignore most occurring byte sequence when normalizing data. It will also replace pixel at that place to magenta")
	flag.Parse() // TODO: gnu style flags

	var args = flag.Args()
	if len(args) == 0 {
		args = append(args, os.Stdin.Name())
	}
	tmpd, err := os.MkdirTemp(".", "out_")
	if err != nil {
		panic(err)
	}
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	done = make(chan int, len(args))
	for i, fn := range args {
		go byteCounter(*flagBetterNormalization, *flagIgnoreMostFrequent, fn, fmt.Sprintf("%s/%d.png", tmpd, i))
	}
	for len(done) != len(args) {
	}
	total := 0
	for {
		if len(done) == 0 {
			break
		}
		total += <-done
	}
	fmt.Printf("total bytes read: %d, output files are in %s\n", total, tmpd)
}
