package config

import (
	"os"

	"github.com/maraloon/datepicker"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	FileColor        string
	FileFormat       string
	DiaryDir         string
	DatepickerConfig datepicker.Config
}

func ValidateFlags() (Config, error) {

	var fileColor string
	var fileFormat string
	var diaryDir string
	var sunday bool
	var monday bool

	viper.SetConfigName("ttydiary.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("HOME"))
	viper.AddConfigPath(os.Getenv("HOME") + "/.config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	flag.StringVar(&fileColor, "file-color", "6", "Color of dates, which have notes")
	flag.StringVar(&fileFormat, "file-format", "md", "Format of note files")
	flag.StringVar(&diaryDir, "diary-dir", os.Getenv("HOME")+"/code/util/notes/diary", "Root dir of notes")
	flag.BoolVarP(&monday, "monday", "m", false, "Monday as first day of week")
	flag.BoolVarP(&sunday, "sunday", "s", true, "Sunday as first day of week")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	datepickerConfig := datepicker.DefaultConfig()

	if flag.Lookup("monday").Changed { // if -m/--monday
		datepickerConfig.FirstWeekdayIsMo = viper.GetBool("monday")
	} else if flag.Lookup("sunday").Changed { // if -s/--sunday
		datepickerConfig.FirstWeekdayIsMo = !viper.GetBool("sunday")
	} else { // get value from config or get default
		datepickerConfig.FirstWeekdayIsMo = viper.GetBool("monday") || !viper.GetBool("sunday")
	}

	config := Config{
		FileColor:        viper.GetString("file-color"),
		FileFormat:       viper.GetString("file-format"),
		DiaryDir:         viper.GetString("diary-dir"),
		DatepickerConfig: datepickerConfig,
	}

	return config, nil
}
