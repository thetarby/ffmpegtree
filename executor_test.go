package ffmpegtree

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

	ov1 := NewOverlayIntoMiddleFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayIntoMiddleFilterNode(s3, s4)

	str := Select(ov2, "out.mp4", nil)
	reg := regexp.MustCompile(`(?P<s1>\[.*])split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])scale=101:101,setsar=1:1(?P<s4>\[.*]);(?P<s2_2>\[.*])scale=100:100,setsar=1:1(?P<s5>\[.*]);(?P<s5_2>\[.*])(?P<s4_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s6>\[.*])(?P<s7>\[.*]);(?P<s7_2>\[.*])scale=103:103,setsar=1:1(?P<s8>\[.*]);(?P<s6_2>\[.*])scale=102:102,setsar=1:1(?P<s9>\[.*]);(?P<s9_2>\[.*])(?P<s8_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2`)
	require.Regexp(t, reg, str.FilterComplex())
	params := getParams(reg, str.FilterComplex())
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

	ov1 := NewOverlayIntoMiddleFilterNode(bg, fg)

	str := Select(ov1, "out.mp4", nil)
	reg := regexp.MustCompile(`(?P<s1>\[.*])scale=100:100,setsar=1:1,colorkey=#00ff00:0\.5(?P<s2>\[.*]);(?P<s3>\[.*])boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=200:200,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2`)
	require.Regexp(t, reg, str)

	params := getParams(reg, str.FilterComplex())
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s4"], params["s4_2"])
}

