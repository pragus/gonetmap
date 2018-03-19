package gonetmap

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "DstType=int,uint16,uint32"

type DstType generic.Type

func ifaceTODstType(i interface{}) DstType {
	var idx DstType
	switch i.(type) {
	case int:
		idx = DstType(i.(int))
	case int16:
		idx = DstType(i.(int16))
	case int32:
		idx = DstType(i.(int32))
	case int64:
		idx = DstType(i.(int64))

	case uint:
		idx = DstType(i.(int))
	case uint16:
		idx = DstType(i.(uint16))
	case uint32:
		idx = DstType(i.(uint32))
	case uint64:
		idx = DstType(i.(uint64))
	}
	return idx
}
