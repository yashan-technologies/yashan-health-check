package define

const (
	SQL_QUERY_CONTROLFILE                         = "select  id, name, bytes/1024/1024 as MBytes from v$controlfile;"
	SQL_QUERY_CONTROLFILE_COUNT                   = "select count(*) as total from v$controlfile;"
	SQL_QUERY_DATAFILE                            = "select * from dba_data_files;"
	SQL_QUERY_BACKUP_SET                          = "select RECID# as RECID, START_TIME, TYPE, decode(COMPLETION_TIME > sysdate, FALSE, TRUE) as SUCCESS from dba_backup_set;"
	SQL_QUERY_FULL_BACKUP_SET_COUNT               = "select count(*) as TOTAL from dba_backup_set where date_add(COMPLETION_TIME , INTERVAL 10 DAY) >= sysdate AND type = 'FULL';"
	SQL_QUERY_BACKUP_SET_PATH                     = "select distinct(PATH) as PATH from dba_backup_set;"
	SQL_QUERY_DATABASE                            = "select database_name, status as database_status, log_mode, open_mode, database_role, protection_mode, protection_level, create_time from v$database;"
	SQL_QUERY_INDEX_BLEVEL                        = "select OWNER, INDEX_NAME, BLEVEL from dba_indexes where BLEVEL>3;"
	SQL_QUERY_INDEX_COLUMN                        = "select INDEX_OWNER, INDEX_NAME, count(*) as column_count from dba_ind_columns group by INDEX_OWNER,INDEX_NAME having count(*) > 10;"
	SQL_QUERY_INDEX_INVISIBLE                     = "select OWNER, INDEX_NAME, TABLE_OWNER, TABLE_NAME FROM dba_indexes where owner<> 'SYS' and VISIBILITY <> 'VISIBLE';"
	SQL_QUERY_INSTANCE                            = "select status as instance_status, version, startup_time from v$instance;"
	SQL_QUERY_LISTEN_ADDR                         = `select VALUE as LISTEN_ADDR from v$parameter where name = 'LISTEN_ADDR';`
	SQL_QUERY_SESSION                             = `select type from v$session`
	SQL_QUERY_SHARE_POOL                          = `select NAME, BYTES from v$sgastat WHERE POOL='SHARE POOL';`
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
	SQL_QUERY_VM_SWAP_RATE              = `SELECT t1.SWAPPED_OUT_BLOCKS / t2.value AS RATE FROM ( SELECT SWAPPED_OUT_BLOCKS FROM v$vm ) t1, ( SELECT value FROM V$SYSSTAT WHERE NAME = 'VM ALLOC' ) t2;`
	SQL_QUERY_YASDB_TOP_SQL_BY_CPU_TIME = `SELECT round(CPU_TIME / 1000, 2) AS CPU_TIME, EXECUTIONS
        , round(ELAPSED_TIME / 1000, 2) AS ALL_ELAPSED_TIME
        , round(ELAPSED_TIME / 1000 / EXECUTIONS, 2) AS AVG_TIME
        , to_char(LAST_ACTIVE_TIME, 'YYYY-MM-DD HH24:MI:SS') AS LAST_TIME, SQL_ID
        , SQL_TEXT
    FROM v$sqlarea
    WHERE EXECUTIONS > 0
    ORDER BY round(ELAPSED_TIME / 1000 / EXECUTIONS, 2) DESC
    LIMIT 10;`
	SQL_QUERY_YASDB_TOP_SQL_BY_BUFFER_GETS = `SELECT BUFFER_GETS, EXECUTIONS
        , round(BUFFER_GETS / EXECUTIONS, 2) AS GETS_PER_EXEC
        , round(ELAPSED_TIME / 1000, 2) AS ALL_ELAPSED_TIME
        , to_char(LAST_ACTIVE_TIME, 'YYYY-MM-DD HH24:MI:SS') AS LAST_TIME, SQL_ID
        , SQL_TEXT
    FROM v$sqlarea
    WHERE EXECUTIONS > 0
    ORDER BY BUFFER_GETS DESC
    LIMIT 10;`
	SQL_QUERY_YASDB_TOP_SQL_BY_DISK_READS = `SELECT DISK_READS, EXECUTIONS
        , round(DISK_READS / EXECUTIONS, 2) AS READS_PER_EXEC
        , round(ELAPSED_TIME / 1000, 2) AS ALL_ELAPSED_TIME
        , to_char(LAST_ACTIVE_TIME, 'YYYY-MM-DD HH24:MI:SS') AS LAST_TIME, SQL_ID
        , SQL_TEXT
    FROM v$sqlarea
    WHERE EXECUTIONS > 0
    ORDER BY DISK_READS DESC
    LIMIT 10;`
	SQL_QUERY_YASDB_TOP_SQL_BY_PARSE_CALLS = `SELECT PARSE_CALLS, EXECUTIONS
        , round(PARSE_CALLS / EXECUTIONS, 2) AS CALLS_PER_EXEC
        , round(ELAPSED_TIME / 1000, 2) AS ALL_ELAPSED_TIME
        , to_char(LAST_ACTIVE_TIME, 'YYYY-MM-DD HH24:MI:SS') AS LAST_TIME, SQL_ID
        , SQL_TEXT
    FROM v$sqlarea
    WHERE EXECUTIONS > 0
    ORDER BY round(ELAPSED_TIME / 1000 / EXECUTIONS, 2) DESC
    LIMIT 10;`
	SQL_QUERY_HIGH_FREQUENCY_SQL = `select SQL_ID, SQL_TEXT, PLSQL_EXEC_TIME, EXECUTIONS from v$sql where EXECUTIONS >= 10000`
	SQL_QUERY_HISTORY_DB_TIME    = `
        WITH dbinfo AS (
            SELECT DISTINCT dbid
            FROM SYS.wRM$_database_instance
            LIMIT 1
        ), 
        t1 AS (
            SELECT snap_id, value
            FROM SYS.wrh$_sysstat, dbinfo
            WHERE SYS.wrh$_sysstat.dbid = dbinfo.dbid
                AND stat_id = 604
        ), 
        t2 AS (
            SELECT snap_id, begin_interval_time + (end_interval_time - begin_interval_time) / 2 AS snap_time
            FROM SYS.wrm$_snapshot, dbinfo
            WHERE SYS.wrm$_snapshot.dbid = dbinfo.dbid
        ), 
        t3 AS (
            SELECT t2_1.snap_time AS prev_snap_time, t2_2.snap_time AS current_snap_time, t1_1.value AS prev_value, t1_2.value AS current_value
            FROM t2 t2_1
                JOIN t2 t2_2 ON t2_1.snap_id + 1 = t2_2.snap_id
                JOIN t1 t1_1 ON t2_1.snap_id = t1_1.snap_id
                JOIN t1 t1_2 ON t2_2.snap_id = t1_2.snap_id
        ), 
        t4 AS (
            SELECT prev_snap_time + (current_snap_time - prev_snap_time) / 2 AS snap_time
                , current_value - prev_value AS db_time_ms
            FROM t3
        )
    SELECT to_char(t4.snap_time, 'YYYY-MM-DD HH24:MI:SS') AS snap_time, t4.db_time_ms
    FROM t4
    ORDER BY snap_time;`
	SQL_QUERY_HISTORY_BUFFER_HIT_RATE = `
    WITH dbinfo AS (
            SELECT DISTINCT dbid
            FROM SYS.wRM$_database_instance
            LIMIT 1
        ), 
        dbstat AS (
            SELECT snap_id, value, stat_id
            FROM SYS.wrh$_sysstat, dbinfo
            WHERE SYS.wrh$_sysstat.dbid = dbinfo.dbid
        ), 
        t1 AS (
            SELECT snap_id, value AS b_cr_get
            FROM dbstat
            WHERE stat_id = 120
        ), 
        t2 AS (
            SELECT snap_id, value AS b_buf_get
            FROM dbstat
            WHERE stat_id = 121
        ), 
        t3 AS (
            SELECT snap_id, value AS e_phy_read
            FROM dbstat
            WHERE stat_id = 131
        ), 
        t4 AS (
            SELECT t1.snap_id
                , (t1.b_cr_get + t2.b_buf_get) / (t1.b_cr_get + t2.b_buf_get + t3.e_phy_read) * 100 AS hit_rate
            FROM t1
                JOIN t2 ON t1.snap_id = t2.snap_id
                JOIN t3 ON t1.snap_id = t3.snap_id
        ), 
        t5 AS (
            SELECT snap_id, begin_interval_time + (end_interval_time - begin_interval_time) / 2 AS snap_time
            FROM SYS.wrm$_snapshot, dbinfo
            WHERE SYS.wrm$_snapshot.dbid = dbinfo.dbid
        )
    SELECT to_char(t5.snap_time, 'YYYY-MM-DD HH24:MI:SS') AS snap_time, t4.hit_rate
    FROM t4
        JOIN t5 ON t4.snap_id = t5.snap_id
    ORDER BY t5.snap_time;`
	SQL_QUERY_BUFFER_HIT_RATE          = `select (sum(decode(NAME, 'BUFFER GETS', VALUE, 0)) + sum(decode(NAME, 'BUFFER CR GETS', VALUE, 0)) - sum(decode(NAME, 'DISK READS', VALUE, 0))) / (sum(decode(NAME, 'BUFFER GETS', VALUE, 0)) + sum(decode(NAME, 'BUFFER CR GETS', VALUE, 0))) * 100 AS HIT_RATE FROM v$sysstat;`
	SQL_QUERY_TABLE_LOCK_WAIT          = `select count(*) as TOTAL from v$lock lo where REQUEST in ('TS','TX');`
	SQL_QUERY_ROW_LOCK_WAIT            = `select count(*) as TOTAL from v$lock lo where REQUEST in ('ROW');`
	SQL_QUERY_LONG_RUNNING_TRANSACTION = `select t.XID, to_char(t.START_DATE, 'yyyy-mm-dd hh24:mi:ss') as START_DATE, t.STATUS , t.RESIDUAL, s.USERNAME, t.SID, t.USED_UBLK from v$transaction t, v$session s where t.START_DATE > sysdate - 3 / 24 and t.SID = s.SID;`
	SQL_QUERY_REPLICATION_STATUS       = "select connection, status, peer_role, peer_addr, transport_lag, apply_lag from v$replication_status;"
	SQL_QUERY_PARAMETER                = "select name, value from v$parameter where value is not null;"
	SQL_QUERY_TOTAL_OBJECT             = "select count(*) as total_count from dba_objects;"
	SQL_QUERY_OWNER_OBJECT             = `SELECT owner, object_type, COUNT(*) AS owner_object_count FROM dba_objects
    WHERE owner NOT IN ('SYS', 'SYSTEM') AND object_type NOT LIKE 'BIN$%'
    GROUP BY owner, object_type
    ORDER BY owner,object_type;`
	SQL_QUERY_TABLESPACE_OBJECT = `SELECT tablespace_name, COUNT(*) AS tablespace_object_count FROM dba_segments
    WHERE segment_type IN ('TABLE', 'INDEX', 'VIEW', 'SEQUENCE')
    GROUP BY tablespace_name
    ORDER BY tablespace_name;`
	SQL_QUERY_INVALID_OBJECT                                = `select OBJECT_ID, OWNER, OBJECT_NAME, OBJECT_TYPE, STATUS from dba_objects where STATUS = 'INVALID';`
	SQL_QUERY_INVISIBLE_INDEX                               = `select OWNER,INDEX_NAME,VISIBILITY from dba_indexes where VISIBILITY !='VISIBLE';`
	SQL_QUERY_DISABLED_CONSTRAINT                           = `select OWNER,CONSTRAINT_NAME,CONSTRAINT_TYPE,STATUS from dba_constraints where STATUS ='DISABLED';`
	SQL_QUERY_TABLE_WITH_TOO_MUCH_COLUMNS                   = `select OWNER,TABLE_NAME,count(*) as COLUMN_COUNT from dba_tab_cols group by OWNER,TABLE_NAME having count(*)>80;`
	SQL_QUERY_TABLE_WITH_TOO_MUCH_INDEXES                   = `select TABLE_OWNER,TABLE_NAME,count(*) as INDEX_COUNT from dba_indexes group by TABLE_OWNER,TABLE_NAME having count(*) >8;`
	SQL_QUERY_PARTITIONED_TABLE_WITHOUT_PARTITIONED_INDEXES = `
    SELECT b.OWNER ,b.name,a.PARTITIONING_TYPE ,b.tab_cols from
    (SELECT owner,TABLE_NAME,PARTITIONING_TYPE FROM DBA_PART_TABLES) a,
    (SELECT OWNER ,NAME ,LISTAGG(COLUMN_NAME,',') WITHIN group(ORDER BY COLUMN_POSITION) AS tab_cols FROM DBA_PART_KEY_COLUMNS WHERE OBJECT_TYPE ='TABLE' GROUP BY OWNER ,NAME ) b
    WHERE a.OWNER = b.owner AND a.TABLE_NAME =b.name AND a.owner<>'SYS'
    minus
    (
    WITH t1 AS
    (
    SELECT b.OWNER ,b.name,a.PARTITIONING_TYPE ,b.tab_cols
    from
    (SELECT owner,TABLE_NAME,PARTITIONING_TYPE FROM DBA_PART_TABLES) a,
    (SELECT OWNER ,NAME ,LISTAGG(COLUMN_NAME,',') WITHIN group(ORDER BY COLUMN_POSITION) AS tab_cols FROM DBA_PART_KEY_COLUMNS WHERE OBJECT_TYPE ='TABLE' GROUP BY OWNER ,NAME ) b
    WHERE a.OWNER = b.owner AND a.TABLE_NAME =b.name
    ),
    t2 AS
    (
    SELECT t2.INDEX_OWNER ,t2.INDEX_NAME,t2.TABLE_OWNER ,t2.TABLE_NAME,t2.ind_cols
    FROM
    (SELECT OWNER ,INDEX_NAME ,TABLE_OWNER ,TABLE_NAME FROM DBA_INDEXES WHERE PARTITIONED ='Y') t1,
    (SELECT INDEX_OWNER ,INDEX_NAME,TABLE_OWNER ,TABLE_NAME,LISTAGG(COLUMN_NAME,',') WITHIN group(ORDER BY COLUMN_POSITION) AS ind_cols
    FROM DBA_IND_COLUMNS WHERE INDEX_OWNER <>'SYS1' GROUP BY INDEX_OWNER ,INDEX_NAME,TABLE_OWNER ,TABLE_NAME) t2
    WHERE t1.OWNER=t2.INDEX_OWNER AND t1.INDEX_NAME=t2.INDEX_NAME AND t1.TABLE_OWNER=t2.TABLE_OWNER AND t1.TABLE_NAME=t2.TABLE_NAME
    )
    SELECT t1.OWNER ,t1.name,t1.PARTITIONING_TYPE ,t1.tab_cols
    FROM t1,t2
    WHERE t1.OWNER=t2.TABLE_OWNER AND t1.name=t2.TABLE_NAME AND t1.tab_cols=t2.ind_cols
    );`
	SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TAB_COLUMNS = `
    SELECT a.OWNER, a.TABLE_NAME
    FROM (
        SELECT OWNER, TABLE_NAME, MAX(DATA_LENGTH) AS MAX_DL
        FROM DBA_TAB_COLUMNS
        WHERE DATA_TYPE NOT LIKE '%LOB'
        GROUP BY OWNER, TABLE_NAME
    ) a, (
            SELECT VALUE
            FROM v$parameter
            WHERE NAME = 'DB_BLOCK_SIZE'
        ) b
    WHERE a.max_dl > b.value;`
	SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TABLES                 = "select OWNER,TABLE_NAME from dba_tables where NUM_ROWS < BLOCKS AND NUM_ROWS > 0;"
	SQL_QUERY_YASDB_PARTITIONED_TABLE_WITH_NUMBER_OF_HASH_PARTITIONS_IS_NOT_A_POWER_OF_TWO = "select OWNER,TABLE_NAME,PARTITIONING_TYPE,PARTITION_COUNT from dba_part_tables where PARTITIONING_TYPE ='HASH' and abs(floor(log(2, PARTITION_COUNT)))!=log(2, PARTITION_COUNT) or log(2, PARTITION_COUNT)='Nan';"
	SQL_QUERY_YASDB_FOREIGN_KEYS_WITHOUT_INDEXES                                           = `WITH t1 AS
    (SELECT owner,CONSTRAINT_NAME,TABLE_NAME,LISTAGG(COLUMN_NAME,',') WITHIN group(ORDER BY posi) AS col_lst FROM
    (SELECT b.OWNER ,b.CONSTRAINT_NAME ,b.TABLE_NAME ,b.COLUMN_NAME ,b.posi FROM
    (SELECT OWNER ,CONSTRAINT_NAME ,TABLE_NAME FROM DBA_CONSTRAINTS WHERE CONSTRAINT_TYPE ='R') a,
    (SELECT b.OWNER ,b.CONSTRAINT_NAME ,b.TABLE_NAME ,b.COLUMN_NAME ,b."POSITION" AS posi FROM ALL_CONS_COLUMNS b) b
    WHERE a.owner=b.OWNER AND a.CONSTRAINT_NAME=b.CONSTRAINT_NAME AND a.table_name=b.TABLE_NAME
    ) GROUP BY owner,CONSTRAINT_NAME,TABLE_NAME),
    t2 AS
    (SELECT INDEX_OWNER ,INDEX_NAME,TABLE_OWNER ,TABLE_NAME,LISTAGG(COLUMN_NAME,',') WITHIN group(ORDER BY COLUMN_POSITION) AS ind_lst
    FROM DBA_IND_COLUMNS GROUP BY INDEX_OWNER ,INDEX_NAME,TABLE_OWNER ,TABLE_NAME)
    SELECT DISTINCT t1.owner,t1.CONSTRAINT_NAME,t1.TABLE_NAME,t1.col_lst FROM t1,t2
    WHERE t1.owner<>'SYS' AND t1.OWNER = t2.TABLE_OWNER AND t1.TABLE_NAME = t2.TABLE_NAME AND t1.col_lst <> t2.ind_lst;`
	SQL_QUERY_YASDB_FOREIGN_KEYS_WITH_IMPLICIT_DATA_TYPE_CONVERSION = `WITH t1 AS
    (SELECT b.OWNER ,b.CONSTRAINT_NAME ,b.TABLE_NAME ,b.COLUMN_NAME AS CHD_COL,b.posi,c.DATA_TYPE AS CHD_TYP FROM
    (SELECT OWNER ,CONSTRAINT_NAME ,TABLE_NAME FROM DBA_CONSTRAINTS WHERE CONSTRAINT_TYPE ='R') a,
    (SELECT b.OWNER ,b.CONSTRAINT_NAME ,b.TABLE_NAME ,b.COLUMN_NAME ,b."POSITION" AS posi FROM DBA_CONS_COLUMNS b) b ,
    (SELECT OWNER ,TABLE_NAME ,COLUMN_NAME ,DATA_TYPE FROM DBA_TAB_COLUMNS) c
    WHERE a.owner=b.OWNER AND a.CONSTRAINT_NAME=b.CONSTRAINT_NAME AND a.table_name=b.TABLE_NAME AND
    b.OWNER =c.OWNER AND b.TABLE_NAME =c.TABLE_NAME AND b.COLUMN_NAME =c.COLUMN_NAME ),
    t2 AS
    (SELECT distinct A.FK_OWNER,A.FK_CON_NAME,A.CHD_TAB,B.PRT_OWNER,B.PRT_CON_NAME,B.PRT_TAB,B.COLUMN_NAME AS PRT_COL,C.DATA_TYPE AS PRT_TYP,B.posi FROM
    (SELECT OWNER AS FK_OWNER, CONSTRAINT_NAME AS FK_CON_NAME,TABLE_NAME AS CHD_TAB,R_OWNER ,R_CONSTRAINT_NAME FROM DBA_CONSTRAINTS WHERE CONSTRAINT_TYPE ='R' ) a,
    (SELECT owner AS PRT_OWNER, CONSTRAINT_NAME AS PRT_CON_NAME, TABLE_NAME AS PRT_TAB,COLUMN_NAME,"POSITION" AS posi FROM DBA_CONS_COLUMNS) b,
    (SELECT OWNER ,TABLE_NAME ,COLUMN_NAME ,DATA_TYPE FROM DBA_TAB_COLUMNS) c
    WHERE a.R_OWNER=b.PRT_OWNER AND a.R_CONSTRAINT_NAME=b.PRT_CON_NAME AND
    B.PRT_OWNER = c.OWNER AND B.PRT_TAB=C.TABLE_NAME AND B.COLUMN_NAME=C.COLUMN_NAME)
    SELECT t2.FK_OWNER,t2.FK_CON_NAME,t2.CHD_TAB,t1.CHD_COL,t1.CHD_TYP,t2.PRT_OWNER,t2.PRT_CON_NAME,t2.PRT_TAB,t2.PRT_COL,T2.PRT_TYP
    FROM t1,t2
    WHERE t1.OWNER=t2.FK_OWNER AND t1.CONSTRAINT_NAME=t2.FK_CON_NAME AND t1.TABLE_NAME=t2.CHD_TAB AND t1.posi=t2.posi AND t1.CHD_TYP<>t2.PRT_TYP;`
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
