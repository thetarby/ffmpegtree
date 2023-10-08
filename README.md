# FfmpegTree
FfmpegTree is a library to generate ffmpeg complex filter scripts from go.

## Usage

To apply vintage filter starting from third second until the end;
```go
// select input file
i := NewInputNode("./test_assets/test-vid.mp4", nil, nil)

// apply curves filter 
c := NewCurvesFilter(i, "vintage")

// apply timeline editing so that filter starts from third second
c.Since(3)

// apply scale filter
res := NewScaleFilterNode(c, 100, 100, true)

// select out streams that you want in the final container
args := Select([]INode{res}, "out.mp4", nil)
```

More complex example, that negates color of the video in the middle and reverts back to original while displaying a text for each change. 
```go
i := NewInputNode("./test_assets/test-vid.mp4", nil, nil)
text1 := NewDrawTextFilter(i, "now it is normal", "", "(w-text_w)/2", "0", 40, 30)
text1.Until(3)

c := NewCurvesFilter(text1, "negative")
c.Since(3)
c.Until(5)

text2 := NewDrawTextFilter(c, "now it is negative", "", "(w-text_w)/2", "0", 40, 30)
text2.Since(3)
text2.Until(5)

text3 := NewDrawTextFilter(text2, "now it is back to normal", "", "(w-text_w)/2", "0", 40, 20)
text3.Since(5)

res := NewScaleFilterNode(text3, 350, 350, true)

args := Select([]INode{res}, "out.mp4", nil)
```

Two input files are used. Second video is scaled to 50x50, rotated and placed at the 4 corners of the first video as an overlay;
```go
i := NewSelectStreamNode(NewInputNode("./test_assets/test-vid.mp4", nil, nil), 0)
s := NewScaleFilterNode(i, 200, -2, true)
b := NewBoxBlurFilter(s, "min(w\\,h)/5", "min(cw\\,ch)/5", 1)
sb := NewScaleFilterNode(b, 200, 250, true)

ov := NewOverlayIntoMiddleFilterNode(sb, s)
i2 := NewSelectStreamNode(NewInputNode("./test_assets/test-vid-2.mp4", nil, nil), 0)
s2 := NewScaleFilterNode(i2, 50, 50, true)

var res INode = NewOverlayFilterNode(ov, s2, "0", "0")                                             // top left
res = NewOverlayFilterNode(res, s2, "main_w-overlay_w", "0")                                       // top right
res = NewOverlayFilterNode(res, NewRotateFilter(s2, "PI"), "0", "main_h-overlay_h")                // bottom left
res = NewOverlayFilterNode(res, NewRotateFilter(s2, "PI"), "main_w-overlay_w", "main_h-overlay_h") // bottom right

args := Select([]INode{res}, "out.mp4", []string{"-shortest"})
```

Tests can be examined for other usage examples. 

To add other filters simply add another filter node (as in filter_node.go) that either embeds BaseFilterNode or TimelineAcceptingFilterNode if it supports timeline editing.