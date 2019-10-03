package msg

import (
	"github.com/tendermint/go-amino"
)

var MsgCdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {

	cdc.RegisterInterface((*Msg)(nil), nil)
}

func init() {
	RegisterCodec(MsgCdc)
}
