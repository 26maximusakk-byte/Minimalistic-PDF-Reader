Репозиторий Minimalistic PDF Reader
Описание
Minimalistic PDF Reader – это набор консольных утилит на 8 языках программирования для чтения и извлечения информации из PDF-файлов. Программа позволяет:

Извлекать текст со всех страниц (или выборочно).

Показывать метаданные (автор, заголовок, создатель, дата создания и модификации).

Выводить количество страниц.

Искать текст по документу (с подсветкой результатов).

Конвертировать PDF в текстовый файл (опционально).

Работать с защищёнными паролем PDF-файлами (ввод пароля через аргумент или интерактивно).

Проект создан как тестовый репозиторий для демонстрации реализации одной задачи на разных языках с использованием продвинутых возможностей: цветной вывод, обработка ошибок, поддержка больших файлов, кэширование результатов поиска.

Структура репозитория
text
pdf-reader/
├── README.md
├── pdf_reader.py          (Python)
├── pdf_reader.js          (JavaScript / Node.js)
├── pdf_reader.ts          (TypeScript)
├── pdf_reader.go          (Go)
├── PdfReader.java         (Java)
├── PdfReader.cs           (C#)
├── pdf_reader.php         (PHP)
└── pdf_reader.rb          (Ruby)
Установка и запуск
Для каждого языка требуются соответствующие библиотеки. Установите их перед запуском.

Язык	Зависимости	Команда запуска
Python	pypdf (или PyPDF2)	pip install pypdf
python pdf_reader.py [опции] <файл.pdf>
JavaScript	pdf-parse, commander, chalk	npm install pdf-parse commander chalk
node pdf_reader.js [опции] <файл.pdf>
TypeScript	pdf-parse, commander, chalk, @types/node	npm install -g ts-node
ts-node pdf_reader.ts [опции] <файл.pdf>
Go	github.com/unidoc/unidoc (или github.com/ledongthuc/pdf)	go get github.com/unidoc/unidoc
go run pdf_reader.go [опции] <файл.pdf>
Java	Apache PDFBox (скачать .jar или через Maven)	javac -cp pdfbox.jar PdfReader.java
java -cp .:pdfbox.jar PdfReader [опции] <файл.pdf>
C#	PdfPig (через NuGet)	dotnet add package PdfPig
dotnet run -- [опции] <файл.pdf>
PHP	smalot/pdfparser (через Composer)	composer require smalot/pdfparser
php pdf_reader.php [опции] <файл.pdf>
Ruby	pdf-reader gem	gem install pdf-reader
ruby pdf_reader.rb [опции] <файл.pdf>
Использование
text
pdf_reader [options] <pdf_file>
Опции
Опция	Описание
-t, --text	Извлечь и показать весь текст (по умолчанию, если не указано другое).
-m, --metadata	Показать метаданные документа.
-p, --pages	Показать только количество страниц.
-s, --search <строка>	Найти указанную строку в тексте и показать контекст.
-o, --output <файл>	Сохранить извлечённый текст в файл.
--password <пароль>	Пароль для зашифрованного PDF.
--page <номер>	Извлечь текст только с указанной страницы (нумерация с 1).
-h, --help	Показать справку.
Примеры
bash
# Извлечь весь текст и вывести в консоль
python pdf_reader.py document.pdf

# Показать метаданные
python pdf_reader.py -m document.pdf

# Найти слово "контракт" и показать контекст
python pdf_reader.py -s "контракт" document.pdf

# Сохранить текст в файл
python pdf_reader.py -o output.txt document.pdf
Особенности реализаций
Все версии поддерживают базовые операции: извлечение текста, метаданные, подсчёт страниц.

Цветной вывод (ANSI) для улучшения читаемости.

Обработка ошибок (неверный файл, пароль, повреждённый PDF).

Использование современных библиотек для работы с PDF.

Поддержка больших файлов (постраничное чтение).

Реализован поиск с учётом регистра (опционально).

Лицензия
MIT
