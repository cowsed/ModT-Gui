package modt


type PrinterStatus struct {
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