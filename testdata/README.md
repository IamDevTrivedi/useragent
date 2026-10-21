# Test data

`conformance.json` contains User-Agent strings and expected parser results. User-Agent
strings are factual protocol data. The fixtures cover common desktop, mobile, tablet,
TV, console, wearable, bot, AI crawler, HTTP client, app, and WebView clients. They are
intentionally stored as data so that future contributors can add cases without editing
assertions.

`bench/baseline.txt` is a local benchmark snapshot for `benchstat` comparisons. It is
not a release gate and should be refreshed when fixtures or the hot path change.

See the repository `LICENSE` for the project copyright and license notice.
