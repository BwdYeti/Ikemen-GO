package main

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// The global, platform-specific rendering backend
var gfx = &Renderer{}

// Blend constants
type BlendFunc int

const (
	BlendOne = BlendFunc(iota)
	BlendZero
	BlendSrcAlpha
	BlendOneMinusSrcAlpha
)

type BlendEquation int

const  (
	BlendAdd = BlendEquation(iota)
	BlendReverseSubtract
)

// Rotation holds rotation parameters
type Rotation struct {
	angle, xangle, yangle float32
}

func (r *Rotation) IsZero() bool {
	return r.angle == 0 && r.xangle == 0 && r.yangle == 0
}

// Tiling holds tiling parameters
type Tiling struct {
	x, y, sx, sy int32
}

var notiling = Tiling{}

// RenderParams holds the common data for all sprite rendering functions
type RenderParams struct {
	// Sprite texture and palette texture
	tex    *Texture
	paltex *Texture
	// Size, position, tiling, scaling and rotation
	size     [2]uint16
	x, y     float32
	tile     Tiling
	xts, xbs float32
	ys, vs   float32
	rxadd    float32
	rot      Rotation
	// Transparency, masking and palette effects
	tint  uint32 // Sprite tint for shadows
	trans int32  // Mugen transparency blending
	mask  int32  // Mask for transparency
	pfx   *PalFX
	// Clipping
	window *[4]int32
	// Rotation center
	rcx, rcy float32
	// Perspective projection
	projectionMode int32
	fLength        float32
	xOffset        float32
	yOffset        float32
}

func (rp *RenderParams) IsValid() bool {
	return rp.tex.IsValid() && IsFinite(rp.x+rp.y+rp.xts+rp.xbs+rp.ys+rp.vs+
		rp.rxadd+rp.rot.angle+rp.rcx+rp.rcy)
}

func drawQuads(modelview mgl.Mat4, x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	gfx.SetUniformMatrix("modelview", modelview[:])
	gfx.SetUniformF("x1x2x4x3", x1, x2, x4, x3) // this uniform is optional
	gfx.SetVertexData(
		x2, y2, 1, 1,
		x3, y3, 1, 0,
		x1, y1, 0, 1,
		x4, y4, 0, 0)

	gfx.RenderQuad()
}

// Render a quad with optional horizontal tiling
func rmTileHSub(modelview mgl.Mat4, x1, y1, x2, y2, x3, y3, x4, y4, width float32,
	tl Tiling, rcx float32) {
	//            p3
	//    p4 o-----o-----o- - -o
	//      /      |      \     ` .
	//     /       |       \       `.
	//    o--------o--------o- - - - o
	//   p1         p2
	topdist := (x3 - x4) * (1 + float32(tl.sx) / width)
	botdist := (x2 - x1) * (1 + float32(tl.sx) / width)
	if AbsF(topdist) >= 0.01 {
		db := (x4 - rcx) * (botdist - topdist) / AbsF(topdist)
		x1 += db
		x2 += db
	}

	// Compute left/right tiling bounds (or right/left when topdist < 0)
	xmax := float32(sys.scrrect[2])
	left, right := int32(0), int32(1)
	if topdist >= 0.01 {
		left = 1 - int32(math.Ceil(float64(MaxF(x3 / topdist, x2 / botdist))))
		right = int32(math.Ceil(float64(MaxF((xmax - x4) / topdist, (xmax - x1) / botdist))))
	} else if topdist <= -0.01 {
		left = 1 - int32(math.Ceil(float64(MaxF((xmax - x3) / -topdist, (xmax - x2) / -botdist))))
		right = int32(math.Ceil(float64(MaxF(x4 / -topdist, x1 / -botdist))))
	}

	if tl.x != 1 {
		left = 0
		right = Min(right, Max(tl.x, 1))
	}

	// Draw all quads in one loop
	for n := left; n < right; n++ {
		x1d, x2d := x1 + float32(n) * botdist, x2 + float32(n) * botdist
		x3d, x4d := x3 + float32(n) * topdist, x4 + float32(n) * topdist
		drawQuads(modelview, x1d, y1, x2d, y2, x3d, y3, x4d, y4)
	}
}

