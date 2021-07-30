package clip

import (
	"git.wow.st/gmp/clip/ns"
)

var pb *ns.NSPasteboard

func Clear() {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	pb.ClearContents()
}

func Set(x string) bool {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	pb.ClearContents()
	return pb.SetString(x)
}

func Get() string {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	ret := pb.GetString()
	if ret.Ptr() == nil {
		return ""
	} else {
		return ret.String()
	}
}
