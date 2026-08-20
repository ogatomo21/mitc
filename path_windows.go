//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	hkeyCurrentUser   = syscall.Handle(0x80000001)
	keyQueryValue     = 0x0001
	keySetValue       = 0x0002
	regSz             = 1
	regExpandSz       = 2
	errorFileNotFound = 2
	errorSuccess      = 0
	hwndBroadcast     = 0xffff
	wmSettingChange   = 0x001a
	smtoAbortIfHung   = 0x0002
	smtoBlock         = 0x0001
	// A broadcast waits once per responsive top-level window. Keep this bounded
	// so a slow application cannot make `mitc path` appear to hang.
	environmentChangeTimeoutMilliseconds = 200
)

var (
	advapi32            = syscall.NewLazyDLL("advapi32.dll")
	regCreateKeyExW     = advapi32.NewProc("RegCreateKeyExW")
	regQueryValueExW    = advapi32.NewProc("RegQueryValueExW")
	regSetValueExW      = advapi32.NewProc("RegSetValueExW")
	regCloseKey         = advapi32.NewProc("RegCloseKey")
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	expandEnvironmentW  = kernel32.NewProc("ExpandEnvironmentStringsW")
	user32              = syscall.NewLazyDLL("user32.dll")
	sendMessageTimeoutW = user32.NewProc("SendMessageTimeoutW")
)

func executableDirectory() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not determine mitc's executable path: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(executable); resolveErr == nil {
		executable = resolved
	}
	directory, err := filepath.Abs(filepath.Dir(executable))
	if err != nil {
		return "", fmt.Errorf("could not determine mitc's executable directory: %w", err)
	}
	return filepath.Clean(directory), nil
}

func updateUserPath(directory string, action pathAction) (bool, error) {
	if err := validatePathAction(action); err != nil {
		return false, err
	}

	key, err := openEnvironmentKey()
	if err != nil {
		return false, err
	}
	defer regCloseKey.Call(uintptr(key))

	path, valueType, err := readPathValue(key)
	if err != nil {
		return false, err
	}
	parts := strings.Split(path, ";")
	entryIndex := make([]int, 0, len(parts))
	for index, part := range parts {
		if sameWindowsPath(part, directory) {
			entryIndex = append(entryIndex, index)
		}
	}

	changed := false
	switch action {
	case pathActionAdd:
		if len(entryIndex) == 0 {
			if path == "" {
				path = directory
			} else if strings.HasSuffix(path, ";") {
				path += directory
			} else {
				path += ";" + directory
			}
			changed = true
		}
	case pathActionRemove:
		if len(entryIndex) != 0 {
			filtered := make([]string, 0, len(parts)-len(entryIndex))
			for _, part := range parts {
				if !sameWindowsPath(part, directory) {
					filtered = append(filtered, part)
				}
			}
			path = strings.Join(filtered, ";")
			changed = true
		}
	}

	if !changed {
		return false, nil
	}
	if err := writePathValue(key, path, valueType); err != nil {
		return false, err
	}
	return true, nil
}

func openEnvironmentKey() (syscall.Handle, error) {
	name, err := syscall.UTF16PtrFromString(`Environment`)
	if err != nil {
		return 0, err
	}
	var key syscall.Handle
	result, _, _ := regCreateKeyExW.Call(
		uintptr(hkeyCurrentUser),
		uintptr(unsafe.Pointer(name)),
		0,
		0,
		0,
		keyQueryValue|keySetValue,
		0,
		uintptr(unsafe.Pointer(&key)),
		0,
	)
	if result != errorSuccess {
		return 0, windowsStatusError("open HKCU\\Environment", result)
	}
	return key, nil
}

