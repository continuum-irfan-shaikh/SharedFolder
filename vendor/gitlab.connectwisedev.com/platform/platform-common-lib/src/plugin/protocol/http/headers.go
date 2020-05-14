package http

import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/protocol"
)

func getHTTPHeader(hdr protocol.HeaderKey) string {
	return string(hdr)
}

func getProtocolHeader(httpHeader string) protocol.HeaderKey {
	return protocol.HeaderKey(httpHeader)
}
