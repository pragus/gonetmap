package gonetmap

import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "DstType=int,uint16,uint32"

type DstType generic.Type

func ifaceTODstType(i interface{}) DstType {
	var idx DstType
	switch i := i.(type) {
	case int:
		idx = DstType(i)
	case int16:
		idx = DstType(i)
	case int32:
		idx = DstType(i)
	case int64:
		idx = DstType(i)

	case uint:
		idx = DstType(i)
	case uint16:
		idx = DstType(i)
	case uint32:
		idx = DstType(i)
	case uint64:
		idx = DstType(i)
	}
	return idx
}
