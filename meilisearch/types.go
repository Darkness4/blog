package meilisearch

import "encoding/json"

type (
	TaskStatus       string
	MatchingStrategy string
)

const (
	// TaskStatusUnknown is the default TaskStatus, should not exist
	TaskStatusUnknown TaskStatus = "unknown"
	// TaskStatusEnqueued the task request has been received and will be processed soon
	TaskStatusEnqueued TaskStatus = "enqueued"
	// TaskStatusProcessing the task is being processed
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusSucceeded the task has been successfully processed
	TaskStatusSucceeded TaskStatus = "succeeded"
	// TaskStatusFailed a failure occurred when processing the task, no changes were made to the database
	TaskStatusFailed TaskStatus = "failed"
	// TaskStatusCanceled the task was canceled
	TaskStatusCanceled TaskStatus = "canceled"
)

const (
	// Last returns documents containing all the query terms first. If there are not enough results containing all
	// query terms to meet the requested limit, Meilisearch will remove one query term at a time,
	// starting from the end of the query.
	Last MatchingStrategy = "last"
	// All only returns documents that contain all query terms. Meilisearch will not match any more documents even
	// if there aren't enough to meet the requested limit.
	All MatchingStrategy = "all"
	// Frequency returns documents containing all the query terms first. If there are not enough results containing
	//all query terms to meet the requested limit, Meilisearch will remove one query term at a time, starting
	//with the word that is the most frequent in the dataset. frequency effectively gives more weight to terms
	//that appear less frequently in a set of results.
	Frequency MatchingStrategy = "frequency"
)

type SubmittedTaskResponse struct {
	TaskUID  int        `json:"taskUid"`
	IndexUID string     `json:"indexUid"`
	Status   TaskStatus `json:"status"`
	Type     string     `json:"type"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	Link    string `json:"link"`
}

type GetTaskResponse struct {
	UID        int           `json:"uid"`
	IndexUID   string        `json:"indexUid"`
	Status     TaskStatus    `json:"status"`
	Type       string        `json:"type"`
	CanceledBy int           `json:"canceledBy"`
	Details    any           `json:"details"`
	Error      ErrorResponse `json:"error"`
}

type CreateKeyRequest struct {
	Actions     []string `json:"actions"`
	Indexes     []string `json:"indexes"`
	ExpiresAt   string   `json:"expiresAt"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	UID         string   `json:"uid,omitempty"`
}

type CreateKeyResponse struct {
	CreateKeyRequest `       json:",inline"`
	Key              string `json:"key"`
}

// SearchRequest is the request url param needed for a search query.
// This struct will be converted to url param before sent.
//
// Documentation: https://www.meilisearch.com/docs/reference/api/search#search-parameters
type SearchRequest struct {
	Offset                  int64                    `json:"offset,omitempty"`
	Limit                   int64                    `json:"limit,omitempty"`
	AttributesToRetrieve    []string                 `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchOn    []string                 `json:"attributesToSearchOn,omitempty"`
	AttributesToCrop        []string                 `json:"attributesToCrop,omitempty"`
	CropLength              int64                    `json:"cropLength,omitempty"`
	CropMarker              string                   `json:"cropMarker,omitempty"`
	AttributesToHighlight   []string                 `json:"attributesToHighlight,omitempty"`
	HighlightPreTag         string                   `json:"highlightPreTag,omitempty"`
	HighlightPostTag        string                   `json:"highlightPostTag,omitempty"`
	MatchingStrategy        MatchingStrategy         `json:"matchingStrategy,omitempty"`
	Filter                  interface{}              `json:"filter,omitempty"`
	ShowMatchesPosition     bool                     `json:"showMatchesPosition,omitempty"`
	ShowRankingScore        bool                     `json:"showRankingScore,omitempty"`
	ShowRankingScoreDetails bool                     `json:"showRankingScoreDetails,omitempty"`
	Facets                  []string                 `json:"facets,omitempty"`
	Sort                    []string                 `json:"sort,omitempty"`
	Vector                  []float32                `json:"vector,omitempty"`
	HitsPerPage             int64                    `json:"hitsPerPage,omitempty"`
	Page                    int64                    `json:"page,omitempty"`
	IndexUID                string                   `json:"indexUid,omitempty"`
	Query                   string                   `json:"q"`
	Distinct                string                   `json:"distinct,omitempty"`
	Hybrid                  *SearchRequestHybrid     `json:"hybrid"`
	RetrieveVectors         bool                     `json:"retrieveVectors,omitempty"`
	RankingScoreThreshold   float64                  `json:"rankingScoreThreshold,omitempty"`
	FederationOptions       *SearchFederationOptions `json:"federationOptions,omitempty"`
	Locales                 []string                 `json:"locales,omitempty"`
	Media                   map[string]any           `json:"media,omitempty"`
}

type SearchRequestHybrid struct {
	SemanticRatio float64 `json:"semanticRatio,omitempty"`
	Embedder      string  `json:"embedder"`
}

type SearchFederationOptions struct {
	Weight float64 `json:"weight,omitempty"`
	Remote string  `json:"remote,omitempty"`
}

// SearchResponse is the response body for search method
type SearchResponse struct {
	Hits               Hits            `json:"hits"`
	EstimatedTotalHits int64           `json:"estimatedTotalHits,omitempty"`
	Offset             int64           `json:"offset,omitempty"`
	Limit              int64           `json:"limit,omitempty"`
	ProcessingTimeMs   int64           `json:"processingTimeMs"`
	Query              string          `json:"query"`
	FacetDistribution  json.RawMessage `json:"facetDistribution,omitempty"`
	TotalHits          int64           `json:"totalHits,omitempty"`
	HitsPerPage        int64           `json:"hitsPerPage,omitempty"`
	Page               int64           `json:"page,omitempty"`
	TotalPages         int64           `json:"totalPages,omitempty"`
	FacetStats         json.RawMessage `json:"facetStats,omitempty"`
	IndexUID           string          `json:"indexUid,omitempty"`
	QueryVector        *[]float32      `json:"queryVector,omitempty"`
}
