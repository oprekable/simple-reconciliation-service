package progressbarhelper

import "github.com/schollz/progressbar/v3"

func BarDescribe(bar *progressbar.ProgressBar, description string) {
	if bar != nil {
		bar.Describe(description)
	}
}
