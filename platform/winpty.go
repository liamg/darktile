// +build windows,cgo

package platform

import (
	"errors"
	"syscall"
	"time"

	"fmt"

	"github.com/MaxRis/w32"
)

// #include "windows.h"
//
//  /*    Until we can specify the platform SDK and target version for Windows.h      *
//   * without breaking our ability to gracefully display an error message, these     *
//   * definitions will be copied from the platform SDK headers and made to work.     *
//   */
//
//  typedef HRESULT (* CreatePseudoConsoleProcType)( COORD, HANDLE, HANDLE, DWORD, uintptr_t * );
//  typedef HRESULT (* ResizePseudoConsoleProcType)( uintptr_t, COORD );
//  typedef HRESULT (* ClosePseudoConsoleProcType)( uintptr_t );
//
//
//  CreatePseudoConsoleProcType pfnCreatePseudoConsole = NULL;
//  ResizePseudoConsoleProcType pfnResizePseudoConsole = NULL;
//  ClosePseudoConsoleProcType pfnClosePseudoConsole = NULL;
//
//  HMODULE hLibKernel32_Kern = NULL;
//
//  DWORD initPtyKernFuncs()
//  {
//      hLibKernel32_Kern = LoadLibrary( "kernel32.dll" );
//      if( hLibKernel32_Kern == NULL )
//      {
//          return -1;
//      }
//
//      pfnCreatePseudoConsole = (CreatePseudoConsoleProcType) GetProcAddress(hLibKernel32_Kern, "CreatePseudoConsole" );
//      if( pfnCreatePseudoConsole == NULL )
//      {
//          return -1;
//      }
//
//      pfnResizePseudoConsole = (ResizePseudoConsoleProcType) GetProcAddress(hLibKernel32_Kern, "ResizePseudoConsole" );
//      if( pfnResizePseudoConsole == NULL )
//      {
//          return -1;
//      }
//
//      pfnClosePseudoConsole = (ClosePseudoConsoleProcType) GetProcAddress(hLibKernel32_Kern, "ClosePseudoConsole" );
//      if( pfnClosePseudoConsole == NULL )
//      {
//          return -1;
//      }
//
//      return 0;
//  }
//
//  DWORD createPtyHelper( int xSize, int ySize, HANDLE input, HANDLE output, DWORD flags, uintptr_t * phPC )
//  {
//      COORD size;
//      size.X = xSize;
//      size.Y = ySize;
//      return (DWORD) (*pfnCreatePseudoConsole)( size, input, output, flags, phPC );
//  }
//
//  DWORD resizePtyHelper( uintptr_t hpc, int xSize, int ySize )
//  {
//      COORD size;
//      size.X = xSize;
//      size.Y = ySize;
//      return (DWORD) (*pfnResizePseudoConsole)( hpc, size );
//  }
//
//  DWORD closePtyHelper( uintptr_t hpc )
//  {
//      return (DWORD) (*pfnClosePseudoConsole)( hpc );
//  }
//
//  int hr_succeeded( DWORD hResult )
//  {
//      return SUCCEEDED( hResult );
//  }
import "C"

var ptyInitSucceeded = false

func init() {
	ret := int(C.initPtyKernFuncs())
	ptyInitSucceeded = (ret == 0)
}

type winConPty struct {
	inPipe                    syscall.Handle
	outPipe                   syscall.Handle
	innerInPipe               syscall.Handle
	innerOutPipe              syscall.Handle
	hcon                      uintptr
	processID                 uint32 // required to obtain old-style console handle (for standard console functions)
	platformDependentSettings PlatformDependentSettings
}

func (pty *winConPty) Read(p []byte) (n int, err error) {
	return syscall.Read(pty.inPipe, p)
}

func (pty *winConPty) Write(p []byte) (n int, err error) {
	return syscall.Write(pty.outPipe, p)
}

func (pty *winConPty) Close() error {
	C.closePtyHelper(C.uintptr_t(pty.hcon))

	err := syscall.CloseHandle(pty.inPipe)
	if err != nil {
		return err
	}
	err = syscall.CloseHandle(pty.outPipe)
	if err != nil {
		return err
	}
	err = syscall.CloseHandle(pty.innerInPipe)
	if err != nil {
		return err
	}
	err = syscall.CloseHandle(pty.innerOutPipe)
	if err != nil {
		return err
	}

	pty.hcon = 0

	return nil
}

func (pty *winConPty) CreateGuestProcess(imagePath string) (Process, error) {
	process, err := createPtyChildProcess(imagePath, pty.hcon)
	if err != nil {
		return nil, err
	}

	pty.processID = process.processID

	err = setupChildConsole(C.DWORD(process.processID), C.STD_OUTPUT_HANDLE, C.ENABLE_PROCESSED_OUTPUT|C.ENABLE_WRAP_AT_EOL_OUTPUT)
	if err != nil {
		process.Close()
		return nil, err
	}

	return process, err
}

