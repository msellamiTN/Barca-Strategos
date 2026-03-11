use crate::core::*;
use crate::gui::*;
use axum::{
    extract::{Path, Query, State, WebSocketUpgrade},
    http::StatusCode,
    response::{Html, IntoResponse, Json},
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use tower::ServiceBuilder;
use tower_http::{
    cors::CorsLayer,
    services::ServeDir,
    trace::TraceLayer,
};
use uuid::Uuid;

/// Web Server for GUI System
/// Provides REST API and WebSocket connections for the cognitive collaboration interface

pub struct GUIWebServer {
    gui_system: Arc<RwLock<PhoenixGUISystem>>,
    server_config: WebServerConfig,
    user_sessions: Arc<RwLock<HashMap<String, UserWebSocketSession>>>,
}

impl GUIWebServer {
    pub fn new(gui_system: PhoenixGUISystem, config: WebServerConfig) -> Self {
        Self {
            gui_system: Arc::new(RwLock::new(gui_system)),
            server_config: config,
            user_sessions: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn start(&self) -> Result<(), GUIError> {
        let app = self.create_app().await?;

        let listener = tokio::net::TcpListener::bind(format!("{}:{}", 
            self.server_config.host, self.server_config.port))
            .await
            .map_err(|e| GUIError::ConfigurationError(e.to_string()))?;

        println!("🌐 Phoenix GUI Web Server starting on http://{}:{}", 
            self.server_config.host, self.server_config.port);

        axum::serve(listener, app)
            .await
            .map_err(|e| GUIError::ConfigurationError(e.to_string()))?;

        Ok(())
    }

    async fn create_app(&self) -> Result<Router, GUIError> {
        let gui_system = self.gui_system.clone();
        let user_sessions = self.user_sessions.clone();

        let app = Router::new()
            // Serve static files
            .nest_service("/static", ServeDir::new("static"))
            
            // Main GUI routes
            .route("/", get(serve_index_html))
            .route("/dashboard", get(get_dashboard_handler))
            .route("/dashboard/:user_id", get(get_user_dashboard_handler))
            .route("/cognitive-workspace", get(get_cognitive_workspace_handler))
            .route("/cognitive-workspace/:user_id", get(get_user_cognitive_workspace_handler))
            .route("/compliance-center", get(get_compliance_center_handler))
            .route("/compliance-center/:user_id", get(get_user_compliance_center_handler))
            .route("/agent-interaction", get(get_agent_interaction_handler))
            .route("/agent-interaction/:user_id", get(get_user_agent_interaction_handler))
            .route("/monitoring", get(get_monitoring_handler))
            .route("/monitoring/:user_id", get(get_user_monitoring_handler))
            
            // API routes
            .route("/api/user/session", post(create_user_session_handler))
            .route("/api/user/:user_id/preferences", get(get_user_preferences_handler))
            .route("/api/user/:user_id/preferences", post(update_user_preferences_handler))
            .route("/api/dashboard/:user_id", get(get_dashboard_api_handler))
            .route("/api/cognitive/:user_id/start-collaboration", post(start_collaboration_handler))
            .route("/api/agents/:user_id/status", get(get_agent_status_handler))
            .route("/api/agents/:user_id/message", post(send_agent_message_handler))
            .route("/api/compliance/:user_id/overview", get(get_compliance_overview_handler))
            .route("/api/monitoring/:user_id/metrics", get(get_monitoring_metrics_handler))
            .route("/api/analytics/:user_id", get(get_analytics_handler))
            .route("/api/system/health", get(get_system_health_handler))
            
            // WebSocket endpoint
            .route("/ws", get(websocket_handler))
            
            // State management
            .with_state(AppState {
                gui_system,
                user_sessions,
            })
            
            // Middleware
            .layer(
                ServiceBuilder::new()
                    .layer(TraceLayer::new_for_http())
                    .layer(CorsLayer::permissive())
            );

        Ok(app)
    }
}

// Application state shared across handlers
#[derive(Clone)]
struct AppState {
    gui_system: Arc<RwLock<PhoenixGUISystem>>,
    user_sessions: Arc<RwLock<HashMap<String, UserWebSocketSession>>>,
}

// HTTP Handlers

async fn serve_index_html() -> Html<&'static str> {
    Html(include_str!("../static/index.html"))
}

async fn get_dashboard_handler() -> Html<&'static str> {
    Html(include_str!("../static/dashboard.html"))
}

async fn get_user_dashboard_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<UnifiedDashboard>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let dashboard = gui_system.get_unified_dashboard(&user_id).await?;
    Ok(Json(dashboard))
}

