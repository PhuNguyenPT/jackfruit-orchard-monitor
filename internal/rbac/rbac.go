package rbac

type Resource string
type Action string
type GroupName string

const (
	ResourceSoilCalibration Resource = "soil_calibration"
	ResourceSensors         Resource = "sensors"
)

const (
	ActionRead  Action = "r"
	ActionWrite Action = "w"
)

const (
	GroupAdmin GroupName = "admin"
	GroupUser  GroupName = "user"
)
