# Contributing

Thanks for taking the time to look at this project.

If this project is useful to you, a star on the repo would be appreciated.

## A note on participation

Issues, discussions, and pull requests on this repository are for humans.

If you read the code, hit a problem, or have an idea, please open an issue or a PR.

Please do not file AI-generated bug reports, machine-summarized feature requests, or autonomously-opened pull requests.
If a model produced the text and no human read it before posting, it does not belong here. AI assistance as a tool is fine; an AI as the participant is not.

Issues and PRs that look automated will be closed without discussion.

## Development

```bash
make check   # lint, build, test, govulncheck
```

Tests replay recorded HTTP cassettes, so no credentials are needed. To
re-record against a live account (PII is stripped as each cassette is written):

```bash
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./internal/record --missing
```

A weekly [smoke test](.github/workflows/smoke.yml) runs `internal/smoke` against
the live API to catch drift the cassettes can't see.

## License

By submitting a contribution you agree that it will be licensed under the
project's [Apache 2.0 license](LICENSE).