async fn get_cognitive_workspace_handler() -> Html<&'static str> {
    Html(include_str!("../static/cognitive_workspace.html"))
}

async fn get_user_cognitive_workspace_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<CognitiveWorkspaceState>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let workspace = gui_system.cognitive_workspace.get_workspace_state(&user_id).await?;
    Ok(Json(workspace))
}

async fn get_compliance_center_handler() -> Html<&'static str> {
    Html(include_str!("../static/compliance_center.html"))
}

async fn get_user_compliance_center_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<ComplianceOverview>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let overview = gui_system.compliance_center.get_overview(&user_id).await?;
    Ok(Json(overview))
}

async fn get_agent_interaction_handler() -> Html<&'static str> {
    Html(include_str!("../static/agent_interaction.html"))
}

async fn get_user_agent_interaction_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<AgentStatusSummary>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let status = gui_system.agent_hub.get_agent_status(&user_id).await?;
    Ok(Json(status))
}

async fn get_monitoring_handler() -> Html<&'static str> {
    Html(include_str!("../static/monitoring.html"))
}

async fn get_user_monitoring_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<MonitoringMetrics>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let metrics = gui_system.monitoring_center.get_metrics(&user_id).await?;
    Ok(Json(metrics))
}

// API Handlers

#[derive(Deserialize)]
struct CreateUserSessionRequest {
    user: User,
    preferences: UserPreferences,
}

async fn create_user_session_handler(
    State(state): State<AppState>,
    Json(request): Json<CreateUserSessionRequest>,
) -> Result<Json<UserSession>, GUIError> {
    let mut gui_system = state.gui_system.write().await;
    let session = gui_system.create_user_session(request.user, request.preferences).await?;
    Ok(Json(session))
}

async fn get_user_preferences_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<UserPreferences>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let sessions = gui_system.user_sessions.read().await;
    let session = sessions.get(&user_id)
        .ok_or(GUIError::SessionNotFound(user_id))?;
    Ok(Json(session.preferences.clone()))
}

#[derive(Deserialize)]
struct UpdateUserPreferencesRequest {
    preferences: UserPreferences,
}

async fn update_user_preferences_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
    Json(request): Json<UpdateUserPreferencesRequest>,
) -> Result<Json<()>, GUIError> {
    let mut gui_system = state.gui_system.write().await;
    let mut sessions = gui_system.user_sessions.write().await;
    let session = sessions.get_mut(&user_id)
        .ok_or(GUIError::SessionNotFound(user_id))?;
    session.preferences = request.preferences;
    Ok(Json(()))
}

async fn get_dashboard_api_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<UnifiedDashboard>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let dashboard = gui_system.get_unified_dashboard(&user_id).await?;
    Ok(Json(dashboard))
}

#[derive(Deserialize)]
struct StartCollaborationRequest {
    participants: CollaborationParticipants,
}

async fn start_collaboration_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
    Json(request): Json<StartCollaborationRequest>,
) -> Result<Json<String>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let session_id = gui_system.cognitive_workspace
        .start_collaboration_session(&user_id, request.participants).await?;
    Ok(Json(session_id))
}

async fn get_agent_status_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<AgentStatusSummary>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let status = gui_system.agent_hub.get_agent_status(&user_id).await?;
    Ok(Json(status))
}

#[derive(Deserialize)]
struct SendAgentMessageRequest {
    interaction_id: String,
    message: AgentMessage,
}

async fn send_agent_message_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
    Json(request): Json<SendAgentMessageRequest>,
) -> Result<Json<AgentResponse>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let response = gui_system.agent_hub
        .send_agent_message(&user_id, &request.interaction_id, request.message).await?;
    Ok(Json(response))
}

async fn get_compliance_overview_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<ComplianceOverview>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let overview = gui_system.compliance_center.get_overview(&user_id).await?;
    Ok(Json(overview))
}

async fn get_monitoring_metrics_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<MonitoringMetrics>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let metrics = gui_system.monitoring_center.get_metrics(&user_id).await?;
    Ok(Json(metrics))
}

#[derive(Deserialize)]
struct AnalyticsQuery {
    time_range: Option<TimeRange>,
}

async fn get_analytics_handler(
    Path(user_id): Path<String>,
    State(state): State<AppState>,
    Query(query): Query<AnalyticsQuery>,
) -> Result<Json<ComplianceAnalytics>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let time_range = query.time_range.unwrap_or(TimeRange {
        start: chrono::Utc::now() - chrono::Duration::days(7),
        end: chrono::Utc::now(),
    });
    let analytics = gui_system.compliance_center.get_analytics(&user_id, time_range).await?;
    Ok(Json(analytics))
}