func rmTileSub(modelview mgl.Mat4, rp RenderParams) {
	x1, y1 := rp.x+rp.rxadd*rp.ys*float32(rp.size[1]), rp.rcy+((rp.y-rp.ys*float32(rp.size[1]))-rp.rcy)*rp.vs
	x2, y2 := x1+rp.xbs*float32(rp.size[0]), y1
	x3, y3 := rp.x+rp.xts*float32(rp.size[0]), rp.rcy+(rp.y-rp.rcy)*rp.vs
	x4, y4 := rp.x, y3
	//var pers float32
	//if AbsF(rp.xts) < AbsF(rp.xbs) {
	//	pers = AbsF(rp.xts) / AbsF(rp.xbs)
	//} else {
	//	pers = AbsF(rp.xbs) / AbsF(rp.xts)
	//}
	if !rp.rot.IsZero() {
		//	kaiten(&x1, &y1, float64(agl), rcx, rcy, vs)
		//	kaiten(&x2, &y2, float64(agl), rcx, rcy, vs)
		//	kaiten(&x3, &y3, float64(agl), rcx, rcy, vs)
		//	kaiten(&x4, &y4, float64(agl), rcx, rcy, vs)
		if rp.vs != 1 {
			y1 = rp.rcy + ((rp.y - rp.ys*float32(rp.size[1])) - rp.rcy)
			y2 = y1
			y3 = rp.rcy + (rp.y - rp.rcy)
			y4 = y3
		}
		if rp.projectionMode == 0 {
			modelview = modelview.Mul4(mgl.Translate3D(rp.rcx, rp.rcy, 0))
		} else if rp.projectionMode == 1 {
			//This is the inverse of the orthographic projection matrix
			matrix := mgl.Mat4{float32(sys.scrrect[2] / 2.0), 0, 0, 0, 0, float32(sys.scrrect[3] / 2), 0, 0, 0, 0, -65535, 0, float32(sys.scrrect[2] / 2), float32(sys.scrrect[3] / 2), 0, 1}
			modelview = modelview.Mul4(mgl.Translate3D(0, -float32(sys.scrrect[3]), rp.fLength))
			modelview = modelview.Mul4(matrix)
			modelview = modelview.Mul4(mgl.Frustum(-float32(sys.scrrect[2])/2/rp.fLength, float32(sys.scrrect[2])/2/rp.fLength, -float32(sys.scrrect[3])/2/rp.fLength, float32(sys.scrrect[3])/2/rp.fLength, 1.0, 65535))
			modelview = modelview.Mul4(mgl.Translate3D(-float32(sys.scrrect[2])/2.0, float32(sys.scrrect[3])/2.0, -rp.fLength))
			modelview = modelview.Mul4(mgl.Translate3D(rp.rcx, rp.rcy, 0))
		} else if rp.projectionMode == 2 {
			matrix := mgl.Mat4{float32(sys.scrrect[2] / 2.0), 0, 0, 0, 0, float32(sys.scrrect[3] / 2), 0, 0, 0, 0, -65535, 0, float32(sys.scrrect[2] / 2), float32(sys.scrrect[3] / 2), 0, 1}
			//modelview = modelview.Mul4(mgl.Translate3D(0, -float32(sys.scrrect[3]), 2048))
			modelview = modelview.Mul4(mgl.Translate3D(rp.rcx-float32(sys.scrrect[2])/2.0-rp.xOffset, rp.rcy-float32(sys.scrrect[3])/2.0+rp.yOffset, rp.fLength))
			modelview = modelview.Mul4(matrix)
			modelview = modelview.Mul4(mgl.Frustum(-float32(sys.scrrect[2])/2/rp.fLength, float32(sys.scrrect[2])/2/rp.fLength, -float32(sys.scrrect[3])/2/rp.fLength, float32(sys.scrrect[3])/2/rp.fLength, 1.0, 65535))
			modelview = modelview.Mul4(mgl.Translate3D(rp.xOffset, -rp.yOffset, -rp.fLength))
		}

		modelview = modelview.Mul4(mgl.Scale3D(1, rp.vs, 1))
		modelview = modelview.Mul4(
			mgl.Rotate3DX(-rp.rot.xangle * math.Pi / 180.0).Mul3(
				mgl.Rotate3DY(rp.rot.yangle * math.Pi / 180.0)).Mul3(
				mgl.Rotate3DZ(rp.rot.angle * math.Pi / 180.0)).Mat4())
		modelview = modelview.Mul4(mgl.Translate3D(-rp.rcx, -rp.rcy, 0))

		drawQuads(modelview, x1, y1, x2, y2, x3, y3, x4, y4)
		return
	}
	if rp.tile.y == 1 && rp.xbs != 0 {
		x1d, y1d, x2d, y2d, x3d, y3d, x4d, y4d := x1, y1, x2, y2, x3, y3, x4, y4
		for {
			x1d, y1d = x4d, y4d+rp.ys*rp.vs*float32(rp.tile.sy)
			x2d, y2d = x3d, y1d
			x3d = x4d - rp.rxadd*rp.ys*float32(rp.size[1]) + (rp.xts/rp.xbs)*(x3d-x4d)
			y3d = y2d + rp.ys*rp.vs*float32(rp.size[1])
			x4d = x4d - rp.rxadd*rp.ys*float32(rp.size[1])
			if AbsF(y3d-y4d) < 0.01 {
				break
			}
			y4d = y3d
			if rp.ys*(float32(rp.size[1])+float32(rp.tile.sy)) < 0 {
				if y1d <= float32(-sys.scrrect[3]) && y4d <= float32(-sys.scrrect[3]) {
					break
				}
			} else if y1d >= 0 && y4d >= 0 {
				break
			}
			if (0 > y1d || 0 > y4d) &&
				(y1d > float32(-sys.scrrect[3]) || y4d > float32(-sys.scrrect[3])) {
				rmTileHSub(modelview, x1d, y1d, x2d, y2d, x3d, y3d, x4d, y4d,
					float32(rp.size[0]), rp.tile, rp.rcx)
			}
		}
	}
	if rp.tile.y == 0 || rp.xts != 0 {
		n := rp.tile.y
		for {
			if rp.ys*(float32(rp.size[1])+float32(rp.tile.sy)) > 0 {
				if y1 <= float32(-sys.scrrect[3]) && y4 <= float32(-sys.scrrect[3]) {
					break
				}
			} else if y1 >= 0 && y4 >= 0 {
				break
			}
			if (0 > y1 || 0 > y4) &&
				(y1 > float32(-sys.scrrect[3]) || y4 > float32(-sys.scrrect[3])) {
				rmTileHSub(modelview, x1, y1, x2, y2, x3, y3, x4, y4,
					float32(rp.size[0]), rp.tile, rp.rcx)
			}
			if rp.tile.y != 1 && n != 0 {
				n--
			}
			if n == 0 {
				break
			}
			x4, y4 = x1, y1-rp.ys*rp.vs*float32(rp.tile.sy)
			x3, y3 = x2, y4
			x2 = x1 + rp.rxadd*rp.ys*float32(rp.size[1]) + (rp.xbs/rp.xts)*(x2-x1)
			y2 = y3 - rp.ys*rp.vs*float32(rp.size[1])
			x1 = x1 + rp.rxadd*rp.ys*float32(rp.size[1])
			if AbsF(y1-y2) < 0.01 {
				break
			}
			y1 = y2
		}
	}
}

