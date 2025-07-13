# Nyla Analytics - API Specification (Core)

## Overview

The Nyla Analytics Core API is a hypermedia-driven service that handles event collection and HTML-based user interfaces for single-site analytics. It's designed for simplicity, performance, and privacy in self-hosted environments.

## API Endpoints

### Event Collection

#### GET /v1/collect

Collects a single event (typically a pageview or simple custom event) via HTTP GET. This is designed for maximum compatibility (e.g., email open tracking, image beacons, JS tracker fallback).

**Query Parameters:**
- `site_id` (string, optional): Site identifier. Must be "default" if provided. If omitted, defaults to "default"
- `type` (string, required): Event type (e.g., `pageview`, `event`)
- `url` (string, optional): Page URL
- `title` (string, optional): Page title
- `referrer` (string, optional): Referrer URL
- `timestamp` (string, optional): ISO8601 timestamp (defaults to server time if omitted)
- `metadata` (string, optional): Base64-encoded JSON or URL-encoded key-value pairs for custom metadata

**Error Responses:**
- `400 Bad Request`: If `site_id` is provided but not equal to "default"
  ```json
  {
    "error": "Invalid site_id. This instance only supports site_id='default'"
  }
  ```

**Authentication:**
- API key via `Authorization: Bearer ...` header (optional for GET, but recommended for production)

**Response:**
- By default, returns a 1x1 transparent GIF (for use as a tracking pixel)
- If `Accept: application/json` is sent, returns:
  ```json
  { "success": true }
  ```
- On error, returns a minimal error GIF or JSON error (if requested)

**Use Cases:**
- Email open tracking
- JS tracker fallback
- Image/script beacons

#### POST /v1/collect

Collects one or more events (batch or complex events) via HTTP POST. Accepts a JSON body for richer payloads and batch processing.

**Request Body:**
```json
{
  "events": [
    {
      "type": "pageview",
      "url": "https://app.getnyla.app/dashboard",
      "title": "Analytics Dashboard",
      "referrer": "https://getnyla.app",
      "timestamp": "2024-03-14T15:09:26Z",
      "metadata": {
        "screen_size": "1920x1080",
        "language": "en-US"
      }
    }
  ],
  "site_id": "default"
}
```

**Authentication:**
- API key via `Authorization: Bearer ...` header (required for POST)

**Response:**
```json
{
  "success": true,
  "processed": 1
}
```

**Use Cases:**
- JavaScript tracker (batch mode)
- Server-to-server event ingestion
- Bulk uploads

### Analytics Interface

#### GET /dashboard

Returns the main dashboard HTML interface.

Response:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Nyla Analytics Dashboard</title>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <script src="https://unpkg.com/hyperscript.org@0.9.12"></script>
</head>
<body>
    <div hx-get="/api/stats/realtime" 
         hx-trigger="load, every 30s"
         hx-swap="innerHTML">
        <!-- Initial stats will be loaded here -->
    </div>
</body>
</html>
```

#### GET /api/stats/realtime

Returns a partial HTML fragment with real-time statistics.

Response:
```html
<div class="stats-container">
    <div class="stat-card">
        <h3>Active Visitors</h3>
        <p class="value">42</p>
    </div>
    <div class="stat-card">
        <h3>Pageviews (30m)</h3>
        <p class="value">156</p>
    </div>
    <div hx-get="/api/top-pages"
         hx-trigger="revealed"
         hx-swap="innerHTML">
        <div class="loading-spinner"></div>
    </div>
</div>
```

#### GET /api/stats/historical

Returns historical analytics data as an HTML table or chart.

Query Parameters:
- `from`: ISO timestamp
- `to`: ISO timestamp
- `resolution`: hour|day|week|month
- `format`: table|chart (defaults to table)

Response:
```html
<table class="analytics-table">
    <thead>
        <tr>
            <th>Date</th>
            <th>Pageviews</th>
            <th>Visitors</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>2024-03-01</td>
            <td>1,234</td>
            <td>567</td>
        </tr>
        <!-- Additional rows -->
    </tbody>
</table>
```

### Site Settings (Core)

#### GET /settings

Display single-site settings interface (core supports one site only).

Response:
```html
<div class="settings-container">
    <div class="site-settings">
        <h3>Site Configuration</h3>
        <form hx-put="/settings" hx-swap="innerHTML" hx-target="#status">
            <label>
                <span>Site Name:</span>
                <input type="text" name="name" value="My Site" required>
            </label>
            <fieldset>
                <legend>Privacy Settings</legend>
                <label>
                    <input type="checkbox" name="ip_anonymization" checked>
                    IP Anonymization
                </label>
                <label>
                    <input type="number" name="retention_days" value="90">
                    Data Retention (days)
                </label>
            </fieldset>
            <button type="submit">Save Settings</button>
        </form>

    </div>
    <div id="status"></div>
</div>
```

## Real-time Updates

### Server-Sent Events

Connect to `/api/updates` for real-time analytics updates:

```html
<div hx-sse="connect:/api/updates">
    <div hx-sse="swap:visitor_count">
        Active Visitors: 0
    </div>
    <div hx-sse="swap:pageviews">
        Pageviews: 0
    </div>
</div>
```

Event Format:
```
event: visitor_count
data: {"count": 42}

event: pageview
data: {"url": "/dashboard", "timestamp": "2024-03-14T15:09:26Z"}
```

## Error Handling

Errors are returned as both HTML fragments and JSON:

HTML Response (default):
```html
<div class="error-message" role="alert">
    <h4>Invalid Request</h4>
    <p>The timestamp format is invalid in events[0]</p>
    <button onclick="this.parentElement.remove()">Dismiss</button>
</div>
```

JSON Response (with Accept: application/json):
```json
{
  "error": {
    "code": "invalid_request",
    "message": "Invalid event format",
    "details": {
      "field": "events[0].timestamp",
      "reason": "invalid datetime format"
    }
  }
}
```

## Rate Limiting

- Collection endpoints: 100 requests per minute per IP
- HTML endpoints: 60 requests per minute per IP
- SSE connections: 10 concurrent connections per IP

Headers included in responses:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`

## Authentication

Session-based authentication for HTML interfaces:

```html
<form hx-post="/login" hx-target="body">
    <input type="email" name="email" required>
    <input type="password" name="password" required>
    <button type="submit">Log In</button>
</form>
```

API key authentication for event collection:
```
Authorization: Bearer nyla_key_123...
```

## Progressive Enhancement

All features degrade gracefully:
1. Base: Server-rendered HTML forms and links
2. Enhancement: HTMX for dynamic updates
3. Optional: Hyperscript for client-side interactions
4. Fallback: Standard form submissions when JavaScript is disabled

## Content Security Policy

```
Content-Security-Policy: 
    default-src 'self';
    script-src 'self' unpkg.com;
    style-src 'self' 'unsafe-inline';
    connect-src 'self' api.getnyla.app;
    frame-ancestors 'none';
```

## Versioning

- API versioned via URL prefix (/v1/)
- HTML interfaces maintain backward compatibility
- Semantic versioning for breaking changes
- Deprecation notices via response headers
- Minimum 6 months notice for breaking changes 