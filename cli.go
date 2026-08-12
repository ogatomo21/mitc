// CLI behavior for mitc.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultUser = "John Doe"

type options struct {
	year              int
	hasYear           bool
	userToSave        string
	temporaryUser     string
	filename          string
	filenameSpecified bool
	printOnly         bool
	showHelp          bool
	showVersion       bool
}

// Run executes mitc with the supplied streams and returns a process exit code.
func Run(args []string, in io.Reader, out, errOut io.Writer, version string) int {
	opts, err := parseOptions(args)
	if err != nil {
		fmt.Fprintf(errOut, "mitc: %s\nTry 'mitc --help' for more information.\n", err)
		return 2
	}

	if opts.showHelp {
		printHelp(out)
		return 0
	}
	if opts.showVersion {
		fmt.Fprintln(out, version)
		return 0
	}
	if opts.userToSave != "" {
		if opts.hasYear || opts.temporaryUser != "" || opts.filenameSpecified || opts.printOnly {
			return argumentError(errOut, "--set-user cannot be used together with generation options.")
		}
		path, saveErr := saveUser(opts.userToSave)
		if saveErr != nil {
			return operationError(errOut, saveErr)
		}
		fmt.Fprintf(out, "Default user set to '%s'.\nSaved to %s\n", opts.userToSave, path)
		return 0
	}
	if opts.printOnly && opts.filenameSpecified {
		return argumentError(errOut, "--filename cannot be used together with --print.")
	}

	year := opts.year
	if !opts.hasYear {
		year = time.Now().Year()
	}
	user := opts.temporaryUser
	if user == "" {
		var loadErr error
		user, loadErr = loadUser()
		if loadErr != nil {
			return operationError(errOut, loadErr)
		}
		if user == "" {
			user = defaultUser
		}
	}
	license := mitLicense(year, user)
	if opts.printOnly {
		_, _ = io.WriteString(out, license)
		return 0
	}

	if _, statErr := os.Stat(opts.filename); statErr == nil {
		fmt.Fprintf(out, "'%s' already exists. Overwrite? [y/N]: ", opts.filename)
		answer, readErr := bufio.NewReader(in).ReadString('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return operationError(errOut, readErr)
		}
		if !strings.EqualFold(strings.TrimSpace(answer), "y") {
			fmt.Fprintln(out, "Cancelled.")
			return 0
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return operationError(errOut, statErr)
	}

	if writeErr := os.WriteFile(opts.filename, []byte(license), 0o644); writeErr != nil {
		return operationError(errOut, writeErr)
	}
	fmt.Fprintf(out, "Created %s\n", opts.filename)
	return 0
}

func parseOptions(args []string) (options, error) {
	opts := options{filename: "LICENSE"}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		value := func(option string) (string, error) {
			if i+1 >= len(args) {
				return "", fmt.Errorf("option '%s' requires a value", option)
			}
			i++
			return args[i], nil
		}
		switch arg {
		case "-h", "--help":
			opts.showHelp = true
		case "-v", "--version":
			opts.showVersion = true
		case "-p", "--print":
			opts.printOnly = true
		case "-y", "--year":
			v, err := value("--year")
			if err != nil {
				return opts, err
			}
			year, err := parseYear(v)
			if err != nil {
				return opts, err
			}
			opts.year, opts.hasYear = year, true
		case "-u", "--user":
			v, err := value("--user")
			if err != nil {
				return opts, err
			}
			if err := validateUser(v); err != nil {
				return opts, err
			}
			opts.temporaryUser = v
		case "--set-user":
			v, err := value("--set-user")
			if err != nil {
				return opts, err
			}
			if err := validateUser(v); err != nil {
				return opts, err
			}
			opts.userToSave = v
		case "-f", "--filename":
			v, err := value("--filename")
			if err != nil {
				return opts, err
			}
			if err := validateFilename(v); err != nil {
				return opts, err
			}
			opts.filename, opts.filenameSpecified = v, true
		default:
			name, v, hasValue := strings.Cut(arg, "=")
			if !hasValue {
				return opts, fmt.Errorf("unknown argument '%s'", arg)
			}
			switch name {
			case "--year":
				year, err := parseYear(v)
				if err != nil {
					return opts, err
				}
				opts.year, opts.hasYear = year, true
			case "--user":
				if err := validateUser(v); err != nil {
					return opts, err
				}
				opts.temporaryUser = v
			case "--set-user":
				if err := validateUser(v); err != nil {
					return opts, err
				}
				opts.userToSave = v
			case "--filename":
				if err := validateFilename(v); err != nil {
					return opts, err
				}
				opts.filename, opts.filenameSpecified = v, true
			default:
				return opts, fmt.Errorf("unknown argument '%s'", arg)
			}
		}
	}
	return opts, nil
}

func argumentError(out io.Writer, message string) int {
	fmt.Fprintf(out, "mitc: %s\nTry 'mitc --help' for more information.\n", message)
	return 2
}

func operationError(out io.Writer, err error) int {
	fmt.Fprintf(out, "mitc: %s\n", err)
	return 1
}

func parseYear(value string) (int, error) {
	if len(value) != 4 {
		return 0, errors.New("year must be a four-digit number between 0001 and 9999")
	}
	year, err := strconv.Atoi(value)
	if err != nil || year < 1 || year > 9999 {
		return 0, errors.New("year must be a four-digit number between 0001 and 9999")
	}
	return year, nil
}

func validateUser(user string) error {
	if strings.TrimSpace(user) == "" {
		return errors.New("user name cannot be empty")
	}
	if strings.ContainsAny(user, "\r\n") {
		return errors.New("user name cannot contain a line break")
	}
	return nil
}

func validateFilename(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("file name cannot be empty")
	}
	return nil
}

func mitLicense(year int, user string) string {
	return fmt.Sprintf(`MIT License

Copyright (c) %04d %s

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
`, year, user)
}

func configPath() (string, error) {
	home := os.Getenv("MITC_CONFIG_HOME")
	if strings.TrimSpace(home) == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", errors.New("could not determine the user home directory")
		}
	}
	if strings.TrimSpace(home) == "" {
		return "", errors.New("could not determine the user home directory")
	}
	return filepath.Join(home, ".mitc.toml"), nil
}

func saveUser(user string) (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}
	content := "user = \"" + escapeTOML(user) + "\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func loadUser() (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "user" {
			continue
		}
		user, err := parseTOMLString(strings.TrimSpace(value))
		if err != nil {
			return "", fmt.Errorf("invalid user value in configuration file '%s': %w", path, err)
		}
		if err := validateUser(user); err != nil {
			return "", fmt.Errorf("invalid user value in configuration file '%s'", path)
		}
		return user, nil
	}
	return "", nil
}

func escapeTOML(value string) string {
	return strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\t", "\\t").Replace(value)
}

func parseTOMLString(value string) (string, error) {
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return "", errors.New("the 'user' setting must be a quoted TOML string")
	}
	var b strings.Builder
	for i := 1; i < len(value)-1; i++ {
		if value[i] != '\\' {
			b.WriteByte(value[i])
			continue
		}
		i++
		if i >= len(value)-1 {
			return "", errors.New("invalid escape sequence in the 'user' setting")
		}
		switch value[i] {
		case '\\', '"':
			b.WriteByte(value[i])
		case 't':
			b.WriteByte('\t')
		default:
			return "", errors.New("unsupported escape sequence in the 'user' setting")
		}
	}
	return b.String(), nil
}

func printHelp(out io.Writer) {
	_, _ = io.WriteString(out, `Usage: mitc [options]

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
`)
}
