# Go-Web Crawler

A concurrent web crawler and keyword search engine built from scratch in Go. Point it at a website, give it a keyword — it crawls hundreds of pages in parallel, scores each one by relevance using a priority queue, stores everything in a SQLite database, and returns ranked results in under a second.

Built to understand concurrent systems deeply: goroutines, channels, mutexes, deadlocks, and the architecture decisions that come with them.

---

## Demo

```bash
# Crawl a site, prioritizing pages about "mystery"
go run cmd/crawler/main.go --target "mystery"

# Search everything that was crawled
go run cmd/crawler/main.go --search "mystery"
```

**Output:**

http://books.toscrape.com/catalogue/category/books/mystery_3/index.html
http://books.toscrape.com/catalogue/the-murder-of-roger-ackroyd_684/index.html
http://books.toscrape.com/catalogue/and-then-there-were-none_679/index.html

---

## Performance

| Mode | Pages Crawled | Time |
|---|---|---|
| Single-threaded BFS | 100 | 10.4s |
| Worker pool (10 goroutines) | 100 | 2.4s |

**4.3x faster.** Network wait time — the bottleneck in any crawler — gets filled by other workers instead of wasted sitting idle.

---

## How It Evolved

This project went through three distinct versions, each solving a real problem the previous one had.

### Version 1: Single-threaded BFS

The first version was a plain breadth-first search crawler. One goroutine, one page at a time.
seed URL → enqueue → dequeue → fetch → tokenize → enqueue new URLs → repeat

It worked but it was slow. Every page fetch blocks until the server responds, so the program spends most of its time sitting idle waiting on the network. 100 pages took around 10.4 seconds.

### Version 2: Concurrent worker pool

To fix the idle time problem, the natural next step was parallelism — fetch multiple pages at the same time. But the first attempt caused a deadlock.

The original design had workers both reading URLs from the jobs channel and writing new URLs back to the same channel:
jobs channel (buffered, size 20)
workers: read URL → fetch → push new URLs back to same channel

With 10 workers all finishing pages at the same time, all 10 try to send new URLs into a full channel simultaneously. Every worker blocks trying to send. But nothing is reading from the channel because all workers are blocked trying to write. The program hangs forever.

The fix was architectural — separate the dispatch channel from the pending-work queue entirely. Workers push new URLs onto an in-memory queue (never blocks). A dedicated feeder goroutine reads from the queue and sends into the channel at a rate workers can consume. This dropped crawl time from 10.4 seconds to 2.4 seconds across 100 pages.

### Version 3: Priority queue and inverted index

Plain BFS visits pages in discovery order. That is fine for completeness but irrelevant to keyword search — you might crawl 50 unrelated pages before hitting a relevant one just because they were discovered first.

Replacing the plain queue with a binary max-heap priority queue meant URLs get scored by keyword relevance before being enqueued. The most relevant page is always crawled next, regardless of when it was found. This turned the crawler from a link-follower into something closer to a search engine.

The final piece was persistence — storing everything in a SQLite-backed inverted index so results survive the crawl and can be queried in sub-second time.
Version 1: single-threaded BFS → 10.4s / 100 pages
Version 2: concurrent worker pool → 2.4s / 100 pages (4.3x faster)
Version 3: priority queue + index → ranked, persistent, searchable results

---

## Architecture
                    ┌─────────────────┐
                    │    main.go      │
                    └────────┬────────┘
                             │
           ┌─────────────────┴──────────────────┐
           │                                    │
      --target                             --search
           │                                    │
           ▼                                    ▼
   ┌───────────────┐                  ┌──────────────────┐
   │  WorkerPool() │                  │  SearchTerms()   │
   └───────┬───────┘                  │  (SQLite query)  │
           │                          └──────────────────┘
 ┌─────────┴──────────┐
 │                    │
 ▼                    ▼
┌──────────┐ ┌───────────────┐
│ Feeder │ │ Priority │◄──── new URLs with scores
│ goroutine│◄─────│ Queue │
└────┬─────┘ │ (max-heap) │
│ └───────────────┘
│ jobs channel (buffered chan string)
▼
┌─────────────────────────────────────────┐
│ Worker 1 │ Worker 2 │ ... │ Worker N │ (concurrent)
└──────────────────────┬──────────────────┘
│
┌────────────┼────────────┐
▼ ▼ ▼
FetchData TokenizerURL ExtractWordCounts
│
▼
┌─────────────────┐
│ SQLite Database│
├─────────────────┤
│ pages table │
│ index_table │
└─────────────────┘

---

## How It Works

### Crawling
The worker pool launches N goroutines that all pull URLs from a shared jobs channel simultaneously. While one worker waits on a network response, the others are already fetching different pages. A separate feeder goroutine manages the flow between the priority queue and the jobs channel.

### Prioritization
Every URL gets scored by how many times the target keyword appears in it. URLs are pushed onto a max-heap priority queue — the highest scoring URL always gets crawled next. This converts plain BFS into best-first search.

### Indexing
After fetching each page, the crawler tokenizes the visible text and builds a word frequency map. Every word and its count gets stored in SQLite as an inverted index. Searching for a term is a single indexed SQL query, not a scan through every page.

