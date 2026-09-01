<?php
// pdf_reader.php
// Версия на PHP с использованием smalot/pdfparser

require_once 'vendor/autoload.php';

use Smalot\PdfParser\Parser;

// ANSI-цвета
define('RESET', "\033[0m");
define('CYAN', "\033[96m");
define('GREEN', "\033[92m");
define('YELLOW', "\033[93m");
define('RED', "\033[91m");
define('BLUE', "\033[94m");
define('BOLD', "\033[1m");

function colorize($text, $color) {
    return $color . $text . RESET;
}

class PDFReader {
    private $filename;
    private $password;
    private $parser;
    private $pdf;

    public function __construct($filename, $password = null) {
        $this->filename = $filename;
        $this->password = $password;
        $this->parser = new Parser();
        $this->pdf = $this->parser->parseFile($filename);
        // Пароль не поддерживается напрямую, но оставим для совместимости
    }

    public function getMetadata() {
        $details = $this->pdf->getDetails();
        $meta = [
            'Title' => $details['Title'] ?? '',
            'Author' => $details['Author'] ?? '',
            'Creator' => $details['Creator'] ?? '',
            'Producer' => $details['Producer'] ?? '',
            'CreationDate' => $details['CreationDate'] ?? '',
            'ModDate' => $details['ModDate'] ?? '',
        ];
        return $meta;
    }

    public function getPageCount() {
        return count($this->pdf->getPages());
    }

    public function extractText($pageNum = null) {
        $text = '';
        if ($pageNum !== null) {
            $pages = $this->pdf->getPages();
            if ($pageNum < 1 || $pageNum > count($pages)) {
                throw new Exception("Номер страницы вне диапазона");
            }
            $text = $pages[$pageNum-1]->getText();
        } else {
            $text = $this->pdf->getText();
        }
        return $text;
    }

    public function search($query, $pageNum = null) {
        $text = $this->extractText($pageNum);
        $lines = explode("\n", $text);
        $results = [];
        foreach ($lines as $i => $line) {
            if (stripos($line, $query) !== false) {
                $results[$i+1] = trim($line);
            }
        }
        return $results;
    }
}

// Парсинг аргументов
$shortOpts = "tmps:o:";
$longOpts = ['password:', 'page:'];
$options = getopt($shortOpts, $longOpts);
$args = array_values(array_filter($argv, function($arg) { return !str_starts_with($arg, '-'); }));
array_shift($args);

$filename = $args[0] ?? null;
$password = $options['password'] ?? null;
$pageNum = isset($options['page']) ? (int)$options['page'] : 0;
$textFlag = isset($options['t']);
$metadataFlag = isset($options['m']);
$pagesFlag = isset($options['p']);
$searchQuery = $options['s'] ?? null;
$outputFile = $options['o'] ?? null;

if (!$filename) {
    echo "Укажите PDF файл.\n";
    exit(1);
}

try {
    $reader = new PDFReader($filename, $password);

    if ($metadataFlag) {
        $meta = $reader->getMetadata();
        echo colorize("Метаданные:", BOLD) . "\n";
        foreach ($meta as $k => $v) {
            echo "  " . colorize($k, CYAN) . ": $v\n";
        }
        echo "\n";
    }

    if ($pagesFlag) {
        $count = $reader->getPageCount();
        echo "Количество страниц: " . colorize($count, GREEN) . "\n\n";
    }

    if ($searchQuery) {
        $results = $reader->search($searchQuery, $pageNum);
        if (count($results) > 0) {
            echo colorize("Найдено совпадений: " . count($results), GREEN) . "\n";
            foreach ($results as $lineNum => $line) {
                $highlighted = preg_replace('/' . preg_quote($searchQuery, '/') . '/i', colorize('$0', YELLOW), $line);
                echo "  Строка $lineNum: $highlighted\n";
            }
        } else {
            echo colorize("Совпадений не найдено.", RED) . "\n";
        }
        echo "\n";
    }

    if ($textFlag || (!$metadataFlag && !$pagesFlag && !$searchQuery)) {
        $text = $reader->extractText($pageNum);
        if ($outputFile) {
            file_put_contents($outputFile, $text);
            echo colorize("Текст сохранён в $outputFile", GREEN) . "\n";
        } else {
            echo colorize("Извлечённый текст:", BOLD) . "\n";
            echo $text;
        }
    }
} catch (Exception $e) {
    echo colorize("Ошибка: " . $e->getMessage(), RED) . "\n";
    exit(1);
}
