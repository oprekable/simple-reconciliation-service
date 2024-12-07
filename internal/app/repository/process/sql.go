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
;
`
	QueryInsertTableSystemTrx = `
-- QueryInsertTableSystemTrx
INSERT INTO system_trx (TrxID, Amount, Type, TransactionTime, FilePath)
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
INSERT INTO bank_trx (UniqueIdentifier, Date, Type, FilePath, Bank, Amount)
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
)
