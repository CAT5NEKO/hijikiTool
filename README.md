# hijikiTool

Misskey/Firefish向けの自動投稿スケジューラー。

## セットアップ

```bash
go build -o hijiki ./cmd/hijiki
```

### 2. 認証情報（.env）

```bash
cp .env.example .env
```

| 変数 | 必須 | 説明 |
|------|------|------|
| `MISSKEY_HOST` | 必須 | インスタンスのホスト名 |
| `MISSKEY_TOKEN` | 必須 | APIトークン |
| `MISSKEY_VISIBILITY` | - | 公開範囲（デフォルト: `home`） |
| `MISSKEY_LOCAL_ONLY` | - | ローカル限定（デフォルト: `false`） |

### 3. スケジュール設定（config.json）

```bash
cp config.example.json config.json
```

| フィールド | 必須 | 説明 |
|------------|------|------|
| `id` | 必須 | 一意識別子 |
| `type` | 必須 | `daily` / `weekly` / `monthly` / `yearly` |
| `hour` | 必須 | 時（0-23） |
| `minute` | 必須 | 分（0-59） |
| `dayOfWeek` | weekly | 曜日（0=日〜6=土） |
| `dayOfMonth` | monthly/yearly | 日（1-31） |
| `month` | yearly | 月（1-12） |
| `content` | 必須 | 投稿内容 |

投稿タイミング: スケジュール時刻から1分以内に投稿されます
（許容時間: 1分）

## systemd（Linux）

```ini
[Unit]
Description=hijikiTool Scheduler
After=network.target

[Service]
Type=simple
WorkingDirectory=/path/to/hijikiTool
ExecStart=/path/to/hijikiTool/hijiki
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
