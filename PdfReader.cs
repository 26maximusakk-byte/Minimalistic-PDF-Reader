// PdfReader.cs
// Версия на C# с использованием PdfPig, System.CommandLine (или простой парсинг)

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using UglyToad.PdfPig;
using UglyToad.PdfPig.Content;

namespace PdfReader
{
    class Program
    {
        // ANSI-цвета
        const string Reset = "\u001b[0m";
        const string Cyan = "\u001b[96m";
        const string Green = "\u001b[92m";
        const string Yellow = "\u001b[93m";
        const string Red = "\u001b[91m";
        const string Blue = "\u001b[94m";
        const string Bold = "\u001b[1m";

        static string Colorize(string text, string color) => color + text + Reset;

        static void Main(string[] args)
        {
            string filename = null;
            string password = null;
            bool textFlag = false;
            bool metadataFlag = false;
            bool pagesFlag = false;
            string searchQuery = null;
            string outputFile = null;
            int pageNum = 0;

            for (int i = 0; i < args.Length; i++)
            {
                switch (args[i])
                {
                    case "-t": textFlag = true; break;
                    case "-m": metadataFlag = true; break;
                    case "-p": pagesFlag = true; break;
                    case "-s": if (i + 1 < args.Length) searchQuery = args[++i]; break;
                    case "-o": if (i + 1 < args.Length) outputFile = args[++i]; break;
                    case "--password": if (i + 1 < args.Length) password = args[++i]; break;
                    case "--page": if (i + 1 < args.Length) pageNum = int.Parse(args[++i]); break;
                    default:
                        if (!args[i].StartsWith("-")) filename = args[i];
                        break;
                }
            }

            if (filename == null)
            {
                Console.WriteLine("Укажите PDF файл.");
                return;
            }

            try
            {
                using (var pdf = PdfDocument.Open(filename, new ParsingOptions { Password = password }))
                {
                    if (metadataFlag)
                    {
                        Console.WriteLine(Colorize("Метаданные:", Bold));
                        var info = pdf.Information;
                        var meta = new Dictionary<string, string>
                        {
                            ["Title"] = info.Title,
                            ["Author"] = info.Author,
                            ["Creator"] = info.Creator,
                            ["Producer"] = info.Producer,
                            ["CreationDate"] = info.CreationDate?.ToString() ?? "",
                            ["ModDate"] = info.ModifiedDate?.ToString() ?? ""
                        };
                        foreach (var kv in meta)
                            Console.WriteLine($"  {Colorize(kv.Key, Cyan)}: {kv.Value}");
                        Console.WriteLine();
                    }

                    if (pagesFlag)
                    {
                        Console.WriteLine($"Количество страниц: {Colorize(pdf.NumberOfPages.ToString(), Green)}");
                        Console.WriteLine();
                    }

                    if (searchQuery != null)
                    {
                        var results = new Dictionary<int, List<string>>();
                        int startPage = pageNum > 0 ? pageNum : 1;
                        int endPage = pageNum > 0 ? pageNum : pdf.NumberOfPages;
                        for (int p = startPage; p <= endPage; p++)
                        {
                            var page = pdf.GetPage(p);
                            string text = page.Text;
                            var lines = text.Split('\n');
                            for (int i = 0; i < lines.Length; i++)
                            {
                                if (lines[i].ToLower().Contains(searchQuery.ToLower()))
                                {
                                    if (!results.ContainsKey(p)) results[p] = new List<string>();
                                    results[p].Add(lines[i].Trim());
                                }
                            }
                        }
                        if (results.Count > 0)
                        {
                            Console.WriteLine(Colorize($"Найдено совпадений: {results.Values.Sum(l => l.Count)}", Green));
                            foreach (var kv in results)
                            {
                                foreach (var line in kv.Value)
                                {
                                    var highlighted = Regex.Replace(line, searchQuery, m => Colorize(m.Value, Yellow), RegexOptions.IgnoreCase);
                                    Console.WriteLine($"  Страница {kv.Key}, строка: {highlighted}");
                                }
                            }
                        }
                        else
                        {
                            Console.WriteLine(Colorize("Совпадений не найдено.", Red));
                        }
                        Console.WriteLine();
                    }

                    if (textFlag || (!metadataFlag && !pagesFlag && searchQuery == null))
                    {
                        var sb = new StringBuilder();
                        int start = pageNum > 0 ? pageNum : 1;
                        int end = pageNum > 0 ? pageNum : pdf.NumberOfPages;
                        for (int p = start; p <= end; p++)
                        {
                            var page = pdf.GetPage(p);
                            sb.Append(page.Text);
                        }
                        string fullText = sb.ToString();
                        if (outputFile != null)
                        {
                            File.WriteAllText(outputFile, fullText);
                            Console.WriteLine(Colorize($"Текст сохранён в {outputFile}", Green));
                        }
                        else
                        {
                            Console.WriteLine(Colorize("Извлечённый текст:", Bold));
                            Console.WriteLine(fullText);
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(Colorize($"Ошибка: {ex.Message}", Red));
            }
        }
    }
}