func readPathValue(key syscall.Handle) (string, uint32, error) {
	name, err := syscall.UTF16PtrFromString("Path")
	if err != nil {
		return "", 0, err
	}
	var valueType uint32
	var byteCount uint32
	result, _, _ := regQueryValueExW.Call(
		uintptr(key),
		uintptr(unsafe.Pointer(name)),
		0,
		uintptr(unsafe.Pointer(&valueType)),
		0,
		uintptr(unsafe.Pointer(&byteCount)),
	)
	if result == errorFileNotFound {
		return "", regExpandSz, nil
	}
	if result != errorSuccess {
		return "", 0, windowsStatusError("read user PATH", result)
	}
	if byteCount == 0 {
		return "", valueType, nil
	}

	buffer := make([]byte, byteCount)
	result, _, _ = regQueryValueExW.Call(
		uintptr(key),
		uintptr(unsafe.Pointer(name)),
		0,
		uintptr(unsafe.Pointer(&valueType)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&byteCount)),
	)
	if result != errorSuccess {
		return "", 0, windowsStatusError("read user PATH", result)
	}
	if valueType != regSz && valueType != regExpandSz {
		return "", 0, fmt.Errorf("user PATH has unsupported registry type %d", valueType)
	}
	return utf16BytesToString(buffer[:byteCount]), valueType, nil
}

func writePathValue(key syscall.Handle, path string, valueType uint32) error {
	name, err := syscall.UTF16PtrFromString("Path")
	if err != nil {
		return err
	}
	value, err := syscall.UTF16FromString(path)
	if err != nil {
		return err
	}
	if valueType != regSz && valueType != regExpandSz {
		valueType = regExpandSz
	}
	result, _, _ := regSetValueExW.Call(
		uintptr(key),
		uintptr(unsafe.Pointer(name)),
		0,
		uintptr(valueType),
		uintptr(unsafe.Pointer(&value[0])),
		uintptr(len(value)*2),
	)
	if result != errorSuccess {
		return windowsStatusError("write user PATH", result)
	}
	return nil
}

func sameWindowsPath(left, right string) bool {
	left = normalizeWindowsPath(left)
	right = normalizeWindowsPath(right)
	return left != "" && strings.EqualFold(left, right)
}

func normalizeWindowsPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if expanded, err := expandWindowsEnvironment(value); err == nil {
		value = expanded
	}
	value = filepath.Clean(value)
	if len(value) > 3 {
		value = strings.TrimRight(value, `\/`)
	}
	return value
}

func expandWindowsEnvironment(value string) (string, error) {
	input, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return "", err
	}
	needed, _, callErr := expandEnvironmentW.Call(uintptr(unsafe.Pointer(input)), 0, 0)
	if needed == 0 {
		return "", callErr
	}
	buffer := make([]uint16, needed)
	result, _, callErr := expandEnvironmentW.Call(
		uintptr(unsafe.Pointer(input)),
		uintptr(unsafe.Pointer(&buffer[0])),
		needed,
	)
	if result == 0 {
		return "", callErr
	}
	return syscall.UTF16ToString(buffer), nil
}

func reloadPath() error {
	environment, err := syscall.UTF16PtrFromString("Environment")
	if err != nil {
		return err
	}
	var result uintptr
	ret, _, callErr := sendMessageTimeoutW.Call(
		hwndBroadcast,
		wmSettingChange,
		0,
		uintptr(unsafe.Pointer(environment)),
		smtoAbortIfHung|smtoBlock,
		environmentChangeTimeoutMilliseconds,
		uintptr(unsafe.Pointer(&result)),
	)
	if ret == 0 {
		if callErr == nil {
			callErr = errors.New("SendMessageTimeoutW failed")
		}
		return fmt.Errorf("notify Windows about the updated PATH: %w", callErr)
	}
	return nil
}

func utf16BytesToString(value []byte) string {
	words := make([]uint16, len(value)/2)
	for index := range words {
		words[index] = uint16(value[index*2]) | uint16(value[index*2+1])<<8
	}
	return strings.TrimSuffix(string(utf16.Decode(words)), "\x00")
}

func windowsStatusError(operation string, status uintptr) error {
	return fmt.Errorf("%s: %w", operation, syscall.Errno(status))
}
