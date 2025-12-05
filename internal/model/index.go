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
	Segments  IndexSegments   `json:"segments"` // 新增：分段信息
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

// IndexSegments 索引分段信息（新增）
type IndexSegments struct {
	Count int `json:"count"`
	MemoryInBytes int64 `json:"memory_in_bytes"`
	TermsMemoryInBytes int64 `json:"terms_memory_in_bytes"`
	StoredFieldsMemoryInBytes int64 `json:"stored_fields_memory_in_bytes"`
	TermVectorsMemoryInBytes int64 `json:"term_vectors_memory_in_bytes"`
	NormsMemoryInBytes int64 `json:"norms_memory_in_bytes"`
	PointsMemoryInBytes int64 `json:"points_memory_in_bytes"`
	DocValuesMemoryInBytes int64 `json:"doc_values_memory_in_bytes"`
	IndexWriterMemoryInBytes int64 `json:"index_writer_memory_in_bytes"`
	VersionMapMemoryInBytes int64 `json:"version_map_memory_in_bytes"`
	FixedBitSetMemoryInBytes int64 `json:"fixed_bit_set_memory_in_bytes"`
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
