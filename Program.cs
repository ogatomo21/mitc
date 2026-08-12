using System.Globalization;
using System.Text;

namespace Mitc;

internal static class Program
{
    private const string DefaultUser = "John Doe";

    public static int Main(string[] args)
    {
        Console.OutputEncoding = new UTF8Encoding(encoderShouldEmitUTF8Identifier: false);

        try
        {
            return Run(args);
        }
        catch (ArgumentException exception)
        {
            Console.Error.WriteLine($"mitc: {exception.Message}");
            Console.Error.WriteLine("Try 'mitc --help' for more information.");
            return 2;
        }
        catch (IOException exception)
        {
            Console.Error.WriteLine($"mitc: {exception.Message}");
            return 1;
        }
        catch (UnauthorizedAccessException exception)
        {
            Console.Error.WriteLine($"mitc: {exception.Message}");
            return 1;
        }
    }

    private static int Run(string[] args)
    {
        int? year = null;
        string? userToSave = null;
        string? temporaryUser = null;
        var fileName = "LICENSE";
        var fileNameSpecified = false;
        var printOnly = false;

        for (var index = 0; index < args.Length; index++)
        {
            var argument = args[index];

            switch (argument)
            {
                case "-h":
                case "--help":
                    PrintHelp();
                    return 0;

                case "--version":
                case "-v":
                    Console.WriteLine(
                        typeof(Program).Assembly
                            .GetCustomAttributes(typeof(System.Reflection.AssemblyInformationalVersionAttribute), inherit: false)
                            .OfType<System.Reflection.AssemblyInformationalVersionAttribute>()
                            .FirstOrDefault()?.InformationalVersion
                        ?? "v1.0");
                    return 0;

                case "--year":
                case "-y":
                    year = ParseYear(ReadOptionValue(args, ref index, "--year"));
                    break;

                case "--set-user":
                    userToSave = ReadOptionValue(args, ref index, "--set-user");
                    ValidateUser(userToSave);
                    break;

                case "--user":
                case "-u":
                    temporaryUser = ReadOptionValue(args, ref index, "--user");
                    ValidateUser(temporaryUser);
                    break;

                case "--filename":
                case "-f":
                    fileName = ReadOptionValue(args, ref index, "--filename");
                    ValidateFileName(fileName);
                    fileNameSpecified = true;
                    break;

                case "--print":
                case "-p":
                    printOnly = true;
                    break;

                default:
                    if (argument.StartsWith("--year=", StringComparison.Ordinal))
                    {
                        year = ParseYear(argument["--year=".Length..]);
                    }
                    else if (argument.StartsWith("--set-user=", StringComparison.Ordinal))
                    {
                        userToSave = argument["--set-user=".Length..];
                        ValidateUser(userToSave);
                    }
                    else if (argument.StartsWith("--user=", StringComparison.Ordinal))
                    {
                        temporaryUser = argument["--user=".Length..];
                        ValidateUser(temporaryUser);
                    }
                    else if (argument.StartsWith("--filename=", StringComparison.Ordinal))
                    {
                        fileName = argument["--filename=".Length..];
                        ValidateFileName(fileName);
                        fileNameSpecified = true;
                    }
                    else
                    {
                        throw new ArgumentException($"Unknown argument '{argument}'.");
                    }

                    break;
            }
        }

        if (userToSave is not null)
        {
            if (year is not null || temporaryUser is not null || fileNameSpecified || printOnly)
            {
                throw new ArgumentException("--set-user cannot be used together with generation options.");
            }

            var path = UserConfig.Save(userToSave);
            Console.WriteLine($"Default user set to '{userToSave}'.");
            Console.WriteLine($"Saved to {path}");
            return 0;
        }

        var selectedYear = year ?? DateTime.Now.Year;
        var user = temporaryUser ?? UserConfig.Load() ?? DefaultUser;
        var license = MitLicense.Create(selectedYear, user);

        if (printOnly)
        {
            if (fileNameSpecified)
            {
                throw new ArgumentException("--filename cannot be used together with --print.");
            }

            Console.Write(license);
            return 0;
        }

        if (File.Exists(fileName) && !ConfirmOverwrite(fileName))
        {
            Console.WriteLine("Cancelled.");
            return 0;
        }

        File.WriteAllText(fileName, license, new UTF8Encoding(encoderShouldEmitUTF8Identifier: false));
        Console.WriteLine($"Created {fileName}");
        return 0;
    }

    private static string ReadOptionValue(string[] args, ref int index, string option)
    {
        if (index + 1 >= args.Length)
        {
            throw new ArgumentException($"Option '{option}' requires a value.");
        }

        index++;
        return args[index];
    }

    private static int ParseYear(string value)
    {
        if (value.Length != 4 ||
            !int.TryParse(value, NumberStyles.None, CultureInfo.InvariantCulture, out var year) ||
            year is < 1 or > 9999)
        {
            throw new ArgumentException("Year must be a four-digit number between 0001 and 9999.");
        }

        return year;
    }