---

## Data Structures

All three built from scratch.

**PriorityQueue** — binary max-heap backed by a slice. Push and Pop are both O(log n). Parent/child relationships maintained with index math: parent at `(i-1)/2`, children at `2i+1` and `2i+2`. Protected with `sync.RWMutex` — read operations use `RLock` for concurrent reads, writes use full `Lock`.

**Queue** — slice-backed FIFO used by the single-threaded BFS crawler. O(1) enqueue, O(1) dequeue via slice re-pointing.

**Set** — `map[string]bool` for O(1) URL deduplication. Prevents the crawler from revisiting pages it has already processed.

---

## Database Schema

```sql
CREATE TABLE pages (
    url        TEXT PRIMARY KEY,
    title      TEXT,
    crawled_at TEXT,
    score      INTEGER
);

CREATE TABLE index_table (
    term      TEXT,
    url       TEXT,
    frequency INTEGER,
    PRIMARY KEY (term, url)
);
```

Search query:
```sql
SELECT url FROM index_table
WHERE term = ?
ORDER BY frequency DESC
```

Pages where the search term appears most frequently rank highest. After crawling 100 pages: **4,000+ term-URL mappings** stored in the index.

---

## Design Decisions

### Why a priority queue instead of a plain queue?

Plain BFS visits pages in discovery order — whatever order links happen to appear in HTML. If you are looking for pages about "mystery", you do not want to crawl 50 unrelated pages first just because they were discovered earlier. The priority queue ensures the crawler always visits the most keyword-relevant page next, regardless of discovery order.

### The goroutine deadlock

The first version of the worker pool had workers sending new URLs back onto the same channel they were reading from. With 10 workers all finishing at the same time, all 10 tried to send into a full channel simultaneously and blocked. Nothing was reading from the channel because all workers were blocked trying to write. Total deadlock, confirmed by the program hanging indefinitely.

The fix was separating concerns: workers write to the priority queue, the feeder reads from the queue and writes to the channel. Verified race-free with `go run -race`.

### Why SQLite and not a flat file?

A flat file requires loading the entire index into memory to search it. SQLite lets you run indexed queries against disk-backed storage — search stays fast as the index grows. The `ORDER BY frequency DESC` ranking is one SQL clause. Without a database, that requires loading and sorting everything in memory on every search.

### Thread safety

The priority queue and set are accessed by multiple goroutines simultaneously. Both use `sync.Mutex` (priority queue uses `sync.RWMutex` to allow concurrent reads). The database uses `database/sql`'s built-in connection pool which handles concurrent access safely. All verified with `go test -race`.

---

## Tech Stack

- **Golang** — standard library for HTTP, concurrency, CLI flags
- **golang.org/x/net/html** — HTML tokenizer for link and text extraction
- **modernc.org/sqlite** — pure-Go SQLite driver, no CGO required
- **SQLite** — local database for page metadata and inverted index

No web framework. No ORM. No external crawler library.

---

## Setup

Requirements: Go 1.21+

```bash
git clone https://github.com/yourusername/GO-WebCrawler
cd GO-WebCrawler
go mod tidy
```

**Crawl a site:**
```bash
go run cmd/crawler/main.go --target "mystery"
```

**Search what you crawled:**
```bash
go run cmd/crawler/main.go --search "mystery"
```

**Run tests:**
```bash
go test ./...
```

**Check for race conditions:**
```bash
go test -race ./...
```

---

## Testing

Table-driven unit tests for all data structures:

```bash
go test ./internal/datastructures/...
```

Covers enqueue/dequeue ordering, set deduplication and lookup, priority queue heap ordering across multiple pushes, and pop returning items in score-descending order. All tests pass with `go test -race`.

---

## Project Structure
GO-WebCrawler/
├── cmd/
│ └── crawler/
│ └── main.go # CLI entry point — flags and wiring only
├── internal/
│ ├── datastructures/
│ │ ├── queue.go # Slice-backed FIFO queue
│ │ ├── queue_test.go
│ │ ├── set.go # map[string]bool set
│ │ ├── set_test.go
│ │ ├── priorityqueue.go # Binary max-heap
│ │ └── priorityqueue_test.go
│ ├── scraper/
│ │ ├── fetcher.go # HTTP client with timeout
│ │ ├── parser.go # HTML tokenizer, word counter, title extractor
│ │ ├── crawler.go # Single-threaded BFS crawler
│ │ └── workerpool.go # Concurrent worker pool
│ └── index/
│ ├── store.go # SQLite init, save page, save terms
│ └── search.go # Ranked keyword search
├── scratch/ # Practice programs used during development
├── go.mod
├── go.sum
├── README.md
└── LICENSE

---

## What I Would Add Next

- **robots.txt compliance** — parse and respect Disallow rules before crawling any domain
- **Per-domain rate limiting** — avoid sending too many concurrent requests to a single server
- **TF-IDF ranking** — replace raw term frequency with proper TF-IDF scoring for better search relevance
- **REST API** — expose search behind `GET /search?q=term` so other services can query the index
- **Distributed crawling** — replace the in-memory priority queue with a Redis-backed queue so multiple crawler instances can share work across machines