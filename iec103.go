package iec103

import (
	"fmt"
)

// proto address limit
const (
	AddressBroadCast = 0
	AddressMin       = 1
	addressMax       = 247
)

// AddressMax proto address max limit
// you can change with SetSpecialAddressMax,
// when your device have address upon addressMax
var AddressMax byte = addressMax

const (
	pduMinSize = 1   // funcCode(1)
	pduMaxSize = 253 // funcCode(1) + data(252)

	rtuAduMinSize = 4   // address(1) + funcCode(1) + crc(2)
	rtuAduMaxSize = 256 // address(1) + PDU(253) + crc(2)

	asciiAduMinSize       = 3
	asciiAduMaxSize       = 256
	asciiCharacterMaxSize = 513

	tcpProtocolIdentifier = 0x0000

	tcpHeaderMbapSize = 7 // MBAP header
	tcpAduMinSize     = 8 // MBAP + funcCode
	tcpAduMaxSize     = 260
)

// proto register limit
const (
	// Bits
	ReadBitsQuantityMin  = 1    // 0x0001
	ReadBitsQuantityMax  = 2000 // 0x07d0
	WriteBitsQuantityMin = 1    // 1
	WriteBitsQuantityMax = 1968 // 0x07b0
	// 16 Bits
	ReadRegQuantityMin             = 1   // 1
	ReadRegQuantityMax             = 125 // 0x007d
	WriteRegQuantityMin            = 1   // 1
	WriteRegQuantityMax            = 123 // 0x007b
	ReadWriteOnReadRegQuantityMin  = 1   // 1
	ReadWriteOnReadRegQuantityMax  = 125 // 0x007d
	ReadWriteOnWriteRegQuantityMin = 1   // 1
	ReadWriteOnWriteRegQuantityMax = 121 // 0x0079
)

// Function Code
const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadFIFOQueue              = 24
	FuncCodeOtherReportSlaveID         = 17
	// FuncCodeDiagReadException          = 7
	// FuncCodeDiagDiagnostic             = 8
	// FuncCodeDiagGetComEventCnt         = 11
	// FuncCodeDiagGetComEventLog         = 12
)

// Exception Code
const (
	ExceptionCodeIllegalFunction                    = 1
	ExceptionCodeIllegalDataAddress                 = 2
	ExceptionCodeIllegalDataValue                   = 3
	ExceptionCodeServerDeviceFailure                = 4
	ExceptionCodeAcknowledge                        = 5
	ExceptionCodeServerDeviceBusy                   = 6
	ExceptionCodeNegativeAcknowledge                = 7
	ExceptionCodeMemoryParityError                  = 8
	ExceptionCodeGatewayPathUnavailable             = 10
	ExceptionCodeGatewayTargetDeviceFailedToRespond = 11
)

// ExceptionError implements error interface.
type ExceptionError struct {
	ExceptionCode byte
}

func (e *ExceptionError) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodeIllegalFunction:
		name = "illegal function"
	case ExceptionCodeIllegalDataAddress:
		name = "illegal data address"
	case ExceptionCodeIllegalDataValue:
		name = "illegal data value"
	case ExceptionCodeServerDeviceFailure:
		name = "server device failure"
	case ExceptionCodeAcknowledge:
		name = "acknowledge"
	case ExceptionCodeServerDeviceBusy:
		name = "server device busy"
	case ExceptionCodeNegativeAcknowledge:
		name = "Negative Acknowledge"
	case ExceptionCodeMemoryParityError:
		name = "memory parity error"
	case ExceptionCodeGatewayPathUnavailable:
		name = "gateway path unavailable"
	case ExceptionCodeGatewayTargetDeviceFailedToRespond:
		name = "gateway target device failed to respond"
	default:
		name = "unknown"
	}
	return fmt.Sprintf("ieccon: exception '%v' (%s)", e.ExceptionCode, name)
}

// ProtocolDataUnit (PDU) is independent of underlying communication layers.
type ProtocolDataUnit struct {
	FuncCode byte
	Data     []byte
}

// protocolFrame 帧结构用于底层对象缓冲池
type protocolFrame struct {
	adu []byte
}

// ClientProvider is the interface implements underlying methods.
type ClientProvider interface {
	// Connect try to connect the remote server
	Connect() error
	// IsConnected returns a bool signifying whether
	// the client is connected or not.
	IsConnected() bool
	// SetAutoReconnect set auto reconnect count
	// if cnt == 0, disable auto reconnect
	// if cnt > 0 ,enable auto reconnect,but max 6
	SetAutoReconnect(cnt byte)
	// LogMode set enable or diable log output when you has set logger
	LogMode(enable bool)
	// SetLogProvider set logger provider
	SetLogProvider(p LogProvider)
	// Close disconnect the remote server
	Close() error
	// Send request to the remote server,it implements on SendRawFrame
	Send(slaveID byte, request ProtocolDataUnit) (ProtocolDataUnit, error)
	// SendPdu send pdu request to the remote server
	SendPdu(slaveID byte, pduRequest []byte) (pduResponse []byte, err error)
	// SendRawFrame send raw frame to the remote server
	SendRawFrame(request string) (response string, err error)
}

// LogProvider RFC5424 log message levels only Debug and Error
type LogProvider interface {
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}
