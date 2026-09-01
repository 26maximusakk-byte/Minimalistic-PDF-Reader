// pdf_reader.go
// Версия на Go с использованием github.com/unidoc/unidoc (или ledongthuc/pdf)
// Для упрощения используем ledongthuc/pdf (легче в установке)

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ANSI-цвета
const (
	reset  = "\033[0m"
	cyan   = "\033[96m"
	green  = "\033[92m"
	yellow = "\033[93m"
	red    = "\033[91m"
	blue   = "\033[94m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

type PDFReader struct {
	filename string
	password string
}

func NewPDFReader(filename, password string) *PDFReader {
	return &PDFReader{filename: filename, password: password}
}

func (r *PDFReader) Open() (*pdf.Reader, error) {
	f, err := os.Open(r.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader, err := pdf.NewReader(f, r.password)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (r *PDFReader) GetMetadata() (map[string]string, error) {
	reader, err := r.Open()
	if err != nil {
		return nil, err
	}
	meta := reader.Trailer()
	result := make(map[string]string)
	if meta != nil {
		for _, key := range []string{"Title", "Author", "Creator", "Producer", "CreationDate", "ModDate"} {
			if val := meta.Key(key).String(); val != "" {
				result[key] = val
			}
		}
	}
	return result, nil
}

func (r *PDFReader) GetPageCount() (int, error) {
	reader, err := r.Open()
	if err != nil {
		return 0, err
	}
	return reader.NumPage(), nil
}

func (r *PDFReader) ExtractText(pageNum int) (string, error) {
	reader, err := r.Open()
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	if pageNum > 0 {
		if pageNum < 1 || pageNum > reader.NumPage() {
			return "", fmt.Errorf("номер страницы вне диапазона")
		}
		page := reader.Page(pageNum)
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		builder.WriteString(text)
	} else {
		for i := 1; i <= reader.NumPage(); i++ {
			page := reader.Page(i)
			text, err := page.GetPlainText(nil)
			if err != nil {
				return "", err
			}
			builder.WriteString(text)
		}
	}
	return builder.String(), nil
}

func (r *PDFReader) Search(query string, pageNum int) (map[int][]string, error) {
	text, err := r.ExtractText(pageNum)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(text, "\n")
	results := make(map[int][]string)
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			results[i+1] = append(results[i+1], strings.TrimSpace(line))
		}
	}
	return results, nil
}

func main() {
	var (
		textFlag      bool
		metadataFlag  bool
		pagesFlag     bool
		searchQuery   string
		outputFile    string
		password      string
		pageNum       int
	)
	flag.BoolVar(&textFlag, "t", false, "Извлечь текст")
	flag.BoolVar(&metadataFlag, "m", false, "Показать метаданные")
	flag.BoolVar(&pagesFlag, "p", false, "Показать количество страниц")
	flag.StringVar(&searchQuery, "s", "", "Искать текст")
	flag.StringVar(&outputFile, "o", "", "Сохранить текст в файл")
	flag.StringVar(&password, "password", "", "Пароль для PDF")
	flag.IntVar(&pageNum, "page", 0, "Номер страницы (1-based)")
	flag.Usage = func() {
		fmt.Println("Использование: go run pdf_reader.go [опции] <pdf_file>")
		fmt.Println("  -t            Извлечь текст")
		fmt.Println("  -m            Показать метаданные")
		fmt.Println("  -p            Показать количество страниц")
		fmt.Println("  -s <query>    Искать текст")
		fmt.Println("  -o <file>     Сохранить текст в файл")
		fmt.Println("  --password    Пароль для PDF")
		fmt.Println("  --page <num>  Номер страницы")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	reader := NewPDFReader(filename, password)

	if metadataFlag {
		meta, err := reader.GetMetadata()
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		fmt.Println(colorize("Метаданные:", bold))
		for k, v := range meta {
			fmt.Printf("  %s: %s\n", colorize(k, cyan), v)
		}
		fmt.Println()
	}

	if pagesFlag {
		count, err := reader.GetPageCount()
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		fmt.Printf("Количество страниц: %s\n", colorize(fmt.Sprintf("%d", count), green))
		fmt.Println()
	}

	if searchQuery != "" {
		results, err := reader.Search(searchQuery, pageNum)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		if len(results) > 0 {
			fmt.Printf(colorize("Найдено совпадений: %d\n", green), len(results))
			for lineNum, lines := range results {
				for _, line := range lines {
					highlighted := strings.ReplaceAll(line, searchQuery, colorize(searchQuery, yellow))
					fmt.Printf("  Строка %d: %s\n", lineNum, highlighted)
				}
			}
		} else {
			fmt.Println(colorize("Совпадений не найдено.", red))
		}
		fmt.Println()
	}

	if textFlag || (!metadataFlag && !pagesFlag && searchQuery == "") {
		text, err := reader.ExtractText(pageNum)
		if err != nil {
			fmt.Println(colorize("Ошибка: "+err.Error(), red))
			os.Exit(1)
		}
		if outputFile != "" {
			err = os.WriteFile(outputFile, []byte(text), 0644)
			if err != nil {
				fmt.Println(colorize("Ошибка записи: "+err.Error(), red))
				os.Exit(1)
			}
			fmt.Println(colorize("Текст сохранён в "+outputFile, green))
		} else {
			fmt.Println(colorize("Извлечённый текст:", bold))
			fmt.Println(text)
		}
	}
}
