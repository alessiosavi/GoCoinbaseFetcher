# go-coinbase-fetcher

#### A simple go library delegated to store the history from coinbase.com for any given currency

## Functionalities:

- The `Downloader` have the following parameter: `Granularity`, `Pair`, `LimitDate`:
    - Granularity: 60, 300, 900, 3600, 21600, 86400;
    - Pair: coinbase pair;
    - Limit date: stop date (download all from now since `limit date`);
- The `Downloader` create a `csv` file using the input parameter, an
  example: `BTC-USD-86400-2015-01-01T00:00:00Z-2022-02-28T00:39:03Z.csv`;
- The `csv` file created is processed in streaming. Once that a `batch` of 300 rows have been downloaded, the result is
  store into the file; serialized as `csv`;

Example:

```go
// ...
var opts = datastructure.DownloadOpts{
Granularity: datastructure.GRANULARITY_DAY,
Pair:        "BTC-USD",
LimitDate:   time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC),
}
// ...
```

#### Once initialized, the downloader is delegated to download the data for the given pair in the given period

```go
download := opts.New()
filename := download.Download(nil)
```

#### A progressbar will show the number estimated finish time and other information:

```text
2021-08-09T00:30:00Z   7% |█████████████████                              | (65/836, 5 it/s) [15s:2m34s]
```

### The manager

The `Manager` is delegated to perform few operation on the file:

- Sort
- Drop duplicates

```go
var manager datastructure.Manager
manager.Sort(filename)
manager.DropDuplicates(filename)
```