# portwatch

Lightweight CLI to monitor and alert on open ports and service changes on a host.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring open ports on the local host and get alerted when changes are detected:

```bash
portwatch watch
```

Specify a scan interval and alert on any new or closed ports:

```bash
portwatch watch --interval 30s --alert stdout
```

Scan a specific range of ports:

```bash
portwatch watch --ports 1-1024 --interval 60s
```

Check the current snapshot of open ports:

```bash
portwatch scan
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `60s` | How often to scan for port changes |
| `--ports` | `1-65535` | Port range to monitor |
| `--alert` | `stdout` | Alert destination (`stdout`, `log`, `webhook`) |
| `--host` | `localhost` | Target host to monitor |

## Example Output

```
[2024-01-15 10:32:01] NEW port open: 8080 (http-alt)
[2024-01-15 10:45:12] CLOSED port: 3000
[2024-01-15 11:00:00] No changes detected.
```

## License

MIT © 2024 portwatch contributors