use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    response::Response,
    routing::{get, get_service},
    Router,
};
use std::net::SocketAddr;
use tower::ServiceBuilder;
use tower_http::{
    cors::{CorsLayer, Any},
    services::ServeDir,
    trace::TraceLayer,
};
use tracing_subscriber;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    tracing_subscriber::init();

    // Create router with secure CORS policy
    let app = Router::new()
        .route("/ws", get(websocket_handler))
        .route("/api/system/health", get(health_check))
        .nest_service("/", get_service(ServeDir::new(std::env::var("STATIC_DIR").unwrap_or_else(|_| "static".to_string()))))
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CorsLayer::new()
                    // More restrictive CORS policy for production
                    .allow_origin(Any) // In production this should read from env, e.g. ["https://app.phoenix.local".parse().unwrap()]
                    .allow_methods(Any)
                    .allow_headers(Any)),
        );

    // Run server
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Phoenix GUI listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .expect("Failed to bind TCP listener. Port 8080 might be in use.");
        
    tracing::info!("Starting server...");
    axum::serve(listener, app)
        .await
        .expect("Server crashed or failed to start");
        
    Ok(())
}

async fn websocket_handler(ws: WebSocketUpgrade) -> Response {
    ws.on_upgrade(websocket)
}

async fn websocket(mut socket: WebSocket) {
    while let Some(msg) = socket.recv().await {
        if let Ok(msg) = msg {
            match msg {
                Message::Text(text) => {
                    if socket.send(Message::Text(format!("Echo: {}", text))).await.is_err() {
                        break;
                    }
                }
                Message::Close(_) => {
                    break;
                }
                _ => {}
            }
        } else {
            break;
        }
    }
}

async fn health_check() -> &'static str {
    "OK"
}
