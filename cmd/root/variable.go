package root

import "embed"

const (
	FlagSystemTRXPath                        string = "systemtrxpath"
	FlagSystemTRXPathShort                   string = "s"
	FlagSystemTRXPathUsage                   string = "Path location of System Transaction directory"
	FlagBankTRXPath                          string = "banktrxpath"
	FlagBankTRXPathShort                     string = "b"
	FlagBankTRXPathUsage                     string = "Path location of Bank Transaction directory"
	FlagReportTRXPath                        string = "reportpath"
	FlagReportTRXPathShort                   string = "r"
	FlagReportTRXPathUsage                   string = "Path location of Archive directory"
	FlagListBank                             string = "listbank"
	FlagListBankShort                        string = "l"
	FlagListBankUsage                        string = "List bank accepted"
	FlagTimeZone                             string = "time_zone"
	FlagTimeZoneShort                        string = "z"
	FlagTimeZoneUsage                        string = `time zone settings`
	FlagFromDate                             string = "from"
	FlagFromDateShort                        string = "f"
	FlagFromDateUsage                        string = `from date (YYYY-MM-DD)`
	FlagToDate                               string = "to"
	FlagToDateShort                          string = "t"
	FlagToDateUsage                          string = `to date (YYYY-MM-DD)`
	FlagTotalDataSampleToGenerate            string = "amountdata"
	FlagTotalDataSampleToGenerateShort       string = "a"
	FlagTotalDataSampleToGenerateUsage       string = `amount system trx data sample to generate, bank trx will be 2 times of this amount`
	FlagPercentageMatchSampleToGenerate      string = "percentagematch"
	FlagPercentageMatchSampleToGenerateShort string = "p"
	FlagPercentageMatchSampleToGenerateUsage string = `percentage of matched trx for data sample to generate`
)

var EmbedFS *embed.FS
var SampleUsageFlags = "--systemtrxpath=/tmp/system --banktrxpath=/tmp/bank --archivepath= /tmp/archive --listbank=bca,mandiri,bri,danamon"
var ProcessUsageFlags = "--systemtrxpath=/tmp/system --banktrxpath=/tmp/bank --archivepath= /tmp/archive --listbank=bca,mandiri,bri,danamon --from=2024/11/10 --from=2024/11/11"
var FlagSystemTRXPathValue string
var FlagBankTRXPathValue string
var FlagReportTRXPathValue string
var FlagListBankValue []string
var FlagTZValue string
var FlagFromDateValue string
var FlagToDateValue string
var DefaultListBank = []string{"bca", "mandiri", "bri", "danamon"}
var FlagTotalDataSampleToGenerateValue int64
var DefaultTotalDataSampleToGenerate int64 = 1000
var FlagPercentageMatchSampleToGenerateValue int
var DefaultPercentageMatchSampleToGenerate = 100
