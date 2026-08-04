package caedral

import "fmt"

// DefaultEmbeddingModel is the default Caedral E1 Small embedding model.
const DefaultEmbeddingModel = "caedral-embed-e1-small-v1"

// DefaultEmbeddingDimensions is the native vector size for DefaultEmbeddingModel.
const DefaultEmbeddingDimensions = 384

var supportedEmbeddingModels = map[string]struct{}{
	DefaultEmbeddingModel: {},
}

var supportedEmbeddingDimensions = map[int]struct{}{
	DefaultEmbeddingDimensions: {},
}

func validateEmbeddingModel(model string) error {
	if _, ok := supportedEmbeddingModels[model]; !ok {
		return fmt.Errorf("unsupported embedding model: %s", model)
	}
	return nil
}

func validateEmbeddingDimensions(dimensions int) error {
	if _, ok := supportedEmbeddingDimensions[dimensions]; !ok {
		return fmt.Errorf("unsupported embedding dimensions: %d", dimensions)
	}
	return nil
}
