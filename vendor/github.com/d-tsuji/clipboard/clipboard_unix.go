// The MIT License (MIT)
// Copyright (c) 2016 Alessandro Arzilli
// https://github.com/aarzilli/nucular/blob/master/LICENSE

// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"golang.org/x/xerrors"
)

const debugClipboardRequests = false

var (
	x             *xgb.Conn
	win           xproto.Window
	clipboardText string
	selnotify     chan bool

	clipboardAtom, primaryAtom, textAtom, targetsAtom, atomAtom xproto.Atom
	targetAtoms                                                 []xproto.Atom
	clipboardAtomCache                                          = map[xproto.Atom]string{}

	doneCh = make(chan interface{}, 1)
)

func start() error {
	var err error
	xServer := os.Getenv("DISPLAY")
	if xServer == "" {
		return xerrors.New("could not identify xserver")
	}
	x, err = xgb.NewConnDisplay(xServer)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	selnotify = make(chan bool, 1)

	win, err = xproto.NewWindowId(x)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	setup := xproto.Setup(x)
	s := setup.DefaultScreen(x)
	err = xproto.CreateWindowChecked(x, s.RootDepth, win, s.Root, 100, 100, 1, 1, 0, xproto.WindowClassInputOutput, s.RootVisual, 0, []uint32{}).Check()
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	clipboardAtom = internAtom(x, "CLIPBOARD")
	primaryAtom = internAtom(x, "PRIMARY")
	textAtom = internAtom(x, "UTF8_STRING")
	targetsAtom = internAtom(x, "TARGETS")
	atomAtom = internAtom(x, "ATOM")

	targetAtoms = []xproto.Atom{targetsAtom, textAtom}

	go eventLoop()

	return nil
}

func set(text string) error {
	if err := start(); err != nil {
		return xerrors.Errorf("init clipboard: %w", err)
	}
	clipboardText = text
	ssoc := xproto.SetSelectionOwnerChecked(x, win, clipboardAtom, xproto.TimeCurrentTime)
	if err := ssoc.Check(); err != nil {
		return xerrors.Errorf("setting clipboard: %w", err)
	}
	return nil
}

func get() (string, error) {
	if err := start(); err != nil {
		return "", xerrors.Errorf("init clipboard: %w", err)
	}
	return getSelection(clipboardAtom)
}

func getSelection(selAtom xproto.Atom) (string, error) {
	csc := xproto.ConvertSelectionChecked(x, win, selAtom, textAtom, selAtom, xproto.TimeCurrentTime)
	err := csc.Check()
	if err != nil {
		return "", xerrors.Errorf("convert selection check: %w", err)
	}

	select {
	case r := <-selnotify:
		if !r {
			return "", nil
		}
		gpc := xproto.GetProperty(x, true, win, selAtom, textAtom, 0, 5*1024*1024)
		gpr, err := gpc.Reply()
		if err != nil {
			return "", xerrors.Errorf("grp reply: %w", err)
		}
		if gpr.BytesAfter != 0 {
			return "", xerrors.New("clipboard too large")
		}
		return string(gpr.Value[:gpr.ValueLen]), nil
	case <-time.After(1 * time.Second):
		return "", xerrors.New("clipboard retrieval failed, timeout")
	}
}

func pollForEvent(X *xgb.Conn, events chan<- xgb.Event) {
	for {
		select {
		case <-doneCh:
			return
		default:
			ev, err := X.PollForEvent()
			if err != nil {
				fmt.Println("wait for event:", err)
			}
			events <- ev
		}
	}
}

func eventLoop() {
	eventCh := make(chan xgb.Event, 1)
	go pollForEvent(x, eventCh)
	for {
		select {
		case event := <-eventCh:
			switch e := event.(type) {
			case xproto.SelectionRequestEvent:
				if debugClipboardRequests {
					tgtname := lookupAtom(e.Target)
					propname := lookupAtom(e.Property)
					fmt.Println("SelectionRequest", e, textAtom, tgtname, propname, "isPrimary:", e.Selection == primaryAtom, "isClipboard:", e.Selection == clipboardAtom)
				}
				t := clipboardText

				switch e.Target {
				case textAtom:
					if debugClipboardRequests {
						fmt.Println("Sending as text")
					}
					cpc := xproto.ChangePropertyChecked(x, xproto.PropModeReplace, e.Requestor, e.Property, textAtom, 8, uint32(len(t)), []byte(t))
					err := cpc.Check()
					if err == nil {
						sendSelectionNotify(e)
					} else {
						fmt.Println(err)
					}

				case targetsAtom:
					if debugClipboardRequests {
						fmt.Println("Sending targets")
					}
					buf := make([]byte, len(targetAtoms)*4)
					for i, atom := range targetAtoms {
						xgb.Put32(buf[i*4:], uint32(atom))
					}

					err := xproto.ChangePropertyChecked(x, xproto.PropModeReplace, e.Requestor, e.Property, atomAtom, 32, uint32(len(targetAtoms)), buf).Check()
					if err == nil {
						sendSelectionNotify(e)
					} else {
						fmt.Println(err)
					}

				default:
					if debugClipboardRequests {
						fmt.Println("Skipping")
					}
					e.Property = 0
					sendSelectionNotify(e)
				}

			case xproto.SelectionNotifyEvent:
				selnotify <- (e.Property == clipboardAtom) || (e.Property == primaryAtom)
			}
		case <-doneCh:
			return
		}
	}
}

func lookupAtom(at xproto.Atom) string {
	if s, ok := clipboardAtomCache[at]; ok {
		return s
	}

	reply, err := xproto.GetAtomName(x, at).Reply()
	if err != nil {
		panic(err)
	}

	// If we're here, it means we didn't have ths ATOM id cached. So cache it.
	atomName := string(reply.Name)
	clipboardAtomCache[at] = atomName
	return atomName
}

func sendSelectionNotify(e xproto.SelectionRequestEvent) {
	sn := xproto.SelectionNotifyEvent{
		Time:      e.Time,
		Requestor: e.Requestor,
		Selection: e.Selection,
		Target:    e.Target,
		Property:  e.Property}
	sec := xproto.SendEventChecked(x, false, e.Requestor, 0, string(sn.Bytes()))
	err := sec.Check()
	if err != nil {
		fmt.Println(err)
	}
}

func internAtom(conn *xgb.Conn, n string) xproto.Atom {
	iac := xproto.InternAtom(conn, true, uint16(len(n)), n)
	iar, err := iac.Reply()
	if err != nil {
		panic(err)
	}
	return iar.Atom
}
