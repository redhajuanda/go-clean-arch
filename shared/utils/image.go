package utils

import (
	"image"
	"io/ioutil"
	"net/http"

	"github.com/nfnt/resize"
)

// ResizeImage resizes image based on given size
func ResizeImage(img image.Image, x, y uint) image.Image {

	// var images []image.Image
	// images = append(images, img) // append original image
	newImage := resize.Resize(x, y, img, resize.Lanczos2)

	return newImage
}

// GetByteImage gets image from url and return []byte of image
func GetByteImage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return b
	// fmt.Println(url)
	// fmt.Println(path.Base(url))
	// buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	// _, err = resp.Body.Read(buff)
	// if err != nil {
	// 	panic(err)
	// }
	// return buff
}
