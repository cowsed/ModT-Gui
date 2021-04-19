package modt

import (
	"fmt"
	g "github.com/AllenDang/giu"
)

func (p PrinterInformation) Build() g.Widget {
	return g.Group().Layout(
		g.TreeNode("Printer Information").Layout(
			g.TreeNode("Metadata").Layout(
				g.Label(fmt.Sprintf("Version: %v", p.Metadata.Version)),
				g.Label(fmt.Sprintf("Type: %v", p.Metadata.Type)),
			).Flags(g.TreeNodeFlagsFramed),
			g.TreeNode("Printer").Layout(
				g.Label(fmt.Sprintf("Model: %v", p.Printer.ModelName)),
				g.Label(fmt.Sprintf("CPUID: %v", p.Printer.CPUID)),
				g.TreeNode("Firmware:").Layout(
					g.Label(fmt.Sprintf("Name: %v",p.Printer.Firmware.Name)),
					g.Label(fmt.Sprintf("Version: %v",p.Printer.Firmware.Version)),
				).Flags(g.TreeNodeFlagsFramed),
				g.Label(fmt.Sprintf("Accept Version: %v", p.Printer.AcceptVersion)),
			).Flags(g.TreeNodeFlagsFramed),

			g.TreeNode("Status").Layout(
				g.Label(fmt.Sprintf("State: %v",p.Status.State)),
				g.Label(fmt.Sprintf("Plate: %v",p.Status.BuildPlate)),
				g.Label(fmt.Sprintf("Filament: %v",p.Status.Filament)),
				g.Label(fmt.Sprintf("Extruder Temp: %vC",p.Status.ExtruderTemperature)),
				g.Label(fmt.Sprintf("Extruder Target: %vC",p.Status.ExtruderTargetTemperature)),
			).Flags(g.TreeNodeFlagsFramed),
			g.TreeNode("Job").Layout(
				g.Label(fmt.Sprint(p.Job)),
			).Flags(g.TreeNodeFlagsFramed),
			g.TreeNode("Time").Layout(
				g.Label(fmt.Sprintf("Boot Time: %d:%d", p.Time.Boot/60,p.Time.Boot%60)),
				g.Label(fmt.Sprintf("Idle Time: %d:%d", p.Time.Idle/60,p.Time.Idle%60)),
			).Flags(g.TreeNodeFlagsFramed),
		).Flags(g.TreeNodeFlagsFramed),
	)

}

type PrinterInformation struct {
	Metadata Metadata `json:"metadata"`
	Printer  Printer  `json:"printer"`
	Error    int      `json:"error"`
	Status   Status   `json:"status"`
	Job      Job      `json:"job"`
	Time     Time     `json:"time"`
}
type Metadata struct {
	Version int    `json:"version"`
	Type    string `json:"type"`
}
type Firmware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type Printer struct {
	ModelName     string   `json:"model_name"`
	CPUID         string   `json:"cpu_id"`
	Firmware      Firmware `json:"firmware"`
	AcceptVersion int      `json:"accept_version"`
}
type Status struct {
	State                     string  `json:"state"`
	BuildPlate                string  `json:"build_plate"`
	Filament                  string  `json:"filament"`
	ExtruderTemperature       float64 `json:"extruder_temperature"`
	ExtruderTargetTemperature float64 `json:"extruder_target_temperature"`
}
type Job struct {
	ID                 string `json:"id"`
	Source             string `json:"source"`
	Progress           int    `json:"progress"`
	RxProgress         int    `json:"rx_progress"`
	CurrentLineNumber  int    `json:"current_line_number"`
	CurrentGcodeNumber int    `json:"current_gcode_number"`
	FileSize           int    `json:"file_size"`
	File               string `json:"file"`
}
type Time struct {
	Boot int `json:"boot"`
	Idle int `json:"idle"`
}
