# pdf_reader.py
# Версия на Python с использованием pypdf, argparse, цветной вывод

import sys
import argparse
import os
from pypdf import PdfReader
from pypdf.errors import PdfReadError, FileNotDecryptedError
import re

# ANSI-цвета
class Colors:
    RESET = '\033[0m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BLUE = '\033[94m'
    BOLD = '\033[1m'

def colorize(text, color):
    return f"{color}{text}{Colors.RESET}"

class PDFReader:
    def __init__(self, filename, password=None):
        self.filename = filename
        self.password = password
        self.reader = None
        self._load()

    def _load(self):
        try:
            self.reader = PdfReader(self.filename)
            if self.reader.is_encrypted:
                if self.password:
                    self.reader.decrypt(self.password)
                else:
                    raise FileNotDecryptedError("Файл зашифрован. Укажите пароль.")
        except Exception as e:
            raise Exception(f"Не удалось открыть PDF: {e}")

    def get_metadata(self):
        meta = self.reader.metadata
        if meta is None:
            return {}
        return {
            'title': meta.get('/Title', ''),
            'author': meta.get('/Author', ''),
            'creator': meta.get('/Creator', ''),
            'producer': meta.get('/Producer', ''),
            'creation_date': meta.get('/CreationDate', ''),
            'mod_date': meta.get('/ModDate', ''),
        }

    def get_page_count(self):
        return len(self.reader.pages)

    def extract_text(self, page_num=None):
        if page_num is not None:
            if page_num < 1 or page_num > len(self.reader.pages):
                raise ValueError("Номер страницы вне диапазона")
            pages = [self.reader.pages[page_num-1]]
        else:
            pages = self.reader.pages
        text = ''
        for page in pages:
            text += page.extract_text()
        return text

    def search(self, query, page_num=None):
        text = self.extract_text(page_num)
        lines = text.split('\n')
        results = []
        for i, line in enumerate(lines):
            if query.lower() in line.lower():
                results.append((i+1, line.strip()))
        return results

def main():
    parser = argparse.ArgumentParser(description='Minimalistic PDF Reader (Python)')
    parser.add_argument('file', help='PDF файл')
    parser.add_argument('-t', '--text', action='store_true', help='Извлечь текст')
    parser.add_argument('-m', '--metadata', action='store_true', help='Показать метаданные')
    parser.add_argument('-p', '--pages', action='store_true', help='Показать количество страниц')
    parser.add_argument('-s', '--search', help='Искать текст')
    parser.add_argument('-o', '--output', help='Сохранить текст в файл')
    parser.add_argument('--password', help='Пароль для зашифрованного PDF')
    parser.add_argument('--page', type=int, help='Номер страницы (1-based)')
    args = parser.parse_args()

    try:
        reader = PDFReader(args.file, args.password)
    except Exception as e:
        print(colorize(f"Ошибка: {e}", Colors.RED))
        sys.exit(1)

    if args.metadata:
        meta = reader.get_metadata()
        print(colorize("Метаданные:", Colors.BOLD))
        for k, v in meta.items():
            print(f"  {colorize(k, Colors.CYAN)}: {v}")
        print()

    if args.pages:
        print(f"Количество страниц: {colorize(str(reader.get_page_count()), Colors.GREEN)}")
        print()

    if args.search:
        results = reader.search(args.search, args.page)
        if results:
            print(colorize(f"Найдено совпадений: {len(results)}", Colors.GREEN))
            for line_num, line in results:
                # Подсветка искомого слова
                highlighted = re.sub(f"({re.escape(args.search)})", colorize(r'\1', Colors.YELLOW), line, flags=re.IGNORECASE)
                print(f"  Строка {line_num}: {highlighted}")
        else:
            print(colorize("Совпадений не найдено.", Colors.RED))
        print()

    if args.text or (not args.metadata and not args.pages and not args.search):
        text = reader.extract_text(args.page)
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(text)
            print(colorize(f"Текст сохранён в {args.output}", Colors.GREEN))
        else:
            print(colorize("Извлечённый текст:", Colors.BOLD))
            print(text)

if __name__ == '__main__':
    main()
