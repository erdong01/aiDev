package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func normalizeJSONPayload(raw []byte) json.RawMessage {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return json.RawMessage([]byte("null"))
	}

	if json.Valid([]byte(trimmed)) {
		return json.RawMessage([]byte(trimmed))
	}

	encoded, err := json.Marshal(trimmed)
	if err != nil {
		return json.RawMessage([]byte("null"))
	}
	return json.RawMessage(encoded)
}

func decodeUpstreamResponseBody(raw string, isBase64 bool) ([]byte, error) {
	payload := []byte(raw)
	if isBase64 {
		decoded, err := decodeBase64String(raw)
		if err != nil {
			return nil, err
		}
		payload = decoded
	}

	if len(payload) >= 2 && payload[0] == 0x1f && payload[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		decoded, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return decoded, nil
	}

	return payload, nil
}

func decodeRequestPayload(raw json.RawMessage, isBase64 bool) ([]byte, error) {
	trimmedRaw := bytes.TrimSpace(raw)
	if len(trimmedRaw) == 0 || bytes.Equal(trimmedRaw, []byte("null")) {
		return []byte("{}"), nil
	}

	var asString string
	if err := json.Unmarshal(trimmedRaw, &asString); err == nil {
		if strings.TrimSpace(asString) == "" {
			return []byte("{}"), nil
		}

		if isBase64 {
			decoded, err := decodeBase64String(asString)
			if err != nil {
				return nil, fmt.Errorf("request body base64 decode failed: %w", err)
			}

			trimmedDecoded := bytes.TrimSpace(decoded)
			if len(trimmedDecoded) == 0 {
				return []byte("{}"), nil
			}
			return trimmedDecoded, nil
		}

		// Some callbacks omit body_base64 even though the body is still a base64-encoded JSON string.
		// When the decoded content is valid JSON, prefer the decoded payload over the encoded string.
		if decoded, err := decodeBase64String(asString); err == nil {
			trimmedDecoded := bytes.TrimSpace(decoded)
			if len(trimmedDecoded) == 0 {
				return []byte("{}"), nil
			}
			if json.Valid(trimmedDecoded) {
				return trimmedDecoded, nil
			}
		}

		return []byte(asString), nil
	}

	if isBase64 {
		// Some upstreams already pass a JSON object here even when the flag is set.
		// Keeping the raw payload is safer than discarding the request log.
		return trimmedRaw, nil
	}

	return trimmedRaw, nil
}

func decodeBase64String(raw string) ([]byte, error) {
	var firstErr error
	for _, encoding := range []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	} {
		decoded, err := encoding.DecodeString(raw)
		if err == nil {
			return decoded, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	return nil, firstErr
}
