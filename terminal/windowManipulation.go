package terminal

import (
	"fmt"
	"strconv"
)

func getOptionalIntegerParam(params []string, paramNo int, defValue int) (int, error) {
	result := defValue
	if len(params) >= paramNo+1 {
		var err error
		result, err = strconv.Atoi(params[paramNo])
		if err != nil {
			return defValue, err
		}
	}

	return result, nil
}

func getMandatoryIntegerParam(params []string, paramNo int) (int, error) {
	if len(params) < paramNo+1 {
		return 0, fmt.Errorf("no mandatory parameter")
	}

	result, err := strconv.Atoi(params[paramNo])
	if err != nil {
		return 0, err
	}

	return result, nil
}

func csiWindowManipulation(params []string, terminal *Terminal) error {
	if terminal.WindowManipulation == nil {
		return fmt.Errorf("Handler for CSI window manipulation commands is not set")
	}

	operation, err := getMandatoryIntegerParam(params, 0)
	if err != nil {
		return fmt.Errorf("CSI t ignored: %s", err.Error())
	}

	switch operation {
	case 1:
		terminal.logger.Debug("De-iconify window")
		return terminal.WindowManipulation.RestoreWindow(terminal)

	case 2:
		terminal.logger.Debug("Iconify window")
		return terminal.WindowManipulation.IconifyWindow(terminal)

	case 3:
		terminal.logger.Debug("Move window")
		{
			x, err := getMandatoryIntegerParam(params, 1)
			if err != nil {
				return err
			}
			y, err := getMandatoryIntegerParam(params, 2)
			if err != nil {
				return err
			}

			return terminal.WindowManipulation.MoveWindow(terminal, x, y)
		}

	case 4:
		terminal.logger.Debug("Resize the window in pixels")
		{
			height, err := getMandatoryIntegerParam(params, 1)
			if err != nil {
				return err
			}
			width, err := getMandatoryIntegerParam(params, 2)
			if err != nil {
				return err
			}

			return terminal.WindowManipulation.ResizeWindowByPixels(terminal, height, width)
		}

	case 5:
		terminal.logger.Debug("Raise the window to the front")
		return terminal.WindowManipulation.BringWindowToFront(terminal)

	case 6:
		return fmt.Errorf("Lowering the window to the bottom is not implemented")

	case 7:
		// NB: On Windows this sequence seem handled by the system
		return fmt.Errorf("Refreshing the window is not implemented")

	case 8:
		terminal.logger.Debug("Resize the text area in characters")
		{
			height, err := getMandatoryIntegerParam(params, 1)
			if err != nil {
				return err
			}
			width, err := getMandatoryIntegerParam(params, 2)
			if err != nil {
				return err
			}
			return terminal.WindowManipulation.ResizeWindowByChars(terminal, height, width)
		}

	case 9:
		{
			p, err := getMandatoryIntegerParam(params, 1)
			if err != nil {
				return err
			}
			if p == 0 {
				terminal.logger.Debug("Restore maximized window")
				return terminal.WindowManipulation.RestoreWindow(terminal)
			} else if p == 1 {
				terminal.logger.Debug("Maximize window")
				return terminal.WindowManipulation.MaximizeWindow(terminal)
			}
		}

	case 11:
		terminal.logger.Debug("Report the window state")
		return terminal.WindowManipulation.ReportWindowState(terminal)

	case 13:
		terminal.logger.Debug("Report the window position")
		return terminal.WindowManipulation.ReportWindowPosition(terminal)

	case 14:
		terminal.logger.Debug("Report the window size in pixels")
		return terminal.WindowManipulation.ReportWindowSizeInPixels(terminal)

	case 18:
		terminal.logger.Debug("Report the window size in characters as CSI 8")
		return terminal.WindowManipulation.ReportWindowSizeInChars(terminal)

	case 19:
		return fmt.Errorf("Reporting the screen size in characters is not implemented")

	case 20:
		return fmt.Errorf("Reporting the window icon label is not implemented")

	case 21:
		return fmt.Errorf("Reporting the window title is not implemented")

	default:
		return fmt.Errorf("not supported CSI t")
	}

	return nil
}
