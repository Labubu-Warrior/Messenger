package networkutil

import (
	"bufio"
	"strings"
	"testing"
)

const (
	TEST_DELIMITER      = '\n'
	TEST_MAX_LINE_BYTES = 8
	TEST_BUFFER_BYTES   = 4
)

func TestReadDelimitedLineAcceptsValidLine(t *testing.T) {
	reader := bufio.NewReaderSize(strings.NewReader("hello\n"), TEST_BUFFER_BYTES)

	line, tooLong, err := ReadDelimitedLine(reader, TEST_DELIMITER, TEST_MAX_LINE_BYTES)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if tooLong {
		t.Fatal("gültige Zeile wurde als zu lang erkannt")
	}
	if line != "hello\n" {
		t.Fatalf("unerwartete Zeile: %q", line)
	}
}

func TestReadDelimitedLineDiscardsLongLineAndKeepsNextRequest(t *testing.T) {
	reader := bufio.NewReaderSize(strings.NewReader("1234567890\nok\n"), TEST_BUFFER_BYTES)

	_, tooLong, err := ReadDelimitedLine(reader, TEST_DELIMITER, TEST_MAX_LINE_BYTES)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if !tooLong {
		t.Fatal("zu lange Zeile wurde nicht erkannt")
	}

	line, tooLong, err := ReadDelimitedLine(reader, TEST_DELIMITER, TEST_MAX_LINE_BYTES)
	if err != nil {
		t.Fatalf("zweite Zeile konnte nicht gelesen werden: %v", err)
	}
	if tooLong || line != "ok\n" {
		t.Fatalf("zweite Anfrage wurde nicht sauber erhalten: %q", line)
	}
}
