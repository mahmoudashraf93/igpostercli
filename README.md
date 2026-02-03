# igpost

igpost is a Go CLI for posting photos, reels, and carousels to Instagram using the official Graph API. It uploads local media to https://uguu.se for temporary hosting, then publishes via the Instagram Graph API.

## Installation

### Go install

```bash
go install github.com/mahmoud/igpostercli/cmd/igpost@latest
```

### Build from source

```bash
make build
./bin/igpost --help
```

## Setup (Meta + Instagram)

1. Create a Meta app in the Meta Developers dashboard and add the **Instagram Graph API** product.
2. Link a Facebook Page to your Instagram Business/Creator account.
3. Generate a short-lived user token with the required permissions (typically: `instagram_basic`, `instagram_content_publish`, `pages_show_list`).
4. Exchange the short-lived token for a long-lived token using `igpost token exchange`.
5. Use `igpost account` to fetch your Instagram Business/User ID.

## Usage

Global flags:

- `--profile`: Profile name (selects stored config + keychain token).
- `--user-id`: Override `IG_USER_ID` at runtime.
- `--page-id`: Override `IG_PAGE_ID` at runtime.
- `--business-id`: Override `IG_BUSINESS_ID` at runtime.

Do you need BOTH `IG_USER_ID` and `IG_PAGE_ID` to post?

- For posting, you only need `IG_USER_ID` + `IG_ACCESS_TOKEN`.
- `IG_PAGE_ID` is only needed to lookup the IG user ID (via `igpost account`).

So the typical flow is:

1. Get `IG_PAGE_ID` (one time).
2. Run:

```bash
igpost --page-id <PAGE_ID> account
```

3. It prints `IG_USER_ID=...`
4. From then on, use `IG_USER_ID` (you can keep `IG_PAGE_ID` around, but it’s not required for posting).

### Post a photo

```bash
igpost photo --file path/to/photo.jpg --caption "hello"
```

```bash
igpost photo --url https://example.com/photo.jpg --caption "hello"
```

### Post a reel

```bash
igpost reel --file path/to/video.mp4 --caption "hello"
```

### Post a carousel

```bash
igpost carousel --files img1.jpg img2.jpg --caption "hello"
```

### Token utilities

```bash
igpost token exchange --short-token "<short_token>"
igpost token debug
```

### Account utilities

```bash
igpost account
igpost owned-pages --business-id <BUSINESS_ID>
```

### Profile management (keychain-backed)

Profiles store non-secret values in `~/.config/igpostercli/config.json`, while access tokens are stored in the OS keychain.

Resolution order:

1. CLI flags (e.g. `--user-id`, `--page-id`, `--business-id`)
2. Selected profile (`--profile` or `IG_PROFILE`)
3. Environment variables / `.env`

```bash
igpost profile set brand-a --access-token "<token>" --user-id <IG_USER_ID> --page-id <PAGE_ID> --business-id <BUSINESS_ID>
igpost profile show brand-a
igpost profile list
igpost profile delete brand-a
```

## Environment variables

Set these in `.env` (see `.env.example`) or export them in your shell.

- `IG_APP_ID`: Meta app ID.
- `IG_APP_SECRET`: Meta app secret.
- `IG_ACCESS_TOKEN`: Long-lived Instagram User access token.
- `IG_PROFILE`: Profile name (default: `default`).
- `IG_PAGE_ID`: Facebook Page ID connected to your Instagram account.
- `IG_BUSINESS_ID`: Meta Business ID (for listing owned pages).
- `IG_USER_ID`: Instagram Business/User ID.
- `IG_GRAPH_VERSION`: Graph API version (default: `v19.0`).
- `IG_POLL_INTERVAL`: Polling interval for media processing (default: `5s`).
- `IG_POLL_TIMEOUT`: Polling timeout for media processing (default: `300s`).

## Token refresh

Long-lived tokens expire. When yours is near expiry:

1. Generate a new short-lived token in the Meta dashboard.
2. Run:

```bash
igpost token exchange --short-token "<short_token>"
```

3. Update `IG_ACCESS_TOKEN` with the new value.
