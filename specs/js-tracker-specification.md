# Nyla Analytics - JavaScript Tracker Specification (Core)

## Overview

The Nyla JavaScript tracker is a lightweight, privacy-focused analytics script that collects essential web analytics data for single-site deployments. It's designed to be minimal, fast, and respectful of user privacy.



## Key Features

- Tiny bundle size (<5KB gzipped)
- No external dependencies
- Automatic pageview tracking
- Custom event support
- Privacy-first design
- CSP compatible
- Async by default
- TypeScript support

## Installation

### Script Tag

```html
<!-- Core self-hosted installation -->
<script async src="https://example.com/collect.js"></script>
<script>
  window.nyla = window.nyla || function(...args) {
    (window.nyla.q = window.nyla.q || []).push(args);
  };
  nyla('init', { 
    site: 'default',  // Always 'default' in core
    endpoint: 'https://example.com/v1/collect'
  });
</script>
```

### NPM Package

```bash
npm install @nyla/collect
```

```typescript
import { init } from '@nyla/collect';

init({ 
  site: 'default',
  endpoint: 'https://example.com/v1/collect'
});
```

## Configuration

```typescript
interface NylaConfig {
  site?: string;          // Site identifier (always 'default' in core)
  endpoint?: string;      // Collection endpoint (required for self-hosted)
  debug?: boolean;        // Enable debug logging
  privacy?: {
    respectDNT?: boolean; // Honor Do Not Track
    anonymizeIP?: boolean;// Anonymize IP addresses
    maskPII?: boolean;    // Mask potential PII in URLs
  };
  sampling?: {
    rate?: number;        // Sample rate (0-1)
  };
  meta?: {
    [key: string]: any;   // Custom metadata
  };
}


```

## API Reference

### Core Methods

#### init(config: NylaConfig)

Initialize the tracker with configuration.

```javascript
nyla('init', {
  site: 'abc123',
  debug: true,
  privacy: {
    respectDNT: true
  }
});
```

#### pageview([options: PageviewOptions])

Track a pageview event.

```javascript
nyla('pageview', {
  url: '/custom-path',
  title: 'Custom Title',
  referrer: 'https://example.com'
});
```

#### event(name: string, properties?: object)

Track a custom event.

```javascript
nyla('event', 'button_click', {
  button_id: 'signup',
  position: 'header'
});
```

#### identify(userId: string, traits?: object)

Associate events with a user ID (optional).

```javascript
nyla('identify', 'user123', {
  plan: 'pro',
  company: 'Acme Inc'
});
```

### Advanced Methods

#### setMeta(key: string, value: any)

Set persistent metadata for all subsequent events.

```javascript
nyla('setMeta', 'campaign', 'spring_2024');
```

#### consent(options: ConsentOptions)

Manage tracking consent.

```javascript
nyla('consent', {
  analytics: true,
  marketing: false
});
```

## Implementation Details

### Event Queue

Events are queued before tracker initialization:

```typescript
interface EventQueue {
  push(args: any[]): void;
  flush(): Promise<void>;
  clear(): void;
}
```

### Batch Processing

Events are batched for efficient transmission:

```typescript
interface BatchProcessor {
  maxSize: number;        // Maximum batch size
  flushInterval: number;  // Flush interval in ms
  add(event: Event): void;
  flush(): Promise<void>;
}
```

### Automatic Data Collection

Automatically collected data points:

- Page URL
- Page title
- Referrer
- Screen dimensions
- Viewport size
- User language
- Connection type
- Device type
- Operating system
- Browser information

### Error Handling

```typescript
interface ErrorHandler {
  handle(error: Error): void;
  log(message: string, level: LogLevel): void;
}
```

## Privacy Features

### Data Minimization

- No cookies by default
- No fingerprinting
- No cross-site tracking
- Minimal automatic collection

### PII Protection

- URL parameter masking
- Email/phone number detection
- Configurable PII patterns
- Data scrubbing options

### Compliance Helpers

- GDPR consent management
- CCPA opt-out support
- Cookie law compatibility
- Privacy policy helpers

## Performance

### Loading Strategy

1. Async script loading
2. Minimal initial payload
3. Feature lazy loading
4. Resource prioritization

### Bandwidth Usage

- Efficient event batching
- Compression
- Cache optimization
- Minimal payload size

### CPU/Memory Impact

- Throttled event processing
- Memory-efficient queue
- Background processing
- Resource cleanup

## Browser Support

### Modern Browsers

- Full feature support
- Performance optimizations
- Modern API usage

### Legacy Support

- Core functionality works in IE11+
- Graceful degradation
- Polyfill strategy

## Development

### Build System

- TypeScript
- Rollup
- Terser optimization
- Source maps

### Testing

- Jest for unit tests
- Cypress for E2E
- Browser compatibility
- Performance benchmarks

### Documentation

- TypeScript types
- JSDoc comments
- Usage examples
- Integration guides

## Security

### CSP Compatibility

```html
Content-Security-Policy: 
  script-src 'self' cdn.getnyla.app;
  connect-src api.getnyla.app;
```

### Request Signing

- Optional request signing
- Timestamp validation
- Nonce generation

### Error Handling

- Safe error logging
- PII scrubbing
- Rate limiting

## Integration

### Popular Frameworks

- React integration
- Vue plugin
- Angular service
- Next.js/Nuxt.js modules

### Build Tools

- Webpack plugin
- Vite plugin
- ESBuild configuration
- Source maps

### CMS Plugins

- WordPress plugin
- Shopify app
- Webflow integration
- Ghost integration 