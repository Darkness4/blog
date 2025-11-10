package meilisearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/Darkness4/blog/web/gen/index"
	"github.com/rs/zerolog/log"
)

const (
	documentsEndpoint = "%s/indexes/%s/documents"
	tasksEndpoint     = "%s/tasks/%d"
	searchEndpoint    = "%s/indexes/%s/search"
)

type Client struct {
	*http.Client
	URL       string
	MasterKey string
	IndexUID  string
}

func NewClient(client *http.Client, url, masterKey, indexUID string) *Client {
	if client == nil {
		panic("client is nil")
	}
	return &Client{
		Client:    client,
		URL:       url,
		MasterKey: masterKey,
		IndexUID:  indexUID,
	}
}

func (c *Client) BuildIndex(ctx context.Context, index [][]index.Index) error {
	records := slices.Collect(IndexToRecords(index))

	log.Info().Int("records", len(records)).Msg("building index")

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&records); err != nil {
		return fmt.Errorf("failed to encode records to json: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf(documentsEndpoint, c.URL, c.IndexUID),
		buf,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.MasterKey)
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("primaryKey", "objectID")
	req.URL.RawQuery = q.Encode()
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 202 {
		body, _ := io.ReadAll(res.Body)
		err := fmt.Errorf("failed to build index: %v", res.Status)
		log.Err(err).Str("body", string(body)).Msg("failed to build index")
		return err
	}

	log.Info().Msg("index built")

	var parsed SubmittedTaskResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return c.waitForSuccess(ctx, parsed.TaskUID)
}

func (c *Client) waitForSuccess(ctx context.Context, uid int) (err error) {
	var task GetTaskResponse
	for task.Status != TaskStatusSucceeded && task.Status != TaskStatusCanceled {
		task, err = c.getTask(ctx, uid)
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}
		switch task.Status {
		case TaskStatusEnqueued, TaskStatusProcessing:
		case TaskStatusCanceled:
			log.Warn().Msg("index build canceled")
		case TaskStatusFailed:
			log.Error().Str("message", task.Error.Message).Msg("index build failed")
			return fmt.Errorf("index build failed: %v", task.Error.Message)
		case TaskStatusSucceeded:
			log.Info().Msg("index build succeeded")
		}
	}
	return nil
}

func (c *Client) getTask(ctx context.Context, uid int) (GetTaskResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(tasksEndpoint, c.URL, uid),
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.MasterKey)
	res, err := c.Do(req)
	if err != nil {
		return GetTaskResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		err := fmt.Errorf("failed to get task: %v", res.Status)
		log.Err(err).Str("body", string(body)).Msg("failed to get task")
		return GetTaskResponse{}, err
	}

	var parsed GetTaskResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return GetTaskResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return parsed, nil
}

func (c *Client) ClearIndex(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf(documentsEndpoint, c.URL, c.IndexUID),
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.MasterKey)
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 202 {
		body, _ := io.ReadAll(res.Body)
		err := fmt.Errorf("failed to clear index: %v", res.Status)
		log.Err(err).Str("body", string(body)).Msg("failed to clear index")
		return err
	}

	var parsed SubmittedTaskResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return c.waitForSuccess(ctx, parsed.TaskUID)
}

func (c *Client) Search(ctx context.Context, query string) (SearchResponse, error) {
	reqBody := SearchRequest{
		AttributesToHighlight: []string{"*"},
		AttributesToCrop:      []string{"content"},
		CropLength:            30,
		Query:                 query,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&reqBody); err != nil {
		return SearchResponse{}, fmt.Errorf("failed to encode Search body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf(searchEndpoint, c.URL, c.IndexUID),
		&buf,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.MasterKey)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		err := fmt.Errorf("failed to search: %v", res.Status)
		log.Err(err).Str("body", string(body)).Msg("failed to search")
		return SearchResponse{}, err
	}

	var parsed SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return SearchResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return parsed, nil
}
