export async function runParallel(count, task) {
  const n = Math.max(1, Number(count) || 1)
  const results = await Promise.allSettled(Array.from({ length: n }, () => task()))
  const ok = results.filter((r) => r.status === 'fulfilled').length
  return { ok, fail: results.length - ok, total: results.length }
}
