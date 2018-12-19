// +build windows,cgo

package platform

// #include "Windows.h"
//
//  /*    Until we can specify the platform SDK and target version for Windows.h      *
//   * without breaking our ability to gracefully display an error message, these     *
//   * definitions will be copied from the platform SDK headers and made to work.     *
//   */
//
//  typedef BOOL (* InitializeProcThreadAttributeListProcType)(LPPROC_THREAD_ATTRIBUTE_LIST, DWORD, DWORD, PSIZE_T);
//  typedef BOOL (* UpdateProcThreadAttributeProcType)(
//            LPPROC_THREAD_ATTRIBUTE_LIST lpAttributeList,
//            DWORD                        dwFlags,
//            DWORD_PTR                    Attribute,
//            PVOID                        lpValue,
//            SIZE_T                       cbSize,
//            PVOID                        lpPreviousValue,
//            PSIZE_T                      lpReturnSize
//      );
//
//  InitializeProcThreadAttributeListProcType pfnInitializeProcThreadAttributeList = NULL;
//  UpdateProcThreadAttributeProcType pfnUpdateProcThreadAttribute = NULL;
//
//  #define ProcThreadAttributePseudoConsole_copy 22
//
//  #define PROC_THREAD_ATTRIBUTE_NUMBER_copy    0x0000FFFF
//  #define PROC_THREAD_ATTRIBUTE_THREAD_copy    0x00010000  // Attribute may be used with thread creation
//  #define PROC_THREAD_ATTRIBUTE_INPUT_copy     0x00020000  // Attribute is input only
//  #define PROC_THREAD_ATTRIBUTE_ADDITIVE_copy  0x00040000  // Attribute may be "accumulated," e.g. bitmasks, counters, etc.
//
//  #define ProcThreadAttributeValue_copy(Number, Thread, Input, Additive) \
//      (((Number) & PROC_THREAD_ATTRIBUTE_NUMBER_copy) | \
//      ((Thread != FALSE) ? PROC_THREAD_ATTRIBUTE_THREAD_copy : 0) | \
//      ((Input != FALSE) ? PROC_THREAD_ATTRIBUTE_INPUT_copy : 0) | \
//      ((Additive != FALSE) ? PROC_THREAD_ATTRIBUTE_ADDITIVE_copy : 0))
//
//  #define PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE_copy ProcThreadAttributeValue_copy (ProcThreadAttributePseudoConsole_copy, FALSE, TRUE, FALSE)
//
//  typedef struct _STARTUPINFOEXW_copy {
//      STARTUPINFOW                 StartupInfo;
//      LPPROC_THREAD_ATTRIBUTE_LIST lpAttributeList;
//  } STARTUPINFOEXW_copy, *LPSTARTUPINFOEXW_copy;
//
//  HMODULE hLibKernel32_Proc = NULL;
//
//  DWORD initProcKernFuncs()
//  {
//      hLibKernel32_Proc = LoadLibrary( "kernel32.dll" );
//      if( hLibKernel32_Proc == NULL )
//      {
//          return -1;
//      }
//
//      pfnInitializeProcThreadAttributeList = (InitializeProcThreadAttributeListProcType) GetProcAddress(hLibKernel32_Proc, "InitializeProcThreadAttributeList" );
//      if( pfnInitializeProcThreadAttributeList == NULL )
//      {
//          return -1;
//      }
//
//      pfnUpdateProcThreadAttribute = (UpdateProcThreadAttributeProcType) GetProcAddress(hLibKernel32_Proc, "UpdateProcThreadAttribute" );
//      if( pfnUpdateProcThreadAttribute == NULL )
//      {
//          return -1;
//      }
//
//      return 0;
//  }
//
//  DWORD createGuestProcHelper( uintptr_t hpc, LPCWSTR imagePath, uintptr_t * hProcess, DWORD * dwProcessID )
//  {
//      STARTUPINFOEXW_copy si;
//      ZeroMemory( &si, sizeof(si) );
//      si.StartupInfo.cb = sizeof(si);
//
//      size_t bytesRequired;
//      (*pfnInitializeProcThreadAttributeList)( NULL, 1, 0, &bytesRequired );
//
//      si.lpAttributeList = (PPROC_THREAD_ATTRIBUTE_LIST)HeapAlloc(GetProcessHeap(), 0, bytesRequired);
//      if( !si.lpAttributeList )
//      {
//          return E_OUTOFMEMORY;
//      }
//
//      if (!(*pfnInitializeProcThreadAttributeList)(si.lpAttributeList, 1, 0, &bytesRequired))
//      {
//          HeapFree(GetProcessHeap(), 0, si.lpAttributeList);
//          return HRESULT_FROM_WIN32(GetLastError());
//      }
//
//      if (!(*pfnUpdateProcThreadAttribute)(si.lpAttributeList,
//              0,
//              PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE_copy,
//              (PVOID) hpc,
//              sizeof(hpc),
//              NULL,
//              NULL))
//      {
//          HeapFree(GetProcessHeap(), 0, si.lpAttributeList);
//          return HRESULT_FROM_WIN32(GetLastError());
//      }
//
//      bytesRequired = (wcslen(imagePath) + 1) * sizeof(wchar_t); // +1 null terminator
//      PWSTR cmdLineMutable = (PWSTR)HeapAlloc(GetProcessHeap(), 0, bytesRequired);
//
//      if (!cmdLineMutable)
//      {
//          HeapFree(GetProcessHeap(), 0, si.lpAttributeList);
//          return E_OUTOFMEMORY;
//      }
//
//      wcscpy_s(cmdLineMutable, bytesRequired, imagePath);
//
//      PROCESS_INFORMATION pi;
//      ZeroMemory(&pi, sizeof(pi));
//
//      if (!CreateProcessW(NULL,
//              cmdLineMutable,
//              NULL,
//              NULL,
//              FALSE,
//              EXTENDED_STARTUPINFO_PRESENT,
//              NULL,
//              NULL,
//              &si.StartupInfo,
//              &pi))
//      {
//          HeapFree(GetProcessHeap(), 0, si.lpAttributeList);
//          HeapFree(GetProcessHeap(), 0, cmdLineMutable);
//          return HRESULT_FROM_WIN32(GetLastError());
//      }
//
//      *hProcess = (uintptr_t) pi.hProcess;
//      *dwProcessID = pi.dwProcessId;
//
//      HeapFree(GetProcessHeap(), 0, si.lpAttributeList);
//      HeapFree(GetProcessHeap(), 0, cmdLineMutable);
//      return S_OK;
//  }
//
//  int hr_succeeded( DWORD hResult );
import "C"
import (
	"errors"
	"syscall"
	"unicode/utf16"
)

var procsInitSucceeded = false

func init() {
	ret := int(C.initProcKernFuncs())
	procsInitSucceeded = (ret == 0)
}

type winProcess struct {
	hproc     uintptr
	processID uint32
}

func createPtyChildProcess(imagePath string, hcon uintptr) (*winProcess, error) {
	path16 := utf16.Encode([]rune(imagePath))

	cpath16 := C.calloc(C.size_t(len(path16)+1), 2)
	pp := (*[1 << 30]uint16)(cpath16)
	copy(pp[:], path16)

	hproc := C.uintptr_t(0)
	dwProcessID := C.DWORD(0)

	hr := C.createGuestProcHelper(C.uintptr_t(hcon), (C.LPCWSTR)(cpath16), &hproc, &dwProcessID)

	C.free(cpath16)

	if int(C.hr_succeeded(hr)) == 0 {
		return nil, errors.New("Failed to create process: " + imagePath)
	}

	return &winProcess{
		hproc:     uintptr(hproc),
		processID: uint32(dwProcessID),
	}, nil
}

func (process *winProcess) Close() error {
	return syscall.CloseHandle(syscall.Handle(process.hproc))
}
