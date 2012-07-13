package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var pat = regexp.MustCompile(`^/img/grid-(\d+)x(\d+)-(\d+)-(\d+)x(\d+)-([a-z]+).png$`)

var COLORS = map[string]color.NRGBA{
	"red":    {0xef, 0x29, 0x29, 0},
	"orange": {0xf5, 0x79, 0x00, 0},
	"yellow": {0xed, 0xd4, 0x00, 0},
	"green":  {0x73, 0xd2, 0x16, 0},
	"blue":   {0x34, 0x65, 0xa4, 0},
	"purple": {0xad, 0x7f, 0xa8, 0},
	"brown":  {0xc1, 0x7d, 0x11, 0},
	"gray":   {0xa8, 0xaa, 0xa5, 0},
}

// each region has a different transparency
const (
	SuperGutter = 0x11
	Gutter      = 0x3c
	LastLine    = 0x44
	Bg          = 0x55
	Line        = 0x88
)

//# cw   - column width - the width of each column (in pixels)
//# ls   - line spacing - the height of each line (in pixels)
//# gw   - gutter width - the width of the gutter between columns (in pixels)
//# cg   - column group - the number of columns in each "super column"
//# lg   - line group - the number of lines in each group
//# wash - one of: red orange yellow green blue purple brown gray
func grid(cw, ls, gw, cg, lg int, wash string) (m *image.RGBA) {
	color := COLORS[wash]
	width := cw*cg + gw*cg
	height := lg * ls
	m = image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch {
			case ls > 1 && (y%ls) == (ls-1):
				color.A = Line
			case inGutter(x, cw*cg+gw*(cg-1), gw):
				color.A = SuperGutter
			case inGutter(x, cw, gw):
				color.A = Gutter
			case y%(ls*lg) >= (ls * (lg - 1)):
				color.A = LastLine
			default:
				color.A = Bg
			}
			m.Set(x, y, color)
		}
	}
	return
}

func inGutter(x, c, g int) bool {
	return x%(c+g) >= c
}

func atoi(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

// grid image png
// example: /img/grid-80x24-10-2x6-red.png
func img(w http.ResponseWriter, r *http.Request) {
	m := pat.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Expires", "Sat, 11-Dec-2110 01:01:01 GMT")
	g := grid(atoi(m[1]), atoi(m[2]), atoi(m[3]), atoi(m[4]), atoi(m[5]), m[6])
	err := png.Encode(w, g)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/img/", http.HandlerFunc(img))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
