package config

type UserAction string

const (
	ActionCopy        UserAction = "copy"
	ActionPaste       UserAction = "paste"
	ActionSearch      UserAction = "search"
	ActionReportBug   UserAction = "report"
	ActionToggleDebug UserAction = "debug"
	ActionToggleSlomo UserAction = "slomo"
)
