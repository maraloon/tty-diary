package filer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Filer struct {
	diaryDir   string
	fileFormat string
}

func NewFiler(diaryDir, fileFormat string) *Filer {
	return &Filer{
		diaryDir:   diaryDir,
		fileFormat: fileFormat,
	}
}

func (f *Filer) GetDatesWithFiles(startYear, endYear int) []string {
	var dates []string

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.
				UTC).Day()
			for day := 1; day <= daysInMonth; day++ {
				date := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
				if f.FileExistsAndNotEmpty(date) {
					dates = append(dates, date)
				}
			}
		}
	}

	return dates
}

func (f *Filer) FileExistsAndNotEmpty(date string) bool {
	path := f.Filepath(date)
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return false
		}

		return len(string(data)) > 0
	}
	return false
}

func (f *Filer) Filepath(date string) string {
	return filepath.Join(f.diaryDir, date+"."+f.fileFormat)
}
