package libkb

import (
	"strings"
	"testing"

	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"github.com/stretchr/testify/require"
)

// requireErrorHasSuffix makes sure that err's string has
// expectedErrSuffix's string as a suffix. This is necessary as
// go-codec prepends stuff to any errors it catches.
func requireErrorHasSuffix(t *testing.T, expectedErrSuffix, err error) {
	t.Helper()
	require.Error(t, err)
	require.True(t, strings.HasSuffix(err.Error(), expectedErrSuffix.Error()), "Expected %q to have %q as a suffix", err, expectedErrSuffix)
}

// Test that various ways in which OuterLinkV2WithMetadata may be
// encoded/decoded all fail.

func TestOuterLinkV2WithMetadataEncode(t *testing.T) {
	var o OuterLinkV2WithMetadata
	_, err := MsgpackEncode(o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
	_, err = MsgpackEncode(&o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
}

func TestOuterLinkV2WithMetadataDecode(t *testing.T) {
	var o OuterLinkV2WithMetadata
	err := MsgpackDecode(&o, []byte{0x1, 0x2})
	requireErrorHasSuffix(t, errCodecDecodeSelf, err)
}

type outerLinkV2WithMetadataEmbedder struct {
	OuterLinkV2WithMetadata
}

func TestOuterLinkV2WithMetadataEmbedderEncode(t *testing.T) {
	var o outerLinkV2WithMetadataEmbedder
	_, err := MsgpackEncode(o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
}

func TestOuterLinkV2WithMetadataEmbedderDecode(t *testing.T) {
	var o outerLinkV2WithMetadataEmbedder
	err := MsgpackDecode(&o, []byte{0x1, 0x2})
	requireErrorHasSuffix(t, errCodecDecodeSelf, err)
}

type outerLinkV2WithMetadataPointerEmbedder struct {
	*OuterLinkV2WithMetadata
}

func TestOuterLinkV2WithMetadataPointerEmbedderEncode(t *testing.T) {
	o := outerLinkV2WithMetadataPointerEmbedder{
		OuterLinkV2WithMetadata: &OuterLinkV2WithMetadata{},
	}
	_, err := MsgpackEncode(o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
}

func TestOuterLinkV2WithMetadataPointerEmbedderDecode(t *testing.T) {
	var o outerLinkV2WithMetadataPointerEmbedder
	err := MsgpackDecode(&o, []byte{0x1, 0x2})
	requireErrorHasSuffix(t, errCodecDecodeSelf, err)
}

type outerLinkV2WithMetadataContainer struct {
	O OuterLinkV2WithMetadata
}

func TestOuterLinkV2WithMetadataContainerEncode(t *testing.T) {
	var o outerLinkV2WithMetadataContainer
	_, err := MsgpackEncode(o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
}

func TestOuterLinkV2WithMetadataContainerDecode(t *testing.T) {
	bytes, err := MsgpackEncode(map[string]interface{}{
		"O": []byte{0x01, 0x02},
	})
	require.NoError(t, err)

	var o outerLinkV2WithMetadataContainer
	err = MsgpackDecode(&o, bytes)
	requireErrorHasSuffix(t, errCodecDecodeSelf, err)
}

type outerLinkV2WithMetadataPointerContainer struct {
	O *OuterLinkV2WithMetadata
}

func TestOuterLinkV2WithMetadataPointerContainerEncode(t *testing.T) {
	o := outerLinkV2WithMetadataPointerContainer{
		O: &OuterLinkV2WithMetadata{},
	}
	_, err := MsgpackEncode(o)
	requireErrorHasSuffix(t, errCodecEncodeSelf, err)
}

func TestOuterLinkV2WithMetadataPointerContainerDecode(t *testing.T) {
	bytes, err := MsgpackEncode(map[string]interface{}{
		"O": []byte{0x01, 0x02},
	})
	require.NoError(t, err)

	var o outerLinkV2WithMetadataPointerContainer
	err = MsgpackDecode(&o, bytes)
	requireErrorHasSuffix(t, errCodecDecodeSelf, err)
}

func serdeOuterLink(t *testing.T, tc TestContext, hPrevSeqno *keybase1.Seqno, hPrevHash LinkID) OuterLinkV2 {
	m := NewMetaContextForTest(tc)

	ModifyFeatureForTest(m, FeatureAllowHighSkips, true, 100000)
	ModifyFeatureForTest(m, FeatureRequireHighSkips, false, 100000)

	var hPrevInfo *HPrevInfo
	hPrevInfo = nil
	if hPrevSeqno != nil {
		hPrevInfoPre := NewHPrevInfo(*hPrevSeqno, hPrevHash)
		hPrevInfo = &hPrevInfoPre
	}
	encoded, err := encodeOuterLink(m, LinkTypeLeave, keybase1.Seqno(20), nil, nil, false, keybase1.SeqType_PUBLIC, false, hPrevInfo)
	require.NoError(t, err)

	o2 := OuterLinkV2{}
	err = MsgpackDecode(&o2, encoded)
	require.NoError(t, err)
	return o2
}

func hPrevInfoMatch(ol OuterLinkV2, hPrevInfo *HPrevInfo) error {
	return ol.AssertFields(ol.Version, ol.Seqno, ol.Prev, ol.Curr, ol.LinkType, ol.SeqType, ol.IgnoreIfUnsupported, hPrevInfo)
}

func TestHPrevBackwardsCompatibility(t *testing.T) {
	tc := SetupTest(t, "test_hprev_backwards_compatibility", 1)
	defer tc.Cleanup()

	o := serdeOuterLink(t, tc, nil, nil)
	require.True(t, o.HPrevSeqno == nil)
	require.True(t, o.HPrevHash == nil)
	require.NoError(t, hPrevInfoMatch(o, nil))

	seqno := keybase1.Seqno(3)
	o = serdeOuterLink(t, tc, &seqno, nil)
	require.True(t, *o.HPrevSeqno == 3)
	require.True(t, o.HPrevHash == nil)
	hPrevInfo := NewHPrevInfo(keybase1.Seqno(3), nil)
	badHPrevInfo1 := NewHPrevInfo(keybase1.Seqno(3), []byte{0, 5, 2})
	badHPrevInfo2 := NewHPrevInfo(keybase1.Seqno(4), nil)
	require.NoError(t, hPrevInfoMatch(o, &hPrevInfo))
	require.Error(t, hPrevInfoMatch(o, &badHPrevInfo1))
	require.Error(t, hPrevInfoMatch(o, &badHPrevInfo2))

	hPrevHash := []byte{0, 6, 2, 42, 123}
	o = serdeOuterLink(t, tc, &seqno, hPrevHash)
	hPrevInfo = NewHPrevInfo(keybase1.Seqno(3), hPrevHash)
	badHPrevInfo1 = NewHPrevInfo(keybase1.Seqno(3), []byte{0, 5, 2})
	badHPrevInfo2 = NewHPrevInfo(keybase1.Seqno(3), nil)
	require.NoError(t, hPrevInfoMatch(o, &hPrevInfo))
	require.Error(t, hPrevInfoMatch(o, &badHPrevInfo1))
	require.Error(t, hPrevInfoMatch(o, &badHPrevInfo2))
}
