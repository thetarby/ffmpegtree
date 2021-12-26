package ffmpegtree

import (
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

// TODO: variable names still seem to be random so tests fails time to time

func TestSplitNode(t *testing.T) {
	i1 := NewSelectStreamNode(NewInputNode("input_1.mp4", nil, nil), 0)
	s1 := NewScaleFilterNode(i1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayFilterNode(s3, s4)

	str := Select(ov2, "out.mp4")
	println(str)
	require.Regexp(t, regexp.MustCompile(`-i input_1.mp4\s*-filter_complex '\[0:0]split\[var_8_0]\[var_8_1];\[var_8_1]scale=101:101,setsar=1:1\[var_2];\[var_8_0]scale=100:100,setsar=1:1\[var_1];\[var_1]\[var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split\[var_7_0]\[var_7_1];\[var_7_1]scale=103:103,setsar=1:1\[var_5];\[var_7_0]scale=102:102,setsar=1:1\[var_4];\[var_4]\[var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out\.mp4`), str)
	//require.Equal(t, `-i input_1.mp4  -filter_complex '[0:0]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out.mp4`, str)
}

func TestMoreThanOneInput(t *testing.T) {
	i1, i2 := NewInputNode("input_1.mp4", nil, nil), NewInputNode("input_2.mp4", nil, nil)
	bg := NewScaleFilterNode(NewBoxBlurFilter(i1, "min(w\\,h)/5", "min(cw\\,ch)/5", 1), 200, 200, true)
	fg := NewChromaFilterNode(NewScaleFilterNode(i2, 100, 100, true), "#00ff00", 0.5)

	ov1 := NewOverlayFilterNode(bg, fg)

	str := Select(ov1, "out.mp4")
	println(str)
	require.Regexp(t, regexp.MustCompile(`-i .*\.mp4 -i .*\.mp4\s*-filter_complex\s*'\[1:0]scale=100:100,setsar=1:1,colorkey=#00ff00:0\.5\[var_4];\[0:0]boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=200:200,setsar=1:1\[var_2];\[var_2]\[var_4]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out\.mp4`), str)
	//require.Equal(t, `-i input_1.mp4 -i input_2.mp4  -filter_complex '[1:0]scale=100:100,setsar=1:1,colorkey=#00ff00:0.5[var_4];[0:0]boxblur=luma_radius=min(w\,h)/5:chroma_radius=min(cw\,ch)/5:luma_power=1,scale=200:200,setsar=1:1[var_2];[var_2][var_4]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out.mp4`, str)
}

func TestSelectStream_1(t *testing.T) {
	i := NewInputNode("input_1.mp4", nil, nil)
	i1 := NewSelectStreamNode(i, 1)
	s1 := NewScaleFilterNode(i1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayFilterNode(s3, s4)

	str := Select(ov2, "out.mp4", NewMap(ov2))
	println(str)
	require.Regexp(t, regexp.MustCompile(`-i input_1\.mp4\s*-filter_complex\s*'\[0:1]split\[var_8_0]\[var_8_1];\[var_8_1]scale=101:101,setsar=1:1\[var_2];\[var_8_0]scale=100:100,setsar=1:1\[var_1];\[var_1]\[var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split\[var_7_0]\[var_7_1];\[var_7_1]scale=103:103,setsar=1:1\[var_5];\[var_7_0]scale=102:102,setsar=1:1\[var_4];\[var_4]\[var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2\[var_6]' -map '\[var_6]' out\.mp4`), str)
	//require.Equal(t, `-i input_1.mp4 -filter_complex '[0:1]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2[var_6]' -map '[var_6]' out.mp4`, str)
}

func TestSelectStream_2(t *testing.T) {
	i1 := NewInputNode("input_1.mp4", nil, nil)
	i1Stream1 := NewSelectStreamNode(i1, 1)
	s1 := NewScaleFilterNode(i1Stream1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayFilterNode(s3, s4)

	str := Select(ov2, "out.mp4")
	println(str)
	//require.Equal(t, `-i input_1.mp4 -filter_complex '[0:1]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2[var_6]' -map '[var_6]' out.mp4`, str)
}

func TestMap(t *testing.T) {
	length := time.Second * 5
	i1 := NewInputNode("vid.mp4", &length, nil)
	s1 := NewScaleFilterNode(i1, 400, 400, true)
	s2 := NewDrawBoxFilter(s1, 0, 0, 400, 70, "#00ff00", "fill")

	str := Select(s2, "out.mp4")
	println(str)
	//require.Equal(t, `-i input_1.mp4 -filter_complex '[0:1]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2[var_6]' -map '[var_6]' out.mp4`, str)
}

func TestCurves(t *testing.T) {
	x := time.Second * 10
	var in INode = NewInputNode("vid.mp4", &x, nil)
	vs := 1.1
	in = NewVideoSpeedFilter(in, float32(1.0/vs))

	scaled := NewScaleFilterNode(in, 1200, -2, true)
	blurred := NewBoxBlurFilter(scaled, "min(w\\,h)/5", "min(cw\\,ch)/5", 1)
	scaledBlurred := NewScaleFilterNode(blurred, 1200, 1600, true)

	o := NewCurvesFilter(NewOverlayFilterNode(scaledBlurred, scaled), "vintage")
	str := Select(o, "out.mp4")
	println(str)
	reg := regexp.MustCompile(`-t 00:00:10 -i vid.mp4\s*-filter_complex '\[0:0]setpts=0.90909094\*PTS,scale=1200:-2,setsar=1:1,split\[var_7_0]\[var_7_1];\[var_7_1]boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=1200:1600,setsar=1:1\[var_4];\[var_4]\[var_7_0]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,curves=preset=vintage' out.mp4`)
	require.Regexp(t, reg, str)
}
