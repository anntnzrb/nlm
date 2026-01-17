# Troubleshooting

## Authentication fails with "no valid browser profiles found"

Try the full profile scan and debug output:

```bash
nlm auth --all --notebooks --debug
```

If your browser keeps locking its profile, use the original profile directory:

```bash
NLM_USE_ORIGINAL_PROFILE=1 nlm auth --all --notebooks --debug
```

Ensure you have a supported browser installed (Chrome, Chromium, Brave, Edge; Safari on macOS).

## Browser opens but auth is not detected

1. Log into NotebookLM in the opened browser window.
2. Re-run `nlm auth --all --notebooks --debug`.
3. If you know the profile, pin it with `nlm auth --profile "Profile 1"`.

## Tokens/cookies look stale

Remove the stored env file and re-auth:

```bash
rm -f ~/.nlm/env
nlm auth
```

## Still stuck?

Gather debug output and open an issue with:

```bash
nlm auth --debug
nlm -debug list
```
