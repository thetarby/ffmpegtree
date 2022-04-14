package ffmpegtree

import "fmt"

type Stream interface {
	GetName() string
}

type IFilterNode interface {
	INode
	GetOutput() Stream
	GetOutStreamName() string
	FilterString() string
}

type BaseFilterNode struct {
	BaseNode
	OutStreamName string
}

func (b *BaseFilterNode) GetOutput() Stream {
	panic("implement me")
}

func (b *BaseFilterNode) GetOutStreamName() string {
	return b.OutStreamName
}

func (b *BaseFilterNode) FilterString() string {
	panic("implement me")
}

func NewBaseFilterNode(children []INode, outStreamName string) *BaseFilterNode {
	return &BaseFilterNode{
		BaseNode:      NewBaseNode(children),
		OutStreamName: outStreamName,
	}
}

type ScaleFilterNode struct {
	BaseFilterNode
	W, H   int
	SetSar bool
}

func (b *ScaleFilterNode) FilterString() string {
	if b.SetSar {
		return fmt.Sprintf("scale=%v:%v,setsar=1:1", b.W, b.H)
	}
	return fmt.Sprintf("scale=%v:%v", b.W, b.H)
}

func NewScaleFilterNode(input INode, w, h int, setsar bool) *ScaleFilterNode {
	return &ScaleFilterNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		W:              w,
		H:              h,
		SetSar:         setsar,
	}
}

/* Overlay Filter */
type OverlayIntoMiddleFilterNode struct {
	BaseFilterNode
}

func (n *OverlayIntoMiddleFilterNode) FilterString() string {
	return "overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2"
}

func NewOverlayIntoMiddleFilterNode(input1, input2 INode) *OverlayIntoMiddleFilterNode {
	return &OverlayIntoMiddleFilterNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input1, input2}, randStr()),
	}
}

type OverlayFilterNode struct {
	BaseFilterNode
	x, y string
}

func (n *OverlayFilterNode) FilterString() string {
	return fmt.Sprintf("overlay=%v:%v", n.x, n.y)
}

func NewOverlayFilterNode(input1, input2 INode, x, y string) *OverlayFilterNode {
	return &OverlayFilterNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input1, input2}, randStr()),
		x:              x,
		y:              y,
	}
}

type ChromaFilterNode struct {
	BaseFilterNode
	Color string
	Sim   float32
}

func (n *ChromaFilterNode) FilterString() string {
	return fmt.Sprintf("colorkey=%v:%v", n.Color, n.Sim)
}

func NewChromaFilterNode(input INode, color string, sim float32) *ChromaFilterNode {
	return &ChromaFilterNode{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		Color:          color,
		Sim:            sim,
	}
}

type VideoSpeedFilter struct {
	BaseFilterNode
	PresentationTimeStamps float32
}

func (s *VideoSpeedFilter) FilterString() string {
	return fmt.Sprintf("setpts=%v*PTS", s.PresentationTimeStamps)
}

func NewVideoSpeedFilter(input INode, presentationTimeStamps float32) *VideoSpeedFilter {
	return &VideoSpeedFilter{
		BaseFilterNode:         *NewBaseFilterNode([]INode{input}, randStr()),
		PresentationTimeStamps: presentationTimeStamps,
	}
}

type DrawBoxFilter struct {
	BaseFilterNode
	X      int
	Y      int
	Width  int
	Height int
	Color  string
	Type   string
}

func (f *DrawBoxFilter) FilterString() string {
	return fmt.Sprintf("drawbox=x=%v:y=%v:w=%v:h=%v:color=%v:t=%v", f.X, f.Y, f.Width, f.Height, f.Color, f.Type)
}

func NewDrawBoxFilter(input INode, x, y, w, h int, color, t string) *DrawBoxFilter {
	return &DrawBoxFilter{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		X:              x,
		Y:              y,
		Width:          w,
		Height:         h,
		Color:          color,
		Type:           t,
	}
}

type BoxBlurFilter struct {
	lumaRadius, chromaRadius string
	lumaPower                int
	BaseFilterNode
}

func (f *BoxBlurFilter) FilterString() string {
	return fmt.Sprintf("boxblur=luma_radius=%v:chroma_radius=%v:luma_power=%v", f.lumaRadius, f.chromaRadius, f.lumaPower)
}

func NewBoxBlurFilter(input INode, lumaRadius string, chromaRadius string, lumaPower int) *BoxBlurFilter {
	return &BoxBlurFilter{
		lumaRadius:     lumaRadius,
		chromaRadius:   chromaRadius,
		lumaPower:      lumaPower,
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
	}
}

type CurvesFilter struct {
	BaseFilterNode
	preset string
}

func (f *CurvesFilter) FilterString() string {
	return fmt.Sprintf("curves=preset=%v", f.preset)
}

func NewCurvesFilter(input INode, preset string) *CurvesFilter {
	return &CurvesFilter{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		preset:         preset,
	}
}

type RotateFilter struct {
	BaseFilterNode
	rotateExpr string
}

func (f *RotateFilter) FilterString() string {
	return fmt.Sprintf("rotate=%v", f.rotateExpr)
}

func NewRotateFilter(input INode, rotateExpr string) *RotateFilter {
	return &RotateFilter{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		rotateExpr:     rotateExpr,
	}
}

type AtempoFilter struct {
	BaseFilterNode
	speed float32
}

func (f *AtempoFilter) FilterString() string {
	return fmt.Sprintf("atempo=%.2f", f.speed)
}

func NewAtempoFilter(input INode, speed float32) *AtempoFilter {
	return &AtempoFilter{
		BaseFilterNode: *NewBaseFilterNode([]INode{input}, randStr()),
		speed:          speed,
	}
}
