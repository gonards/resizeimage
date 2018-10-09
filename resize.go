package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

var assocScalers = map[string]draw.Scaler{
	"NearestNeighbor": draw.NearestNeighbor,
	"ApproxBiLinear":  draw.ApproxBiLinear,
	"BiLinear":        draw.BiLinear,
	"CatmullRom":      draw.CatmullRom,
}

func singleImageResize(path, dstFolderName, scalerName string, width, height int) {
	initImgFormat()

	srcImg, format := openImage(path)

	if height == 0 && width == defaultWidth {
		height = getHeight(srcImg)
	}

	resizedImg := generateImg(srcImg, assocScalers[scalerName], width, height)
	dstName := extractFileName(path)
	saveImg(resizedImg, dstFolderName+"/"+dstName, format)
}

func multipleImageResize(path, dstFolderName, scalerName string, width, height int) {
	initImgFormat()

	folderList := strings.Split(path, "/")
	rootFolderName := folderList[len(folderList)-1]

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		fi, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		dstPath := strings.Replace(path, rootFolderName, dstFolderName, -1)

		switch mode := fi.Mode(); {
		case mode.IsDir():
			os.MkdirAll(dstPath, 0755)
		case mode.IsRegular():
			singleImageResize(path, filepath.Dir(dstPath), scalerName, width, height)
			fmt.Println("----------")
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
}

func getHeight(img image.Image) int {
	result := float64(defaultWidth) / float64(img.Bounds().Max.X) * float64(img.Bounds().Max.Y)
	return int(result)
}

func extractFileName(path string) string {
	filename := filepath.Base(path)
	result := strings.Split(filename, ".")
	return result[0]
}

func saveImg(img image.Image, dstName string, format string) {
	filename := dstName + "." + format
	fmt.Print("Saving Image " + filename + " ...")
	dstFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	switch format {
	case "png":
		err = png.Encode(dstFile, img)
	case "jpeg":
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(dstFile, img, &opt)
	}

	// close the file
	defer dstFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" DONE")
}

func generateImg(imgSrc image.Image, scaler draw.Scaler, width, height int) image.Image {
	fmt.Print("Resizing Image ...")
	rect := image.Rect(0, 0, width, height)
	imgDst := image.NewRGBA(rect)
	scaler.Scale(imgDst, rect, imgSrc, imgSrc.Bounds(), draw.Over, nil)
	fmt.Println(" DONE")
	return imgDst
}

func convertStringToInt(srt string) int {
	result, err := strconv.Atoi(srt)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func openImage(pathImg string) (image.Image, string) {
	fmt.Print("Opening Image " + pathImg + " ...")
	file, err := os.Open(pathImg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" DONE")
	return img, format
}

func initImgFormat() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
