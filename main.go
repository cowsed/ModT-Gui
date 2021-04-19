package main

import (
	"github.com/cowsed/ModT-Gui/Scripts"

	g "github.com/AllenDang/giu"
)

var isOpen bool
var styleNum = 1
var message string = "Drop a file"
var filepath string = "Drag a File Here because file dialogs are hard"
var filecontents string = ""
var log string = "- Nothing in the log yet\n"

var recentStatus modt.PrinterInformation

var connected = false
var commChan chan modt.Comm
var currentInfo modt.PrinterInformation

func Connect() {

	if !connected {
		fmt.Println("Connecting...")
		connected = true
		go func() {
			commChan <- modt.COMM_PLAY
			s := modt.GetStatus(&commChan, &currentInfo)
			log += "- " + s + "\n"
			log += "- The printer probably became disconnected"
			connected = false
		}()
	}
}

func Pause() {
	commChan <- modt.COMM_PAUSE
}
func Play() {
	commChan <- modt.COMM_PLAY
}
func End() {
	commChan <- modt.COMM_END
}

func sendGCodeAsync() {

	go func() {
		if filecontents != "" {
			log += "Sending GCODE...\n"
			//Pause status getting
			Pause()
			modt.SendGcode(filecontents)
			//Resume status getting

			Play()
		} else {
			log += "- Can not send empty file"

		}
	}()
}

func loop() {
	controlButtons := g.Group()
	if connected {
		controlButtons = g.Group().Layout(
			//g.Line(
			//	g.Button("Load Filament").OnClick(loadFilamentAsync),
			//	g.Button("Unload Filament"),
			//),
			g.Spacing(),
			g.Line(
				g.Button("Send GCODE").OnClick(sendGCodeAsync),
				g.Button("Stop"),
			),
			g.Line(
				g.Button("Pause").OnClick(Pause),
				g.Button("Play").OnClick(Play),
			),
		)
	} else {
		controlButtons = g.Group().Layout(
			g.Label("Connect to interact with the printer"),
		)
	}
	var connectionButton *g.ButtonWidget = g.Button("Disconnect").OnClick(End)
	if !connected {
		connectionButton = g.Button("Connect").OnClick(Connect)
	}

	g.SingleWindowWithMenuBar("Mod-T-Controller").IsOpen(&isOpen).Layout(
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
						connectionButton,
						g.Dummy(0, 20),
						controlButtons,
						g.Spacing(),

						g.Line(
							g.InputText("##Filename", &filepath).Flags(g.InputTextFlags_ReadOnly).Size(-1.0),
						),
						currentInfo.Build(),
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
	commChan = make(chan modt.Comm, 1)
	modt.SaveUsbData()
	wnd := g.NewMasterWindow("MOD-T", 800, 600, 0, nil)
	wnd.SetDropCallback(onDrop)
	UpdateStyle()

	wnd.Run(loop)
}

/*
func loadFilamentAsync() {
	go func() {
		log += "Loading Filament..."
		//Pause status getting
		commChan <- modt.COMM_PAUSE

		modt.LoadFilamentTemp()
		//Resume status getting
		commChan <-modt.COMM_PLAY
	}()
}
*/