func TestSelectStream_1(t *testing.T) {
	i := NewInputNode("input_1.mp4", nil, nil)
	i1 := NewSelectStreamNode(i, 1)
	s1 := NewScaleFilterNode(i1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayIntoMiddleFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayIntoMiddleFilterNode(s3, s4)

	str := Select(ov2, "out.mp4", nil, NewMap(ov2))
	
	reg := regexp.MustCompile(`(?P<s1>\[.*])split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])scale=101:101,setsar=1:1(?P<s4>\[.*]);(?P<s2_2>\[.*])scale=100:100,setsar=1:1(?P<s5>\[.*]);(?P<s5_2>\[.*])(?P<s4_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s6>\[.*])(?P<s7>\[.*]);(?P<s7_2>\[.*])scale=103:103,setsar=1:1(?P<s8>\[.*]);(?P<s6_2>\[.*])scale=102:102,setsar=1:1(?P<s9>\[.*]);(?P<s9_2>\[.*])(?P<s8_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2(?P<s10>\[.*])`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str.FilterComplex())
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s3"], params["s3_2"])
	require.Equal(t, params["s4"], params["s4_2"])
	require.Equal(t, params["s5"], params["s5_2"])
	require.Equal(t, params["s6"], params["s6_2"])
	require.Equal(t, params["s7"], params["s7_2"])
	require.Equal(t, params["s8"], params["s8_2"])
	require.Equal(t, params["s9"], params["s9_2"])
	//require.Equal(t, `-i input_1.mp4 -filter_complex '[0:1]split[var_8_0][var_8_1];[var_8_1]scale=101:101,setsar=1:1[var_2];[var_8_0]scale=100:100,setsar=1:1[var_1];[var_1][var_2]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split[var_7_0][var_7_1];[var_7_1]scale=103:103,setsar=1:1[var_5];[var_7_0]scale=102:102,setsar=1:1[var_4];[var_4][var_5]overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2[var_6]' -map '[var_6]' out.mp4`, str)
}

func TestSelectStream_2(t *testing.T) {
	i1 := NewInputNode("input_1.mp4", nil, nil)
	i1Stream1 := NewSelectStreamNode(i1, 1)
	s1 := NewScaleFilterNode(i1Stream1, 100, 100, true)
	s2 := NewScaleFilterNode(i1, 101, 101, true)

	ov1 := NewOverlayIntoMiddleFilterNode(s1, s2)

	s3 := NewScaleFilterNode(ov1, 102, 102, true)
	s4 := NewScaleFilterNode(ov1, 103, 103, true)

	ov2 := NewOverlayIntoMiddleFilterNode(s3, s4)

	str := Select(ov2, "out.mp4", nil)
	
	reg := regexp.MustCompile(`(?P<s1>\[.*])scale=101:101,setsar=1:1(?P<s2>\[.*]);(?P<s3>\[.*])scale=100:100,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,split(?P<s5>\[.*])(?P<s6>\[.*]);(?P<s6_2>\[.*])scale=103:103,setsar=1:1(?P<s7>\[.*]);(?P<s5_2>\[.*])scale=102:102,setsar=1:1(?P<s8>\[.*]);(?P<s8_2>\[.*])(?P<s7_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2`)
	require.Regexp(t, reg, str)
	params := getParams(reg, str.FilterComplex())
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

	str := Select(s2, "out.mp4", nil, NewMap(i2, "a"))
	require.Equal(t, "-t", str[0])
	require.Equal(t, "00:00:05", str[1])
	require.Equal(t, "-i", str[2])
	require.Equal(t, "vid.mp4", str[3])
	require.Equal(t, "-t", str[4])
	require.Equal(t, "00:00:05", str[5])
	require.Equal(t, "-i", str[6])
	require.Equal(t, "vid3.mp4", str[7])
	require.Equal(t, "-filter_complex", str[8])
	require.Equal(t, `[0:0]scale=400:400,setsar=1:1,drawbox=x=0:y=0:w=400:h=70:color=#00ff00:t=fill`, str.FilterComplex())
	require.Equal(t, `-map`, str[10])
	require.Equal(t, `1:a`, str[11])
	require.Equal(t, `out.mp4`, str[12])
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

	o := NewCurvesFilter(NewOverlayIntoMiddleFilterNode(scaledBlurred, scaled), "vintage")
	str := Select(o, "out.mp4", nil)

	//-t 00:00:10 -i vid.mp4\s*-filter_complex '

	reg := regexp.MustCompile(`(?P<s1>\[.*])setpts=0.90909094\*PTS,scale=1200:-2,setsar=1:1,split(?P<s2>\[.*])(?P<s3>\[.*]);(?P<s3_2>\[.*])boxblur=luma_radius=min\(w\\,h\)/5:chroma_radius=min\(cw\\,ch\)/5:luma_power=1,scale=1200:1600,setsar=1:1(?P<s4>\[.*]);(?P<s4_2>\[.*])(?P<s2_2>\[.*])overlay=main_w/2-overlay_w/2:main_h/2-overlay_h/2,curves=preset=vintage`)
	
	require.Equal(t, "-t", str[0])
	require.Equal(t, "00:00:10", str[1])
	require.Equal(t, "-i", str[2])
	require.Equal(t, "vid.mp4", str[3])

	require.Regexp(t, reg, str)
	params := getParams(reg, str.FilterComplex())
	require.Equal(t, params["s2"], params["s2_2"])
	require.Equal(t, params["s3"], params["s3_2"])
	require.Equal(t, params["s4"], params["s4_2"])
}

func TestSelectMoreThanOneStream(t *testing.T) {
	x := time.Second * 10
	var in IInputNode = NewInputNode("vid.mp4", &x, nil)
	vs := 2.0
	var inVStream INode = NewVideoSpeedFilter(in, float32(1.0/vs))
	var scaledInVStream INode = NewScaleFilterNode(inVStream, 100, 100, false)
	inVStreamFinal := NewScaleFilterNode(scaledInVStream, 10, 10, false)
	
	var differVStream INode = NewScaleFilterNode(scaledInVStream, 200, 200, false)
	differVStream = NewScaleFilterNode(differVStream, 300, 300, false)


	exec := NewFfmpegExecutor(nil, "out.mp4", nil)
	res := exec.ToFfmpeg(differVStream, inVStreamFinal)

	reg := regexp.MustCompile(`(?P<s1>\[.*])setpts=0.5\*PTS,scale=100:100(?P<s2>\[.*]);(?P<s3>\[.*])scale=10:10;(?P<s4>\[.*])scale=200:200,scale=300:300`)
	// [0:0]setpts=0.5*PTS,scale=100:100[var_2];[var_2]scale=10:10;[var_2]scale=200:200,scale=300:300
	require.Regexp(t, reg, res.FilterComplex())
	params := getParams(reg, res.FilterComplex())
	require.Equal(t, params["s2"], params["s3"])
	require.Equal(t, params["s3"], params["s4"])
	fmt.Println(res)
}

func TestCurvesWithTimeline(t *testing.T) {
	i := NewInputNode("./test_assets/test-vid.mp4", nil, nil)
	c := NewCurvesFilter(i, "vintage")
	c.Since(3)
	res := NewScaleFilterNode(c, 100, 100, true)
	
	exec := NewFfmpegExecutor(nil, "out.mp4", nil)
	args := exec.ToFfmpeg(res)
	
	// [0:0]curves=preset=vintage:enable='between(t, 1,5)',scale=100:100,setsar=1:1
	require.Equal(t, `[0:0]curves=preset=vintage:enable='gte(t, 3.00)',scale=100:100,setsar=1:1`, args.FilterComplex())
}