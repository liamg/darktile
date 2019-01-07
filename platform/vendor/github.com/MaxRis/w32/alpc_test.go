package w32

import (
	"testing"
)

var testPortName = "\\TestAlpcPort"

var basicPortAttr = ALPC_PORT_ATTRIBUTES{
	MaxMessageLength: uint64(SHORT_MESSAGE_MAX_SIZE),
	SecurityQos: SECURITY_QUALITY_OF_SERVICE{
		Length:              SECURITY_QOS_SIZE,
		ContextTrackingMode: SECURITY_DYNAMIC_TRACKING,
		EffectiveOnly:       1,
		ImpersonationLevel:  SecurityAnonymous,
	},
	Flags:          ALPC_PORFLG_ALLOW_LPC_REQUESTS,
	DupObjectTypes: ALPC_SYNC_OBJECT_TYPE,
}

func ObjectAttributes(name string) (oa OBJECT_ATTRIBUTES, e error) {

	sd, e := InitializeSecurityDescriptor(1)
	if e != nil {
		return
	}

	e = SetSecurityDescriptorDacl(sd, nil)
	if e != nil {
		return
	}

	oa, e = InitializeObjectAttributes(name, 0, 0, sd)
	return
}

func Send(
	hPort HANDLE,
	msg *AlpcShortMessage,
	flags uint32,
	pMsgAttrs *ALPC_MESSAGE_ATTRIBUTES,
	timeout *int64,
) (e error) {

	e = NtAlpcSendWaitReceivePort(hPort, flags, msg, pMsgAttrs, nil, nil, nil, timeout)
	return

}

func Recv(
	hPort HANDLE,
	pMsg *AlpcShortMessage,
	pMsgAttrs *ALPC_MESSAGE_ATTRIBUTES,
	timeout *int64,
) (bufLen uint32, e error) {

	bufLen = uint32(pMsg.TotalLength)
	e = NtAlpcSendWaitReceivePort(hPort, 0, nil, nil, pMsg, &bufLen, pMsgAttrs, timeout)
	return

}

// Convenience method to create an ALPC port with a NULL DACL. Requires an
// absolute port name ( where / is the root of the kernel object directory )
func CreatePort(name string) (hPort HANDLE, e error) {

	oa, e := ObjectAttributes(name)
	if e != nil {
		return
	}

	hPort, e = NtAlpcCreatePort(&oa, &basicPortAttr)

	return
}

func ConnectPort(serverName, clientName string, pConnMsg *AlpcShortMessage) (hPort HANDLE, e error) {

	oa, e := InitializeObjectAttributes(clientName, 0, 0, nil)
	if e != nil {
		return
	}

	hPort, e = NtAlpcConnectPort(
		serverName,
		&oa,
		&basicPortAttr,
		ALPC_PORFLG_ALLOW_LPC_REQUESTS,
		nil,
		pConnMsg,
		nil,
		nil,
		nil,
		nil,
	)

	return
}

func Accept(
	hSrv HANDLE,
	context *AlpcPortContext,
	pConnReq *AlpcShortMessage,
	accept bool,
) (hPort HANDLE, e error) {

	oa, _ := InitializeObjectAttributes("", 0, 0, nil)

	var accepted uintptr
	if accept {
		accepted++
	}

	hPort, e = NtAlpcAcceptConnectPort(
		hSrv,
		0,
		&oa,
		&basicPortAttr,
		context,
		pConnReq,
		nil,
		accepted,
	)

	return
}

func TestNtAlpcCreatePort(t *testing.T) {

	hPort, err := CreatePort(testPortName)

	if err != nil {
		t.Errorf("failed to create ALPC port %v: %v", testPortName, err)
	} else {
		t.Logf("[OK] Created ALPC port %v with handle 0x%x", testPortName, hPort)
	}
}
