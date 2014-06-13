package transformer

import (
	"github.com/gographics/imagick/imagick"
	_ "os"
	"fmt"
)

func Resize(data []byte, width uint, height uint) (image []byte, err error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err = mw.ReadImageBlob(data); err != nil { return }

	w := mw.GetImageWidth()
	h := mw.GetImageHeight()

	fmt.Printf("Original file size %dx%d\n", w, h)

	if err = mw.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1); err != nil { return }

	if err = mw.SetImageCompressionQuality(95); err != nil { return }

	image = mw.GetImageBlob()

	return
}
