package root

import "embed"

const (
	FlagSystemTRXPath       string = "systemtrxpath"
	FlagSystemTRXPathShort  string = "s"
	FlagSystemTRXPathUsage  string = "Path location of System Transaction directory"
	FlagBankTRXPath         string = "banktrxpath"
	FlagBankTRXPathShort    string = "b"
	FlagBankTRXPathUsage    string = "Path location of Bank Transaction directory"
	FlagArchiveTRXPath      string = "archivepath"
	FlagArchiveTRXPathShort string = "a"
	FlagArchiveTRXPathUsage string = "Path location of Archive directory"
	FlagListBank            string = "listbank"
	FlagListBankShort       string = "l"
	FlagListBankUsage       string = "Path location of Archive directory"
	FlagTimeZone            string = "time_zone"
	FlagTimeZoneShort       string = "t"
	FlagTimeZoneUsage       string = `time zone settings`
)

var EmbedFS *embed.FS
var SampleUsageFlags = "--systemtrxpath=/tmp/system --banktrxpath=/tmp/bank --archivepath= /tmp/archive --listbank=bca,mandiri,bri,danamon"
var ProcessUsageFlags = "--systemtrxpath=/tmp/system --banktrxpath=/tmp/bank --archivepath= /tmp/archive --listbank=bca,mandiri,bri,danamon --from=2024/11/10 --from=2024/11/11"
var SystemTRXPath string
var BankTRXPath string
var ArchivePath string
var ListBank []string
var FlagTZValue string
var DefaultListBank = []string{"bca", "mandiri", "bri", "danamon"}
