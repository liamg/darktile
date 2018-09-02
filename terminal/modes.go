package terminal

import "fmt"

func csiSetMode(modeStr string, enabled bool, terminal *Terminal) error {

	/*
	   Mouse support

	   		#define SET_X10_MOUSE               9
	        #define SET_VT200_MOUSE             1000
	        #define SET_VT200_HIGHLIGHT_MOUSE   1001
	        #define SET_BTN_EVENT_MOUSE         1002
	        #define SET_ANY_EVENT_MOUSE         1003

	        #define SET_FOCUS_EVENT_MOUSE       1004

	        #define SET_EXT_MODE_MOUSE          1005
	        #define SET_SGR_EXT_MODE_MOUSE      1006
	        #define SET_URXVT_EXT_MODE_MOUSE    1015

	        #define SET_ALTERNATE_SCROLL        1007
	*/

	switch modeStr {
	case "4":
		if enabled { // @todo support replace mode
			terminal.ActiveBuffer().SetInsertMode()
		} else {
			terminal.ActiveBuffer().SetReplaceMode()
		}
	case "?1":
		terminal.modes.ApplicationCursorKeys = enabled
	case "?7":
		// auto-wrap mode
		//DECAWM
		terminal.ActiveBuffer().SetAutoWrap(enabled)
	case "?9":
		if enabled {
			terminal.logger.Infof("Turning on X10 mouse mode")
			terminal.SetMouseMode(MouseModeX10)
		} else {
			terminal.logger.Infof("Turning off X10 mouse mode")
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?12", "?13":
		terminal.modes.BlinkingCursor = enabled
	case "?25":
		terminal.modes.ShowCursor = enabled
	case "?47", "?1047":
		if enabled {
			terminal.UseAltBuffer()
		} else {
			terminal.UseMainBuffer()
		}
	case "?1000", "?1006;1000", "?10061000": // ?10061000 seen from htop
		// enable mouse tracking
		// 1000 refers to ext mode for extended mouse click area - otherwise only x <= 255-31
		if enabled {
			terminal.logger.Infof("Turning on VT200 mouse mode")
			terminal.SetMouseMode(MouseModeVT200)
		} else {
			terminal.logger.Infof("Turning off VT200 mouse mode")
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?1048":
		if enabled {
			terminal.ActiveBuffer().SaveCursor()
		} else {
			terminal.ActiveBuffer().RestoreCursor()
		}
	case "?1049":
		if enabled {
			terminal.UseAltBuffer()
		} else {
			terminal.UseMainBuffer()
		}
	default:
		return fmt.Errorf("Unsupported CSI %sl code", modeStr)
	}

	return nil
}
