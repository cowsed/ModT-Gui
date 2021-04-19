package main

import (
	"github.com/cowsed/ModT-Gui/Scripts"

	g "github.com/AllenDang/giu"

)

var styleNum = 0
var message string = "Drop a file"
var filepath string = "Drag a File Here because file dialogs are hard"
var filecontents string = ""
var log string = "- Nothing in the log yet\n"

var recentStatus modt.PrinterStatus

var connected = false
var commChan chan bool

func Connect() {
	if !connected {
		connected = true
		go func() {
			commChan <- true
			s := modt.GetStatus(&commChan)
			log += "- " + s + "\n"
			log += "- The printer probably became disconnected"
			connected = false
		}()
	}
}

func Pause() {
	commChan <- false
}
func Play() {
	commChan <- true
}

func loadFilamentAsync() {
	go func() {
		log += "Loading Filament..."
		//Pause status getting
		commChan <- false
		connected = false

		modt.LoadFilamentTemp()
		//Resume status getting
		Connect()
	}()
}

func sendGCodeAsync() {

	go func() {
		if filecontents != "" {
			log += "Sending GCODE..."
			//Pause status getting
			commChan <- false
			connected = false
			modt.SendGcodeTemp(filepath)
			//Resume status getting

			Connect()
		} else {
			log += "- Can not send empty file"

		}
	}()
}

func loop() {
	controlButtons := g.Group()
	if connected {
		controlButtons = g.Group().Layout(
			g.Line(
				g.Button("Load Filament").OnClick(loadFilamentAsync),
				g.Button("Unload Filament"),
			),
			g.Spacing(),
			g.Line(
				g.Button("Send GCODE").OnClick(sendGCodeAsync),
				g.Button("Stop"),
			),
		)
	} else {
		controlButtons = g.Group().Layout(
			g.Label("Connect to interact with the printer"),
		)
	}

	g.SingleWindowWithMenuBar("On Drop Demo").Layout(
		g.MenuBar().Layout(
			g.Menu("File").Layout(
				g.MenuItem("Open GCODE"),
				g.MenuItem("Stop"),
				g.MenuItem("View Log"),
			), g.Menu("View").Layout(
				g.Menu("Style").Layout(
					g.RadioButton("Classic", styleNum == 0).OnChange(func() { styleNum = 0; UpdateStyle() }),
					g.RadioButton("Dark", styleNum == 1).OnChange(func() { styleNum = 1; UpdateStyle() }),
					g.RadioButton("Light", styleNum == 2).OnChange(func() { styleNum = 2; UpdateStyle() }),
					g.RadioButton("Bi", styleNum == 3).OnChange(func() { styleNum = 3; UpdateStyle() }),
				),
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
						g.Dummy(0, 20),
						controlButtons,
						g.Spacing(),
						g.Button("Pause").OnClick(Pause),
						g.Button("Play").OnClick(Play),

						g.Line(
							g.InputText("##Filename", &filepath).Flags(g.InputTextFlags_ReadOnly).Size(-1.0),
						),
					},
					g.Layout{
						g.InputTextMultiline("Log", &log).Flags(0).Size(-1, -1),
					},
				),
			},
			g.Layout{
				g.InputTextMultiline("##DroppedFiles", &filecontents).Flags(g.InputTextFlags_ReadOnly).Size(-1, -1),
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

	UpdateStyle()

	wnd.Run(loop)
}
