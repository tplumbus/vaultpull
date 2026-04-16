# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with your Vault instance and run `vaultpull` pointing at a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --path secret/data/myapp --output .env
```

This will pull all key/value pairs from the specified Vault path and write them to a local `.env` file.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to read from | *(required)* |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file without prompting | `false` |
| `--addr` | Vault server address | `$VAULT_ADDR` |

### Example Output

```dotenv
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
API_KEY=supersecretkey
DEBUG=false
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token or supported auth method

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername