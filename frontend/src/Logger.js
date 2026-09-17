import { trace, context } from '@opentelemetry/api';

const OTEL_COLLECTOR_URL = import.meta.env.VITE_OTEL_COLLECTOR_URL || '';

function currentTraceContext() {
    const span = trace.getSpan(context.active());
    if (!span) return {};
    const ctx = span.spanContext();
    return { trace_id: ctx.traceId, span_id: ctx.spanId };
}

export function logEvent(severity, message, attributes = {}) {
    const body = {
        resourceLogs: [{
            resource: { attributes: [{ key: 'service.name', value: { stringValue: 'finflow-frontend' } }] },
            scopeLogs: [{
                logRecords: [{
                    timeUnixNano: String(Date.now() * 1e6),
                    severityText: severity,
                    body: { stringValue: message },
                    attributes: Object.entries({ ...attributes, ...currentTraceContext() })
                        .map(([key, value]) => ({ key, value: { stringValue: String(value) } })),
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