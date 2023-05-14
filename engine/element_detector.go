package engine

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

type Detector struct {
	net            gocv.Net
	outputNames    []string
	footage        *gocv.VideoCapture
	window         *gocv.Window
	classes        []string
	gpuEnabled     bool
	scoreThreshold float32
	nmsThreshold   float32
}

func NewAIDSEngine() *Detector {
	return &Detector{gpuEnabled: false, scoreThreshold: 0.45, nmsThreshold: 0.5}
}

func (d *Detector) Load() {

	var err error

	d.net = gocv.ReadNet("./model/yolov4.weights", "./model/yolov4.cfg")

	d.net.SetPreferableBackend(gocv.NetBackendDefault)
	d.net.SetPreferableTarget(gocv.NetTargetCPU)

	d.outputNames = getOutputsNames(&d.net)

	d.footage, err = gocv.VideoCaptureFile("https://rr2---sn-qxaelnez.googlevideo.com/videoplayback?expire=1684101159&ei=xwNhZIG2JsS-lu8P8vG8iAQ&ip=181.41.206.203&id=o-ALts9GDf2X6bszLvXeSMlwls7s7cMndgnFJ_4DfUUsam&itag=18&source=youtube&requiressl=yes&spc=qEK7B1CpqumqMMR8HA-r_QAk-DNTDN2RRWJEHrI4mg&vprv=1&svpuc=1&mime=video%2Fmp4&ns=uUCgUry2TFq0n7emdmx7iz4N&cnr=14&ratebypass=yes&dur=1801.496&lmt=1665157572495972&fexp=24007246,51000013&c=WEB&txp=4430434&n=jysEnwuFPmcGYA&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cspc%2Cvprv%2Csvpuc%2Cmime%2Cns%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AOq0QJ8wRgIhANzvWLbOz3qKhqG1u6NHOq2yN2Rv-usvIvZwJ4OPZ5wwAiEAwayFFH25nC2kaRsLJzlu0wMYNmU4_250DBDzIR5dDLg%3D&rm=sn-q4feek7z&req_id=4e0682f15f5a3ee&redirect_counter=2&cm2rm=sn-ugp2ax2a5t-qxas7l&cms_redirect=yes&cmsv=e&ipbypass=yes&mh=z3&mip=103.240.235.23&mm=29&mn=sn-qxaelnez&ms=rdu&mt=1684081754&mv=m&mvi=2&pl=24&lsparams=ipbypass,mh,mip,mm,mn,ms,mv,mvi,pl&lsig=AG3C_xAwRAIgVcdf_0RP72tNkb8XlR24jXjbf--buzbY1rKrXyeM9ZgCIAWiARcqYT1pHzdFtvNSqZGF_HY0vFle-fpEXcQtmLQr")

	if err != nil {
		log.Error(err)
	}

	d.window = gocv.NewWindow("Animal Intrusion Detection System")

	d.classes = readCOCO()

}

func (d *Detector) Process() {

	defer d.Close()

	mat := gocv.NewMat()

	for {
		isTrue := d.footage.Read(&mat)

		if mat.Empty() {
			continue
		}

		if isTrue {
			frame, _ := detect(&d.net, mat.Clone(), d.scoreThreshold,
				d.nmsThreshold, d.outputNames, d.classes)

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

func readCOCO() []string {
	var classes []string
	read, _ := os.Open("./model/coco.names")
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
