package sample

const (
	QueryDropTableArguments = `
-- QuerystmtDropTableArguments
DROP TABLE IF EXISTS arguments;
`

	QueryDropTableBanks = `
-- QuerystmtDropTableBanks
DROP TABLE IF EXISTS banks;
`

	QueryDropTableBaseData = `
-- QuerystmtDropTableBaseData
DROP TABLE IF EXISTS base_data;
`

	QueryCreateTableArguments = `
-- QueryCreateTableArguments
CREATE TEMPORARY TABLE IF NOT EXISTS arguments AS
SELECT
    DATETIME(DATE(?)) AS start
    , DATETIME(DATE(?)) AS end
    , ABS(CAST(? AS INTEGER)) AS limit_trx_data
    , ABS(CAST(? AS INTEGER)) AS match_percentage
;
`

	QueryCreateTableBanks = `
-- QueryCreateTableBanks
CREATE TEMPORARY TABLE IF NOT EXISTS banks AS
SELECT
    key AS id
    , value AS bank_name
FROM json_each(
    ?
)
;
`
	QueryCreateTableBaseData = `
-- QueryCreateTableBaseData
CREATE TEMPORARY TABLE IF NOT EXISTS base_data AS
WITH RECURSIVE generate_series(value) AS (
    SELECT
        start AS v
    FROM arguments
    UNION ALL
    SELECT
        DATETIME(value , '+1 second') AS v
    FROM generate_series
    WHERE
        DATETIME(value , '+1 second') <= (SELECT DATETIME(DATETIME(end, '+1 day'), '-1 second') FROM arguments)
),
data AS (
    SELECT
        value AS transactionTime
        , lower(hex(randomblob(16))) AS trxID
        , RANDOM() AS randomData
    FROM generate_series
)
SELECT
    trxID
    , transactionTime
    , CAST(
        FLOOR((ABS(randomData) % (100000 - 1000) + 1000) / 100) * 100 AS FLOAT
      ) AS amount
    , CASE
          WHEN randomData % 2 = 1 THEN 'DEBIT'
          ELSE 'CREDIT'
    END AS type
    , bank_name AS bank
FROM data
LEFT JOIN banks ON id = ABS(randomData % (SELECT COUNT(*) FROM banks))
;
`
	QueryCreateIndexTableBaseData = `
-- QueryCreateIndexTableBaseData
CREATE INDEX IF NOT EXISTS idx_base_data_trx_id ON base_data(trxID);
`

	QueryGetTrxData = `
-- QueryGetTrxData
WITH with_system_trx AS (
    SELECT
        *
        , ROW_NUMBER() OVER win    AS row_number
        , FLOOR((ROW_NUMBER() OVER win/ CAST(COUNT() OVER counter AS REAL)) * 100) AS percentile
    FROM (
         SELECT
             bd1.trxID
              , system_trx.trxID IS NOT NULL AS is_system_trx
         FROM base_data bd1
         LEFT JOIN (
             SELECT
                 bd2.trxID
             FROM base_data bd2
             ORDER BY random()
             LIMIT (SELECT limit_trx_data FROM arguments)
         ) system_trx ON bd1.trxID = system_trx.trxID
    ) tmp_data
    WINDOW
        win AS (ORDER BY tmp_data.is_system_trx DESC)
        , counter AS (PARTITION BY tmp_data.is_system_trx)
    ORDER BY tmp_data.is_system_trx DESC
)
, with_bank_trx_flagged AS (
    SELECT
        wst.trxID
         , wst.is_system_trx
         , wst.row_number
         , CASE
               WHEN a.match_percentage = 100 THEN wst.is_system_trx
               WHEN a.match_percentage < 100 AND a.match_percentage >= 0 THEN
                   CASE
                       WHEN wst.percentile <= a.match_percentage THEN TRUE
                       WHEN wst.percentile > a.match_percentage THEN FALSE
                       ELSE TRUE END
               ELSE FALSE
        END AS is_bank_trx
    FROM with_system_trx wst, arguments a
)
SELECT
    wbtf.trxID
    , lower(bd.bank) || '-' || lower(hex(randomblob(16))) AS uniqueIdentifier
    , wbtf.is_system_trx
    , wbtf.is_bank_trx
    , bd.type
    , bd.bank
    , bd.amount
    , bd.transactionTime
    , DATE(bd.transactionTime)  AS date
FROM with_bank_trx_flagged wbtf
INNER JOIN base_data bd ON bd.trxID = wbtf.trxID
ORDER BY wbtf.row_number
LIMIT (SELECT limit_trx_data * 2 FROM arguments)
;
`
)
