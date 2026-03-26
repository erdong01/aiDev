package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestDecodeRequestPayload_Base64JSONBody(t *testing.T) {
	payload := `{
    "model": "doubao-seedance-1-5-pro-251215",
    "content": [
        {
            "type": "text",
            "text": "女孩抱着狐狸，女孩轻轻抚摸狐狸，温柔地看向镜头，镜头缓缓拉出，女孩的头发被风吹动，可以听到风声"
        },
        {
            "type": "image_url",
            "image_url": {
                "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/i2v_foxrgirl.png"
            }
        }
    ],
    "generate_audio": true,
    "ratio": "adaptive",
    "duration": 5,
    "watermark": false
}`

	raw, err := json.Marshal(base64.StdEncoding.EncodeToString([]byte(payload)))
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	decoded, err := decodeRequestPayload(json.RawMessage(raw), true)
	if err != nil {
		t.Fatalf("decode request payload: %v", err)
	}

	if string(decoded) != payload {
		t.Fatalf("decoded payload mismatch\nwant: %s\n got: %s", payload, string(decoded))
	}

	if got := string(normalizeJSONPayload(decoded)); got == "{}" {
		t.Fatalf("normalized payload unexpectedly fell back to empty JSON")
	}
}

func TestDecodeRequestPayload_Base64PlainText(t *testing.T) {
	raw, err := json.Marshal(base64.StdEncoding.EncodeToString([]byte("hello world")))
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	decoded, err := decodeRequestPayload(json.RawMessage(raw), true)
	if err != nil {
		t.Fatalf("decode request payload: %v", err)
	}

	if string(decoded) != "hello world" {
		t.Fatalf("unexpected decoded payload: %q", string(decoded))
	}

	if got := string(normalizeJSONPayload(decoded)); got != `"hello world"` {
		t.Fatalf("unexpected normalized payload: %s", got)
	}
}

func TestDecodeRequestPayload_Base64JSONBodyWithoutFlag(t *testing.T) {
	raw, err := json.Marshal(base64.StdEncoding.EncodeToString([]byte(`{"hello":"world"}`)))
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	decoded, err := decodeRequestPayload(json.RawMessage(raw), false)
	if err != nil {
		t.Fatalf("decode request payload: %v", err)
	}

	if string(decoded) != `{"hello":"world"}` {
		t.Fatalf("unexpected decoded payload: %s", string(decoded))
	}
}

func TestDecodeRequestPayload_RawJSONObject(t *testing.T) {
	raw := json.RawMessage(`{"hello":"world"}`)

	decoded, err := decodeRequestPayload(raw, false)
	if err != nil {
		t.Fatalf("decode request payload: %v", err)
	}

	if string(decoded) != string(raw) {
		t.Fatalf("unexpected decoded payload: %s", string(decoded))
	}
}

func TestDecodeUpstreamResponseBody_Base64Gzip(t *testing.T) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write([]byte(`{"id":"cgt-20260324211610-sfnfs"}`)); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	decoded, err := decodeUpstreamResponseBody(encoded, true)
	if err != nil {
		t.Fatalf("decode upstream response: %v", err)
	}

	if string(decoded) != `{"id":"cgt-20260324211610-sfnfs"}` {
		t.Fatalf("unexpected decoded response: %s", string(decoded))
	}
}
