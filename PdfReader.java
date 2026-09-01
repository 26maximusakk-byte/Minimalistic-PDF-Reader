// PdfReader.java
// Версия на Java с использованием Apache PDFBox

import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDDocumentInformation;
import org.apache.pdfbox.text.PDFTextStripper;
import org.apache.pdfbox.Loader;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.*;

public class PdfReader {
    private final String filename;
    private final String password;
    private PDDocument document;

    public PdfReader(String filename, String password) throws IOException {
        this.filename = filename;
        this.password = password;
        File file = new File(filename);
        if (password != null && !password.isEmpty()) {
            document = Loader.loadPDF(file, password);
        } else {
            document = Loader.loadPDF(file);
        }
    }

    public void close() throws IOException {
        if (document != null) document.close();
    }

    public Map<String, String> getMetadata() {
        PDDocumentInformation info = document.getDocumentInformation();
        Map<String, String> meta = new LinkedHashMap<>();
        meta.put("Title", info.getTitle());
        meta.put("Author", info.getAuthor());
        meta.put("Creator", info.getCreator());
        meta.put("Producer", info.getProducer());
        meta.put("CreationDate", info.getCreationDate() != null ? info.getCreationDate().toString() : "");
        meta.put("ModDate", info.getModificationDate() != null ? info.getModificationDate().toString() : "");
        return meta;
    }

    public int getPageCount() {
        return document.getPages().getCount();
    }

    public String extractText(int pageNum) throws IOException {
        PDFTextStripper stripper = new PDFTextStripper();
        if (pageNum > 0) {
            stripper.setStartPage(pageNum);
            stripper.setEndPage(pageNum);
        }
        return stripper.getText(document);
    }

    public Map<Integer, List<String>> search(String query, int pageNum) throws IOException {
        String text = extractText(pageNum);
        String[] lines = text.split("\n");
        Map<Integer, List<String>> results = new LinkedHashMap<>();
        for (int i = 0; i < lines.length; i++) {
            if (lines[i].toLowerCase().contains(query.toLowerCase())) {
                results.computeIfAbsent(i+1, k -> new ArrayList<>()).add(lines[i].trim());
            }
        }
        return results;
    }

    // ANSI-цвета (консольные)
    private static String colorize(String text, String color) {
        return color + text + "\u001B[0m";
    }

    public static void main(String[] args) {
        String filename = null;
        String password = null;
        boolean textFlag = false;
        boolean metadataFlag = false;
        boolean pagesFlag = false;
        String searchQuery = null;
        String outputFile = null;
        int pageNum = 0;

        // Простой парсинг аргументов
        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
                case "-t": textFlag = true; break;
                case "-m": metadataFlag = true; break;
                case "-p": pagesFlag = true; break;
                case "-s": if (i+1 < args.length) searchQuery = args[++i]; break;
                case "-o": if (i+1 < args.length) outputFile = args[++i]; break;
                case "--password": if (i+1 < args.length) password = args[++i]; break;
                case "--page": if (i+1 < args.length) pageNum = Integer.parseInt(args[++i]); break;
                default:
                    if (!args[i].startsWith("-")) filename = args[i];
            }
        }

        if (filename == null) {
            System.err.println("Укажите PDF файл.");
            System.exit(1);
        }

        try (PdfReader reader = new PdfReader(filename, password)) {
            if (metadataFlag) {
                Map<String, String> meta = reader.getMetadata();
                System.out.println(colorize("Метаданные:", "\u001B[1m"));
                for (Map.Entry<String, String> e : meta.entrySet()) {
                    System.out.printf("  %s: %s\n", colorize(e.getKey(), "\u001B[96m"), e.getValue());
                }
                System.out.println();
            }

            if (pagesFlag) {
                System.out.printf("Количество страниц: %s\n", colorize(String.valueOf(reader.getPageCount()), "\u001B[92m"));
                System.out.println();
            }

            if (searchQuery != null) {
                Map<Integer, List<String>> results = reader.search(searchQuery, pageNum);
                if (!results.isEmpty()) {
                    System.out.printf("%s\n", colorize("Найдено совпадений: " + results.size(), "\u001B[92m"));
                    for (Map.Entry<Integer, List<String>> e : results.entrySet()) {
                        for (String line : e.getValue()) {
                            String highlighted = line.replaceAll("(?i)" + searchQuery, colorize("$0", "\u001B[93m"));
                            System.out.printf("  Строка %d: %s\n", e.getKey(), highlighted);
                        }
                    }
                } else {
                    System.out.println(colorize("Совпадений не найдено.", "\u001B[91m"));
                }
                System.out.println();
            }

            if (textFlag || (!metadataFlag && !pagesFlag && searchQuery == null)) {
                String text = reader.extractText(pageNum);
                if (outputFile != null) {
                    Files.write(Paths.get(outputFile), text.getBytes());
                    System.out.println(colorize("Текст сохранён в " + outputFile, "\u001B[92m"));
                } else {
                    System.out.println(colorize("Извлечённый текст:", "\u001B[1m"));
                    System.out.println(text);
                }
            }
        } catch (Exception e) {
            System.err.println(colorize("Ошибка: " + e.getMessage(), "\u001B[91m"));
            System.exit(1);
        }
    }
}