async fn get_system_health_handler(
    State(state): State<AppState>,
) -> Result<Json<SystemHealth>, GUIError> {
    let gui_system = state.gui_system.read().await;
    let health = gui_system.get_system_health().await?;
    Ok(Json(health))
}

// WebSocket Handler

async fn websocket_handler(
    ws: WebSocketUpgrade,
    State(state): State<AppState>,
) -> impl IntoResponse {
    ws.on_upgrade(|socket| handle_websocket(socket, state))
}

async fn handle_websocket(
    mut socket: axum::extract::ws::WebSocket,
    state: AppState,
) {
    let mut user_id: Option<String> = None;

    while let Some(msg) = socket.recv().await {
        match msg {
            Ok(axum::extract::ws::Message::Text(text)) => {
                if let Ok(ws_message) = serde_json::from_str::<WebSocketMessage>(&text) {
                    match ws_message.message_type {
                        WebSocketMessageType::Authenticate => {
                            if let Ok(auth_data) = serde_json::from_value::<AuthData>(ws_message.data) {
                                user_id = Some(auth_data.user_id.clone());
                                
                                // Register user session
                                let session = UserWebSocketSession {
                                    user_id: auth_data.user_id.clone(),
                                    socket: socket.clone(),
                                    connected_at: chrono::Utc::now(),
                                };
                                
                                let mut sessions = state.user_sessions.write().await;
                                sessions.insert(auth_data.user_id, session);
                                
                                // Send confirmation
                                let response = WebSocketMessage {
                                    message_type: WebSocketMessageType::Authenticated,
                                    data: serde_json::json!({"status": "success"}),
                                    timestamp: chrono::Utc::now(),
                                };
                                
                                if let Ok(response_text) = serde_json::to_string(&response) {
                                    let _ = socket.send(axum::extract::ws::Message::Text(response_text)).await;
                                }
                            }
                        },
                        WebSocketMessageType::Update => {
                            // Handle real-time updates
                            if let Some(ref uid) = user_id {
                                if let Ok(update) = serde_json::from_value::<GUIUpdate>(ws_message.data) {
                                    let mut gui_system = state.gui_system.write().await;
                                    let _ = gui_system.handle_gui_update(update).await;
                                }
                            }
                        },
                        WebSocketMessageType::Ping => {
                            // Respond to ping
                            let response = WebSocketMessage {
                                message_type: WebSocketMessageType::Pong,
                                data: serde_json::json!({"timestamp": chrono::Utc::now()}),
                                timestamp: chrono::Utc::now(),
                            };
                            
                            if let Ok(response_text) = serde_json::to_string(&response) {
                                let _ = socket.send(axum::extract::ws::Message::Text(response_text)).await;
                            }
                        },
                    }
                }
            },
            Ok(axum::extract::ws::Message::Close(_)) => {
                // Handle disconnect
                if let Some(ref uid) = user_id {
                    let mut sessions = state.user_sessions.write().await;
                    sessions.remove(uid);
                }
                break;
            },
            Err(_) => {
                // Handle error
                break;
            },
            _ => {}
        }
    }
}

// Supporting structures

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WebServerConfig {
    pub host: String,
    pub port: u16,
    pub static_files_path: String,
    pub max_connections: usize,
    pub enable_cors: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserWebSocketSession {
    pub user_id: String,
    pub socket: axum::extract::ws::WebSocket,
    pub connected_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WebSocketMessage {
    pub message_type: WebSocketMessageType,
    pub data: serde_json::Value,
    pub timestamp: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum WebSocketMessageType {
    Authenticate,
    Authenticated,
    Update,
    Ping,
    Pong,
    Error,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuthData {
    pub user_id: String,
    pub token: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TimeRange {
    pub start: chrono::DateTime<chrono::Utc>,
    pub end: chrono::DateTime<chrono::Utc>,
}

// Error handling for HTTP responses
impl IntoResponse for GUIError {
    fn into_response(self) -> axum::response::Response {
        let (status, error_message) = match self {
            GUIError::SessionNotFound(msg) => (StatusCode::NOT_FOUND, msg),
            GUIError::ComponentInitializationFailed(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
            GUIError::UpdateProcessingFailed(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
            GUIError::InterfaceCreationFailed(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
            GUIError::ConfigurationError(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
            GUIError::RenderingError(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
        };

        let body = Json(serde_json::json!({
            "error": error_message,
            "status": status.as_u16()
        }));

        (status, body).into_response()
    }
}
