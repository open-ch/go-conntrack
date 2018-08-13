// +build linux,!386

package conntrack

import (
	"reflect"
	"testing"

	"golang.org/x/net/bpf"
)

func TestConstructFilter(t *testing.T) {
	tests := []struct {
		name     string
		table    CtTable
		filters  []ConnAttr
		rawInstr []bpf.RawInstruction
		err      error
	}{
		// Modified example from libnetfilter_conntrack/utils/conntrack_filter.c
		{name: "conntrack_filter.c", table: Ct, filters: []ConnAttr{
			{Type: AttrOrigL4Proto, Data: []byte{0x11}},                                                                    // TCP
			{Type: AttrOrigL4Proto, Data: []byte{0x06}},                                                                    // UDP
			{Type: AttrTCPState, Data: []byte{0x3}},                                                                        // TCP_CONNTRACK_ESTABLISHED
			{Type: AttrOrigIPv4Src, Data: []byte{0x7F, 0x0, 0x0, 0x1}, Mask: []byte{0xff, 0xff, 0xff, 0xff}, Negate: true}, // SrcIP != 127.0.0.1
			{Type: AttrOrigIPv6Src, Data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, // SrcIP != ::1
				Mask: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Negate: true},
		}, rawInstr: []bpf.RawInstruction{
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0050, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0015, Jt: 0x01, Jf: 0x00, K: 0x00000001},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0000, Jt: 0x00, Jf: 0x00, K: 0x00000014},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0e, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000002},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0a, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x06, Jf: 0x00, K: 0x00000000},
			{Op: 0x0007, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0050, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0015, Jt: 0x02, Jf: 0x00, K: 0x00000011},
			{Op: 0x0050, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0015, Jt: 0x01, Jf: 0x00, K: 0x00000006},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0000, Jt: 0x00, Jf: 0x00, K: 0x00000014},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0c, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x08, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x04, Jf: 0x00, K: 0x00000000},
			{Op: 0x0007, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0050, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0015, Jt: 0x01, Jf: 0x00, K: 0x00000003},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0000, Jt: 0x00, Jf: 0x00, K: 0x00000014},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0e, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0a, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x06, Jf: 0x00, K: 0x00000000},
			{Op: 0x0007, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0040, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0054, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0015, Jt: 0x02, Jf: 0x00, K: 0x7f000001},
			{Op: 0x0005, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0000, Jt: 0x00, Jf: 0x00, K: 0x00000014},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x17, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x13, Jf: 0x00, K: 0x00000000},
			{Op: 0x0004, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0001, Jt: 0x00, Jf: 0x00, K: 0x00000003},
			{Op: 0x0030, Jt: 0x00, Jf: 0x00, K: 0xfffff00c},
			{Op: 0x0015, Jt: 0x0f, Jf: 0x00, K: 0x00000000},
			{Op: 0x0007, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0040, Jt: 0x00, Jf: 0x00, K: 0x00000004},
			{Op: 0x0054, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0015, Jt: 0x00, Jf: 0x0a, K: 0x00000000},
			{Op: 0x0040, Jt: 0x00, Jf: 0x00, K: 0x00000008},
			{Op: 0x0054, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0015, Jt: 0x00, Jf: 0x07, K: 0x00000000},
			{Op: 0x0040, Jt: 0x00, Jf: 0x00, K: 0x0000000c},
			{Op: 0x0054, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0015, Jt: 0x00, Jf: 0x04, K: 0x00000000},
			{Op: 0x0040, Jt: 0x00, Jf: 0x00, K: 0x00000010},
			{Op: 0x0054, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
			{Op: 0x0015, Jt: 0x02, Jf: 0x00, K: 0x00000001},
			{Op: 0x0005, Jt: 0x00, Jf: 0x00, K: 0x00000001},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0x00000000},
			{Op: 0x0006, Jt: 0x00, Jf: 0x00, K: 0xffffffff},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rawInstr, err := constructFilter(tc.table, tc.filters)
			if err != tc.err {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(rawInstr, tc.rawInstr) {
				t.Fatalf("unexpected replies:\n- want: %#v\n-  got: %#v", tc.rawInstr, rawInstr)
			}

		})
	}
}
