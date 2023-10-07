package define

const (
	SQL_QUERY_CONTROLFILE                         = "select  id, name, bytes/1024/1024 as MBytes from v$controlfile;"
	SQL_QUERY_CONTROLFILE_COUNT                   = "select count(*) as total from v$controlfile;"
	SQL_QUERY_DATABASE                            = "select database_name, status as database_status, log_mode, open_mode, database_role, protection_mode, protection_level, create_time from v$database;"
	SQL_QUERY_INDEX_BLEVEL                        = "select OWNER, INDEX_NAME, BLEVEL from dba_indexes where BLEVEL>3;"
	SQL_QUERY_INDEX_COLUMN                        = "select INDEX_OWNER, INDEX_NAME, count(*) as column_count from dba_ind_columns group by INDEX_OWNER,INDEX_NAME having count(*) > 10;"
	SQL_QUERY_INDEX_INVISIBLE                     = "select OWNER, INDEX_NAME, TABLE_OWNER, TABLE_NAME FROM dba_indexes where owner<> 'SYS' and VISIBILITY <> 'VISIBLE';"
	SQL_QUERY_INSTANCE                            = "select status as instance_status, version, startup_time from v$instance;"
	SQL_QUERY_LISTEN_ADDR                         = `select VALUE as LISTEN_ADDR from v$parameter where name = 'LISTEN_ADDR';`
	SQL_QUERY_SESSION                             = `select type from v$session`
	SQL_QUERY_TABLESPACE                          = `SELECT TABLESPACE_NAME, CONTENTS, STATUS, ALLOCATION_TYPE , TOTAL_BYTES - USER_BYTES AS USED_BYTES, TOTAL_BYTES, (TOTAL_BYTES - USER_BYTES) / TOTAL_BYTES * 100 AS USED_RATE FROM SYS.DBA_TABLESPACES;`
	SQL_QUERY_TABLESPACE_DATA_PERCENTAGE_FORMATER = `SELECT A.TABLESPACE_NAME, A.B1/B.B2*100 AS DATA_PERCENTAGE FROM 
    (SELECT TABLESPACE_NAME,SUM(BYTES) AS B1 FROM dba_segments WHERE SEGMENT_TYPE LIKE 'TABLE%%' GROUP BY TABLESPACE_NAME ) A,
    (SELECT TABLESPACE_NAME,TOTAL_BYTES AS B2 FROM DBA_TABLESPACES) B WHERE (A.TABLESPACE_NAME=B.TABLESPACE_NAME AND A.TABLESPACE_NAME ='%s');`
	SQL_QUERY_WAIT_EVENT = `SELECT count(s.WAIT_EVENT) current_waits FROM sys.v$system_event se, sys.v$session s WHERE se.EVENT = s.WAIT_EVENT
    AND se.event not in ('SQL*Net message from client',
    'SQL*Net more data from client',
    'pmon timer',
    'rdbms ipc message',
    'rdbms ipc reply',
    'smon timer');`
	SQL_QUERY_REPLICATION_STATUS = "select connection, status, peer_role, peer_addr, transport_lag, apply_lag from v$replication_status;"
	SQL_QUERY_PARAMETER          = "select name, value from v$parameter where value is not null;"
	SQL_QUERY_TOTAL_OBJECT       = "select count(*) as total_count from dba_objects;"
	SQL_QUERY_OWNER_OBJECT       = `SELECT owner, object_type, COUNT(*) AS owner_object_count FROM dba_objects
    WHERE owner NOT IN ('SYS', 'SYSTEM') AND object_type NOT LIKE 'BIN$%'
    GROUP BY owner, object_type
    ORDER BY owner,object_type;`
	SQL_QUERY_TABLESPACE_OBJECT = `SELECT tablespace_name, COUNT(*) AS tablespace_object_count FROM dba_segments
    WHERE segment_type IN ('TABLE', 'INDEX', 'VIEW', 'SEQUENCE')
    GROUP BY tablespace_name
    ORDER BY tablespace_name;`
	SQL_QUERY_LOGFILE       = "select ID, NAME, STATUS, BLOCK_SIZE, BLOCK_COUNT, USED_BLOCKS, SEQUENCE# AS SEQUENCE from v$logfile;"
	SQL_QUERY_LOGFILE_COUNT = `select count(*) as total_count, SUM(CASE WHEN STATUS = 'CURRENT' THEN 1 ELSE 0 END) AS current_count,
    SUM(CASE WHEN STATUS = 'ACTIVE' THEN 1 ELSE 0 END) AS active_count, SUM(CASE WHEN STATUS = 'INACTIVE' THEN 1 ELSE 0 END) AS inactive_count
    from v$logfile;`
)
