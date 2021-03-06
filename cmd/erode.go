package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(erodeCmd)
}

var currentErodeShape int

var erodeCmd = &cobra.Command{
	Use:   "erode",
	Short: "Erode video images",
	Long: `Erode video images.
	
Key commands:
  Use 'z' and 'x' keys to page through structuring element shapes.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleErodeCmd()
	},
}

func handleErodeCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(erodeWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	trackerX := window.CreateTrackbar("ksize X", 25)
	trackerX.SetMin(1)
	trackerX.SetPos(12)

	trackerY := window.CreateTrackbar("ksize Y", 25)
	trackerY.SetMin(1)
	trackerY.SetPos(12)

	fmt.Printf("Start reading video: %v\n", videoSource)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// Erode image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentErodeShape), image.Pt(trackerX.GetPos(), trackerY.GetPos()))
		gocv.Erode(img, &processed, kernel)
		kernel.Close()

		// Display the processed image?
		if pause {
			window.IMShow(img)
		} else {
			window.IMShow(processed)
		}

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case zKey:
			currentErodeShape = prevShape(currentErodeShape)
			window.SetWindowTitle(erodeWindowTitle())
		case xKey:
			currentErodeShape = nextShape(currentErodeShape)
			window.SetWindowTitle(erodeWindowTitle())
		case gKey:
			erodeGoCodeFragment(getCurrentMorphShapeDescription(currentErodeShape), trackerX.GetPos(), trackerY.GetPos())
		case pKey:
			erodePythonCodeFragment(currentErodeShape, trackerX.GetPos(), trackerY.GetPos())
		case space:
			handlePause(erodeWindowTitle())
		case esc:
			return
		}
	}
}

func erodeWindowTitle() string {
	return "Erode - " + getCurrentMorphShapeDescription(currentErodeShape) + "- CVscope"
}

func erodeGoCodeFragment(morphType string, x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("\nkernel := gocv.GetStructuringElement(gocv.%s, image.Pt(%d, %d))\n", morphType, x, y)
	fmt.Printf("gocv.Erode(src, &dest, kernel)\n\n")
}

func erodePythonCodeFragment(morphType, x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
