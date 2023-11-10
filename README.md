# FGK_PASMAS_backend

## Wichtiger Hinweis

Bitte beachten Sie, dass dieses Projekt noch in der Entwicklungsphase ist. Einige Funktionen sind möglicherweise noch nicht vollständig implementiert oder können Änderungen unterliegen.

## Voraussetzungen

Bevor Sie mit der Installation beginnen, stellen Sie sicher, dass Sie folgende Voraussetzungen erfüllt haben:

- Go 

## Installation

Folgen Sie diesen Schritten, um das Projekt zu installieren:

1. Klonen Sie das Repository:

   ```bash
   git clone https://github.com/MetaEMK/FGK_PASMAS_backend
   ```

3. Wechseln sie in den docker Ordner:
     ```bash
     cd FGK_PASMAS_backend/docker
     ```
4. Kopieren sie die .env.example Datei und passen Sie an:
    ```bash
    cp .env.example .env
    ```
    **Bitte passen Sie alle Werte wenn nötig an**

5. Starten der Datenbank:
    ```bash
    docker compose up -d

3. Wechseln Sie in das Source Verzeichnis:

   ```bash
   cd ../src
   ```

4. Installieren Sie die erforderlichen Abhängigkeiten:

   ```bash
   go mod tidy
   ```

5. Bauen Sie das Projekt oder führen Sie es testweise aus:

   ```bash
   go run main.go
   ```

   ODER

   ```bash
   go build
   ```
