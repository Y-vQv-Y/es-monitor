package model

// IndexStats 索引统计
type IndexStats struct {
	Indices map[string]IndexStat `json:"indices"`
}

// IndexStat 单个索引统计
type IndexStat struct {
	Primaries IndexShardStats `json:"primaries"`
	Total     IndexShardStats `json:"total"`
	UUID      string          `json:"uuid"`
	Health    string          `json:"health"`
	Status    string          `json:"status"`
}

// IndexShardStats 索引分片统计
type IndexShardStats struct {
	Docs struct {
		Count   int `json:"count"`
		Deleted int `json:"deleted"`
	} `json:"docs"`
	Store struct {
		SizeInBytes int64 `json:"size_in_bytes"`
	} `json:"store"`
	Indexing struct {
		IndexTotal        int   `json:"index_total"`
		IndexTimeInMillis int64 `json:"index_time_in_millis"`
		IndexCurrent      int   `json:"index_current"`
	} `json:"indexing"`
	Search struct {
		QueryTotal        int   `json:"query_total"`
		QueryTimeInMillis int64 `json:"query_time_in_millis"`
		QueryCurrent      int   `json:"query_current"`
	} `json:"search"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}
