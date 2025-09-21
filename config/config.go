package config

import (
	"github.com/maraloon/datepicker"
	"github.com/spf13/pflag"
)

type Config struct {
	NotesColor       string
	FileFormat       string
	DiaryDir         string
	DatepickerConfig datepicker.Config
}

func ValidateFlags() Config {
	var notesColor string
	var fileFormat string
	var diaryDir string
	var sunday bool
	var monday bool
	pflag.StringVar(&notesColor, "color", "6", "Color of dates, which have notes")
	pflag.StringVar(&fileFormat, "file-format", "md", "Format of note files")
	pflag.StringVar(&diaryDir, "diary-dir", "/code/util/notes/diary", "Root dir of notes")
	pflag.BoolVarP(&sunday, "sunday", "s", true, "Sunday as first day of week")
	pflag.BoolVarP(&monday, "monday", "m", false, "Monday as first day of week")
	pflag.Parse()

	datepickerConfig := datepicker.DefaultConfig()
	datepickerConfig.FirstWeekdayIsMo = monday || !sunday

	config := Config{
		NotesColor: notesColor,
		FileFormat: fileFormat,
		DiaryDir:   diaryDir,
		DatepickerConfig: datepickerConfig,
	}

	return config
}
