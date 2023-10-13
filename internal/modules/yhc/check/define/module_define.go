package define

const (
	// level 1 modules
	MODULE_OVERVIEW ModuleName = "overview"
	MODULE_HOST     ModuleName = "host_check"
	MODULE_YASDB    ModuleName = "yasdb_check"
	MODULE_OBJECT   ModuleName = "object_check"
	MODULE_SECURITY ModuleName = "security_check"
	MODULE_LOG      ModuleName = "log_analysis"
	MODULE_CUSTOM   ModuleName = "custom_check"

	// the followings are level 2 modules
	// parent module: MN_OVERVIEW
	MODULE_OVERVIEW_HOST  ModuleName = "overview_host"
	MODULE_OVERVIEW_YASDB ModuleName = "overview_yasdb"

	// parent module: MN_HOST
	MODULE_HOST_WORKLOAD ModuleName = "host_workload_check"

	// parent module: MN_YASDB
	MODULE_YASDB_STANDBY     ModuleName = "yasdb_standby_check"
	MODULE_YASDB_CONFIG      ModuleName = "yasdb_config_check"
	MODULE_YASDB_TABLESPACE  ModuleName = "yasdb_tablespace_check"
	MODULE_YASDB_CONTROLFILE ModuleName = "yasdb_controlfile_check"
	MODULE_YASDB_BACKUP      ModuleName = "yasdb_backup_check"
	MODULE_YASDB_WORKLOAD    ModuleName = "yasdb_workload_check"
	MODULE_YASDB_PERFORMANCE ModuleName = "yasdb_performance_analysis"

	// parent module: MN_OBJECT
	MODULE_OBJECT_NUMBER     ModuleName = "object_number_count"
	MODULE_OBJECT_STATUS     ModuleName = "object_status_check"
	MODULE_OBJECT_TABLE      ModuleName = "object_table_check"
	MODULE_OBJECT_CONSTRAINT ModuleName = "object_constraint_check"
	MODULE_OBJECT_INDEX      ModuleName = "object_index_check"
	MODULE_OBJECT_SEQUENCE   ModuleName = "object_sequence_check"
	MODULE_OBJECT_TASK       ModuleName = "object_task_check"
	MODULE_OBJECT_PACKAGE    ModuleName = "object_package_check"

	MODULE_SECURITY_LOGIN      ModuleName = "security_login_config"
	MODULE_SECURITY_PERMISSION ModuleName = "security_permission_check"
	MODULE_SECURITY_AUDIT      ModuleName = "security_audit_check"

	// parent module: MN_LOG
	MODULE_LOG_RUN   ModuleName = "log_run_analysis"
	MODULE_LOG_REDO  ModuleName = "log_redo_analysis"
	MODULE_LOG_UNDO  ModuleName = "log_undo_analysis"
	MODULE_LOG_ERROR ModuleName = "log_error_analysis"

	MODULE_CUSTOM_BASH ModuleName = "custom_check_bash"
	MODULE_CUSTOM_SQL  ModuleName = "custom_check_sql"
)

type ModuleName string

const (
	METRIC_YASDB_INSTANCE            MetricName = "yasdb_instance"
	METRIC_YASDB_DATABASE            MetricName = "yasdb_database"
	METRIC_YASDB_FILE_PERMISSION     MetricName = "yasdb_file_permission"
	METRIC_YASDB_LISTEN_ADDR         MetricName = "yasdb_listen_address"
	METRIC_YASDB_OS_AUTH             MetricName = "yasdb_os_auth"
	METRIC_HOST_INFO                 MetricName = "host_info"
	METRIC_HOST_FIREWALLD            MetricName = "host_firewalld"
	METRIC_HOST_IPTABLES             MetricName = "host_iptables"
	METRIC_HOST_CPU_INFO             MetricName = "host_cpu_info"
	METRIC_HOST_DISK_INFO            MetricName = "host_disk_info"
	METRIC_HOST_DISK_BLOCK_INFO      MetricName = "host_disk_block_info"
	METRIC_HOST_BIOS_INFO            MetricName = "host_bios_info"
	METRIC_HOST_MEMORY_INFO          MetricName = "host_memory_info"
	METRIC_HOST_NETWORK_INFO         MetricName = "host_network_info"
	METRIC_HOST_HISTORY_CPU_USAGE    MetricName = "host_history_cpu_usage"
	METRIC_HOST_CURRENT_CPU_USAGE    MetricName = "host_current_cpu_usage"
	METRIC_HOST_HISTORY_DISK_IO      MetricName = "host_history_disk_io"
	METRIC_HOST_CURRENT_DISK_IO      MetricName = "host_current_disk_io"
	METRIC_HOST_HISTORY_MEMORY_USAGE MetricName = "host_history_memory_usage"
	METRIC_HOST_CURRENT_MEMORY_USAGE MetricName = "host_current_memory_usage"
	METRIC_HOST_HISTORY_NETWORK_IO   MetricName = "host_history_network_io"
	METRIC_HOST_CURRENT_NETWORK_IO   MetricName = "host_current_network_io"
	METRIC_YASDB_REPLICATION_STATUS  MetricName = "yasdb_replication_status"
	METRIC_YASDB_PARAMETER           MetricName = "yasdb_parameter"
	METRIC_YASDB_TABLESPACE          MetricName = "yasdb_tablespace"
	METRIC_YASDB_CONTROLFILE_COUNT   MetricName = "yasdb_controlfile_count"
	METRIC_YASDB_CONTROLFILE         MetricName = "yasdb_controlfile"
	METRIC_YASDB_DATAFILE            MetricName = "yasdb_datafile"
	METRIC_YASDB_SESSION             MetricName = "yasdb_session"
	METRIC_YASDB_WAIT_EVENT          MetricName = "yasdb_wait_event"
	METRIC_YASDB_OBJECT_COUNT        MetricName = "yasdb_object_count"
	METRIC_YASDB_OBJECT_OWNER        MetricName = "yasdb_object_owner"
	METRIC_YASDB_OBJECT_TABLESPACE   MetricName = "yasdb_object_tablespace"
	METRIC_YASDB_INDEX_BLEVEL        MetricName = "yasdb_index_blevel"
	METRIC_YASDB_INDEX_COLUMN        MetricName = "yasdb_index_column"
	METRIC_YASDB_INDEX_INVISIBLE     MetricName = "yasdb_index_invisible"
	METRIC_YASDB_REDO_LOG            MetricName = "yasdb_redo_log"
	METRIC_YASDB_REDO_LOG_COUNT      MetricName = "yasdb_redo_log_count"
	METRIC_YASDB_RUN_LOG_ERROR       MetricName = "yasdb_run_log_error"
)

type MetricName string
