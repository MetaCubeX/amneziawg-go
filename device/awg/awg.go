package awg

import (
	"bytes"
	"sync"
	"sync/atomic"
)

type aSecCfgType struct {
	IsSet                      bool
	JunkPacketCount            int
	JunkPacketMinSize          int
	JunkPacketMaxSize          int
	InitHeaderJunkSize         int
	ResponseHeaderJunkSize     int
	CookieReplyHeaderJunkSize  int
	TransportHeaderJunkSize    int
	InitPacketMagicHeader      uint32
	ResponsePacketMagicHeader  uint32
	UnderloadPacketMagicHeader uint32
	TransportPacketMagicHeader uint32
	// InitPacketMagicHeader      Limit
	// ResponsePacketMagicHeader  Limit
	// UnderloadPacketMagicHeader Limit
	// TransportPacketMagicHeader Limit
}

type Protocol struct {
	IsASecOn atomic.Bool
	// TODO: revision the need of the mutex
	ASecMux     sync.RWMutex
	ASecCfg     aSecCfgType
	JunkCreator junkCreator

	HandshakeHandler SpecialHandshakeHandler
}

func (protocol *Protocol) CreateInitHeaderJunk() ([]byte, error) {
	return protocol.createHeaderJunk(protocol.ASecCfg.InitHeaderJunkSize)
}

func (protocol *Protocol) CreateResponseHeaderJunk() ([]byte, error) {
	return protocol.createHeaderJunk(protocol.ASecCfg.ResponseHeaderJunkSize)
}

func (protocol *Protocol) CreateCookieReplyHeaderJunk() ([]byte, error) {
	return protocol.createHeaderJunk(protocol.ASecCfg.CookieReplyHeaderJunkSize)
}

func (protocol *Protocol) CreateTransportHeaderJunk(packetSize int) ([]byte, error) {
	return protocol.createHeaderJunk(protocol.ASecCfg.TransportHeaderJunkSize, packetSize)
}

func (protocol *Protocol) createHeaderJunk(junkSize int, optExtraSize ...int) ([]byte, error) {
	extraSize := 0
	if len(optExtraSize) == 1 {
		extraSize = optExtraSize[0]
	}

	var junk []byte
	protocol.ASecMux.RLock()
	if junkSize != 0 {
		buf := make([]byte, 0, junkSize+extraSize)
		writer := bytes.NewBuffer(buf[:0])
		err := protocol.JunkCreator.AppendJunk(writer, junkSize)
		if err != nil {
			protocol.ASecMux.RUnlock()
			return nil, err
		}
		junk = writer.Bytes()
	}
	protocol.ASecMux.RUnlock()

	return junk, nil
}
