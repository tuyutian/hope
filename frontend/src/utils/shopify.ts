

export function redirectRemote(url: string, newContext = false) {
  newContext ? window.open(url, "_top") : (window.location.href = url);
}
