
use std::env;
use std::sync::Arc;

use aws_config::BehaviorVersion;
use aws_sdk_eventbridge::{
    types::PutEventsRequestEntry,
    Client as EventBridgeClient,
};

use hmac::{Hmac, Mac};
use lambda_http::{
    service_fn,
    Body,
    Error,
    IntoResponse,
    Request,
    Response,
};
use sha2::Sha256;

type HmacSha256 = Hmac<Sha256>;

const SIGNATURE_HEADER: &str = "x-razorpay-signature";

#[derive(Clone)]
struct AppState {
    eventbridge: EventBridgeClient,
    event_bus_name: String,
    webhook_secret: Arc<str>,
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    /*
     * Everything here runs during Lambda initialization.
     *
     * The EventBridge client is created once for this execution
     * environment and then captured by the handler.
     *
     * Warm invocations reuse it.
     */
    let aws_config = aws_config::defaults(BehaviorVersion::latest())
        .load()
        .await;

    let eventbridge = EventBridgeClient::new(&aws_config);

    let event_bus_name = env::var("EVENT_BUS_NAME")
        .unwrap_or_else(|_| "alshifa-bus".to_string());

    let webhook_secret = env::var("RAZORPAY_WEBHOOK_SECRET")
        .expect("RAZORPAY_WEBHOOK_SECRET environment variable is required");

    let state = AppState {
        eventbridge,
        event_bus_name,
        webhook_secret: Arc::from(webhook_secret),
    };

    /*
     * The state is created outside the request handler.
     *
     * Lambda can reuse this execution environment for subsequent
     * requests, allowing the AWS SDK client to be reused.
     */
    lambda_http::run(service_fn(move |request| {
        let state = state.clone();

        async move {
            handle_webhook(request, state).await
        }
    }))
    .await?;

    Ok(())
}

async fn handle_webhook(
    request: Request,
    state: AppState,
) -> Result<impl IntoResponse, Error> {
    /*
     * Razorpay sends POST requests.
     */
    if request.method() != http::Method::POST {
        return Ok(response(
            405,
            "method not allowed",
        ));
    }

    /*
     * Get Razorpay's signature.
     */
    let signature = match request.headers().get(SIGNATURE_HEADER) {
        Some(value) => match value.to_str() {
            Ok(value) if !value.is_empty() => value,
            _ => {
                return Ok(response(
                    401,
                    "invalid signature",
                ));
            }
        },

        None => {
            return Ok(response(
                401,
                "missing signature",
            ));
        }
    };

    /*
     * IMPORTANT:
     *
     * Do not deserialize and re-serialize the JSON before
     * calculating the HMAC.
     *
     * Razorpay's signature is calculated against the original
     * webhook body.
     */
    let raw_body = request.body().as_ref();

    /*
     * Verify Razorpay HMAC signature before doing anything
     * with the webhook.
     */
    if !verify_signature(
        raw_body,
        signature,
        &state.webhook_secret,
    ) {
        return Ok(response(
            401,
            "invalid signature",
        ));
    }

    /*
     * Optional sanity check:
     *
     * The webhook is expected to contain JSON.
     *
     * We do NOT modify the body.
     */
    if serde_json::from_slice::<serde_json::Value>(raw_body).is_err() {
        return Ok(response(
            400,
            "invalid json",
        ));
    }

    /*
     * Convert the verified payload into UTF-8.
     *
     * Razorpay webhooks are JSON, so this should normally
     * succeed.
     */
    let detail = match std::str::from_utf8(raw_body) {
        Ok(body) => body.to_owned(),

        Err(_) => {
            return Ok(response(
                400,
                "invalid request body",
            ));
        }
    };

    /*
     * Publish the verified webhook to EventBridge.
     *
     * EventBridge requires Source, DetailType and Detail.
     */
    let entry = PutEventsRequestEntry::builder()
        .event_bus_name(&state.event_bus_name)
        .source("alshifa.payment")
        .detail_type("PaymentWebhookReceived")
        .detail(detail)
        .build();

    let result = state
        .eventbridge
        .put_events()
        .entries(entry)
        .send()
        .await;

    match result {
        Ok(output) => {
            /*
             * PutEvents can succeed at the API level while an
             * individual entry fails.
             *
             * Check FailedEntryCount.
             */
            if output.failed_entry_count() > 0 {
                eprintln!(
                    "EventBridge rejected webhook event: failed_entry_count={}",
                    output.failed_entry_count()
                );

                return Ok(response(
                    500,
                    "failed to publish event",
                ));
            }

            /*
             * Only acknowledge Razorpay after the event was
             * successfully accepted by EventBridge.
             */
            Ok(response(
                200,
                "ok",
            ))
        }

        Err(error) => {
            eprintln!(
                "EventBridge PutEvents failed: {error:?}"
            );

            /*
             * Return 500 so Razorpay does not receive a false
             * success response when our event was not accepted.
             */
            Ok(response(
                500,
                "failed to publish event",
            ))
        }
    }
}

