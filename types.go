package caedral

// ChatMessage is a single chat completion message.
type ChatMessage struct {
	Role    string  `json:"role"`
	Content *string `json:"content"`
	Name    string  `json:"name,omitempty"`
}

// NotreOptions configures optional Notre platform optimization.
type NotreOptions struct {
	Mode      string `json:"mode,omitempty"`
	Telemetry *bool  `json:"telemetry,omitempty"`
}

// NotrePublicMetadata is customer-facing Notre data on a chat completion.
type NotrePublicMetadata struct {
	Enabled      bool    `json:"enabled"`
	Mode         string  `json:"mode"`
	Intervened   bool    `json:"intervened"`
	FallbackUsed bool    `json:"fallback_used"`
	InputBefore  *int    `json:"input_before,omitempty"`
	InputSent    *int    `json:"input_sent,omitempty"`
	InputSaved   *int    `json:"input_saved,omitempty"`
	ValueUSD     *float64 `json:"value_usd,omitempty"`
	Result       string  `json:"result,omitempty"`
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
	Notre            *NotreOptions `json:"notre,omitempty"`
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
	Notre   *NotrePublicMetadata   `json:"notre,omitempty"`
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
	Enabled        bool `json:"enabled"`
	LimitCents     *int `json:"limitCents"`
	UsedCents      int  `json:"usedCents"`
	RemainingCents *int `json:"remainingCents"`
}

// UsagePool is an included quota pool from GET /v1/usage.
type UsagePool struct {
	UsedMilli      int     `json:"usedMilli"`
	LimitMilli     int     `json:"limitMilli"`
	UsedFormatted  string  `json:"usedFormatted,omitempty"`
	LimitFormatted string  `json:"limitFormatted,omitempty"`
	PercentUsed    float64 `json:"percentUsed"`
	Available      bool    `json:"available"`
}

// UsagePlan is the current commercial plan.
type UsagePlan struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Interval   string `json:"interval"`
	Status     string `json:"status"`
	PriceCents int    `json:"priceCents,omitempty"`
}

// UsageOnDemand is optional usage after included quota.
type UsageOnDemand struct {
	Mode             string `json:"mode"`
	Allowed          bool   `json:"allowed"`
	Blocked          bool   `json:"blocked,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	AccruedMilli     int    `json:"accruedMilli"`
	SpentMilli       int    `json:"spentMilli,omitempty"`
	AccruedFormatted string `json:"accruedFormatted,omitempty"`
	SpentFormatted   string `json:"spentFormatted,omitempty"`
}

// UsageSummary is the account usage response from GET /v1/usage.
type UsageSummary struct {
	AccountStatus  string         `json:"accountStatus"`
	Plan           UsagePlan      `json:"plan"`
	BillingPeriod  map[string]any `json:"billingPeriod,omitempty"`
	Pools          map[string]UsagePool `json:"pools"`
	OnDemand       UsageOnDemand  `json:"onDemand"`
	Quota          map[string]any `json:"quota,omitempty"`
}

// EmbeddingCreateRequest configures an embeddings call.
type EmbeddingCreateRequest struct {
	Model          string `json:"model"`
	Input          any    `json:"input"`
	Dimensions     *int   `json:"dimensions,omitempty"`
	InputType      string `json:"input_type,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
}

// EmbeddingData is one embedding vector.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingCreateResponse is the embeddings API response.
type EmbeddingCreateResponse struct {
	Object string           `json:"object"`
	Model  string           `json:"model"`
	Data   []EmbeddingData  `json:"data"`
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
	Model string           `json:"model"`
	Data  []ImageData      `json:"data"`
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
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
}

// RerankCreateResponse is the rerank API response.
type RerankCreateResponse struct {
	Model   string         `json:"model"`
	Results []RerankResult `json:"results"`
}
