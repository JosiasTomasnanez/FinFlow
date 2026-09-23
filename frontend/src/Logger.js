const OTEL_COLLECTOR_URL = import.meta.env.VITE_OTEL_COLLECTOR_URL || '';

export function logEvent(severity, attributes = {}, span) {
    const service = 'finflow-frontend';
    const { traceId: trace_id, spanId: span_id } = span?.spanContext() ?? {};
    const jsonLogBody = {
        level: severity,
        service: service,
        time: new Date().toISOString(),
        trace_id,
        span_id,
        ...attributes,
    };

    const body = {
        resourceLogs: [{
            resource: { attributes: [{ key: 'service.name', value: { stringValue: service } }] },
            scopeLogs: [{
                logRecords: [{
                    timeUnixNano: String(Date.now() * 1e6),
                    severityText: severity,
                    body: { stringValue: JSON.stringify(jsonLogBody) },
                    attributes: []
                }],
            }],
        }],
    };

    fetch(`${OTEL_COLLECTOR_URL}/v1/logs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
        keepalive: true,
    }).catch(() => { });
}