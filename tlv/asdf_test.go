package tlv_test

import (
	"bytes"
	crand "crypto/rand"
	"io"
	"testing"

	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/invoices"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/tlv"
)

func TestASDF(t *testing.T) {

	var preImage [32]byte
	_, err := io.ReadFull(crand.Reader, preImage[:])
	if err != nil {
		t.Fatalf("ERR: %v", err.Error())
	}

	// Now that we have our pre-image, we'll encode it as a TLV
	// record into a temporary bytes buffer that will become our
	// final TLV stream.
	var b bytes.Buffer
	tlvStream, err := tlv.NewStream(
		tlv.MakePrimitiveRecord(
			invoices.PreimageTLV, &preImage,
		),
		// tlv.MakeSentinelRecord(),
	)
	if err := tlvStream.Encode(&b); err != nil {
		t.Fatalf("ERR: %v", err.Error())
	}

	eob := b.Bytes()
	t.Logf("TLV bytes: %x\n", eob)

	/* ========== DECODE ========== */

	var (
		preImage2     [32]byte
		blankPreImage [32]byte
	)
	tlvStream2, err := tlv.NewStream(
		tlv.MakePrimitiveRecord(invoices.PreimageTLV, &preImage2),
	) // instead of this was "PreImageTLV" = 128
	err = tlvStream2.Decode(
		bytes.NewReader(eob),
	)
	t.Logf("decoded TLV: %+v\n", tlvStream2)

	if err != nil {
		t.Fatalf("ERR: %v", err.Error())
	}

	var rHash lntypes.Hash
	paymentSecret := lntypes.Preimage(preImage2)
	if preImage != blankPreImage && paymentSecret.Matches(rHash) {

		t.Log("settled")

		// TODO(roasbeef): record settled spontaneous payment
		// in DB

		var circuitKey channeldb.CircuitKey
		var acceptHeight int32
		hodlEvent := &invoices.HodlEvent{
			CircuitKey:   circuitKey,
			AcceptHeight: acceptHeight,
			Preimage:     &paymentSecret,
		}
		t.Logf("HODL EVENT: %+v\n", hodlEvent)
	}
}
