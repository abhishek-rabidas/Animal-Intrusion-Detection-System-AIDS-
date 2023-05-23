package engine

import (
	"AIDS/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

type Detector struct {
	net         gocv.Net
	outputNames []string
	footage     *gocv.VideoCapture
	window      *gocv.Window
	classes     []string
	config      *config.Config
}

func InitializeDetector(config *config.Config) *Detector {
	return &Detector{config: config}
}

func (d *Detector) Load() {

	var err error

	d.net = gocv.ReadNet(d.config.Model, d.config.Cfg)

	d.net.SetPreferableBackend(gocv.NetBackendType(gocv.NetBackendDefault))
	d.net.SetPreferableTarget(gocv.NetTargetType(gocv.NetTargetCPU))

	d.outputNames = getOutputsNames(&d.net)

	d.footage, err = gocv.VideoCaptureFile(d.config.Feed)

	if err != nil {
		log.Error(err)
	}

	d.window = gocv.NewWindow("Animal Intrusion Detection System")

	d.classes = readCOCO(d.config.Classnames)

}

func (d *Detector) Process() {

	mat := gocv.NewMat()

	for {
		isTrue := d.footage.Read(&mat)

		if mat.Empty() {
			continue
		}

		if isTrue {
			frame, _ := detect(&d.net, mat.Clone(), d.config.ScoreThreshold,
				d.config.NmsThreshold, d.outputNames, d.classes)

			d.window.IMShow(frame)
			key := d.window.WaitKey(1)
			if key == 113 {
				break
			}
		} else {
			return
		}

	}
}

func (d *Detector) Close() {
	d.net.Close()
	d.footage.Close()
	d.window.Close()
	log.Info("Process Completed")
}

func detect(net *gocv.Net, src gocv.Mat, scoreThreshold float32, nmsThreshold float32, OutputNames []string, classes []string) (gocv.Mat, []string) {
	img := src.Clone()
	img.ConvertTo(&img, gocv.MatTypeCV32F)
	blob := gocv.BlobFromImage(img, 1/255.0, image.Pt(416, 416), gocv.NewScalar(0, 0, 0, 0), true, false)
	net.SetInput(blob, "")
	probs := net.ForwardLayers(OutputNames)
	boxes, confidences, classIds := postProcess(img, &probs)
	indices := make([]int, 100)
	if len(boxes) == 0 { // No Classes
		return src, []string{}
	}
	gocv.NMSBoxes(boxes, confidences, scoreThreshold, nmsThreshold, indices)

	return drawRect(src, boxes, classes, classIds, indices)
}

func postProcess(frame gocv.Mat, outs *[]gocv.Mat) ([]image.Rectangle, []float32, []int) {
	var classIds []int
	var confidences []float32
	var boxes []image.Rectangle
	for _, out := range *outs {

		data, _ := out.DataPtrFloat32()
		for i := 0; i < out.Rows(); i, data = i+1, data[out.Cols():] {

			scoresCol := out.RowRange(i, i+1)

			scores := scoresCol.ColRange(5, out.Cols())
			_, confidence, _, classIDPoint := gocv.MinMaxLoc(scores)
			if confidence > 0.5 {

				centerX := int(data[0] * float32(frame.Cols()))
				centerY := int(data[1] * float32(frame.Rows()))
				width := int(data[2] * float32(frame.Cols()))
				height := int(data[3] * float32(frame.Rows()))

				left := centerX - width/2
				top := centerY - height/2
				classIds = append(classIds, classIDPoint.X)
				confidences = append(confidences, float32(confidence))
				boxes = append(boxes, image.Rect(left, top, width, height))
			}
		}
	}
	return boxes, confidences, classIds
}

func drawRect(img gocv.Mat, boxes []image.Rectangle, classes []string, classIds []int, indices []int) (gocv.Mat, []string) {
	var detectClass []string
	for _, idx := range indices {
		if idx == 0 {
			continue
		}
		gocv.Rectangle(&img, image.Rect(boxes[idx].Max.X, boxes[idx].Max.Y, boxes[idx].Max.X+boxes[idx].Min.X, boxes[idx].Max.Y+boxes[idx].Min.Y), color.RGBA{255, 0, 0, 0}, 2)
		gocv.PutText(&img, classes[classIds[idx]], image.Point{boxes[idx].Max.X, boxes[idx].Max.Y + 30}, gocv.FontHersheySimplex, 1, color.RGBA{0, 0, 255, 0}, 1)
		detectClass = append(detectClass, classes[classIds[idx]])
	}
	return img, detectClass
}

func getOutputsNames(net *gocv.Net) []string {
	var outputLayers []string
	for _, i := range net.GetUnconnectedOutLayers() {
		layer := net.GetLayer(i)
		layerName := layer.GetName()
		if layerName != "_input" {
			outputLayers = append(outputLayers, layerName)
		}
	}
	return outputLayers
}

func readCOCO(path string) []string {
	var classes []string
	read, _ := os.Open(path)
	defer read.Close()
	for {
		var t string
		_, err := fmt.Fscan(read, &t)
		if err != nil {
			break
		}
		classes = append(classes, t)
	}
	return classes
}
