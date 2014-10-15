# Change Log

All notable changes to this project will be documented in this file.

## Soon

- netstat poller. Polls values from `/proc/net/netstat`

## 0.5.0 - 2014-10-09

### Added

- Units. Measurements now have units which the Librato outputter takes
  advantage of. No other outputters currently take advantage of this.

### Changed (Breaking)

- sockstat poller now uses lowercase protocol names in emitted
  metrics. Previously, it broke convention and used
  uppercase. (i.e. 'UDP' is now 'udp')
- The percentage calculations for the "mem" and "df" pollers resulted
  in values between 0 and 1. CPU percentages are between 0-100. They
  now all use 0-100 instead of 0-1.

## Prior Versions

Sadly, we didn't keep a proper changelog for previous versions. :(