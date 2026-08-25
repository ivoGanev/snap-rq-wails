#!/usr/bin/env python3
"""
Inject mock HTTP requests into the Wails SQLite database.

Use this to populate the app database with realistic data of varying sizes
so you can test the UI for layout issues, performance, and edge cases.

Examples:
    # Create 50 mock requests in the default app database
    python scripts/mock_requests.py

    # Create 200 requests in a custom database file
    python scripts/mock_requests.py --count 200 --db ./my-test.db

    # Small data only (good for quick smoke tests)
    python scripts/mock_requests.py --count 10 --min-len 5 --max-len 30

    # Large data only (good for stress testing the UI)
    python scripts/mock_requests.py --count 25 --min-len 500 --max-len 4000
"""

import argparse
import os
import random
import sqlite3
import string
import sys
from pathlib import Path

SCHEMA = """
CREATE TABLE IF NOT EXISTS http_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    body TEXT,
    request_headers TEXT,
    status_code INTEGER,
    response_id TEXT
);
CREATE INDEX IF NOT EXISTS idx_http_requests_name ON http_requests(name);
"""

METHODS = ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"]
STATUS_CODES = [200, 201, 204, 301, 302, 400, 401, 403, 404, 422, 500, 502, 503]


def default_db_path() -> Path:
    """Return the default Wails app database path for the current OS."""
    if sys.platform == "win32":
        base = Path(os.environ.get("LOCALAPPDATA", Path.home() / "AppData" / "Local"))
    elif sys.platform == "darwin":
        base = Path.home() / "Library" / "Application Support"
    else:
        base = Path(os.environ.get("XDG_DATA_HOME", Path.home() / ".local" / "share"))
    return base / "snap-rq-wails-v3" / "app.db"


def random_text(min_len: int, max_len: int) -> str:
    """Generate a random ASCII string of variable length."""
    length = random.randint(min_len, max_len)
    chars = string.ascii_letters + string.digits + string.punctuation + " \n"
    return "".join(random.choices(chars, k=length))


def random_name(min_len: int, max_len: int) -> str:
    """Generate a readable random name."""
    length = random.randint(min_len, max_len)
    words = ["user", "order", "product", "search", "login", "checkout", "profile", "api", "v1", "v2", "webhook", "auth", "token", "refresh", "list", "detail", "create", "update", "delete"]
    parts = []
    while len(" ".join(parts)) < length:
        parts.append(random.choice(words))
    name = " ".join(parts)
    if len(name) > max_len:
        name = name[:max_len].rstrip()
    return name


def random_url(min_len: int, max_len: int) -> str:
    """Generate a random URL, sometimes with query params and fragments."""
    hosts = ["api.example.com", "localhost:8080", "test.service.io", "dev.internal.net", "staging.gateway.org"]
    paths = ["users", "orders", "products", "search", "auth/login", "webhooks", "v1/resources", "v2/items"]
    url = f"https://{random.choice(hosts)}/{random.choice(paths)}"

    if random.random() < 0.7:
        url += f"/{random.randint(1, 99999)}"
    if random.random() < 0.5:
        params = "&".join(
            f"{random.choice(['q', 'id', 'page', 'limit', 'filter', 'sort'])}={random.randint(1, 9999)}"
            for _ in range(random.randint(1, 5))
        )
        url += f"?{params}"
    if random.random() < 0.2:
        url += "#section"

    # Pad or trim to requested size range
    if len(url) < min_len:
        url += "/" + "x" * (min_len - len(url) - 1)
    if len(url) > max_len:
        url = url[:max_len]
    return url


def random_body(min_len: int, max_len: int) -> str:
    """Generate a random JSON-like or plain text body."""
    choice = random.random()
    if choice < 0.3:
        # Empty / tiny body
        return "" if random.random() < 0.5 else "{}"
    if choice < 0.7:
        # JSON-like object
        size = random.randint(min_len, max_len)
        obj = {
            "id": random.randint(1, 999999),
            "name": random_name(5, 40),
            "active": random.choice([True, False]),
            "tags": [random_name(3, 12) for _ in range(random.randint(0, 10))],
            "payload": random_text(max(0, size - 100), max(0, size - 50)),
        }
        return str(obj).replace("'", '"')
    # Plain text / random payload
    return random_text(min_len, max_len)


