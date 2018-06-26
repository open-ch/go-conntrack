//+build linux

package conntrack

import (
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

// Nfct represents a conntrack handler
type Nfct struct {
	con *netlink.Conn
}

// Conn contains all the information of a connection
type Conn struct {
}

// CtType specifies the subsystem of conntrack
type CtType int

// Supported conntrack subsystems
const (
	Ct         CtType = unix.NFNL_SUBSYS_CTNETLINK
	CtExpected CtType = unix.NFNL_SUBSYS_CTNETLINK_EXP
	CtTimeout  CtType = unix.NFNL_SUBSYS_CTNETLINK_TIMEOUT
)

const (
	ipctnlMsgCtNew            = iota
	ipctnlMsgCtGet            = iota
	ipctnlMsgCtDelete         = iota
	ipctnlMsgCtGetCtrZero     = iota
	ipctnlMsgCtGetStatsCPU    = iota
	ipctnlMsgCtGetStats       = iota
	ipctnlMsgCtGetDying       = iota
	ipctnlMsgCtGetUnconfirmed = iota

	ipctnlMsgMax = iota
)

// CtQuery specifies the type of the query
type CtQuery int

// Supported query types
const (
	CtQCreate          CtQuery = iota
	CtQUpdate          CtQuery = iota
	CtQDestroy         CtQuery = iota
	CtQGet             CtQuery = iota
	CtQFlush           CtQuery = iota
	CtQDump            CtQuery = iota
	CtQDumpReset       CtQuery = iota
	CtQCreateUpdate    CtQuery = iota
	CtQDumpFilter      CtQuery = iota
	CtQDumpFilterReset CtQuery = iota
)

// CtFamily specifies the network family
type CtFamily uint8

// Supported family types
const (
	CtIPv6 CtFamily = unix.AF_INET6
	CtIPv4 CtFamily = unix.AF_INET
)

// Open a connection to the given conntrack subsystem
func Open() (*Nfct, error) {
	var nfct Nfct

	con, err := netlink.Dial(unix.NETLINK_NETFILTER, nil)
	if err != nil {
		return nil, err
	}
	nfct.con = con

	return &nfct, nil
}

// Close the connection to the conntrack subsystem
func (nfct *Nfct) Close() error {
	return nfct.con.Close()
}

// Flush a conntrack subsystem
func (nfct *Nfct) Flush(f CtFamily) error {
	data := putExtraHeader(uint8(f), unix.NFNETLINK_V0, 0)
	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType((Ct << 8) | ipctnlMsgCtDelete),
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
		},
		Data: data,
	}
	return nfct.execute(req)
}

// Dump a conntrack subsystem
func (nfct *Nfct) Dump(f CtFamily) ([]*Conn, error) {
	data := putExtraHeader(uint8(f), unix.NFNETLINK_V0, 0)
	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType((Ct << 8) | ipctnlMsgCtGet),
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
		},
		Data: data,
	}

	verify, err := nfct.con.Send(req)
	if err != nil {
		return nil, err
	}
	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	reply, err := nfct.con.Receive()
	if err != nil {
		return nil, err
	}

	var conn []*Conn
	for _, msg := range reply {
		c, err := parseConnectionMsg(msg.Data[4:])
		if err != nil {
			return nil, err
		}
		conn = append(conn, c)
	}

	return conn, nil
}

// ErrMsg as defined in nlmsgerr
type ErrMsg struct {
	Code  int
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

func unmarschalErrMsg(b []byte) (ErrMsg, error) {
	var msg ErrMsg

	msg.Code = int(nlenc.Uint32(b[0:4]))
	msg.Len = nlenc.Uint32(b[4:8])
	msg.Type = nlenc.Uint16(b[8:10])
	msg.Flags = nlenc.Uint16(b[10:12])
	msg.Seq = nlenc.Uint32(b[12:16])
	msg.Pid = nlenc.Uint32(b[16:20])

	return msg, nil
}

func (nfct *Nfct) execute(req netlink.Message) error {
	reply, e := nfct.con.Execute(req)
	if e != nil {
		return e
	}
	if e := netlink.Validate(req, reply); e != nil {
		return e
	}
	for _, msg := range reply {
		errMsg, err := unmarschalErrMsg(msg.Data)
		if err != nil {
			return err
		}
		if errMsg.Code != 0 {
			return fmt.Errorf("%#v", errMsg)
		}
	}
	return nil
}

func putExtraHeader(familiy, version uint8, resid uint16) []byte {
	buf := make([]byte, 2)
	nlenc.PutUint16(buf, resid)
	return append([]byte{familiy, version}, buf...)
}

func parseConnectionMsg(b []byte) (*Conn, error) {
	return nil, fmt.Errorf("Not implemented yet")
}