package model

// NodeStats 节点统计信息
type NodeStats struct {
	Nodes map[string]NodeStat `json:"nodes"`
}

// NodeStat 单个节点统计
type NodeStat struct {
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	IP        string    `json:"ip"`
	Roles     []string  `json:"roles"`
	Timestamp int64     `json:"timestamp"`
	JVM       JVMStats  `json:"jvm"`
	OS        OSStats   `json:"os"`
	Process   Process   `json:"process"`
	FS        FSStats   `json:"fs"`
	Transport Transport `json:"transport"`
	HTTP      HTTP      `json:"http"`
	Indices   Indices   `json:"indices"`
}

// JVMStats JVM 统计
type JVMStats struct {
	Timestamp      int64 `json:"timestamp"`
	UptimeInMillis int64 `json:"uptime_in_millis"`
	Mem            struct {
		HeapUsedInBytes         int64 `json:"heap_used_in_bytes"`
		HeapUsedPercent         int   `json:"heap_used_percent"`
		HeapCommittedInBytes    int64 `json:"heap_committed_in_bytes"`
		HeapMaxInBytes          int64 `json:"heap_max_in_bytes"`
		NonHeapUsedInBytes      int64 `json:"non_heap_used_in_bytes"`
		NonHeapCommittedInBytes int64 `json:"non_heap_committed_in_bytes"`
	} `json:"mem"`
	Threads struct {
		Count     int `json:"count"`
		PeakCount int `json:"peak_count"`
	} `json:"threads"`
	GC struct {
		Collectors struct {
			Young struct {
				CollectionCount        int `json:"collection_count"`
				CollectionTimeInMillis int `json:"collection_time_in_millis"`
			} `json:"young"`
			Old struct {
				CollectionCount        int `json:"collection_count"`
				CollectionTimeInMillis int `json:"collection_time_in_millis"`
			} `json:"old"`
		} `json:"collectors"`
	} `json:"gc"`
}

// OSStats 操作系统统计
type OSStats struct {
	Timestamp int64 `json:"timestamp"`
	CPU       struct {
		Percent     int `json:"percent"`
		LoadAverage struct {
			OneMinute      float64 `json:"1m"`
			FiveMinutes    float64 `json:"5m"`
			FifteenMinutes float64 `json:"15m"`
		} `json:"load_average"`
	} `json:"cpu"`
	Mem struct {
		TotalInBytes int64 `json:"total_in_bytes"`
		FreeInBytes  int64 `json:"free_in_bytes"`
		UsedInBytes  int64 `json:"used_in_bytes"`
		FreePercent  int   `json:"free_percent"`
		UsedPercent  int   `json:"used_percent"`
	} `json:"mem"`
	Swap struct {
		TotalInBytes int64 `json:"total_in_bytes"`
		FreeInBytes  int64 `json:"free_in_bytes"`
		UsedInBytes  int64 `json:"used_in_bytes"`
	} `json:"swap"`
}

// Process 进程统计
type Process struct {
	Timestamp           int64 `json:"timestamp"`
	OpenFileDescriptors int   `json:"open_file_descriptors"`
	MaxFileDescriptors  int   `json:"max_file_descriptors"`
	CPU                 struct {
		Percent       int   `json:"percent"`
		TotalInMillis int64 `json:"total_in_millis"`
	} `json:"cpu"`
	Mem struct {
		TotalVirtualInBytes int64 `json:"total_virtual_in_bytes"`
	} `json:"mem"`
}

// FSStats 文件系统统计
type FSStats struct {
	Timestamp int64 `json:"timestamp"`
	Total     struct {
		TotalInBytes     int64 `json:"total_in_bytes"`
		FreeInBytes      int64 `json:"free_in_bytes"`
		AvailableInBytes int64 `json:"available_in_bytes"`
	} `json:"total"`
	Data []struct {
		Path             string `json:"path"`
		Mount            string `json:"mount"`
		Type             string `json:"type"`
		TotalInBytes     int64  `json:"total_in_bytes"`
		FreeInBytes      int64  `json:"free_in_bytes"`
		AvailableInBytes int64  `json:"available_in_bytes"`
	} `json:"data"`
	IOStats struct {
		Total struct {
			Operations int64 `json:"operations"`
			ReadOps    int64 `json:"read_operations"`
			WriteOps   int64 `json:"write_operations"`
			ReadKB     int64 `json:"read_kilobytes"`
			WriteKB    int64 `json:"write_kilobytes"`
		} `json:"total"`
	} `json:"io_stats"`
}

// Transport 传输层统计
type Transport struct {
	ServerOpen    int   `json:"server_open"`
	RxCount       int64 `json:"rx_count"`
	RxSizeInBytes int64 `json:"rx_size_in_bytes"`
	TxCount       int64 `json:"tx_count"`
	TxSizeInBytes int64 `json:"tx_size_in_bytes"`
}

// HTTP HTTP 统计
type HTTP struct {
	CurrentOpen int `json:"current_open"`
	TotalOpened int `json:"total_opened"`
}

// Indices 索引统计
type Indices struct {
	Docs struct {
		Count   int `json:"count"`
		Deleted int `json:"deleted"`
	} `json:"docs"`
	Store struct {
		SizeInBytes             int64 `json:"size_in_bytes"`
		ReservedInBytes         int64 `json:"reserved_in_bytes"`
		TotalDataSetSizeInBytes int64 `json:"total_data_set_size_in_bytes"`
	} `json:"store"`
	Indexing struct {
		IndexTotal           int   `json:"index_total"`
		IndexTimeInMillis    int64 `json:"index_time_in_millis"`
		IndexCurrent         int   `json:"index_current"`
		IndexFailed          int   `json:"index_failed"`
		DeleteTotal          int   `json:"delete_total"`
		DeleteTimeInMillis   int64 `json:"delete_time_in_millis"`
		DeleteCurrent        int   `json:"delete_current"`
		ThrottleTimeInMillis int64 `json:"throttle_time_in_millis"`
	} `json:"indexing"`
	Search struct {
		OpenContexts       int   `json:"open_contexts"`
		QueryTotal         int   `json:"query_total"`
		QueryTimeInMillis  int64 `json:"query_time_in_millis"`
		QueryCurrent       int   `json:"query_current"`
		FetchTotal         int   `json:"fetch_total"`
		FetchTimeInMillis  int64 `json:"fetch_time_in_millis"`
		FetchCurrent       int   `json:"fetch_current"`
		ScrollTotal        int   `json:"scroll_total"`
		ScrollTimeInMillis int64 `json:"scroll_time_in_millis"`
		ScrollCurrent      int   `json:"scroll_current"`
	} `json:"search"`
	Merges struct {
		Current            int   `json:"current"`
		CurrentDocs        int   `json:"current_docs"`
		CurrentSizeInBytes int64 `json:"current_size_in_bytes"`
		Total              int   `json:"total"`
		TotalTimeInMillis  int64 `json:"total_time_in_millis"`
		TotalDocs          int   `json:"total_docs"`
		TotalSizeInBytes   int64 `json:"total_size_in_bytes"`
	} `json:"merges"`
	Refresh struct {
		Total             int   `json:"total"`
		TotalTimeInMillis int64 `json:"total_time_in_millis"`
	} `json:"refresh"`
	Flush struct {
		Total             int   `json:"total"`
		TotalTimeInMillis int64 `json:"total_time_in_millis"`
	} `json:"flush"`
}