def random_headers(min_len: int, max_len: int) -> str:
    """Generate random HTTP headers as a single string."""
    if random.random() < 0.2:
        return ""
    headers = []
    header_pool = [
        ("Content-Type", random.choice(["application/json", "application/xml", "text/plain", "multipart/form-data"])),
        ("Authorization", f"Bearer {random_text(20, 60)}"),
        ("X-Request-ID", random_text(10, 40)),
        ("Accept", "application/json"),
        ("User-Agent", "MockClient/1.0"),
        ("X-Custom-Header", random_text(5, 30)),
    ]
    random.shuffle(header_pool)
    for name, value in header_pool[: random.randint(1, len(header_pool))]:
        headers.append(f"{name}: {value}")
    text = "\n".join(headers)
    if len(text) < min_len:
        text += "\n" + random_text(max(0, min_len - len(text) - 1), max(0, max_len - len(text) - 1))
    if len(text) > max_len:
        text = text[:max_len].rstrip()
    return text


def ensure_schema(conn: sqlite3.Connection) -> None:
    """Create the http_requests table if it does not exist."""
    conn.executescript(SCHEMA)
    conn.commit()


def insert_mock_requests(conn: sqlite3.Connection, count: int, min_len: int, max_len: int) -> int:
    """Insert N mock requests with varied field sizes."""
    cursor = conn.cursor()
    inserted = 0
    for i in range(count):
        # Vary each field independently around the requested range to produce
        # a realistic mix of small, medium, and large values.
        field_min = random.randint(min_len, max(1, max_len // 4))
        field_max = random.randint(field_min, max(field_min, max_len))

        name = random_name(field_min, field_max)
        url = random_url(field_min, field_max)
        method = random.choice(METHODS)
        body = random_body(field_min, field_max)
        headers = random_headers(field_min, field_max)
        status_code = random.choice(STATUS_CODES)
        response_id = random_text(8, 64)

        cursor.execute(
            """
            INSERT INTO http_requests (name, url, method, body, request_headers, status_code, response_id)
            VALUES (?, ?, ?, ?, ?, ?, ?)
            """,
            (name, url, method, body, headers, status_code, response_id),
        )
        inserted += 1
        if i % 100 == 0:
            conn.commit()
            print(f"  inserted {i + 1}/{count}...")

    conn.commit()
    return inserted


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Inject mock HTTP requests into the Wails app SQLite database.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s --count 50
  %(prog)s --count 200 --db ./test.db
  %(prog)s --count 10 --min-len 5 --max-len 30
  %(prog)s --count 25 --min-len 500 --max-len 4000
""",
    )
    parser.add_argument(
        "--db",
        type=Path,
        default=default_db_path(),
        help="Path to the SQLite database file (default: Wails app data directory).",
    )
    parser.add_argument(
        "--count",
        type=int,
        default=50,
        help="Number of mock requests to insert (default: 50).",
    )
    parser.add_argument(
        "--min-len",
        type=int,
        default=5,
        help="Minimum length for variable text fields (default: 5).",
    )
    parser.add_argument(
        "--max-len",
        type=int,
        default=1000,
        help="Maximum length for variable text fields (default: 1000).",
    )
    parser.add_argument(
        "--seed",
        type=int,
        default=None,
        help="Random seed for reproducible data (default: None).",
    )
    parser.add_argument(
        "--clear",
        action="store_true",
        help="Clear existing http_requests rows before inserting new ones.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()

    if args.count < 0:
        print("error: --count must be >= 0", file=sys.stderr)
        return 1
    if args.min_len < 0 or args.max_len < 0 or args.min_len > args.max_len:
        print("error: --min-len and --max-len must be >= 0 and min <= max", file=sys.stderr)
        return 1

    if args.seed is not None:
        random.seed(args.seed)

    db_path = args.db.resolve()
    print(f"Database: {db_path}")
    print(f"Creating {args.count} mock request(s) with field lengths {args.min_len}-{args.max_len} chars...")

    # Ensure parent directory exists
    db_path.parent.mkdir(parents=True, exist_ok=True)

    conn = sqlite3.connect(db_path)
    try:
        ensure_schema(conn)

        if args.clear:
            conn.execute("DELETE FROM http_requests")
            conn.commit()
            print("Cleared existing http_requests rows.")

        inserted = insert_mock_requests(conn, args.count, args.min_len, args.max_len)
        print(f"Done. Inserted {inserted} mock request(s).")
    finally:
        conn.close()

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
