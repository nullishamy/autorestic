package main

import (
	"encoding/json"
)

func UnmarshalStatusMessage(data []byte) (StatusMessage, error) {
	var r StatusMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StatusMessage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StatusMessage struct {
	MessageType  string   `json:"message_type"`
	PercentDone  float64  `json:"percent_done"`
	TotalFiles   int64    `json:"total_files"`
	FilesDone    int64    `json:"files_done"`
	TotalBytes   int64    `json:"total_bytes"`
	BytesDone    int64    `json:"bytes_done"`
	CurrentFiles []string `json:"current_files"`
}

func UnmarshalSummaryMessage(data []byte) (SummaryMessage, error) {
	var r SummaryMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SummaryMessage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SummaryMessage struct {
	MessageType         string  `json:"message_type"`
	FilesNew            int64   `json:"files_new"`
	FilesChanged        int64   `json:"files_changed"`
	FilesUnmodified     int64   `json:"files_unmodified"`
	DirsNew             int64   `json:"dirs_new"`
	DirsChanged         int64   `json:"dirs_changed"`
	DirsUnmodified      int64   `json:"dirs_unmodified"`
	DataBlobs           int64   `json:"data_blobs"`
	TreeBlobs           int64   `json:"tree_blobs"`
	DataAdded           int64   `json:"data_added"`
	DataAddedPacked     int64   `json:"data_added_packed"`
	TotalFilesProcessed int64   `json:"total_files_processed"`
	TotalBytesProcessed int64   `json:"total_bytes_processed"`
	TotalDuration       float64 `json:"total_duration"`
	SnapshotID          string  `json:"snapshot_id"`
}

func UnmarshalStatsMessage(data []byte) (StatsMessage, error) {
	var r StatsMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StatsMessage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StatsMessage struct {
	TotalSize      int64 `json:"total_size"`
	TotalFileCount int64 `json:"total_file_count"`
	SnapshotsCount int64 `json:"snapshots_count"`
}

func UnmarshalUntypedMessage(data []byte) (UntypedMessage, error) {
	var r UntypedMessage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UntypedMessage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type UntypedMessage struct {
	MessageType string `json:"message_type"`
}
