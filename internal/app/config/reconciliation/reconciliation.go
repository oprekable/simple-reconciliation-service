package reconciliation

// Reconciliation ..
type Reconciliation struct {
	Action        string   `default:"-" mapstructure:"action"`
	SystemTRXPath string   `default:"-" mapstructure:"system_trx_path"`
	BankTRXPath   string   `default:"-" mapstructure:"bank_trx_path"`
	ArchivePath   string   `default:"-" mapstructure:"archive_path"`
	ListBank      []string `default:"-" mapstructure:"list_bank"`
}
