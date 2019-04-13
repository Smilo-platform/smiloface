package server

import (
	"bytes"
	"encoding/base64"
	"github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
)

func init() {
}

func FacialRecognition(w http.ResponseWriter, r *http.Request) {

	//define endpoint permissions, block all except OPTIONS
	switch r.Method {
	case "OPTIONS":
		AllowOrigin(w, r)
		w.WriteHeader(http.StatusOK)
		return
	case "PUT":
	case "DELETE":
	case "HEAD":
	case "TRACE":
	case "CONNECT":
		w.WriteHeader(http.StatusForbidden)
		return
	}

	AllowOrigin(w, r)

	cascadeFile, err := ioutil.ReadFile("./facefinder")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	p := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Error("Could not process img")
		return
	}

	frameBuffer := new(bytes.Buffer)
	err = jpeg.Encode(frameBuffer, img, nil)
	if err != nil {
		log.Error("[ERROR] encoding frame buffer ", err)
		return
	}

	src := pigo.ImgToNRGBA(img)
	frame := pigo.RgbToGrayscale(src)

	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     *minSize,
		MaxSize:     *maxSize,
		ShiftFactor: *shiftFactor,
		ScaleFactor: *scaleFactor,
		ImageParams: pigo.ImageParams{
			Pixels: frame,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, *angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0)

	dc = gg.NewContext(cols, rows)
	dc.DrawImage(src, 0, 0)

	buff := new(bytes.Buffer)
	drawMarker(dets, buff, *circleMarker)

	// Encode as MJPEG
	w.Write(buff.Bytes())

}

func AllowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, API-KEY")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, HEAD, POST, PUT, OPTIONS")

	//TODO: add origin validation
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

}
