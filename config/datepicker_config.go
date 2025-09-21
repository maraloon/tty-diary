package datepicker_config

import (
	"github.com/maraloon/datepicker"
	"github.com/spf13/pflag"
)

func ValidateFlags() datepicker.Config {
	var sunday bool
	var monday bool
	pflag.BoolVarP(&sunday, "sunday", "s", true, "Sunday as first day of week")
	pflag.BoolVarP(&monday, "monday", "m", false, "Monday as first day of week")
	pflag.Parse()

	config := datepicker.DefaultConfig()
	config.FirstWeekdayIsMo = monday || !sunday

	return config
}
