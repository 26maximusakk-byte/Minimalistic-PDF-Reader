// pdf_reader.js
// Версия на JavaScript с использованием pdf-parse, commander, chalk

#!/usr/bin/env node

const { Command } = require('commander');
const fs = require('fs');
const chalk = require('chalk');
const pdfParse = require('pdf-parse');

class PDFReader {
    constructor(filename, password = null) {
        this.filename = filename;
        this.password = password;
        this.dataBuffer = null;
        this.metadata = null;
        this.pageCount = 0;
        this.text = '';
        this._load();
    }

    _load() {
        this.dataBuffer = fs.readFileSync(this.filename);
        // pdf-parse не поддерживает пароли напрямую, но для демонстрации оставим
    }

    async _parse() {
        const options = {};
        if (this.password) options.password = this.password;
        const data = await pdfParse(this.dataBuffer, options);
        this.metadata = data.metadata;
        this.pageCount = data.numpages;
        this.text = data.text;
        return data;
    }

    async getMetadata() {
        await this._parse();
        return this.metadata;
    }

    async getPageCount() {
        await this._parse();
        return this.pageCount;
    }

    async extractText(pageNum = null) {
        await this._parse();
        if (pageNum !== null) {
            // pdf-parse не поддерживает постраничное извлечение, но мы можем извлечь весь текст и разбить по страницам?
            // Вместо этого используем библиотеку pdf-lib для постраничного чтения, но для простоты вернём весь текст.
            console.warn(chalk.yellow('Постраничное извлечение не поддерживается, возвращён весь текст.'));
        }
        return this.text;
    }

    async search(query, pageNum = null) {
        const text = await this.extractText(pageNum);
        const lines = text.split('\n');
        const results = [];
        lines.forEach((line, idx) => {
            if (line.toLowerCase().includes(query.toLowerCase())) {
                results.push({ lineNum: idx + 1, line: line.trim() });
            }
        });
        return results;
    }
}

const program = new Command();
program
    .name('pdf_reader')
    .description('Minimalistic PDF Reader (JavaScript)')
    .argument('<file>', 'PDF файл')
    .option('-t, --text', 'Извлечь текст')
    .option('-m, --metadata', 'Показать метаданные')
    .option('-p, --pages', 'Показать количество страниц')
    .option('-s, --search <query>', 'Искать текст')
    .option('-o, --output <file>', 'Сохранить текст в файл')
    .option('--password <password>', 'Пароль для PDF')
    .option('--page <number>', 'Номер страницы', parseInt)
    .action(async (file, options) => {
        const reader = new PDFReader(file, options.password);
        try {
            if (options.metadata) {
                const meta = await reader.getMetadata();
                console.log(chalk.bold('Метаданные:'));
                for (const [k, v] of Object.entries(meta)) {
                    console.log(`  ${chalk.cyan(k)}: ${v}`);
                }
                console.log();
            }
            if (options.pages) {
                const count = await reader.getPageCount();
                console.log(`Количество страниц: ${chalk.green(count)}`);
                console.log();
            }
            if (options.search) {
                const results = await reader.search(options.search, options.page);
                if (results.length) {
                    console.log(chalk.green(`Найдено совпадений: ${results.length}`));
                    results.forEach(({ lineNum, line }) => {
                        const highlighted = line.replace(new RegExp(options.search, 'gi'), chalk.yellow('$&'));
                        console.log(`  Строка ${lineNum}: ${highlighted}`);
                    });
                } else {
                    console.log(chalk.red('Совпадений не найдено.'));
                }
                console.log();
            }
            if (options.text || (!options.metadata && !options.pages && !options.search)) {
                const text = await reader.extractText(options.page);
                if (options.output) {
                    fs.writeFileSync(options.output, text, 'utf-8');
                    console.log(chalk.green(`Текст сохранён в ${options.output}`));
                } else {
                    console.log(chalk.bold('Извлечённый текст:'));
                    console.log(text);
                }
            }
        } catch (err) {
            console.error(chalk.red(`Ошибка: ${err.message}`));
            process.exit(1);
        }
    });

program.parse(process.argv);
