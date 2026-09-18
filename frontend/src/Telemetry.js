import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } from '@opentelemetry/semantic-conventions';

const OTEL_COLLECTOR_URL = import.meta.env.VITE_OTEL_COLLECTOR_URL || '';
const BACKEND_URL = import.meta.env.VITE_API_URL || '';

function getDomain(url) {
    try {
        return new URL(url).hostname;
    } catch {
        return url.replace(/https?:\/\//, '').split('/')[0];
    }
}

export function initTelemetry() {
    const provider = new WebTracerProvider({
        resource: resourceFromAttributes({
            [ATTR_SERVICE_NAME]: 'finflow-frontend',
            [ATTR_SERVICE_VERSION]: '1.0.0',
        }),
        spanProcessors: [
            new BatchSpanProcessor(
                new OTLPTraceExporter({
                    url: `${OTEL_COLLECTOR_URL}/v1/traces`,
                })
            ),
        ],
    });

    provider.register({
        contextManager: new ZoneContextManager(),
    });

    const domain = getDomain(BACKEND_URL);

    registerInstrumentations({
        instrumentations: [
            new DocumentLoadInstrumentation(),
            new FetchInstrumentation({
                propagateTraceHeaderCorsUrls: [
                    new RegExp(`.*${domain}.*`),
                ],
                clearTimingResources: true,
            }),
        ],
    });
}