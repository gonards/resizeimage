package main

import (
	"flag"
	"log"
)

var validSizes = map[string]bool{
	"manual": true,
	"auto":   true,
}

var validModes = map[string]bool{
	"image":  true,
	"folder": true,
}

var validScalers = map[string]bool{
	"NearestNeighbor": true,
	"ApproxBiLinear":  true,
	"BiLinear":        true,
	"CatmullRom":      true,
}

const defaultWidth = 400

func main() {
	// Get command line flags
	modeFlag := flag.String("m", "image", "Mode of resize. Allowed values are : image|folder")
	dstFolderFlag := flag.String("dst", "small", "Define the result folder name")
	widthFlag := flag.Int("w", defaultWidth, "Define the expected width")
	heightFlag := flag.Int("h", 0, "Define the expected height")
	pathFlag := flag.String("p", "./", "Define the path")
	scalerFlag := flag.String("s", "NearestNeighbor", "Define the scaler used to do the resize. Allowed values are : NearestNeighbor|ApproxBiLinear|BiLinear|CatmullRom")
	flag.Parse()

	validateMode(*modeFlag)
	validateScaler(*scalerFlag)

	switch *modeFlag {
	case "image":
		singleImageResize(*pathFlag, *dstFolderFlag, *scalerFlag, *widthFlag, *heightFlag)
	case "folder":
		multipleImageResize(*pathFlag, *dstFolderFlag, *scalerFlag, *widthFlag, *heightFlag)
	}
}

func validateMode(mode string) {
	if !validModes[mode] {
		log.Fatal("Invalid mode. Allowd modes are : image|folder")
	}
}

func validateScaler(scaler string) {
	if !validScalers[scaler] {
		log.Fatal("Invalid scaler. Allowd scalers are : NearestNeighbor|ApproxBiLinear|BiLinear|CatmullRom")
	}
}