func rmInitSub(rp *RenderParams) {
	if rp.vs < 0 {
		rp.vs *= -1
		rp.ys *= -1
		rp.rot.angle *= -1
		rp.rot.xangle *= -1
	}
	if rp.tile.x == 0 {
		rp.tile.sx = 0
	} else if rp.tile.sx > 0 {
		rp.tile.sx -= int32(rp.size[0])
	}
	if rp.tile.y == 0 {
		rp.tile.sy = 0
	} else if rp.tile.sy > 0 {
		rp.tile.sy -= int32(rp.size[1])
	}
	if rp.xts >= 0 {
		rp.x *= -1
	}
	rp.x += rp.rcx
	rp.rcy *= -1
	if rp.ys < 0 {
		rp.y *= -1
	}
	rp.y += rp.rcy
}

func RenderSprite(rp RenderParams) {
	if !rp.IsValid() {
		return
	}

	rmInitSub(&rp)

	neg, grayscale, padd, pmul := false, float32(0), [3]float32{0, 0, 0}, [3]float32{1, 1, 1}
	tint := [4]float32{float32(rp.tint&0xff)/255, float32(rp.tint>>8&0xff)/255,
		float32(rp.tint>>16&0xff)/255, float32(rp.tint>>24&0xff)/255}

	if rp.pfx != nil {
		neg, grayscale, padd, pmul = rp.pfx.getFcPalFx(rp.trans == -2)
		if rp.trans == -2 {
			padd[0], padd[1], padd[2] = -padd[0], -padd[1], -padd[2]
		}
	}

	proj := mgl.Ortho(0, float32(sys.scrrect[2]), 0, float32(sys.scrrect[3]), -65535, 65535)
	modelview := mgl.Translate3D(0, float32(sys.scrrect[3]), 0)

	gfx.Scissor(rp.window[0], rp.window[1], rp.window[2], rp.window[3])

	renderWithBlending(func(eq BlendEquation, src, dst BlendFunc, a float32) {

		gfx.SetPipeline(eq, src, dst)

		gfx.SetUniformMatrix("projection", proj[:])
		gfx.SetTexture("tex", rp.tex)
		if rp.paltex == nil {
			gfx.SetUniformI("isRgba", 1)
		} else {
			gfx.SetTexture("pal", rp.paltex)
			gfx.SetUniformI("isRgba", 0)
			gfx.SetUniformI("mask", int(rp.mask))
		}
		gfx.SetUniformI("isTrapez", int(Btoi(AbsF(AbsF(rp.xts)-AbsF(rp.xbs)) > 0.001)))
		gfx.SetUniformI("isFlat", 0)

		gfx.SetUniformI("neg", int(Btoi(neg)))
		gfx.SetUniformF("gray", grayscale)
		gfx.SetUniformFv("add", padd[:])
		gfx.SetUniformFv("mult", pmul[:])
		gfx.SetUniformFv("tint", tint[:])
		gfx.SetUniformF("alpha", a)

		rmTileSub(modelview, rp)

		gfx.ReleasePipeline()
	}, rp.trans, rp.paltex != nil)

	gfx.DisableScissor()
}

