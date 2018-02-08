package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(blurCmd)
}

var blurCmd = &cobra.Command{
	Use:   "blur",
	Short: "Blur video images",
	Long:  `Blur video images using a normalized box filter`,
	Run: func(cmd *cobra.Command, args []string) {
		handleBlurCmd()
	},
}

func handleBlurCmd() {
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(blurWindowTitle())
	defer window.Close()

	trackerX := window.CreateTrackbar("ksize X", 25)
	trackerX.SetMin(1)
	trackerX.SetPos(12)

	trackerY := window.CreateTrackbar("ksize Y", 25)
	trackerY.SetMin(1)
	trackerY.SetPos(12)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	fmt.Printf("Start reading camera device: %v\n", deviceID)
MainLoop:
	for {
		if ok := webcam.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// Blur image proccessing filter
		gocv.Blur(img, processed, image.Pt(trackerX.GetPos(), trackerY.GetPos()))

		// Display the processed image
		window.IMShow(processed)

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case 103:
			// 'g'
			blurGoCodeFragment(trackerX.GetPos(), trackerY.GetPos())
		case 112:
			// 'p'
			blurPythonCodeFragment(trackerX.GetPos(), trackerY.GetPos())
		case 27:
			// 'ESC'
			break MainLoop
		}
	}
}

func blurWindowTitle() string {
	return "Blur - CV Toolkit"
}

func blurGoCodeFragment(x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("gocv.Blur(src, dest, image.Pt(%d, %d))\n\n", x, y)
}

func blurPythonCodeFragment(x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}