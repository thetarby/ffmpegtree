package ffmpegtree

import (
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

// TODO: variable names still seem to be random so tests fails time to time

/**
 * Parses url with the given regular expression and returns the 
 * group values defined in the expression.
 *
 */
 func getParams(regEx *regexp.Regexp, str string) (paramsMap map[string]string) {
    match := regEx.FindStringSubmatch(str)

    paramsMap = make(map[string]string)
    for i, name := range regEx.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }
    return paramsMap
}

func TestSplitNode(t *testing.T) {
	i1 := NewSelectStreamNode(NewInputNode("input_1.mp4", nil, nil), 0)
	s1 := NewScaleFilterNode(i1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayFilterNode(s3, s4)

	str := Select(ov2, "out.mp4")
	reg := regexp.MustCompile(`-i input_1.mp4\s*-filter_complex '(?P<s1>\[.*])split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])scale=101:101,setsar=1:1(?P<s4>\[.*]);(?P<s2_2>\[.*])scale=100:100,setsar=1:1(?P<s5>\[.*]);(?P<s5_2>\[.*])(?P<s4_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s6>\[.*])(?P<s7>\[.*]);(?P<s7_2>\[.*])scale=103:103,setsar=1:1(?P<s8>\[.*]);(?P<s6_2>\[.*])scale=102:102,setsar=1:1(?P<s9>\[.*]);(?P<s9_2>\[.*])(?P<s8_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out\.mp4`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str)
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s3"], params["s3_2"])
	require.Equal(t, params["s4"], params["s4_2"])
	require.Equal(t, params["s5"], params["s5_2"])
	require.Equal(t, params["s6"], params["s6_2"])
	require.Equal(t, params["s7"], params["s7_2"])
	require.Equal(t, params["s8"], params["s8_2"])
	require.Equal(t, params["s9"], params["s9_2"])
	//require.Equal(t, `-i input_1.mp4  -filter_complex '[0:0]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out.mp4`, str)
}

func TestMoreThanOneInput(t *testing.T) {
	i1, i2 := NewInputNode("input_1.mp4", nil, nil), NewInputNode("input_2.mp4", nil, nil)
	bg := NewScaleFilterNode(NewBoxBlurFilter(i1, "min(w\\,h)/5", "min(cw\\,ch)/5", 1), 200, 200, true)
	fg := NewChromaFilterNode(NewScaleFilterNode(i2, 100, 100, true), "#00ff00", 0.5)

	ov1 := NewOverlayFilterNode(bg, fg)

	str := Select(ov1, "out.mp4")
	reg := regexp.MustCompile(`-i .*\.mp4 -i .*\.mp4\s*-filter_complex\s*'(?P<s1>\[.*])scale=100:100,setsar=1:1,colorkey=#00ff00:0\.5(?P<s2>\[.*]);(?P<s3>\[.*])boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=200:200,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out\.mp4`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str)
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s4"], params["s4_2"])
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
	reg := regexp.MustCompile(`-i input_1\.mp4\s*-filter_complex\s*'(?P<s1>\[.*])split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])scale=101:101,setsar=1:1(?P<s4>\[.*]);(?P<s2_2>\[.*])scale=100:100,setsar=1:1(?P<s5>\[.*]);(?P<s5_2>\[.*])(?P<s4_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s6>\[.*])(?P<s7>\[.*]);(?P<s7_2>\[.*])scale=103:103,setsar=1:1(?P<s8>\[.*]);(?P<s6_2>\[.*])scale=102:102,setsar=1:1(?P<s9>\[.*]);(?P<s9_2>\[.*])(?P<s8_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2(?P<s10>\[.*])' -map '(?P<s10_2>\[.*])' out\.mp4`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str)
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s3"], params["s3_2"])
	require.Equal(t, params["s4"], params["s4_2"])
	require.Equal(t, params["s5"], params["s5_2"])
	require.Equal(t, params["s6"], params["s6_2"])
	require.Equal(t, params["s7"], params["s7_2"])
	require.Equal(t, params["s8"], params["s8_2"])
	require.Equal(t, params["s9"], params["s9_2"])
	require.Equal(t, params["s10"], params["s10_2"])
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
	reg := regexp.MustCompile(`-i input_1.mp4  -filter_complex '(?P<s1>\[.*])scale=101:101,setsar=1:1(?P<s2>\[.*]);(?P<s3>\[.*])scale=100:100,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s5>\[.*])(?P<s6>\[.*]);(?P<s6_2>\[.*])scale=103:103,setsar=1:1(?P<s7>\[.*]);(?P<s5_2>\[.*])scale=102:102,setsar=1:1(?P<s8>\[.*]);(?P<s8_2>\[.*])(?P<s7_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out.mp4`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str)
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s4"], params["s4_2"])
	require.Equal(t, params["s5"], params["s5_2"])
	require.Equal(t, params["s6"], params["s6_2"])
	require.Equal(t, params["s7"], params["s7_2"])
	require.Equal(t, params["s8"], params["s8_2"])
	//-i input_1.mp4  -filter_complex '[0:0]scale=101:101,setsar=1:1[var_2];[0:1]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2' out.mp4
}

func TestMap(t *testing.T) {
	length := time.Second * 5
	i1 := NewInputNode("vid.mp4", &length, nil)
	i2 := NewInputNode("vid3.mp4", &length, nil)
	s1 := NewScaleFilterNode(i1, 400, 400, true)
	s2 := NewDrawBoxFilter(s1, 0, 0, 400, 70, "#00ff00", "fill")

	str := Select(s2, "out.mp4", NewMap(i2, "a"))
	require.Equal(t, `-t 00:00:05 -i vid.mp4 -t 00:00:05 -i vid3.mp4  -filter_complex '[0:0]scale=400:400,setsar=1:1,drawbox=x=0:y=0:w=400:h=70:color=#00ff00:t=fill' -map '1:a' out.mp4`, str)
	// -t 00:00:05 -i vid.mp4 -t 00:00:05 -i vid3.mp4  -filter_complex '[0:0]scale=400:400,setsar=1:1,drawbox=x=0:y=0:w=400:h=70:color=#00ff00:t=fill' -map '1:a' out.mp4
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
	reg := regexp.MustCompile(`-t 00:00:10 -i vid.mp4\s*-filter_complex '(?P<s1>\[.*])setpts=0.90909094\*PTS,scale=1200:-2,setsar=1:1,split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=1200:1600,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,curves=preset=vintage' out.mp4`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str)
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s3"], params["s3_2"])
	require.Equal(t, params["s4"], params["s4_2"])
}
