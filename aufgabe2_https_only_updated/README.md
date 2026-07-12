# Aufgabe 2 – Terminvorschläge aus dem Stundenplan

Diese Version verwendet **keine lokale XML-Datei**. Sowohl der Hauptserver als auch der eigene Wertepaar-Server laden bei jeder fachlichen Anfrage die aktuelle XML-Datei direkt über HTTPS:

```text
https://intern.fh-wedel.de/~tho/api/splan.php
```

## Start

Drei Terminals öffnen:

```bash
go mod tidy
```

```bash
go run ./cmd/valuepairserver
```

```bash
go run ./cmd/server
```

```bash
go run ./cmd/client
```

## Datenfluss

```text
GUI-Client
    -> Hauptserver
       -> aktuelle XML-Datei über HTTPS
       -> eigener Wertepaar-Server
          -> aktuelle XML-Datei über HTTPS
       <- Index-Wertepaare
    <- lesbare Terminvorschläge
```

Q1- und Q2-Veranstaltungen werden nicht zeitlich gefiltert und dadurch wie Langläufer behandelt.
