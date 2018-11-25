package config

type UserAction string

const (
	ActionCopy        UserAction = "copy"
	ActionPaste       UserAction = "paste"
	ActionGoogle      UserAction = "google"
	ActionReportBug   UserAction = "report"
	ActionToggleDebug UserAction = "debug"
	ActionToggleSlomo UserAction = "slomo"
)
