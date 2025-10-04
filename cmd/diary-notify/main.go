package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"tty-diary/config"
	"tty-diary/filer"
)

var (
	dircachefile = os.Getenv("HOME") + "/.cache/tty-diary/"
	cachefile    = dircachefile + "last_check"
)

func main() {
	config, err := config.ValidateFlags()
	filer := filer.NewFiler(config.DiaryDir, config.FileFormat)

	var last_exec time.Time
	cache_file_content, err := os.ReadFile(cachefile)
	if err != nil {
		if os.IsNotExist(err) {
			last_exec = time.Now().Add(-1 * time.Minute)
			os.Mkdir(dircachefile, 0755)
			err := os.WriteFile(cachefile, []byte(time.Now().Format("2006/01/02 15:04")), 0644)
			if err != nil {
				log.Fatalf("error write file: %v", err)
			}
		}
	} else {
		loc, err := time.LoadLocation("Local")
		if err != nil {
			log.Fatalf("Error load location: %v", err)
		}
		last_exec, err = time.ParseInLocation("2006/01/02 15:04", string(cache_file_content), loc)

		if err != nil {
			last_exec, err = time.ParseInLocation("2006/01/02 15:04", string(cache_file_content[:len(cache_file_content)-1]), loc)
			if err != nil {
				log.Fatalf("Error parse date: %v", err)
			}
		}
	}

	last_exec_epoch := last_exec.Unix()
	current_epoch := time.Now().Unix()

	date_iter := last_exec
	is_time_regexp := regexp.MustCompile(`^([0-9]{2}:[0-9]{2})\ (.*)$`)

	for {
		if date_iter.Compare(time.Now()) == 1 {
			break
		}
		notes_file, err := os.Open(filer.Filepath(date_iter.Format("2006/01/02")))

		if err != nil {

			if os.IsNotExist(err) {
				date_iter = date_iter.Add(24 * time.Hour)
				continue
			} else {
				log.Fatalf("Error opening file: %v", err)
			}
		}
		defer notes_file.Close()

		scanner := bufio.NewScanner(notes_file)

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		for scanner.Scan() {
			line := scanner.Text()
			line_with_time := is_time_regexp.FindStringSubmatch(line)

			if line_with_time != nil {

				line_time := line_with_time[1]
				line_text := line_with_time[2]

				loc, err := time.LoadLocation("Local")
				if err != nil {
					log.Fatalf("Error load location: %v", err)
				}

				line_datetime, err := time.ParseInLocation("2006/01/02 15:04", fmt.Sprintf("%s %s", date_iter.Format("2006/01/02"), line_time), loc)
				if err != nil {
					log.Fatalf("Error parse date: %v", err)

				}
				line_epoch := line_datetime.Unix()

				if line_epoch > last_exec_epoch && line_epoch <= current_epoch {
					if int(time.Now().Sub(date_iter).Hours()) > 0 {
						day_and_time := line_datetime.Format("Jan 2 15:04")
						cmd := exec.Command("notify-send", "-c", "diary-not-today", line_text, day_and_time)
						cmd.Run()
					} else {
						if line_time != time.Now().Format("15:04") {
							cmd := exec.Command("notify-send", "-c", "diary-today", line_text, line_time)
							cmd.Run()
						} else {
							cmd := exec.Command("notify-send", "-c", "diary-now", line_text)
							cmd.Run()
						}
					}
				}

			}

		}

		date_iter = date_iter.Add(24 * time.Hour)
	}

	err = os.WriteFile(cachefile, []byte(time.Now().Format("2006/01/02 15:04")), 0644)
	if err != nil {
		log.Fatalf("error write file: %v", err)
	}

}
