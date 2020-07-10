package vcode

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"time"
)

const (
	stdWidth  = 200
	stdHeight = 80
	// Maximum absolute skew factor of a single digit.
	maxSkew = .7
	// Number of background circles.
	circleCount = 20
)

type Image struct {
	*image.Paletted
	numWidth  int
	numHeight int
	dotSize   int
}

func NewImage(digits string) *Image {
	m := new(Image)
	m.Paletted = image.NewPaletted(image.Rect(0, 0, stdWidth, stdHeight), m.getRandomPalette())
	m.calculateSizes(stdWidth, stdHeight, len(digits))
	// Randomly position captcha inside the image.
	maxx := stdWidth - (m.numWidth+m.dotSize)*len(digits) - m.dotSize
	maxy := stdHeight - m.numHeight - m.dotSize*2
	var border int
	if stdWidth > stdHeight {
		border = stdHeight / 5
	} else {
		border = stdWidth / 5
	}
	x := rngint(border, maxx-border)
	y := rngint(border, maxy-border)
	// 绘制数字
	for _, n := range digits {
		m.drawDigit(font[n-48], x, y)
		x += m.numWidth + m.dotSize
	}
	// Draw strike-through line.
	m.strikeThrough()
	// // Apply wave distortion.
	m.distort(rngfloat(5, 10), rngfloat(100, 200))
	// // Fill image with random circles.
	m.fillWithCircles(circleCount, m.dotSize+1) //m.dotSize
	return m
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
func (m *Image) getRandomPalette() color.Palette {
	p := make([]color.Color, circleCount+1)
	p[0] = color.RGBA{0xFF, 0xFF, 0xFF, 0x00}
	// Circle colors.
	for i := 1; i <= circleCount; i++ {
		p[i] = color.RGBA{
			// 限制颜色
			uint8(rand.Intn(128) + 64),
			uint8(rand.Intn(128) + 64),
			uint8(rand.Intn(128) + 64),
			0xFF,
		}
	}
	return p
}

// EncodedPNG encodes an image to PNG and returns
// the result as a byte slice.
func (m *Image) EncodedPNG() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m.Paletted); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

func (m *Image) Base64() string {
	return base64.RawStdEncoding.EncodeToString(m.EncodedPNG())
}

func (m *Image) calculateSizes(stdWidth, stdHeight, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if stdWidth > stdHeight {
		border = stdHeight / 4
	} else {
		border = stdWidth / 4
	}
	// 将所有内容转换为浮点数进行计算。
	w := float64(stdWidth - border*2)
	h := float64(stdHeight - border*2)
	// fw考虑了数字之间的1点间距
	fw := float64(fontWidth + 1)
	fh := float64(fontHeight)
	nc := float64(ncount)
	// 计算单个数字的stdWidth，只考虑图像的stdWidth。
	nw := w / nc
	// 从这个stdWidth计算一个数字的stdHeight。
	nh := nw * fh / fw
	// Digit too high?
	if nh > h {
		// Fit digits based on stdHeight.
		nh = h
		nw = fw / fh * nh
	}
	// Calculate dot size.
	m.dotSize = int(nh / fh)
	if m.dotSize < 1 {
		m.dotSize = 1
	}
	// Save everything, making the actual stdWidth smaller by 1 dot to account
	// for spacing between digits.
	m.numWidth = int(nw) - m.dotSize
	m.numHeight = int(nh)
}

func (m *Image) drawHorizLine(fromX, toX, y int, colorIdx uint8) {
	for x := fromX; x <= toX; x++ {
		m.SetColorIndex(x, y, colorIdx)
	}
}

func (m *Image) drawCircle(x, y, radius int, colorIdx uint8) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xo := 0
	yo := radius

	m.SetColorIndex(x, y+radius, colorIdx)
	m.SetColorIndex(x, y-radius, colorIdx)
	m.drawHorizLine(x-radius, x+radius, y, colorIdx)

	for xo < yo {
		if f >= 0 {
			yo--
			dfy += 2
			f += dfy
		}
		xo++
		dfx += 2
		f += dfx
		m.drawHorizLine(x-xo, x+xo, y+yo, colorIdx)
		m.drawHorizLine(x-xo, x+xo, y-yo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y+xo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y-xo, colorIdx)
	}
}

func (m *Image) fillWithCircles(n, maxradius int) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		colorIdx := uint8(rngint(1, circleCount-1))
		r := rngint(1, maxradius)
		m.drawCircle(rngint(r, maxx-r), rngint(r, maxy-r), r, colorIdx)
	}
}

func (m *Image) strikeThrough() {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	y := rngint(maxy/3, maxy-maxy/3)
	amplitude := rngfloat(5, 20)
	period := rngfloat(80, 180)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxx; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < m.dotSize; yn++ {
			r := rngint(0, m.dotSize)
			m.drawCircle(x+int(xo), y+int(yo)+(yn*m.dotSize), r/2, 1)
		}
	}
}

func (m *Image) drawDigit(digit []byte, x, y int) {
	skf := rngfloat(-maxSkew, maxSkew)
	xs := float64(x)
	r := m.dotSize / 2
	ci := uint8(rand.Intn(circleCount) + 1)
	y += rngint(-r, r)
	for yo := 0; yo < fontHeight; yo++ {
		for xo := 0; xo < fontWidth; xo++ {
			if digit[yo*fontWidth+xo] != blackChar {
				continue
			}
			m.drawCircle(x+xo*m.dotSize, y+yo*m.dotSize, r, ci)
		}
		xs += skf
		x = int(xs)
	}
}

func (m *Image) distort(amplude float64, period float64) {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y

	oldm := m.Paletted
	newm := image.NewPaletted(image.Rect(0, 0, w, h), oldm.Palette)

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			newm.SetColorIndex(x, y, oldm.ColorIndexAt(x+int(xo), y+int(yo)))
		}
	}
	m.Paletted = newm
}

func rngint(from, to int) int {
	return rand.Intn(to+1-from) + from
}
func rngfloat(from, to float64) float64 {
	return rand.Float64()*(to-from) + from
}
