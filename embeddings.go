package caedral

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// EmbeddingsService creates text embeddings.
type EmbeddingsService struct {
	client *Client
}

// Create generates embeddings for the given input.
func (s EmbeddingsService) Create(ctx context.Context, req EmbeddingCreateRequest) (*EmbeddingCreateResponse, error) {
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = DefaultEmbeddingModel
	}
	if err := validateEmbeddingModel(model); err != nil {
		return nil, err
	}

	dimensions := DefaultEmbeddingDimensions
	if req.Dimensions != nil {
		dimensions = *req.Dimensions
	}
	if err := validateEmbeddingDimensions(dimensions); err != nil {
		return nil, err
	}

	body := EmbeddingCreateRequest{
		Model:          model,
		Input:          req.Input,
		Dimensions:     &dimensions,
		InputType:      req.InputType,
		EncodingFormat: req.EncodingFormat,
	}

	var raw struct {
		Object string `json:"object"`
		Model  string `json:"model"`
		Data   []struct {
			Object    string          `json:"object"`
			Index     int             `json:"index"`
			Embedding json.RawMessage `json:"embedding"`
		} `json:"data"`
		Usage *CompletionUsage `json:"usage,omitempty"`
	}
	if err := s.client.doPostJSON(ctx, "/v1/embeddings", body, &raw); err != nil {
		return nil, err
	}

	out := &EmbeddingCreateResponse{
		Object: raw.Object,
		Model:  raw.Model,
		Usage:  raw.Usage,
		Data:   make([]EmbeddingData, len(raw.Data)),
	}
	for i, item := range raw.Data {
		vector, err := decodeEmbeddingVector(item.Embedding, dimensions)
		if err != nil {
			return nil, fmt.Errorf("caedral: decode embedding index %d: %w", item.Index, err)
		}
		out.Data[i] = EmbeddingData{
			Object:    item.Object,
			Index:     item.Index,
			Embedding: vector,
		}
	}
	return out, nil
}

func decodeEmbeddingVector(raw json.RawMessage, dimensions int) ([]float64, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("missing embedding value")
	}

	var floats []float64
	if err := json.Unmarshal(raw, &floats); err == nil {
		return floats, nil
	}

	var encoded string
	if err := json.Unmarshal(raw, &encoded); err != nil {
		return nil, fmt.Errorf("unsupported embedding encoding")
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 embedding: %w", err)
	}
	expectedBytes := dimensions * 4
	if len(decoded) != expectedBytes {
		return nil, fmt.Errorf(
			"base64 payload is %d bytes, expected %d for %d dimensions",
			len(decoded),
			expectedBytes,
			dimensions,
		)
	}

	out := make([]float64, dimensions)
	for i := 0; i < dimensions; i++ {
		bits := binary.LittleEndian.Uint32(decoded[i*4 : i*4+4])
		out[i] = float64(math.Float32frombits(bits))
	}
	return out, nil
}
