/*
 * nyla-collector - GDPR compliant privacy focused web analytics
 * Copyright (C) 2024 Joe Purdy
 * mailto:nyla AT purdy DOT dev
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 3 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

// TypeScript types for config and event
interface NylaConfig {
  site: string;
  endpoint?: string;
  logLevel?: 'none' | 'warn' | 'info' | 'debug';
}

interface PageviewEvent {
  url: string;
  title: string;
  referrer?: string;
  timestamp?: string;
}

// Internal state
let config: NylaConfig | null = null;
let initialized = false;

// Logging utility
function log(level: 'debug' | 'info' | 'warn' | 'error', ...args: any[]) {
  const levels = { none: 0, warn: 1, info: 2, debug: 3 };
  const levelOrder = { error: 1, warn: 1, info: 2, debug: 3 };
  const logLevel = (config && config.logLevel) || 'warn';
  if (logLevel === 'none') return;
  if (levels[logLevel] >= levelOrder[level]) {
    if (level === 'error') {
      console.error(...args);
    } else if (level === 'warn') {
      console.warn(...args);
    } else {
      console.log(...args);
    }
  }
}

function getEndpoint(): string {
  return (config && config.endpoint) || 'https://api.getnyla.app';
}

function getSiteId(): string | null {
  return config && config.site ? config.site : null;
}

function getReferrer(): string {
  return document.referrer || '';
}

function getPageviewEvent(): PageviewEvent {
  return {
    url: window.location.href,
    title: document.title,
    referrer: getReferrer(),
    timestamp: new Date().toISOString(),
  };
}

function sendPageview(event: PageviewEvent) {
  if (!config || !config.site) {
    log('warn', '[nyla] No config or site set, not sending pageview');
    return;
  }

  const params = new URLSearchParams({
    site_id: config.site,
    type: 'pageview',
    url: event.url,
    title: event.title,
    referrer: event.referrer || '',
    timestamp: event.timestamp || new Date().toISOString(),
  });

  const endpoint = getEndpoint();
  const url = `${endpoint}/v1/collect?${params.toString()}`;
  log('debug', '[nyla] About to send pageview:', {
    endpoint,
    params: Object.fromEntries(params.entries()),
    url,
    event
  });

  // Keep a reference to avoid GC
  (window as any)._nylaImgs = (window as any)._nylaImgs || [];
  const img = new Image();
  (window as any)._nylaImgs.push(img);

  img.onload = function() {
    log('info', '[nyla] Pageview image loaded successfully:', url);
    // Remove from array to avoid memory leak
    (window as any)._nylaImgs = (window as any)._nylaImgs.filter((i: any) => i !== img);
  };
  img.onerror = function(e) {
    log('error', '[nyla] Pageview image failed to load:', url, e);
    (window as any)._nylaImgs = (window as any)._nylaImgs.filter((i: any) => i !== img);
  };

  img.src = url;
  log('debug', '[nyla] Image src set:', img.src);
}

function trackPageview() {
  sendPageview(getPageviewEvent());
}

function setupSPANavigation() {
  const his = window.history;
  if (his.pushState) {
    const originalPushState = his.pushState;
    his.pushState = function () {
      originalPushState.apply(this, arguments);
      trackPageview();
    };
    window.addEventListener('popstate', trackPageview);
  }
  window.addEventListener('hashchange', trackPageview, false);
}

function initNyla(userConfig: NylaConfig) {
  if (initialized) {
    log('info', '[nyla] Already initialized');
    return;
  }
  config = userConfig;
  initialized = true;
  log('info', '[nyla] Initialized with config:', config);
  setupSPANavigation();
  trackPageview();
}

// Global queue function
(function (w: any) {
  const q: any[] = [];
  function nyla(...args: any[]) {
    if (args[0] === 'init') {
      initNyla(args[1]);
    } else if (args[0] === 'pageview') {
      if (!initialized) q.push(args);
      else sendPageview({ ...getPageviewEvent(), ...(args[1] || {}) });
    } else {
      // Ignore unsupported calls for MVP
    }
  }
  nyla.q = q;
  w.nyla = nyla;

  // Fallback: auto-init if data-siteid is present and no explicit init
  if (!initialized) {
    const ds = document.currentScript?.dataset;
    if (ds && ds.siteid) {
      nyla('init', { site: ds.siteid });
    }
  }
})(window);