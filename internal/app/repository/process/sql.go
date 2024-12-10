package process

const (
	QueryDropTableArguments = `
-- QuerystmtDropTableArguments
DROP TABLE IF EXISTS arguments;
`

	QueryDropTableBanks = `
-- QuerystmtDropTableBanks
DROP TABLE IF EXISTS banks;
`
	QueryDropTableSystemTrx = `
-- QueryDropTableSystemTrx
DROP TABLE IF EXISTS system_trx;
`
	QueryDropTableBankTrx = `
-- QueryDropTableBankTrx
DROP TABLE IF EXISTS bank_trx;
`
	QueryDropTableReconciliationMap = `
-- QueryDropTableReconciliationMap
DROP TABLE IF EXISTS reconciliation_map;
`
	QueryCreateTableArguments = `
-- QueryCreateTableArguments
CREATE TABLE IF NOT EXISTS arguments AS
SELECT
    DATETIME(DATE(?)) AS start
    , DATETIME(DATE(?)) AS end
;
`
	QueryCreateTableBanks = `
-- QueryCreateTableBanks
CREATE TABLE IF NOT EXISTS banks AS
SELECT
    key AS id
    , LOWER(value) AS bank_name
FROM json_each(
    ?
)
;
`

	QueryCreateTableSystemTrx = `
-- QueryCreateTableSystemTrx
CREATE TABLE IF NOT EXISTS system_trx (
	TrxID TEXT PRIMARY KEY,
	Amount FLOAT,
	Type TEXT,
	TransactionTime DATETIME,
	FilePath TEXT
);
;
`
	QueryCreateTableBankTrx = `
-- QueryCreateTableBankTrx
CREATE TABLE IF NOT EXISTS bank_trx (
	UniqueIdentifier TEXT PRIMARY KEY,
	Amount FLOAT,
	Type TEXT,
	Bank TEXT,
	Date DATE,
	FilePath TEXT
);

CREATE INDEX IF NOT EXISTS bank_trx_Date_Type_Amount_index ON bank_trx (Date, Type, Amount);
`
	QueryCreateTableReconciliationMap = `
-- QueryCreateTableReconciliationMap
CREATE TABLE IF NOT EXISTS reconciliation_map (
	TrxID TEXT PRIMARY KEY,
	UniqueIdentifier TEXT
);

CREATE INDEX IF NOT EXISTS reconciliation_map_TrxID_index ON reconciliation_map (TrxID);
CREATE INDEX IF NOT EXISTS reconciliation_map_UniqueIdentifier_index ON reconciliation_map (UniqueIdentifier);
`
	QueryInsertTableSystemTrx = `
-- QueryInsertTableSystemTrx
INSERT OR IGNORE INTO system_trx (TrxID, Amount, Type, TransactionTime, FilePath)
	SELECT
	json_extract(j.value, '$.TrxID') AS TrxID
	 , json_extract(j.value, '$.Amount') AS Amount
	 , json_extract(j.value, '$.Type') AS Type
	 , json_extract(j.value, '$.TransactionTime') AS TransactionTime
	 , json_extract(j.value, '$.FilePath') AS FilePath
	FROM json_each(
	 ?
	) AS j
;
`

	QueryInsertTableBankTrx = `
-- QueryInsertTableBankTrx
INSERT OR IGNORE INTO bank_trx (UniqueIdentifier, Date, Type, FilePath, Bank, Amount)
	SELECT
	json_extract(j.value, '$.UniqueIdentifier') AS UniqueIdentifier
	 , json_extract(j.value, '$.Date') AS Date
	 , json_extract(j.value, '$.Type') AS Type
	 , json_extract(j.value, '$.FilePath') AS FilePath
	 , json_extract(j.value, '$.Bank') AS Bank
	 , json_extract(j.value, '$.Amount') AS Amount
	FROM json_each(
	 ?
	) AS j
;
`

	QueryInsertTableReconciliationMap = `
-- QueryInsertTableReconciliationMap
WITH base_search AS (
    SELECT
        st.TrxID
         , bt.UniqueIdentifier
         , st.Amount
         , st.Type
         , bt.Bank
         , bt.Date
         , st.TransactionTime
    FROM
        arguments a
       , system_trx st
    INNER JOIN bank_trx bt ON
        bt.Date = DATE(st.TransactionTime)
        AND bt.Type = st.Type
        AND bt.Amount = st.Amount
)
, with_counter AS (
    SELECT
        ROW_NUMBER() OVER (PARTITION BY bs.Date, bs.Amount, bs.Type, bs.TrxID ORDER BY bs.UniqueIdentifier) AS r_system
         , ROW_NUMBER() OVER (PARTITION BY bs.Date, bs.Amount, bs.Type, bs.UniqueIdentifier ORDER BY bs.TrxID) AS r_bank
         , bs.*
    FROM base_search bs
)
, matched_trx AS (
    SELECT
        wc.*
    FROM with_counter wc
    WHERE wc.r_system = wc.r_bank
    ORDER BY wc.Date, wc.Amount, wc.TrxID, wc.UniqueIdentifier, wc.TransactionTime
)
INSERT INTO reconciliation_map(
    TrxID,
    UniqueIdentifier
)
SELECT
    mt.TrxID,
    mt.UniqueIdentifier
FROM matched_trx mt
WHERE NOT EXISTS (SELECT 1 FROM reconciliation_map LIMIT 1)
;
`
)
