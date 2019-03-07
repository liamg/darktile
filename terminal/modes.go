package terminal

import (
	"errors"
	"fmt"
	"strings"
)

func recoverCodeFromEnabled(enabled bool) string {
	code := ""
	if enabled {
		code = "h"
	} else {
		code = "l"
	}
	return code
}

func csiSetModes(modes []string, enabled bool, terminal *Terminal) error {
	if len(modes) == 0 {
		return fmt.Errorf("CSI %s without parameters is not allowed", recoverCodeFromEnabled(enabled))
	}
	if len(modes) == 1 {
		return csiSetMode(modes[0], enabled, terminal)
	}
	// should we propagate DEC prefix?
	const decPrefix = '?'
	isDec := len(modes[0]) > 0 && modes[0][0] == decPrefix

	// iterate through params, propagating DEC prefix to subsequent elements
	errorStrings := make([]string, 0)
	for i, v := range modes {
		updatedMode := v
		if i > 0 && isDec {
			updatedMode = string(decPrefix) + v
		}
		err := csiSetMode(updatedMode, enabled, terminal)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf(strings.Join(errorStrings, "\n"))
	}

	return nil
}

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
			terminal.SetInsertMode()
		} else {
			terminal.SetReplaceMode()
		}
	case "20":
		if enabled {
			terminal.SetNewLineMode()
		} else {
			terminal.SetLineFeedMode()
		}
	case "?1":
		terminal.modes.ApplicationCursorKeys = enabled
	case "?3":
		_, lines := terminal.GetSize()
		if enabled {
			// DECCOLM - COLumn mode, 132 characters per line
			terminal.SetSize(132, uint(lines))
		} else {
			// DECCOLM - 80 characters per line (erases screen)
			terminal.SetSize(80, uint(lines))
		}
		terminal.Clear()
		/*
			case "?4":
				// DECSCLM
				// @todo smooth scrolling / jump scrolling
		*/
	case "?5": // DECSCNM
		terminal.SetScreenMode(enabled)
	case "?6":
		// DECOM
		terminal.SetOriginMode(enabled)
	case "?7":
		// auto-wrap mode
		//DECAWM
		terminal.SetAutoWrap(enabled)
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
	case "?1000", "?10061000": // ?10061000 seen from htop
		// enable mouse tracking
		// 1000 refers to ext mode for extended mouse click area - otherwise only x <= 255-31
		if enabled {
			terminal.logger.Infof("Turning on VT200 mouse mode")
			terminal.SetMouseMode(MouseModeVT200)
		} else {
			terminal.logger.Infof("Turning off VT200 mouse mode")
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?1002":
		if enabled {
			terminal.logger.Infof("Turning on Button Event mouse mode")
			terminal.SetMouseMode(MouseModeButtonEvent)
		} else {
			terminal.logger.Infof("Turning off Button Event mouse mode")
			terminal.SetMouseMode(MouseModeNone)
		}
	case "?1003":
		return errors.New("Any Event mouse mode is not supported")
		/*
			if enabled {
				terminal.logger.Infof("Turning on Any Event mouse mode")
				terminal.SetMouseMode(MouseModeAnyEvent)
			} else {
				terminal.logger.Infof("Turning off Any Event mouse mode")
				terminal.SetMouseMode(MouseModeNone)
			}
		*/
	case "?1005":
		return errors.New("UTF-8 ext mouse mode is not supported")
		/*
			if enabled {
				terminal.logger.Infof("Turning on UTF-8 ext mouse mode")
				terminal.SetMouseExtMode(MouseExtUTF)
			} else {
				terminal.logger.Infof("Turning off UTF-8 ext mouse mode")
				terminal.SetMouseExtMode(MouseExtNone)
			}
		*/
	case "?1006":
		if enabled {
			terminal.logger.Infof("Turning on SGR ext mouse mode")
			terminal.SetMouseExtMode(MouseExtSGR)
		} else {
			terminal.logger.Infof("Turning off SGR ext mouse mode")
			terminal.SetMouseExtMode(MouseExtNone)
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
	case "?2004":
		terminal.SetBracketedPasteMode(enabled)
	default:
		return fmt.Errorf("Unsupported CSI %s%s code", modeStr, recoverCodeFromEnabled(enabled))
	}

	return nil
}
