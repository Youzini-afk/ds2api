'use strict';

const fs = require('fs');
const path = require('path');

const DEFAULT_BASE_HEADERS = Object.freeze({
  Host: 'chat.deepseek.com',
  'User-Agent': 'DeepSeek/1.6.11 Android/35',
  Accept: 'application/json',
  'Content-Type': 'application/json',
  'x-client-platform': 'android',
  'x-client-version': '1.6.11',
  'x-client-locale': 'zh_CN',
  'accept-charset': 'UTF-8',
});

const WEB_HEADER_OVERRIDES = Object.freeze({
  'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36',
  Accept: 'application/json, text/plain, */*',
  'Content-Type': 'application/json',
  Origin: 'https://chat.deepseek.com',
  Referer: 'https://chat.deepseek.com/',
  'accept-language': 'zh-CN,zh;q=0.9,en;q=0.8',
  'x-client-platform': 'web',
  'x-client-version': '1.0.0',
  'x-client-locale': 'zh_CN',
});

const DEFAULT_SKIP_PATTERNS = Object.freeze([
  'quasi_status',
  'elapsed_secs',
  'token_usage',
  'pending_fragment',
  'conversation_mode',
  'fragments/-1/status',
  'fragments/-2/status',
  'fragments/-3/status',
]);

const DEFAULT_SKIP_EXACT_PATHS = Object.freeze([
  'response/search_status',
]);

function loadSharedConstants() {
  const sharedPath = path.resolve(__dirname, '../../internal/deepseek/constants_shared.json');
  try {
    const raw = fs.readFileSync(sharedPath, 'utf8');
    const parsed = JSON.parse(raw);
    const baseHeaders = parsed && typeof parsed.base_headers === 'object' && !Array.isArray(parsed.base_headers)
      ? { ...DEFAULT_BASE_HEADERS, ...parsed.base_headers }
      : { ...DEFAULT_BASE_HEADERS };
    const skipPatterns = Array.isArray(parsed && parsed.skip_contains_patterns)
      ? parsed.skip_contains_patterns.filter((v) => typeof v === 'string' && v !== '')
      : [...DEFAULT_SKIP_PATTERNS];
    const skipExactPaths = Array.isArray(parsed && parsed.skip_exact_paths)
      ? parsed.skip_exact_paths.filter((v) => typeof v === 'string' && v !== '')
      : [...DEFAULT_SKIP_EXACT_PATHS];
    return {
      baseHeaders,
      skipPatterns,
      skipExactPaths,
    };
  } catch (_err) {
    return {
      baseHeaders: { ...DEFAULT_BASE_HEADERS },
      skipPatterns: [...DEFAULT_SKIP_PATTERNS],
      skipExactPaths: [...DEFAULT_SKIP_EXACT_PATHS],
    };
  }
}

const shared = loadSharedConstants();

module.exports = {
  BASE_HEADERS: Object.freeze(shared.baseHeaders),
  WEB_BASE_HEADERS: Object.freeze({ ...shared.baseHeaders, ...WEB_HEADER_OVERRIDES }),
  baseHeadersForProfile(profile) {
    const p = String(profile || '').trim().toLowerCase();
    if (p === 'web') {
      return { ...shared.baseHeaders, ...WEB_HEADER_OVERRIDES };
    }
    return { ...shared.baseHeaders };
  },
  SKIP_PATTERNS: Object.freeze(shared.skipPatterns),
  SKIP_EXACT_PATHS: new Set(shared.skipExactPaths),
};