/*
 * Razorpay webhook verification.
 *
 * HMAC-SHA256(
 *     webhook_secret,
 *     raw_request_body
 * )
 *
 * Razorpay sends the resulting hexadecimal digest in:
 *
 * X-Razorpay-Signature
 */
fn verify_signature(
    body: &[u8],
    received_signature: &str,
    secret: &str,
) -> bool {
    /*
     * Razorpay signatures are hexadecimal strings.
     */
    let received_bytes = match hex::decode(received_signature) {
        Ok(bytes) => bytes,

        Err(_) => {
            return false;
        }
    };

    let mut mac = match HmacSha256::new_from_slice(
        secret.as_bytes(),
    ) {
        Ok(mac) => mac,

        Err(_) => {
            return false;
        }
    };

    mac.update(body);

    /*
     * verify_slice performs the cryptographic verification.
     */
    mac.verify_slice(&received_bytes).is_ok()
}

/*
 * Small helper for HTTP responses.
 */
fn response(
    status: u16,
    message: &'static str,
) -> Response<Body> {
    Response::builder()
        .status(status)
        .header("content-type", "text/plain")
        .body(Body::Text(message.to_string()))
        .expect("failed to build HTTP response")
}
```

### `Cargo.toml`

Use this alongside it:

```toml
[package]
name = "alshifa-razorpay-webhook"
version = "0.1.0"
edition = "2021"

[dependencies]
lambda_http = "1.3"
tokio = { version = "1", features = ["macros"] }

aws-config = { version = "1", features = ["behavior-version-latest"] }
aws-sdk-eventbridge = "1"

serde_json = "1"

hmac = "0.12"
sha2 = "0.10"
hex = "0.4"
```

A few important details in this implementation:

* **EventBridge client is initialized before the handler**, not on every webhook.
* The original request body is used for HMAC verification.
* Signature verification happens **before** publishing.
* The Lambda doesn't connect to MongoDB or perform appointment/payment processing.
* `PutEvents`'s `failed_entry_count` is checked, rather than treating every HTTP-successful API response as successful event publication. The AWS SDK exposes this per-entry failure information.
* The EventBridge event has:

  * `source = alshifa.payment`
  * `detail-type = PaymentWebhookReceived`
  * `detail = original verified Razorpay JSON`

So downstream rules can filter on the event type without your Lambda understanding the payment business logic.

**Lambda environment variables:**

```text
RAZORPAY_WEBHOOK_SECRET=<your Razorpay webhook secret>
EVENT_BUS_NAME=alshifa-bus
```

**IAM permission for this Lambda:**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "events:PutEvents",
      "Resource": "arn:aws:events:ap-south-1:YOUR_ACCOUNT_ID:event-bus/alshifa-bus"
    }
  ]
}
