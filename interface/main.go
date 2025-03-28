package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var MyError bool = false
var templates map[string]*template.Template
var serverAddr string

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("template/index.html", "template/base.html"))
	templates["decoder"] = template.Must(template.ParseFiles("template/decoder.html", "template/base.html"))
}

type Data struct {
	Array      []string
	StatusCode int
}

func main() {
	flag.StringVar(&serverAddr, "server", ":8080", "server address")
	flag.Parse()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MyError = false

		data := Data{}
		data.StatusCode = http.StatusOK
		if r.Method == "GET" {
			renderTemplate(w, "index", data)
			return
		}

		if r.Method == "POST" {
			r.ParseForm()
			input := r.FormValue("input")
			choice := r.FormValue("processMethod")

			var result []string
			var statusCode int

			switch choice {
			case "decode":
				result, statusCode = decodeInput(input)
			case "encode":
				result, statusCode = encodeInput(input)
			default:
				result = nil
				statusCode = http.StatusBadRequest
			}

			data := Data{
				Array:      result,
				StatusCode: statusCode,
			}
			renderTemplate(w, "decoder", data)
			return
		}
	})

	fmt.Printf("Server is running at %s\n", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}

func renderTemplate(w http.ResponseWriter, name string, data Data) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template does not exist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(data.StatusCode)
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func decodeInput(input string) ([]string, int) {
	lines := strings.Split(input, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			result = append(result, processLine(line))
		}
	}
	if MyError {
		return result, http.StatusBadRequest
	}
	return result, http.StatusAccepted
}

func processLine(line string) string {
	var result strings.Builder
	var bracketCount int

	for i := 0; i < len(line); i++ {
		if line[i] == '[' {
			bracketCount++
		} else if line[i] == ']' {
			bracketCount--
			if bracketCount < 0 {
				MyError = true
				return "Error"
			}
		}
	}

	if bracketCount != 0 {
		MyError = true
		return "Error"
	}

	i := 0
	for i < len(line) {
		if line[i] == '[' {
			countStart := i + 1
			for line[i] != ']' {
				i++
				if i == len(line) {
					MyError = true
					return "Error"
				}
			}
			countAndChars := line[countStart:i]
			parts := splitCountAndChars(countAndChars)

			fmt.Println("Parts:", parts)

			if len(parts) != 2 || parts[1] == "" || !startsWithNumber(parts[0]) {
				MyError = true
				return "Error"
			}

			count, err := strconv.Atoi(parts[0])
			if err != nil {
				MyError = true
				return "Error"
			}
			chars := parts[1]

			for j := 0; j < count; j++ {
				result.WriteString(chars)
			}
			i++
		} else {
			result.WriteByte(line[i])
			i++
		}
	}

	return result.String()
}

func encodeInput(input string) ([]string, int) {
	lines := strings.Split(input, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		result = append(result, encodeLine(line))
	}
	return result, http.StatusAccepted
}

func encodeLine(input string) string {
	var result strings.Builder
	count := 1
	prevChar := input[0]

	if len(input) == 0 || input[0] == '\n' {
		return ""
	}

	for i := 1; i < len(input); i++ {
		if input[i] == '\n' {
			continue
		}

		if input[i] == prevChar {
			count++
		} else {
			if count > 1 {
				result.WriteString(fmt.Sprintf("[%d %s]", count, string(prevChar)))
			} else {
				result.WriteString(string(prevChar))
			}
			prevChar = input[i]
			count = 1
		}
	}

	if count > 1 {
		result.WriteString(fmt.Sprintf("[%d %s]", count, string(prevChar)))
	} else {
		result.WriteString(string(prevChar))
	}

	fmt.Println("Encoded:", result.String())

	return result.String()
}

func splitCountAndChars(s string) []string {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) == 1 {
		if strings.TrimSpace(parts[0]) == "" {
			return []string{"", ""}
		}
		return parts
	}
	return parts
}

func startsWithNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
