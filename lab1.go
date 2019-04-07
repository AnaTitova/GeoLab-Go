package main

import "fmt"
import "os"
import "io"
import geojson "github.com/paulmach/go.geojson"
import "github.com/fogleman/gg"

//====================================================================
type XYcoordinates struct {
	X float64
	Y float64
}

//====================================================================
var Coordinates [255]XYcoordinates
var NumofCoordinates = 0
var JSONfile = []byte(readFile("2.geojson"))

//====================================================================
func readFile(GeoJSONfile string) string {
	file, err := os.Open(GeoJSONfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	coordinatesData := make([]byte, 64)
	var data string

	for {
		fc, err := file.Read(coordinatesData)
		if err == io.EOF {
			break
		}
		data += string(coordinatesData[:fc])
	}
	return data
}

//====================================================================
func Decoder() *geojson.FeatureCollection {
	fc, _ := geojson.UnmarshalFeatureCollection(JSONfile)
	return fc
}

//====================================================================
func main() {
	Properties := Decoder()
	CoordinatesJSON := Properties.Features[0].Geometry.Polygon
	NumofCoordinates = len(CoordinatesJSON[0])
	for i := 0; i < NumofCoordinates; i++ {
		Coordinates[i].X = CoordinatesJSON[0][i][0]
		Coordinates[i].Y = CoordinatesJSON[0][i][1]
	}

	var coefW float64 = 0.003 * 1366
	var coefH float64 = 0.003 * 1024
	dc := gg.NewContext(1366, 1024)
	dc.SetHexColor("fff")
	dc.InvertY()
	dc.Clear()

	for i := 0; i < NumofCoordinates; i++ {
		x := (Coordinates[i].X + 70) * coefW
		y := (Coordinates[i].Y + 100) * coefH
		if i == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.ClosePath()
	dc.SetHexColor("f00")
	dc.Fill()

	//контур
	dc.SetHexColor("000")
	dc.SetLineWidth(2)
	for i := 0; i < NumofCoordinates; i++ {
		x := (Coordinates[i].X + 70) * coefW
		y := (Coordinates[i].Y + 100) * coefH
		if i == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.ClosePath()
	dc.Stroke()
	dc.SavePNG("Result.png")
}
