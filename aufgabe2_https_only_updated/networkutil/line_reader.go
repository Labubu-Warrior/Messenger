// Package networkutil enthält kleine, wiederverwendbare Hilfsfunktionen
// für die begrenzte Verarbeitung zeilenbasierter Netzwerkprotokolle.
package networkutil

import (
	"bufio"
	"errors"
	"io"
)

// ReadDelimitedLine liest bis einschließlich delimiter, ohne mehr als maxBytes
// Nutzdaten im Speicher zu sammeln. Ist die Zeile zu lang, wird der Rest bis
// zum delimiter verworfen, damit die nächste Anfrage sauber gelesen werden kann.
func ReadDelimitedLine(reader *bufio.Reader, delimiter byte, maxBytes int) (string, bool, error) {
	if reader == nil {
		return "", false, errors.New("kein Reader angegeben")
	}
	if maxBytes <= 0 {
		return "", false, errors.New("ungültige maximale Zeilenlänge")
	}

	line := make([]byte, 0, maxBytes)

	for {
		fragment, readErr := reader.ReadSlice(delimiter)

		if len(line)+len(fragment) > maxBytes {
			if errors.Is(readErr, bufio.ErrBufferFull) {
				if err := discardUntilDelimiter(reader, delimiter); err != nil {
					return "", true, err
				}
			} else if readErr != nil && !errors.Is(readErr, io.EOF) {
				return "", true, readErr
			}
			return "", true, nil
		}

		line = append(line, fragment...)

		switch {
		case readErr == nil:
			return string(line), false, nil
		case errors.Is(readErr, bufio.ErrBufferFull):
			continue
		default:
			return "", false, readErr
		}
	}
}

func discardUntilDelimiter(reader *bufio.Reader, delimiter byte) error {
	for {
		_, err := reader.ReadSlice(delimiter)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, bufio.ErrBufferFull):
			continue
		default:
			return err
		}
	}
}
