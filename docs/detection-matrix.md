# Detection matrix

This document summarizes the built-in detection surface and the tests that protect it.
It is a guide to current behavior, not a guarantee that every User-Agent variant is
recognized.

## Browsers and clients

| Family | Examples |
|---|---|
| Mainstream desktop and mobile | Chrome, Chromium, Firefox, Safari, Mobile Safari, Edge, Opera |
| Alternative browsers | Brave, Vivaldi, Arc, DuckDuckGo, Whale, Coc Coc, Yandex Browser |
| Regional and vendor browsers | QQ Browser, UC Browser, HeyTapBrowser, Samsung Browser, Miui Browser, Huawei Browser, Baidu Browser, Sogou Browser, Quark, 360 Browser |
| Legacy and embedded | Internet Explorer, NetFront, Konqueror, CFNetwork, AppleWebKit, Android browser |
| Automation and bots | HeadlessChrome, Splash, Scrapy, Googlebot, Bingbot, YandexBot, Baiduspider, DuckDuckBot, Applebot, Bytespider, GoogleOther |
| HTTP clients | curl, Wget, Go-http-client, python-requests, okhttp, Java, axios, node-fetch, undici, PostmanRuntime |

## Operating systems

Windows, Windows Phone, Android, iOS, macOS, Linux, ChromeOS, Tizen, webOS, HarmonyOS,
KaiOS, BlackBerry, FreeBSD, OpenBSD, NetBSD, and Solaris.

Windows versions are reported as they appear in the User-Agent. A `Windows NT 10.0`
token remains `10.0`; the parser does not infer a marketing version.

## Device types

| Type | Coverage |
|---|---|
| Desktop | Default for recognized desktop platforms when no mobile, tablet, TV, console, or wearable marker is present |
| Mobile | Android phones, iPhones, Windows Phone, feature phones, and mobile app containers |
| Tablet | iPads, Android tablets, and tablet-specific vendor models |
| TV | Tizen, webOS, HbbTV, GoogleTV, Samsung TV, and smart-TV markers |
| Console | PlayStation, Xbox, and Nintendo |
| Wearable | Apple Watch, Wear OS, and watch-specific markers |
| Embedded | Available to custom device rules; built-in detection does not broadly classify embedded clients |
| Unknown | Input without a reliable device signal |

## Bots and applications

Bot categories include search crawlers, social previewers, AI crawlers, SEO tools,
monitoring services, security scanners, scrapers, HTTP clients, headless browsers, and
advertising crawlers. Social preview bots can retain the underlying browser.

Embedded app detection covers Facebook, Instagram, TikTok, WeChat, QQ, Seafile,
CFNetwork, and other common app containers. `App.InApp` and `WebView` are independent
signals.

## Test coverage

| Area | Primary tests |
|---|---|
| Behavior fixtures | `TestDetectionRegressions`, `TestEdgeCaseRegressions` |
| Exact JSON conformance cases | `TestConformanceCorpus` |
| Corpus breadth gates | `TestConformanceCoverage` |
| Broad browser and device matrix | `TestBroadCorpus` |
| CPU normalization | `TestCPUArchitectureRegressions` |
| TV fixtures | `TestAdditionalTVFixtures` |
| Rules and custom matching | `TestCustomBrowserRule`, `TestCustomRulePriority`, `TestCustomMatchForms`, `TestCustomRulesOverrideBuiltIns` |
| Client hints | `TestApplyClientHints*`, `TestParseBrandList` |
| HTTP integration | `uahttp` package tests |
| CLI | `cmd/ua` package tests |
| Fuzzing | `FuzzParse`, `FuzzTokenizer`, `FuzzParseVersion`, `FuzzParseBrandList`, `FuzzParseRequest` |

The test names above are the source of truth for local verification. No external
tracker is required to understand or extend the suite.
