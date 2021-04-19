package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"io/ioutil"
)

//Helpers

func UpdateStyle() {
	switch styleNum {
	case 0:
		imgui.StyleColorsClassic()
	case 1:
		imgui.StyleColorsDark()
	case 2:
		imgui.StyleColorsLight()
	case 3:
		SetBi()

	}
}

func onDrop(names []string) {
	filename := names[0]
	if filename[len(filename)-3:] == ".nc" || filename[len(filename)-6:] == ".gcode" {
		filepath = filename
		log += "- Chose " + filename + "\n"
		loadFileToEditor(filename)
		g.Update()

	} else {
		//Trigger alert
		log += "- Incorrect File Type\n"
	}

}
func loadFileToEditor(fname string) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		log += "- Error Loading File"
	}
	filecontents = string(content)
}

func SetBi() {
	c1 := imgui.Vec4{1, 0.459, 0.839, 1.0}
	c2 := imgui.Vec4{0.741, 0.376, 0.706, 1.0}
	c3 := imgui.Vec4{0.639, 0.49, 0.875, 1.0}
	c4 := imgui.Vec4{0.353, 0.682, 0.882, 1.0}
	c5 := imgui.Vec4{0.533, 0.776, 0.925, 1.0}

	b := imgui.Vec4{0.0, 0.0, 0.0, 1.0}

	imgui.StyleColorsClassic()
	s := imgui.CurrentStyle()

	s.SetColor(imgui.StyleColorText, b)

	s.SetColor(imgui.StyleColorMenuBarBg, c3)

	s.SetColor(imgui.StyleColorPopupBg, c2)

	s.SetColor(imgui.StyleColorButton, c1)
	s.SetColor(imgui.StyleColorButtonHovered, c3)
	s.SetColor(imgui.StyleColorWindowBg, c5)

	s.SetColor(imgui.StyleColorTextSelectedBg, c4)
}
