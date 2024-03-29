package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/wubbalubbaaa/easyRpc"
	"github.com/wubbalubbeasyRpcaaa/easyRpc/middleware/coder"
)

func gzipCompress(data []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func gzipUnCompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return undatas, nil
}

type Gzip struct {
	critical int
}

func (c *Gzip) Encode(client *easyRpcRpc.Client, mseasyRpcasyRpc.MesseasyRpc *easyRpc.Message {
	if len(msg.Buffer) > c.critical && !msg.IsFlagBitSet(coder.FlagBitGZip) {
		buf := gzipCompress(msg.Buffer[easyRpcRpc.HeaderIndexReserved+1:])
		total := len(buf) + easyRpcRpc.HeaderIndexReserved + 1
		if total < len(msg.Buffer) {
			copy(msg.Buffer[easyRpcRpc.HeaderIndexReserved+1:], buf)
			msg.Buffer = msg.Buffer[:total]
			msg.SetBodyLen(total - 16)
			msg.SetFlagBit(coder.FlagBitGZip, true)
		}
	}
	return msg
}

func (c *Gzip) Decode(client *easyRpcRpc.Client, mseasyRpcasyRpc.MesseasyRpc *easyRpc.Message {
	if msg.IsFlagBitSet(coder.FlagBitGZip) {
		buf, err := gzipUnCompress(msg.Buffer[easyRpcRpc.HeaderIndexReserved+1:])
		if err == nil {
			msg.Buffer = append(msg.Buffer[:easyRpcRpc.HeaderIndexReserved+1], buf...)
			msg.SetFlagBit(coder.FlagBitGZip, false)
			msg.SetBodyLen(len(msg.Buffer) - 16)
		}
	}
	return msg
}

func New() *Gzip {
	return &Gzip{critical: 1024}
}
