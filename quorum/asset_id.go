package quorum

import (
	"crypto/md5"
	"io"
	"strings"

	"github.com/gofrs/uuid"
)

func MixinAssetID(assetKey string) string {
	h := md5.New()
	_, _ = io.WriteString(h, EthAsset)
	_, _ = io.WriteString(h, strings.ToLower(assetKey))
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}
