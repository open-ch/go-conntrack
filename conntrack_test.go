//+build linux

package conntrack_test

import (
	"testing"

	ct "github.com/florianl/go-conntrack"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nltest"
)

func TestFlush(t *testing.T) {
	tests := []struct {
		name   string
		family ct.CtFamily
		want   []netlink.Message
	}{
		{name: "Flush IPv4", family: ct.CtIPv4, want: []netlink.Message{
			{
				Header: netlink.Header{
					Length: 20,
					// NFNL_SUBSYS_CTNETLINK<<8|IPCTNL_MSG_CT_DELETE
					Type: netlink.HeaderType(1<<8 | 2),
					// NLM_F_REQUEST|NLM_F_ACK
					Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
					// Can and will be ignored
					Sequence: 0,
					// Can and will be ignored
					PID: nltest.PID,
				},
				// nfgen_family=AF_INET, version=NFNETLINK_V0, res_id=htons(0)
				Data: []byte{0x2, 0x0, 0x0, 0x0},
			},
		},
		},
		{name: "Flush IPv6", family: ct.CtIPv6, want: []netlink.Message{
			{
				Header: netlink.Header{
					Length: 20,
					// NFNL_SUBSYS_CTNETLINK<<8|IPCTNL_MSG_CT_DELETE
					Type: netlink.HeaderType(1<<8 | 2),
					// NLM_F_REQUEST|NLM_F_ACK
					Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsAcknowledge,
					// Can and will be ignored
					Sequence: 0,
					// Can and will be ignored
					PID: nltest.PID,
				},
				// nfgen_family=AF_INET6, version=NFNETLINK_V0, res_id=htons(0)
				Data: []byte{0xA, 0x0, 0x0, 0x0},
			},
		},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Fake a netfilter conntrack connection
			nfct := &ct.Nfct{}
			nfct.Con = nltest.Dial(func(reqs []netlink.Message) ([]netlink.Message, error) {
				if len(reqs) == 0 {
					return nil, nil
				}
				if len(reqs) != 1 {
					t.Fatalf("Expected only one request, got %d", len(reqs))
				}

				// To ignore the Sequence number, we set it to the same value
				tc.want[0].Header.Sequence = reqs[0].Header.Sequence

				if len(reqs) != len(tc.want) {
					t.Fatalf("differen length:\n- want: %#v\n- got: %#v\n", tc.want, reqs)
				}

				for i := 0; i < len(reqs); i++ {
					if len(reqs[i].Data) != len(tc.want[i].Data) {
						t.Fatalf("differen length:\n- want: %#v\n- got: %#v\n", tc.want[i], reqs[i])
					}
					if reqs[i].Header.Type != tc.want[i].Header.Type {
						t.Fatalf("differen header types:\n- want: %#v\n- got: %#v\n", tc.want[i].Header.Type, reqs[i].Header.Type)
					}
					if reqs[i].Header.Flags != tc.want[i].Header.Flags {
						t.Fatalf("differen header flags:\n- want: %#v\n- got: %#v\n", tc.want[i].Header.Flags, reqs[i].Header.Flags)
					}
					for j, v := range reqs[i].Data {
						if v != tc.want[i].Data[j] {
							t.Fatalf("unexpected reply:\n- want: %#v\n-  got: %#v", tc.want[i].Data, reqs[i].Data)
						}
					}
				}
				return nil, nil
			})
			defer nfct.Con.Close()

			if err := nfct.Flush(ct.Ct, tc.family); err != nil {
				t.Fatal(err)
			}

		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		attributes []ct.ConnAttr
		want       []netlink.Message
	}{
		{name: "noAttributes", attributes: []ct.ConnAttr{}, want: []netlink.Message{
			{
				Header: netlink.Header{
					Length: 20,
					// NFNL_SUBSYS_CTNETLINK<<8|IPCTNL_MSG_CT_NEW
					Type: netlink.HeaderType(1<<8 | 0),
					// NLM_F_REQUEST|NLM_F_CREATE|NLM_F_ACK|NLM_F_EXCL
					Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsCreate | netlink.HeaderFlagsAcknowledge | netlink.HeaderFlagsExcl,
					// Can and will be ignored
					Sequence: 0,
					// Can and will be ignored
					PID: nltest.PID,
				},
				// nfgen_family=AF_INET, version=NFNETLINK_V0, NFNL_SUBSYS_CTNETLINK
				Data: []byte{0x2, 0x0, 0x0, 0x1},
			},
		}},
		// Example from libnetfilter_conntrack/utils/conntrack_create.c
		{name: "conntrack_create.c", attributes: []ct.ConnAttr{
			{Type: ct.AttrOrigIPv4Src, Data: []byte{0x1, 0x1, 0x1, 0x1}}, // SrcIP
			{Type: ct.AttrOrigIPv4Dst, Data: []byte{0x2, 0x2, 0x2, 0x2}}, // DstIP
			{Type: ct.AttrOrigL4Proto, Data: []byte{0x11}},               // TCP
			{Type: ct.AttrOrigPortSrc, Data: []byte{0x00, 0x14}},         // SrcPort
			{Type: ct.AttrOrigPortDst, Data: []byte{0x00, 0x0A}},         // DstPort
			{Type: ct.AttrTCPState, Data: []byte{0x1}},                   // TCP-State
			{Type: ct.AttrTimeout, Data: []byte{0x00, 0x00, 0x00, 0x64}}, // Timeout
		}, want: []netlink.Message{
			{
				Header: netlink.Header{
					Length: 80,
					// NFNL_SUBSYS_CTNETLINK<<8|IPCTNL_MSG_CT_NEW
					Type: netlink.HeaderType(1<<8 | 0),
					// NLM_F_REQUEST|NLM_F_CREATE|NLM_F_ACK|NLM_F_EXCL
					Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsCreate | netlink.HeaderFlagsAcknowledge | netlink.HeaderFlagsExcl,
					// Can and will be ignored
					Sequence: 0,
					// Can and will be ignored
					PID: nltest.PID,
				},
				// nfgen_family=AF_INET, version=NFNETLINK_V0, NFNL_SUBSYS_CTNETLINKa + netlinkes Attributes
				Data: []byte{0x2, 0x0, 0x0, 0x1, 0x2c, 0x0, 0x1, 0x80, 0x14, 0x0, 0x1, 0x80, 0x8, 0x0, 0x1, 0x0, 0x1, 0x1, 0x1, 0x1, 0x8, 0x0, 0x2, 0x0, 0x2, 0x2, 0x2, 0x2, 0x14, 0x0, 0x2, 0x80, 0x6, 0x0, 0x2, 0x0, 0x0, 0x14, 0x0, 0x0, 0x6, 0x0, 0x3, 0x0, 0x0, 0xa, 0x0, 0x0, 0x5, 0x0, 0x1, 0x0, 0x11, 0x0, 0x0, 0x0, 0x8, 0x0, 0x7, 0x0, 0x0, 0x0, 0x0, 0x64},
			},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nfct := &ct.Nfct{}
			nfct.Con = nltest.Dial(func(reqs []netlink.Message) ([]netlink.Message, error) {
				if len(reqs) == 0 {
					return nil, nil
				}
				if len(reqs) != 1 {
					t.Fatalf("Expected only one request, got %d", len(reqs))
				}
				// To ignore the Sequence number, we set it to the same value
				tc.want[0].Header.Sequence = reqs[0].Header.Sequence

				if len(reqs) != len(tc.want) {
					t.Fatalf("differen length:\n- want: %#v\n- got: %#v\n", tc.want, reqs)
				}

				for i := 0; i < len(reqs); i++ {
					if len(reqs[i].Data) != len(tc.want[i].Data) {
						t.Fatalf("differen length:\n- want: %#v\n- got: %#v\n", tc.want[i], reqs[i])
					}
					if reqs[i].Header.Type != tc.want[i].Header.Type {
						t.Fatalf("differen header types:\n- want: %#v\n- got: %#v\n", tc.want[i].Header.Type, reqs[i].Header.Type)
					}
					if reqs[i].Header.Flags != tc.want[i].Header.Flags {
						t.Fatalf("differen header flags:\n- want: %#v\n- got: %#v\n", tc.want[i].Header.Flags, reqs[i].Header.Flags)
					}
					for j, v := range reqs[i].Data {
						if v != tc.want[i].Data[j] {
							t.Fatalf("unexpected reply:\n- want: %#v\n-  got: %#v", tc.want[i].Data, reqs[i].Data)
						}
					}
				}
				return nil, nil
			})
			defer nfct.Con.Close()

			if err := nfct.Create(ct.Ct, ct.CtIPv4, tc.attributes); err != nil {
				t.Fatal(err)
			}
		})
	}
}
