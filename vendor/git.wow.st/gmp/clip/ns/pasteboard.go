package ns

/*
#cgo CFLAGS: -x objective-c -fno-objc-arc
#cgo LDFLAGS: -framework AppKit -framework Foundation
#pragma clang diagnostic ignored "-Wformat-security"

#import <Foundation/Foundation.h>
#import <AppKit/NSPasteboard.h>

void
NSObject_inst_Release(void* o) {
        @autoreleasepool {
                [(NSObject*)o release];
        }
}

void
NSString_inst_Release(void* o) {
        @autoreleasepool {
                [(NSString*)o release];
        }
}

const void* _Nullable
NSString_inst_UTF8String(void* o) {
        const char* _Nullable ret;
        @autoreleasepool {
                ret = strdup([(NSString*)o UTF8String]);
        }
        return ret;

}

void
NSPasteboard_inst_Release(void* o) {
        @autoreleasepool {
                [(NSPasteboard*)o release];
        }
}

void* _Nullable
NSString_StringWithUTF8String(void* nullTerminatedCString) {
        NSString* _Nullable ret;
        @autoreleasepool {
                ret = [NSString stringWithUTF8String:nullTerminatedCString];
                if(ret != nil) { [ret retain]; }
        }
        return ret;

}

void* _Nonnull
NSPasteboard_GeneralPasteboard() {
        NSPasteboard* _Nonnull ret;
        @autoreleasepool {
                ret = [NSPasteboard generalPasteboard];
        }
        return ret;
}

void
NSPasteboard_inst_ClearContents(void* o) {
        @autoreleasepool {
                [(NSPasteboard*)o clearContents];
        }
}

BOOL
NSPasteboard_inst_SetString(void* o, void* string) {
        BOOL ret;
        @autoreleasepool {
                ret = [(NSPasteboard*)o setString:string forType:NSPasteboardTypeString];
        }
        return ret;
}

void* _Nullable
NSPasteboard_inst_GetString(void* o) {
        NSString* _Nullable ret;
        @autoreleasepool {
                ret = [(NSPasteboard*)o stringForType:NSPasteboardTypeString];
                if (ret != nil && ret != o) { [ret retain]; }
        }
        return ret;

}

*/
import "C"

import (
	"unsafe"
	"runtime"
)

type Id struct {
        ptr unsafe.Pointer
}
func (o *Id) Ptr() unsafe.Pointer { if o == nil { return nil }; return o.ptr }

type NSObject interface {
        Ptr() unsafe.Pointer
}

func (o *Id) Release()  {
        C.NSObject_inst_Release(o.Ptr())
        runtime.KeepAlive(o)
}

func (o *NSPasteboard) Release()  {
        C.NSPasteboard_inst_Release(o.Ptr())
        runtime.KeepAlive(o)
}

func (o *NSString) Release()  {
        C.NSString_inst_Release(o.Ptr())
        runtime.KeepAlive(o)
}

func (c *Char) Free() {
        C.free(unsafe.Pointer(c))
}

type BOOL C.uchar

type NSString struct { Id }
func (o *NSString) Ptr() unsafe.Pointer { if o == nil { return nil }; return o.ptr }
func (o *Id) NSString() *NSString {
        return (*NSString)(unsafe.Pointer(o))
}

func (o *NSString) UTF8String() *Char {
        ret := (*Char)(unsafe.Pointer(C.NSString_inst_UTF8String(o.Ptr())))
        runtime.KeepAlive(o)
        return ret
}

func (o *NSString) String() string {
        utf8 := o.UTF8String()
        ret := utf8.String()
        utf8.Free()
        runtime.KeepAlive(o)
        return ret
}

type NSPasteboard struct { Id }
func (o *NSPasteboard) Ptr() unsafe.Pointer { if o == nil { return nil }; return o.ptr }
func (o *Id) NSPasteboard() *NSPasteboard {
        return (*NSPasteboard)(unsafe.Pointer(o))
}

type Char C.char

func CharWithGoString(s string) *Char {
        return (*Char)(unsafe.Pointer(C.CString(s)))
}

func (c *Char) String() string {
        return C.GoString((*C.char)(c))
}

func NSStringWithUTF8String(nullTerminatedCString *Char) *NSString {
        ret := &NSString{}
        ret.ptr = unsafe.Pointer(C.NSString_StringWithUTF8String(unsafe.Pointer(nullTerminatedCString)))
        if ret.ptr == nil { return ret }
        runtime.SetFinalizer(ret, func(o *NSString) {
                o.Release()
        })
        return ret
}

func NSStringWithGoString(string string) *NSString {
        string_chr := CharWithGoString(string)
        defer string_chr.Free()
        ret := NSStringWithUTF8String(string_chr)
        return ret
}

func NSPasteboardGeneralPasteboard() *NSPasteboard {
        ret := &NSPasteboard{}
        ret.ptr = unsafe.Pointer(C.NSPasteboard_GeneralPasteboard())
        if ret.ptr == nil { return ret }
        return ret
}

func (o *NSPasteboard) ClearContents() {
        C.NSPasteboard_inst_ClearContents(o.Ptr())
        runtime.KeepAlive(o)
}

func (o *NSPasteboard) SetString(s string) bool {
	string := NSStringWithGoString(s)
        ret := (C.NSPasteboard_inst_SetString(o.Ptr(), string.Ptr())) != 0
        runtime.KeepAlive(o)
	runtime.KeepAlive(string)
        return ret
}

func (o *NSPasteboard) GetString() *NSString {
        ret := &NSString{}
        ret.ptr = unsafe.Pointer(C.NSPasteboard_inst_GetString(o.Ptr()))
        if ret.ptr == nil { runtime.KeepAlive(o); return ret }
        if ret.ptr == o.ptr { runtime.KeepAlive(o); return (*NSString)(unsafe.Pointer(o)) }
        runtime.SetFinalizer(ret, func(o *NSString) {
                o.Release()
        })
        runtime.KeepAlive(o)
        return ret
}

