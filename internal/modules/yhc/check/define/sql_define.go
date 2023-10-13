package define

const (
	SQL_QUERY_CONTROLFILE                         = "select  id, name, bytes/1024/1024 as MBytes from v$controlfile;"
	SQL_QUERY_CONTROLFILE_COUNT                   = "select count(*) as total from v$controlfile;"
	SQL_QUERY_DATAFILE                            = "select * from dba_data_files;"
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
	/**对象检查**/
	SQL_QUERY_OVERSIZED_INDEX          = `SELECT ind.OWNER AS ind_owner,ind.SEGMENT_NAME AS ind_name,ind.SEGMENT_TYPE as IND_SEGMENT_TYPE ,tab.SEGMENT_TYPE as TAB_SEGMENT_TYPE,tab.OWNER AS tab_owner ,tab.SEGMENT_NAME AS tab_name,ind.BYTES AS ind_bytes,tab.BYTES AS tab_bytes FROM DBA_SEGMENTS ind,DBA_SEGMENTS tab,DBA_INDEXES di WHERE IND.SEGMENT_TYPE IN ('INDEX','INDEX PARTITION') AND tab.SEGMENT_TYPE IN ('TABLE','TABLE PARTITION') AND ind.OWNER =di.OWNER AND ind.SEGMENT_NAME =di.INDEX_NAME AND tab.OWNER =di.TABLE_OWNER AND tab.SEGMENT_NAME =di.TABLE_NAME AND ind.BYTES > tab.BYTES;`
	SQL_QUERY_TABLE_INDEX_NOT_TOGETHER = `SELECT OWNER,INDEX_NAME ,TABLE_OWNER ,TABLE_NAME FROM dba_indexes WHERE OWNER <> TABLE_OWNER;`
	SQL_QUERY_NO_AVAILABLE_VALUE       = `SELECT SEQUENCE_OWNER ,SEQUENCE_NAME ,MIN_VALUE / MAX_VALUE * 100 as USED_RATE FROM DBA_SEQUENCES ds WHERE MIN_VALUE / MAX_VALUE > 7/10;`
	SQL_QUERY_RUNNING_JOB              = `select OWNER ,JOB_NAME ,JOB_STYLE ,JOB_CREATOR ,JOB_ACTION  from DBA_SCHEDULER_JOBS where STATE='RUNNING';`
	SQL_NO_PACKAGE_PACKAGE_BODY        = `SELECT OWNER ,NAME FROM (SELECT OWNER ,NAME,LISTAGG("TYPE",'-') AS str FROM DBA_SOURCE GROUP by OWNER ,NAME) WHERE str<>'PACKAGE-PACKAGE BODY'`
	/**安全检查**/
	SQL_QUERY_PASSWORD_STRENGTH                            = `SELECT value FROM x$parameter WHERE name ='_CHECK_PASSWORD_COMPLEXITY';`
	SQL_QUERY_MAXIMUM_LOGIN_ATTEMPTS                       = `select PROFILE,RESOURCE_NAME ,RESOURCE_TYPE, LIMIT from DBA_PROFILES where PROFILE<>'DEFAULT' and RESOURCE_NAME='FAILED_LOGIN_ATTEMPTS' and LIMIT<>'UNLIMITED';`
	SQL_QUERY_USER_NO_OPEN                                 = `select username,ACCOUNT_STATUS from dba_users where ACCOUNT_STATUS!='OPEN';`
	SQL_QUERY_USER_WITH_SYSTEM_TABLE_PRIVILEGES            = `select GRANTEE from DBA_TAB_PRIVS where OWNER='SYS' and TYPE='TABLE' and GRANTEE in (select username from dba_users);`
	SQL_QUERY_ALL_USERS_WITH_DBA_ROLE                      = `select GRANTEE from dba_role_privs where GRANTED_ROLE='DBA';`
	SQL_QUERY_ALL_USERS_ALL_PRIVILEGE_OR_SYSTEM_PRIVILEGES = `select GRANTEE from dba_sys_privs where PRIVILEGE='ALL PRIVILEGES' AND GRANTEE IN ( SELECT USERNAME FROM DBA_USERS);`
	SQL_QUERY_USERS_USE_SYSTEM_TABLESPACE                  = `select username,default_tablespace from dba_users;`
	SQL_QUERY_AUDIT_CLEANUP_TASK                           = `select AUDIT_TRAIL,LAST_ARCHIVE_TS,DATABASE_ID from DBA_AUDIT_MGMT_LAST_ARCH_TS;`
	SQL_QUERY_AUDIT_FILE_SIZE                              = `select segment_name ,bytes/1024/1024/1024 as size_gb from dba_segments where segment_name like 'AUD$';`
	/**日志分析**/
	SQL_QUERY_UNDO_LOG_SIZE                 = `SELECT  a.USED_UBLK * b.value /1024/1024/1044 AS SIZE_GB, XID from V$TRANSACTION as a ,(SELECT VALUE FROM V$PARAMETER WHERE NAME='DB_BLOCK_SIZE') AS B ;`
	SQL_QUERY_UNDO_LOG_TOTAL_BLOCK          = `SELECT  SUM(USED_UBLK) as TOTAL_BLOCK from V$TRANSACTION ;`
	SQL_QUERY_UNDO_LOG_RUNNING_TRANSACTIONS = `SELECT XID, SID,XRMID,XEXT, XNODE,XSN,STATUS,RESIDUAL, USED_UBLK, FIRST_UBAFIL,FIRST_UBABLK,FIRST_UBAVER ,FIRST_UBAREC,LAST_UBAFIL,LAST_UBABLK, PTX_XID, START_DATE,ISOLATION_LEVEL from V$TRANSACTION ;`
)