func renderWithBlending(render func(eq BlendEquation, src, dst BlendFunc, a float32), trans int32, correctAlpha bool) {
	blendSourceFactor := BlendSrcAlpha
	if !correctAlpha {
		blendSourceFactor = BlendOne
	}
	switch {
	case trans == -1:
		render(BlendAdd, blendSourceFactor, BlendOne, 1)
	case trans == -2:
		render(BlendReverseSubtract, BlendOne, BlendOne, 1)
	case trans <= 0:
	case trans < 255:
		render(BlendAdd, blendSourceFactor, BlendOneMinusSrcAlpha, float32(trans) / 255)
	case trans < 512:
		render(BlendAdd, blendSourceFactor, BlendOneMinusSrcAlpha, 1)
	default:
		src, dst := trans&0xff, trans>>10&0xff
		if dst < 255 {
			render(BlendAdd, BlendZero, BlendOneMinusSrcAlpha, 1 - float32(dst)/255)
		}
		if src > 0 {
			render(BlendAdd, blendSourceFactor, BlendOne, float32(src) / 255)
		}
	}
}

func FillRect(rect [4]int32, color uint32, trans int32) {
	r := float32(color>>16&0xff) / 255
	g := float32(color>>8&0xff) / 255
	b := float32(color&0xff) / 255

	modelview := mgl.Translate3D(0, float32(sys.scrrect[3]), 0)
	proj := mgl.Ortho(0, float32(sys.scrrect[2]), 0, float32(sys.scrrect[3]), -65535, 65535)

	x1, y1 := float32(rect[0]), -float32(rect[1])
	x2, y2 := float32(rect[0]+rect[2]), -float32(rect[1]+rect[3])

	renderWithBlending(func(eq BlendEquation, src, dst BlendFunc, a float32) {
		gfx.SetPipeline(eq, src, dst)
		gfx.SetVertexData(
			x2, y2, 1, 1,
			x2, y1, 1, 0,
			x1, y2, 0, 1,
			x1, y1, 0, 0)

		gfx.SetUniformMatrix("modelview", modelview[:])
		gfx.SetUniformMatrix("projection", proj[:])
		gfx.SetUniformI("isFlat", 1)
		gfx.SetUniformF("tint", r, g, b, a)
		gfx.RenderQuad()
		gfx.ReleasePipeline()
	}, trans, true)
}