func setupChildConsole(processID C.DWORD, nStdHandle C.DWORD, mode uint) error {
	C.FreeConsole()
	defer C.AttachConsole(^C.DWORD(0)) // attach to parent process console

	// process may not be ready so we'll do retries
	const maxWaitMilliSeconds = 5000
	const waitStepMilliSeconds = 200
	count := maxWaitMilliSeconds / waitStepMilliSeconds

	for {
		if r := C.AttachConsole(processID); r != 0 {
			break // success
		}
		lastError := C.GetLastError()
		if lastError != C.ERROR_GEN_FAILURE || count <= 0 {
			return fmt.Errorf("Was not able to attach to the child prosess' console")
		}

		time.Sleep(time.Millisecond * time.Duration(waitStepMilliSeconds))
		count--
	}

	h := C.GetStdHandle(nStdHandle)
	C.SetConsoleMode(h, C.DWORD(mode))
	C.FreeConsole()

	return nil
}

func (pty *winConPty) Resize(x, y int) error {
	cret := C.resizePtyHelper(C.uintptr_t(pty.hcon), C.int(x), C.int(y))

	if int(C.hr_succeeded(cret)) == 0 {
		return errors.New("Failed to resize ConPTY")
	}

	return nil
}

func (pty *winConPty) GetPlatformDependentSettings() PlatformDependentSettings {
	return pty.platformDependentSettings
}

func (pty *winConPty) Clear() {
	C.FreeConsole()
	defer C.AttachConsole(^C.DWORD(0)) // attach to parent process console

	if C.AttachConsole(C.DWORD(pty.processID)) == 0 {
		return
	}
	defer C.FreeConsole()
	hConsole := C.GetStdHandle(C.STD_OUTPUT_HANDLE)

	var coordScreen C.COORD
	coordScreen.X = 0
	coordScreen.Y = 0
	var cCharsWritten C.DWORD
	var csbi C.CONSOLE_SCREEN_BUFFER_INFO

	// Get the number of character cells in the current buffer.
	if C.GetConsoleScreenBufferInfo(hConsole, &csbi) != C.BOOL(0) {

		dwConSize := C.DWORD(csbi.dwSize.X * csbi.dwSize.Y)

		// Fill the entire screen with blanks.
		C.FillConsoleOutputCharacterA(hConsole, // Handle to console screen buffer
			' ',            // Character to write to the buffer
			dwConSize,      // Number of cells to write
			coordScreen,    // Coordinates of first cell
			&cCharsWritten) // Receive number of characters written

		// Get the current text attribute.
		if C.GetConsoleScreenBufferInfo(hConsole, &csbi) != C.BOOL(0) {
			// Set the buffer's attributes accordingly.

			C.FillConsoleOutputAttribute(hConsole, // Handle to console screen buffer
				csbi.wAttributes, // Character attributes to use
				dwConSize,        // Number of cells to set attribute
				coordScreen,      // Coordinates of first cell
				&cCharsWritten)   // Receive number of characters written
		}
	}

	// Put the cursor at its home coordinates.
	C.SetConsoleCursorPosition(hConsole, coordScreen)
}

// NewPty creates a new instance of a Pty implementation for Windows on a newly allocated ConPTY
func NewPty(x, y int) (pty Pty, err error) {
	if !ptyInitSucceeded {
		w32.MessageBox(0, "Aminal requires APIs that are only available on Windows 10 1809 (October 2018 Update) or above. Please upgrade", "Aminal", 0)
		return nil, errors.New("Windows PseudoConsole API unavailable on this version of Windows")
	}
	pty = nil

	var inputReadSide, inputWriteSide syscall.Handle
	var outputReadSide, outputWriteSide syscall.Handle

	err = syscall.CreatePipe(&inputReadSide, &inputWriteSide, nil, 0)
	if err != nil {
		return
	}

	err = syscall.CreatePipe(&outputReadSide, &outputWriteSide, nil, 0)
	if err != nil {
		return
	}

	var hc C.uintptr_t

	cret := C.createPtyHelper(C.int(x), C.int(y), C.HANDLE(inputReadSide), C.HANDLE(outputWriteSide), 0, &hc)
	ret := int(cret)

	if ret != 0 {
		return nil, errors.New("Failed to allocate a ConPTY instance")
	}

	pty = &winConPty{
		inPipe:       outputReadSide,
		outPipe:      inputWriteSide,
		innerInPipe:  inputReadSide,
		innerOutPipe: outputWriteSide,
		hcon:         uintptr(hc),
		platformDependentSettings: PlatformDependentSettings{
			OSCTerminators: map[rune]struct{}{0x00: {}, 0x07: {}},
		},
	}

	return pty, nil
}
