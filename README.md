# My-Clock

A collection of Go terminal projects for learning terminal graphics, CLI tools, and system monitoring.

## Projects

### 1. ASCII 7-Segment Clock (`cmd/clock`)

Displays the current time (hours:minutes) using ASCII art 7-segment digits with a blinking colon.

```bash
go run ./cmd/clock
```

### 2. Folder Backup Tool (`cmd/backup`)

A CLI utility that creates timestamped backups of a directory.

```bash
go run ./cmd/backup -src ~/myproject -dst ~/backups
```

### 3. Terminal Animations (`cmd/animate`)

Scrolling text and countdown timer animations with ANSI colors.

```bash
# Scrolling text
go run ./cmd/animate -mode scroll -text "Good Morning Dadi"

# Countdown timer
go run ./cmd/animate -mode countdown -seconds 60
```

### 4. Dev Dashboard (`cmd/dashboard`)

A terminal dashboard showing CPU usage, RAM, date/time with color-coded indicators.

```bash
go run ./cmd/dashboard
```

## Running Tests

```bash
go test ./...
```
