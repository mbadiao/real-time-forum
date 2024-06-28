package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"text/template"
	"time"
)

func FileService(str string, w http.ResponseWriter, data any) {
	tmpl, err := template.ParseFiles("./frontend/views/" + str)
	if err != nil {
		fmt.Println("NOTFOUND")
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("error while executing the template")
		return
	}
}

func IsImage(file multipart.File) bool {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return false
	}
	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/svg+xml", "image/jpeg", "image/gif", "image/png":
	default:
		fmt.Println(contentType)
		fmt.Println("type")
		return false
	}
	return true
}

func FormatTimeAgo(creationDate time.Time) string {
	now := time.Now()
	diff := now.Sub(creationDate)

	var res string

	switch {
	case diff.Hours() >= 24*365:
		years := int(diff.Hours() / (24 * 365))
		res = fmt.Sprintf("%d year", years)
		if years > 1 {
			res += "s"
		}
		res += " ago"
	case diff.Hours() >= 24*30:
		months := int(diff.Hours() / (24 * 30))
		res = fmt.Sprintf("%d month", months)
		if months > 1 {
			res += "s"
		}
		res += " ago"
	case diff.Hours() >= 24*7:
		weeks := int(diff.Hours() / (24 * 7))
		res = fmt.Sprintf("%d week", weeks)
		if weeks > 1 {
			res += "s"
		}
		res += " ago"
	case diff.Hours() >= 24:
		days := int(diff.Hours() / 24)
		res = fmt.Sprintf("%d day", days)
		if days > 1 {
			res += "s"
		}
		res += " ago"
	case diff.Hours() >= 1:
		hours := int(diff.Hours())
		res = fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			res += "s"
		}
		res += " ago"
	case diff.Minutes() >= 1:
		minutes := int(diff.Minutes())
		res = fmt.Sprintf("%d minute", minutes)
		if minutes > 1 {
			res += "s"
		}
		res += " ago"
	default:
		seconds := int(diff.Seconds())
		res = fmt.Sprintf("%d second", seconds)
		if seconds != 1 {
			res += "s"
		}
		res += " ago"
	}

	return res
}
