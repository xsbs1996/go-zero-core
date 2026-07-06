package xpacket

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/xsbs1996/go-zero-core/xcode"
	"github.com/xsbs1996/go-zero-core/xreply"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestValidateLength(t *testing.T) {
	t.Parallel()

	if err := ValidateLength(nil); !errors.Is(err, ErrPacketTooShort) {
		t.Fatalf("expected ErrPacketTooShort, got %v", err)
	}

	packet := make([]byte, HeaderLen)
	binary.BigEndian.PutUint32(packet[2:HeaderLen], 1)
	if err := ValidateLength(packet); !errors.Is(err, ErrPacketLengthMismatch) {
		t.Fatalf("expected ErrPacketLengthMismatch, got %v", err)
	}
}

func TestEncodeDecodeBody(t *testing.T) {
	t.Parallel()

	packet, err := EncodeBody(24, []byte("payload"))
	if err != nil {
		t.Fatalf("EncodeBody() error = %v", err)
	}

	action, body, err := DecodeBody(packet)
	if err != nil {
		t.Fatalf("DecodeBody() error = %v", err)
	}
	if action != 24 {
		t.Fatalf("action = %d, want 24", action)
	}
	if string(body) != "payload" {
		t.Fatalf("body = %q, want %q", body, "payload")
	}
}

type businessStringValue struct {
	wrapperspb.StringValue
}

func (v *businessStringValue) ValidateBusinessProto() error {
	if v.Value == "blocked" {
		return fmt.Errorf("blocked proto value")
	}
	return nil
}

func TestEncodeDecodeProto(t *testing.T) {
	t.Parallel()

	packet, err := EncodeProto(7, wrapperspb.String("hello"))
	if err != nil {
		t.Fatalf("EncodeProto() error = %v", err)
	}

	var req wrapperspb.StringValue
	action, err := DecodeProto(packet, &req)
	if err != nil {
		t.Fatalf("DecodeProto() error = %v", err)
	}
	if action != 7 || req.Value != "hello" {
		t.Fatalf("action/value = %d/%q, want 7/%q", action, req.Value, "hello")
	}
}

func TestDecodeProtoCallsProtoBusinessValidator(t *testing.T) {
	t.Parallel()

	packet, err := EncodeProto(1, wrapperspb.String("blocked"))
	if err != nil {
		t.Fatalf("EncodeProto() error = %v", err)
	}

	var req businessStringValue
	_, err = DecodeProto(packet, &req)
	if err == nil || err.Error() != "blocked proto value" {
		t.Fatalf("expected business validation error, got %v", err)
	}
}

type jsonUser struct {
	Name string `json:"name"`
}

type businessJsonUser struct {
	Name string `json:"name"`
}

func (u *businessJsonUser) ValidateBusinessJson() error {
	if u.Name == "" {
		return fmt.Errorf("missing json user name")
	}
	if u.Name == "blocked" {
		return fmt.Errorf("blocked json user")
	}
	return nil
}

func TestEncodeDecodeJson(t *testing.T) {
	t.Parallel()

	packet, err := EncodeJson(9, jsonUser{Name: "alice"})
	if err != nil {
		t.Fatalf("EncodeJson() error = %v", err)
	}

	var user jsonUser
	action, err := DecodeJson(packet, &user)
	if err != nil {
		t.Fatalf("DecodeJson() error = %v", err)
	}
	if action != 9 || user.Name != "alice" {
		t.Fatalf("action/name = %d/%q, want 9/%q", action, user.Name, "alice")
	}
}

func TestDecodeJsonCallsJsonBusinessValidator(t *testing.T) {
	t.Parallel()

	packet, err := EncodeJson(2, businessJsonUser{Name: "blocked"})
	if err != nil {
		t.Fatalf("EncodeJson() error = %v", err)
	}

	var user businessJsonUser
	_, err = DecodeJson(packet, &user)
	if err == nil || err.Error() != "blocked json user" {
		t.Fatalf("expected business validation error, got %v", err)
	}
}

func TestEncodeDecodeJsonResult(t *testing.T) {
	t.Parallel()

	packet, err := EncodeJsonResult(3, xcode.CodeInvalidParam, jsonUser{Name: "alice"}, xcode.Vars{"field": "name"})
	if err != nil {
		t.Fatalf("EncodeJsonResult() error = %v", err)
	}

	action, result, err := DecodeJsonResult[jsonUser](packet)
	if err != nil {
		t.Fatalf("DecodeJsonResult() error = %v", err)
	}
	if action != 3 {
		t.Fatalf("action = %d, want 3", action)
	}
	if result.Code != xcode.CodeInvalidParam || result.Msg != "invalid param: name" || result.Data.Name != "alice" {
		t.Fatalf("unexpected result: %+v", result)
	}

	_, body, err := DecodeBody(packet)
	if err != nil {
		t.Fatalf("DecodeBody() error = %v", err)
	}
	var raw xreply.Result[jsonUser]
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
}
