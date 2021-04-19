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
