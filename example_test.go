package wabin

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/tetratelabs/wabin/binary"
	"github.com/tetratelabs/wabin/wasm"
)

func newExample() *wasm.Module {
	three := wasm.Index(3)
	f32, i32, i64 := wasm.ValueTypeF32, wasm.ValueTypeI32, wasm.ValueTypeI64
	return &wasm.Module{
		TypeSection: []*wasm.FunctionType{
			{Params: []wasm.ValueType{i32, i32}, Results: []wasm.ValueType{i32}},
			{},
			{Params: []wasm.ValueType{i32, i32, i32, i32}, Results: []wasm.ValueType{i32}},
			{Params: []wasm.ValueType{i64}, Results: []wasm.ValueType{i64}},
			{Params: []wasm.ValueType{f32}, Results: []wasm.ValueType{i32}},
			{Params: []wasm.ValueType{i32, i32}, Results: []wasm.ValueType{i32, i32}},
		},
		ImportSection: []*wasm.Import{
			{
				Module: "wasi_snapshot_preview1", Name: "args_sizes_get",
				Type:     wasm.ExternTypeFunc,
				DescFunc: 0,
			}, {
				Module: "wasi_snapshot_preview1", Name: "fd_write",
				Type:     wasm.ExternTypeFunc,
				DescFunc: 2,
			},
		},
		FunctionSection: []wasm.Index{wasm.Index(1), wasm.Index(1), wasm.Index(0), wasm.Index(3), wasm.Index(4), wasm.Index(5)},
		CodeSection: []*wasm.Code{
			{Body: []byte{wasm.OpcodeCall, 3, wasm.OpcodeEnd}},
			{Body: []byte{wasm.OpcodeEnd}},
			{Body: []byte{wasm.OpcodeLocalGet, 0, wasm.OpcodeLocalGet, 1, wasm.OpcodeI32Add, wasm.OpcodeEnd}},
			{Body: []byte{wasm.OpcodeLocalGet, 0, wasm.OpcodeI64Extend16S, wasm.OpcodeEnd}},
			{Body: []byte{
				wasm.OpcodeLocalGet, 0x00,
				wasm.OpcodeMiscPrefix, wasm.OpcodeMiscI32TruncSatF32S,
				wasm.OpcodeEnd,
			}},
			{Body: []byte{wasm.OpcodeLocalGet, 1, wasm.OpcodeLocalGet, 0, wasm.OpcodeEnd}},
		},
		MemorySection: &wasm.Memory{Min: 1, Max: three, IsMaxEncoded: true},
		ExportSection: []*wasm.Export{
			{Name: "AddInt", Type: wasm.ExternTypeFunc, Index: wasm.Index(4)},
			{Name: "", Type: wasm.ExternTypeFunc, Index: wasm.Index(3)},
			{Name: "mem", Type: wasm.ExternTypeMemory, Index: wasm.Index(0)},
			{Name: "swap", Type: wasm.ExternTypeFunc, Index: wasm.Index(7)},
		},
		StartSection: &three,
		NameSection: &wasm.NameSection{
			ModuleName: "example",
			FunctionNames: wasm.NameMap{
				{Index: wasm.Index(0), Name: "wasi.args_sizes_get"},
				{Index: wasm.Index(1), Name: "wasi.fd_write"},
				{Index: wasm.Index(2), Name: "call_hello"},
				{Index: wasm.Index(3), Name: "hello"},
				{Index: wasm.Index(4), Name: "addInt"},
				{Index: wasm.Index(7), Name: "swap"},
			},
			LocalNames: wasm.IndirectNameMap{
				{Index: wasm.Index(1), NameMap: wasm.NameMap{
					{Index: wasm.Index(0), Name: "fd"},
					{Index: wasm.Index(1), Name: "iovs_ptr"},
					{Index: wasm.Index(2), Name: "iovs_len"},
					{Index: wasm.Index(3), Name: "nwritten_ptr"},
				}},
				{Index: wasm.Index(4), NameMap: wasm.NameMap{
					{Index: wasm.Index(0), Name: "value_1"},
					{Index: wasm.Index(1), Name: "value_2"},
				}},
			},
		},
	}
}

func Example_binary() {
	bin := binary.EncodeModule(newExample())

	if mod, err := binary.DecodeModule(bin, wasm.CoreFeaturesV2); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(mod.NameSection.ModuleName)
	}

	// Output:
	// example
}
