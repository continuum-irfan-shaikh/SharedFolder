package http

import (
	"io"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/protocol"
)

type responseSerializerImpl struct{}

func (ser *responseSerializerImpl) Serialize(res *protocol.Response, dst io.Writer) (err error) {
	return responseSerializeRaw(res, dst)
}

func (ser *responseSerializerImpl) Deserialize(src io.Reader) (res *protocol.Response, err error) {
	return responseDeserializeRaw(src)
}
