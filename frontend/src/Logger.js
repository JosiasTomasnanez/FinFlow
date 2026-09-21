import { trace, context } from '@opentelemetry/api';

const OTEL_COLLECTOR_URL = import.meta.env.VITE_OTEL_COLLECTOR_URL || '';

function currentTraceContext() {
    const span = trace.getSpan(context.active());
    if (!span) return {};
    const ctx = span.spanContext();
    return { trace_id: ctx.traceId, span_id: ctx.spanId };
}

export function logEvent(severity, attributes = {}) {
    const service = 'finflow-frontend';
    const jsonLogBody = {
        level: severity,
        service: service,
        time: new Date().toISOString(),
        ...attributes,
        ...currentTraceContext() // Inyecta trace_id y span_id al mismo nivel
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