    private static void ValidateUser(string user)
    {
        if (string.IsNullOrWhiteSpace(user))
        {
            throw new ArgumentException("User name cannot be empty.");
        }

        if (user.Contains('\r') || user.Contains('\n'))
        {
            throw new ArgumentException("User name cannot contain a line break.");
        }
    }

    private static void ValidateFileName(string fileName)
    {
        if (string.IsNullOrWhiteSpace(fileName))
        {
            throw new ArgumentException("File name cannot be empty.");
        }
    }

    private static bool ConfirmOverwrite(string fileName)
    {
        Console.Write($"'{fileName}' already exists. Overwrite? [y/N]: ");
        var answer = Console.ReadLine();
        return answer?.Trim().Equals("y", StringComparison.OrdinalIgnoreCase) == true;
    }

    private static void PrintHelp()
    {
        Console.WriteLine(
            """
            Usage: mitc [options]

            Generates an MIT License and saves it to LICENSE by default.

            Options:
              -y, --year <YEAR>      Set the copyright year (default: current year)
              -u, --user <NAME>      Use a copyright holder for this run only
                  --set-user <NAME>  Save the default copyright holder
              -f, --filename <FILE>  Change the output file name
              -p, --print            Write the license to standard output only
              -h, --help             Show this help
              -v, --version          Show the version

            Examples:
              mitc
              mitc -y 2025 -u "Tomoya Ogawa"
              mitc --filename LICENSE.txt
              mitc --print
              mitc --set-user "Tomoya Ogawa"
            """);
    }
}

internal static class UserConfig
{
    private const string FileName = ".mitc.toml";

    public static string? Load()
    {
        var path = GetPath();
        if (!File.Exists(path))
        {
            return null;
        }

        foreach (var rawLine in File.ReadLines(path, Encoding.UTF8))
        {
            var line = rawLine.Trim();
            if (line.Length == 0 || line.StartsWith('#'))
            {
                continue;
            }

            var separator = line.IndexOf('=');
            if (separator < 0 || !line[..separator].Trim().Equals("user", StringComparison.Ordinal))
            {
                continue;
            }

            var value = line[(separator + 1)..].Trim();
            var user = ParseTomlString(value);
            if (string.IsNullOrWhiteSpace(user) || user.Contains('\r') || user.Contains('\n'))
            {
                throw new IOException($"Invalid user value in configuration file '{path}'.");
            }

            return user;
        }

        return null;
    }

    public static string Save(string user)
    {
        var path = GetPath();
        var content = $"user = \"{EscapeTomlString(user)}\"{Environment.NewLine}";
        File.WriteAllText(path, content, new UTF8Encoding(encoderShouldEmitUTF8Identifier: false));
        return path;
    }

    private static string GetPath()
    {
        // This override also makes the configuration behavior easy to verify without
        // writing to the caller's real home directory.
        var home = Environment.GetEnvironmentVariable("MITC_CONFIG_HOME");
        if (string.IsNullOrWhiteSpace(home))
        {
            home = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
        }

        if (string.IsNullOrWhiteSpace(home))
        {
            home = Environment.GetEnvironmentVariable("USERPROFILE")
                ?? Environment.GetEnvironmentVariable("HOME");
        }

        if (string.IsNullOrWhiteSpace(home))
        {
            throw new IOException("Could not determine the user home directory.");
        }

        return Path.Combine(home, FileName);
    }

    private static string EscapeTomlString(string value) => value
        .Replace("\\", "\\\\", StringComparison.Ordinal)
        .Replace("\"", "\\\"", StringComparison.Ordinal)
        .Replace("\t", "\\t", StringComparison.Ordinal);

    private static string ParseTomlString(string value)
    {
        if (value.Length < 2 || value[0] != '"' || value[^1] != '"')
        {
            throw new IOException("The 'user' setting must be a quoted TOML string.");
        }

        var result = new StringBuilder();
        for (var index = 1; index < value.Length - 1; index++)
        {
            var character = value[index];
            if (character != '\\')
            {
                result.Append(character);
                continue;
            }

            if (++index >= value.Length - 1)
            {
                throw new IOException("Invalid escape sequence in the 'user' setting.");
            }

            result.Append(value[index] switch
            {
                '\\' => '\\',
                '"' => '"',
                't' => '\t',
                _ => throw new IOException("Unsupported escape sequence in the 'user' setting.")
            });
        }

        return result.ToString();
    }
}

internal static class MitLicense
{
    public static string Create(int year, string user) =>
        $$"""
        MIT License

        Copyright (c) {{year.ToString("D4", CultureInfo.InvariantCulture)}} {{user}}

        Permission is hereby granted, free of charge, to any person obtaining a copy
        of this software and associated documentation files (the "Software"), to deal
        in the Software without restriction, including without limitation the rights
        to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
        copies of the Software, and to permit persons to whom the Software is
        furnished to do so, subject to the following conditions:

        The above copyright notice and this permission notice shall be included in all
        copies or substantial portions of the Software.

        THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
        IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
        FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
        AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
        LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
        OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
        SOFTWARE.
        """ + Environment.NewLine;
}
