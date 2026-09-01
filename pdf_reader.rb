# pdf_reader.rb
# Версия на Ruby с использованием gem pdf-reader, OptionParser, colorize

require 'pdf-reader'
require 'optparse'
require 'colorize'

class PDFReader
  attr_reader :filename, :password

  def initialize(filename, password = nil)
    @filename = filename
    @password = password
    @reader = PDF::Reader.new(filename, password: password)
  rescue => e
    raise "Не удалось открыть PDF: #{e.message}"
  end

  def metadata
    info = @reader.info || {}
    {
      'Title' => info[:Title] || '',
      'Author' => info[:Author] || '',
      'Creator' => info[:Creator] || '',
      'Producer' => info[:Producer] || '',
      'CreationDate' => info[:CreationDate] ? info[:CreationDate].to_s : '',
      'ModDate' => info[:ModDate] ? info[:ModDate].to_s : ''
    }
  end

  def page_count
    @reader.page_count
  end

  def extract_text(page_num = nil)
    if page_num
      if page_num < 1 || page_num > page_count
        raise "Номер страницы вне диапазона"
      end
      @reader.page(page_num).text
    else
      @reader.pages.map(&:text).join("\n")
    end
  end

  def search(query, page_num = nil)
    text = extract_text(page_num)
    lines = text.split("\n")
    results = {}
    lines.each_with_index do |line, idx|
      if line.downcase.include?(query.downcase)
        results[idx+1] = line.strip
      end
    end
    results
  end
end

# Парсинг аргументов
options = {}
OptionParser.new do |opts|
  opts.banner = "Использование: ruby pdf_reader.rb [опции] <pdf_file>"
  opts.on('-t', '--text', 'Извлечь текст') { options[:text] = true }
  opts.on('-m', '--metadata', 'Показать метаданные') { options[:metadata] = true }
  opts.on('-p', '--pages', 'Показать количество страниц') { options[:pages] = true }
  opts.on('-s', '--search QUERY', 'Искать текст') { |q| options[:search] = q }
  opts.on('-o', '--output FILE', 'Сохранить текст в файл') { |f| options[:output] = f }
  opts.on('--password PASS', 'Пароль для PDF') { |p| options[:password] = p }
  opts.on('--page NUM', Integer, 'Номер страницы') { |n| options[:page] = n }
end.parse!

filename = ARGV[0]
unless filename
  puts "Укажите PDF файл."
  exit 1
end

begin
  reader = PDFReader.new(filename, options[:password])

  if options[:metadata]
    puts "Метаданные:".bold
    reader.metadata.each do |k, v|
      puts "  #{k.colorize(:cyan)}: #{v}"
    end
    puts
  end

  if options[:pages]
    puts "Количество страниц: #{reader.page_count.to_s.colorize(:green)}"
    puts
  end

  if options[:search]
    results = reader.search(options[:search], options[:page])
    if results.any?
      puts "Найдено совпадений: #{results.size}".colorize(:green)
      results.each do |line_num, line|
        highlighted = line.gsub(/#{Regexp.escape(options[:search])}/i) { |m| m.colorize(:yellow) }
        puts "  Строка #{line_num}: #{highlighted}"
      end
    else
      puts "Совпадений не найдено.".colorize(:red)
    end
    puts
  end

  if options[:text] || (!options[:metadata] && !options[:pages] && !options[:search])
    text = reader.extract_text(options[:page])
    if options[:output]
      File.write(options[:output], text)
      puts "Текст сохранён в #{options[:output]}".colorize(:green)
    else
      puts "Извлечённый текст:".bold
      puts text
    end
  end
rescue => e
  puts "Ошибка: #{e.message}".colorize(:red)
  exit 1
end
