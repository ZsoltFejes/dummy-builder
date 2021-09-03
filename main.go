package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func getSize(fileSize string) int {
	var units = [8]string{"x", "KB", "MB", "GB", "TB", "PB", "EB", "ZB"} // x is there as a place holder, it helps with math later in the function
	var unitsShort = [8]string{"x", "K", "M", "G", "T", "P", "E", "Z"}   // x is there as a place holder, it helps with math later in the function
	var chosenUnitIndex int
	var chosenUnit string
	for i, unit := range units {
		if strings.Contains(fileSize, units[i]) {
			chosenUnitIndex = i
			chosenUnit = unit
			break
		}
		if strings.Contains(fileSize, unitsShort[i]) {
			chosenUnitIndex = i
			chosenUnit = unitsShort[i]
			break
		}
	}
	fileSizeTrimmed := fileSize
	fileSizeTrimmed = strings.TrimSuffix(fileSizeTrimmed, chosenUnit)
	fileSizeInt, err := strconv.Atoi(fileSizeTrimmed)
	if err != nil {
		fmt.Println("Unable to find storage size unit")
		fmt.Println("Accepted storage size units:")
		fmt.Println("	KB|K - Kilobyte")
		fmt.Println("	MB|M - Megabyte")
		fmt.Println("	GB|G - Gigabyte")
		fmt.Println("	TB|T - Terabyte")
		fmt.Println("	PB|P - Petabyte")
		fmt.Println("	EB|E - Exabyte")
		fmt.Println("	ZB|Z - Yottabyte")
		os.Exit(2)
	}
	if chosenUnitIndex == 0 {
		return fileSizeInt
	} else {
		return int((math.Pow(1024, float64(chosenUnitIndex))) * float64(fileSizeInt))
	}
}

func main() {
	fileSize := flag.String("s", "1M", "Set the size of the dummy file you want to create")
	fileName := flag.String("o", "dummy", "Set where you want to output the dummy file")
	flag.Parse()

	bytesToGenerate := getSize(*fileSize)
	buffer := 1048576 // 1 MB
	f, err := os.OpenFile(*fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		print("Unable to create '" + *fileName + "', make sure path exist, check permissions and try again")
	}
	defer f.Close()
	bar := progressbar.NewOptions(bytesToGenerate,
		progressbar.OptionShowBytes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("Generating file '"+*fileName+"' ("+*fileSize+")"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	start := time.Now()
	for {
		if bytesToGenerate > buffer {
			token := make([]byte, buffer)
			rand.Read(token)
			f.Write(token)
			bar.Add(buffer)
			bytesToGenerate -= buffer
		} else {
			token := make([]byte, bytesToGenerate)
			rand.Read(token)
			f.Write(token)
			bar.Add(bytesToGenerate)
			bytesToGenerate = 0
		}
		if bytesToGenerate <= 0 {
			break
		}
	}
	fmt.Println("\n'" + *fileName + "' was generated in " + time.Since(start).String())
}
