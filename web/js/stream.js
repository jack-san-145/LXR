/**
 * stream.js
 *
 * Reusable abstraction for the daemon's hijacked plain-text streaming
 * endpoints (/create, /start, /pull_image). Reads exactly the bytes the
 * server sends via ReadableStream + TextDecoder, line by line, as they
 * arrive — never setInterval(), never simulated progress percentages.
 *
 *   streamRequest(url, fetchOptions, { onLine, onComplete, onError })
 */
async function streamRequest(url, options, { onLine, onComplete, onError } = {}) {
  let res;
  try {
    res = await fetch(url, options);
  } catch (err) {
    onError && onError(new (window.ApiError || Error)('Could not reach the LXR web server.', { cause: err }));
    return;
  }

  if (!res.ok || !res.body) {
    let detail = '';
    try { detail = await res.text(); } catch (_) { /* ignore */ }
    onError && onError(new (window.ApiError || Error)(detail || `Request failed (${res.status})`, { status: res.status }));
    return;
  }

  const reader = res.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop(); // keep the incomplete trailing line for the next chunk
      for (const line of lines) {
        onLine && onLine(line.replace(/\r$/, ''));
      }
    }
    if (buffer.trim().length) {
      onLine && onLine(buffer.replace(/\r$/, ''));
    }
    onComplete && onComplete();
  } catch (err) {
    onError && onError(new (window.ApiError || Error)('The stream ended unexpectedly.', { cause: err }));
  }
}

window.streamRequest = streamRequest;
