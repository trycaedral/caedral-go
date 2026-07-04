package caedral

// ChatMessage is a single chat completion message.
type ChatMessage struct {
	Role    string  `json:"role"`
	Content *string `json:"content"`
	Name    string  `json:"name,omitempty"`
}

// ChatCompletionRequest configures a chat completion call.
type ChatCompletionRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	Stream           bool          `json:"stream,omitempty"`
	Temperature      *float64      `json:"temperature,omitempty"`
	MaxTokens        *int          `json:"max_tokens,omitempty"`
	TopP             *float64      `json:"top_p,omitempty"`
	FrequencyPenalty *float64      `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64      `json:"presence_penalty,omitempty"`
	Stop             any           `json:"stop,omitempty"`
	User             string        `json:"user,omitempty"`
}

// ChatCompletionChoice is one completion choice.
type ChatCompletionChoice struct {
	Index        int            `json:"index"`
	Message      map[string]any `json:"message"`
	FinishReason *string        `json:"finish_reason"`
}

// CompletionUsage reports token usage.
type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletion is a non-streaming chat completion response.
type ChatCompletion struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   *CompletionUsage       `json:"usage,omitempty"`
}

// ChatCompletionChunkChoice is a streaming delta choice.
type ChatCompletionChunkChoice struct {
	Index        int            `json:"index"`
	Delta        map[string]any `json:"delta"`
	FinishReason *string        `json:"finish_reason"`
}

// ChatCompletionChunk is one SSE chunk from a streaming completion.
type ChatCompletionChunk struct {
	ID      string                      `json:"id"`
	Object  string                      `json:"object"`
	Created int64                       `json:"created"`
	Model   string                      `json:"model"`
	Choices []ChatCompletionChunkChoice `json:"choices"`
}

// Model describes a Caedral model.
type Model struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Created       int64  `json:"created"`
	OwnedBy       string `json:"owned_by"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ContextWindow int    `json:"context_window"`
	PricingTier   string `json:"pricing_tier"`
}

// ModelListResponse is the response from Models.List.
type ModelListResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// WeeklyPool summarizes weekly token pool usage.
type WeeklyPool struct {
	Limit     int `json:"limit"`
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
}

// OverageSummary summarizes overage billing.
type OverageSummary struct {
	Enabled         bool `json:"enabled"`
	LimitCents      *int `json:"limitCents"`
	UsedCents       int  `json:"usedCents"`
	RemainingCents  *int `json:"remainingCents"`
}

// UsageSummary is the account usage response.
type UsageSummary struct {
	AccountStatus                  string         `json:"accountStatus"`
	Plan                           string         `json:"plan"`
	PlanStatus                     string         `json:"planStatus"`
	BalanceCents                   int            `json:"balanceCents"`
	WeeklyPool                     WeeklyPool     `json:"weeklyPool"`
	Overage                        OverageSummary `json:"overage"`
	BalanceWeightedUnitsAffordable int            `json:"balanceWeightedUnitsAffordable"`
}

// EmbeddingCreateRequest configures an embeddings call.
type EmbeddingCreateRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"`
}

// EmbeddingData is one embedding vector.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingCreateResponse is the embeddings API response.
type EmbeddingCreateResponse struct {
	Object string          `json:"object"`
	Model  string          `json:"model"`
	Data   []EmbeddingData `json:"data"`
	Usage  *CompletionUsage `json:"usage,omitempty"`
}

// ImageGenerateRequest configures image generation.
type ImageGenerateRequest struct {
	Model  string `json:"model,omitempty"`
	Prompt string `json:"prompt"`
	N      *int   `json:"n,omitempty"`
	Size   string `json:"size,omitempty"`
}

// ImageData holds one generated image.
type ImageData struct {
	URL     string `json:"url,omitempty"`
	B64JSON string `json:"b64_json,omitempty"`
}

// ImageGenerateResponse is the image generation response.
type ImageGenerateResponse struct {
	Model string      `json:"model"`
	Data  []ImageData `json:"data"`
	Usage *CompletionUsage `json:"usage,omitempty"`
}

// AudioGenerateRequest configures speech generation.
type AudioGenerateRequest struct {
	Model string `json:"model,omitempty"`
	Input string `json:"input"`
	Voice string `json:"voice,omitempty"`
}

// AudioGenerateResponse is the speech generation response.
type AudioGenerateResponse struct {
	Model   string           `json:"model"`
	Choices []map[string]any `json:"choices,omitempty"`
	Usage   *CompletionUsage `json:"usage,omitempty"`
}

// RerankCreateRequest configures a rerank call.
type RerankCreateRequest struct {
	Model     string   `json:"model,omitempty"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      *int     `json:"top_n,omitempty"`
}

// RerankResult is one reranked document.
type RerankResult struct {
	Index           int     `json:"index"`
	RelevanceScore  float64 `json:"relevance_score"`
}

// RerankCreateResponse is the rerank API response.
type RerankCreateResponse struct {
	Model   string         `json:"model"`
	Results []RerankResult `json:"results"`
}
