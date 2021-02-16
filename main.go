package main

import (
	"./Scripts"

	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"io/ioutil"
	_ "io/ioutil"
	_ "strings"
)

var message string = "Drop a file"
var filepath string = "Drag a File Here because file dialogs are hard"
var filecontents string = ""
var log string = "- Nothing in the log yet\n"

var recentStatus modt.PrinterStatus

var commChan chan bool

func loadFileToEditor(fname string) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		log += "- Error Loading File"
	}
	filecontents = string(content)
}

func onDrop(names []string) {
	filename := names[0]
	if filename[len(filename)-3:] == ".nc" || filename[len(filename)-6:] == ".gcode" {
		filepath = filename
		log += "- Chose " + filename
		loadFileToEditor(filename)
		g.Update()

	} else {
		//Trigger alert
		log += "- Incorrect File Type\n"
	}

}


func Connect(){
	go func(){
		s:=modt.GetStatus(&commChan)
		log+=s+"\n"
	}()
}


func loadFilamentAsync(*chan bool) {
	go func() {
		//Pause status getting
		commChan <- false

		modt.LoadFilamentTemp()
		//Resume status getting
		commChan <- true

	}()
}

func sendGCodeAsync() {

	go func() {
		//Pause status getting
		commChan <- false
		modt.SendGcodeTemp(filepath)
		//Resume status getting
		commChan <- true

	}()
}

func loop() {
	//g.PushItemSpacing(0,10)
	imgui.StyleColorsClassic()

	g.SingleWindowWithMenuBar("On Drop Demo").Layout(
		g.MenuBar().Layout(
			g.Menu("File").Layout(
				g.MenuItem("Open GCODE"),
				g.MenuItem("Stop"),
				g.MenuItem("View Log"),
			),
			g.Menu("Help").Layout(
				g.MenuItem("Show full status"),
				g.MenuItem("Show Help Text"),
			),
		),

		g.SplitLayout("MainSplit", g.DirectionHorizontal, true, 300,
			g.Layout{
				g.Label("Printer Status"),
				g.SplitLayout("ControlSplit", g.DirectionVertical, true, 300,
					g.Layout{
						g.Button("Connect").OnClick(Connect),
						g.Line(
							g.Button("Load Filament"),
							g.Button("Unload Filament"),
						),
						g.Spacing(),

						g.Line(
							g.InputText("##Filename", &filepath).Flags(g.InputTextFlagsReadOnly).Size(-1.0),
						),
						g.Spacing(),
						g.Line(
							g.Button("Send GCODE").OnClick(sendGCodeAsync),
							g.Button("Stop"),
						),
					},
					g.Layout{
						g.InputTextMultiline("Log", &log).Flags(0).Size(-1, -1),
					},
				),
			},
			g.Layout{
				g.InputTextMultiline("##DroppedFiles", &filecontents).Flags(g.InputTextFlagsReadOnly+g.InputTextFlagsAlwaysInsertMode).Size(-1, -1),
			},
		),
	)
}

func main() {
	//Start monitoring and set up channel to pause that
	//go modt.Status
	commChan = make(chan bool, 1)

	wnd := g.NewMasterWindow("Hello world", 800, 600, 0, nil)
	wnd.SetDropCallback(onDrop)

	wnd.Run(loop)
}
