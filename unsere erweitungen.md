* mm-less druid ist erlaubt 
* DruidHistoricalFreeSpace: prometheus.NewDesc("druid_historical_free_space",
			"Freespace of all historicals per node",
			[]string{"host", "server_type", "ip", "pod"}, nil),

		DruidHistoricalUsagePercent: prometheus.NewDesc("druid_historical_usage_percent",
			"Usage of all historicals per node in Percent",
			[]string{"host", "server_type", "ip", "pod"}, nil),

		DruidHistoricalUsageAbsolute: prometheus.NewDesc("druid_historical_usage_absolute",
			"Absolute Usage of all historicals per node",
			[]string{"host", "server_type", "ip", "pod"}, nil),