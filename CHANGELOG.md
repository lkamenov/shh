# Change Log

All notable changes to this project will be documented in this file.

## 0.9.8 - 2015-01-05

### Changed

- Bugfix: ETS table memory to report in words not bytes.

## 0.9.7 - 2015-01-05

### Changed

- Bugfix: Used memory is a gauge, report it as such.

## 0.9.6 - 2015-01-05

### Changed

- Bugfix: Renamed LinearSliceContainsString to SliceContainsString, removing
    the old implementation, so that we only match on exact matches

## 0.9.5 - 2015-01-05

### Added

- Redis poller.

### Changed

- Update Go Version to 1.4.
- Travis CI Updates

## 0.9.4 - 2014-12-16

### Changed

- Folsom Poller (folsom) now polls for ETS table size and memory.

## 0.9.3 - 2014-11-25

### Changed

- Folsom Poller (folsom) now uses native `folsom_cowboy` API calls to determine
  metric types.

## 0.9.2 - 2014-11-12

### Added

- Folsom Poller (folsom), which adds the ability to fetch metrics from a
  runinning Erlang system. See the [docs](POLLERS.md) for more info.

## 0.9.1 - 2014-11-11

- Retry all low-level delivery errors

### Changed (Breaking)

- Fixed a typo in the per-processes rss metric. byts -> bytes

## 0.9.0 - 2014-11-07

### Changed (Breaking)

- MetricsNameNormalizer was changed so some metrics names may change.
```
    Old # -> _
        - -> _
    New # -> .
        _ -> -
```

## 0.8.9 - 2014-11-06

- Actually catch io.EOFs so we can re-try them. Experience has shown they are
    most likely dropped connections.

## 0.8.8 - 2014-11-06

### Added

- Splunk Search Peers poller (splunksearchpeers), which adds the
  ability to monitor Splunk search clusters. See the
  [docs](POLLERS.md) for more info.

### Changed (Breaking)

- Removed LibratoNetworkTimeout, use NetworkTimeout
  (`SHH_NETWORK_TIMEOUT`) instead.
- Merged listen poller metrics into 3 (all parse errors are treated
  the same) and added prefix `meta` to them all.

### Changed (Non-breaking)

- Refactored listen poller stats counters to avoid mutex
- Removed a bunch of useless logging
- handleListenConnection now part of Listen poller instead of bare.


## 0.8.7 - 2014-11-03

- nagios2stats bugfix

## 0.8.6 - 2014-11-03

- `nagios3stats` poller

## 0.8.5 - 2014-10-31

- Fix reporting of per processes sys / user cpu

## 0.8.4 - 2014-10-30

- Ignore processes w/o names, likely due to a process exiting between
    enumerating the directory entries and reading /proc/<pid>/stat

## 0.8.3 - 2014-10-30

- Generate additional process stats for processes that match
    `SHH_PROCESSES_REGEX`.
- `SHH_PROCESSES_REGEX`: \A\z - Regex of process names to poll and extra
    additional measurements for
- `SHH_TICKS`: 100 - cpu ticks per second. Default should be correct
    for most systems. see `getconf CLK_TCK`. Temporary until we use
    cgo to get it
- `SHH_PAGE_SIZE`: 4096 - kernel page size. Default should be correct
    for most systems. See `getconf PAGESIZE`. Temporary until we use
    cgo to get it.

## 0.8.2 - 2014-10-28

- `SHH_LIBRATO_ROUND`: true - round measurement times to nearest
    interval during submission

## 0.8.1 - 2014-10-24

### Changed (Breaking)

- Some disk poller stats were incorrectly being reported as Gauges.

## 0.8.0 - 2014-10-22

### Added

- `SHH_DF_LOOP` introduced to avoid loopback mounts from showing up in df
    poller (default false)
- `SHH_FULL` introduced to add back full set of metrics from some pollers

### Changed

- Remove `DEFAULT_SELF_POLLER_MODE` in favor of SHH_FULL="self"
- `SHH_CPU_AGGR` defaults to true, eliminating cpu metrics for all cores by
    default
- `SHH_DF_TYPES` removes tmpfs by default
- Default pollers to minimal set. Utilize `SHH_FULL` to get full set of metrics

## 0.7.0 - 2014-10-22

- Update some defaults: `SHH_INTERVAL=60s` & `SHH_LIBRATO_BATCH_TIMEOUT=10s`

## 0.6.4 - 2014-10-22

- Bugfix

## 0.6.3 - 2014-10-22

- SHH_LIBRATO_BATCH_SIZE defaults to 500
- SHH_LIBRATO_NETWORK_TIMEOUT defaults to 5s
- SHH_LIBRATO_BATCH_TIMEOUT defaults to SHH_INTERVAL
- Librato Outputter timeout doesn't start until there is a measurement

## 0.6.2 - 2014-10-22

- Report to librato how many guages / counters are being reported in a batch

## 0.6.1 - 2014-10-22

- `$SHH_LISTEN_TIMEOUT` now controls the timeout on the socket.
- Librato outlet now reports a `User-Agent` header at the request of Librato.
- Timeout errors to the librato api are now reported.
- Better handling of ntpdate sub process error messages.

## 0.6.0 - 2014-10-20

- shh-value cli tool for interacting with the unix socket.
- Latest version of Go (1.3.3) used.
- use github.com/heroku/slog for structured logging (extracted from shh
    originally).
- LISTEN Poller documentation.
- Improved Listen Poller with support for types and units
- Use of Go's logger (over fmt.Println)

## 0.5.0 - 2014-10-09

### Added

- Units. Measurements now have units which the Librato outputter takes
    advantage of. No other outputters currently take advantage of this.

### Changed (Breaking)

- sockstat poller now uses lowercase protocol names in emitted metrics.
    Previously, it broke convention and used uppercase. (i.e. 'UDP' is now
    'udp')
- The percentage calculations for the "mem" and "df" pollers resulted in values
    between 0 and 1. CPU percentages are between 0-100. They now all use 0-100
    instead of 0-1.

## Prior Versions

Sadly, we didn't keep a proper changelog for previous versions. :(
