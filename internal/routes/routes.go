package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barca-strategos/phoenix/internal/routes/api/v1"
	"github.com/barca-strategos/phoenix/internal/websocket"
	"github.com/barca-strategos/phoenix/internal/compliance/iso27001"
	"github.com/barca-strategos/phoenix/internal/compliance/nist"
	"github.com/barca-strategos/phoenix/internal/compliance/gdpr"
	"github.com/barca-strategos/phoenix/internal/compliance/policy"
	"github.com/barca-strategos/phoenix/internal/compliance/vendor"
	"github.com/barca-strategos/phoenix/internal/compliance/soc2"
	"github.com/barca-strategos/phoenix/internal/compliance/pci_dss"
	"github.com/barca-strategos/phoenix/internal/risk"
	"github.com/barca-strategos/phoenix/internal/compliance/evidence"
	"github.com/barca-strategos/phoenix/internal/chat"
	"github.com/barca-strategos/phoenix/internal/agentic/broker"
	"github.com/barca-strategos/phoenix/internal/observability"
	"github.com/barca-strategos/phoenix/internal/middleware"
	"github.com/barca-strategos/phoenix/internal/auth"
	"github.com/barca-strategos/phoenix/internal/case"
	"github.com/barca-strategos/phoenix/internal/playbook"
	"github.com/barca-strategos/phoenix/internal/siem"
	"github.com/barca-strategos/phoenix/internal/asset"
	"github.com/barca-strategos/phoenix/internal/threatintel"
	"github.com/barca-strategos/phoenix/internal/tenant"
	"github.com/barca-strategos/phoenix/internal/reporting"
	"github.com/barca-strategos/phoenix/internal/integrations"
)

// Register wires up the HTTP routes for the API.
func Register(app *fiber.App) {
	// Initialize services
	hub := websocket.NewHub()
	go hub.Run()

	riskSvc := risk.New()
	evidenceSvc := evidence.New()
	caseSvc := case.New()
	playbookEngine := playbook.New()
	siemSvc := siem.New()
	assetSvc := asset.New()
	tiSvc := threatintel.New()
	tenantSvc := tenant.New()
	
	// Compliance services
	iso27001Svc := iso27001.New()
	nistSvc := nist.New()
	gdprSvc := gdpr.New()
	policySvc := policy.New()
	vendorSvc := vendor.New()
	soc2Svc := soc2.New()
	pciDssSvc := pci_dss.New()
	
	brk := broker.New()
	brk.RegisterTool(broker.Tool{
		Name:        "resolve_alert",
		Description: "Resolve a security alert",
		Inputs:      map[string]string{"alert_id": "string"},
		SafetyTier:  "medium",
	})
	slackBot := chat.NewSlackBot(brk)

	// Register pre-built playbooks
	playbookEngine.RegisterPlaybook(&playbook.Playbook{
		ID:          "phishing-response",
		Name:        "Phishing Response",
		Description: "Automated phishing incident response",
		StartNode:   "trigger",
		Nodes: []playbook.Node{
			{ID: "trigger", Type: playbook.NodeTypeTrigger, Config: map[string]interface{}{}, Next: []string{"isolate"}},
			{ID: "isolate", Type: playbook.NodeTypeAction, Config: map[string]interface{}{"action": "isolate_host"}, Next: []string{"delay"}},
			{ID: "delay", Type: playbook.NodeTypeDelay, Config: map[string]interface{}{"seconds": 30}, Next: []string{"notify"}},
			{ID: "notify", Type: playbook.NodeTypeAction, Config: map[string]interface{}{"action": "create_ticket"}, Next: []string{}},
		},
	})

	// Register SIEM connectors
	siemSvc.RegisterConnector("splunk", siem.NewSplunkConnector())

	// Initialize auth
	authSvc := auth.New("super-secret-key")
	ssoHandler := auth.NewSSOHandler(authSvc, "http://localhost:8080", "https://idp.example.com/sso")

	// Initialize metrics
	observability.InitMetrics()

	// Integrations
	serviceNow := integrations.NewServiceNow("https://instance.service-now.com", "admin", "password")
	jira := integrations.NewJira("https://instance.atlassian.net", "admin", "token")

	// Middleware (except SSO endpoints)
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/auth/saml" || c.Path() == "/auth/saml/callback" || c.Path() == "/auth/oidc" || c.Path() == "/auth/oidc/callback" {
			return c.Next()
		}
		return middleware.JWTAuth(authSvc)(c)
	})

	// Metrics endpoint
	app.Get("/metrics", func(c *fiber.Ctx) error {
		observability.MetricsHandler().ServeHTTP(c.Context(), c.Response())
		return nil
	})

	// SSO endpoints (no auth)
	v1.RegisterSSORoutes(app, ssoHandler)

	// WebSocket
	v1.RegisterWebSocketRoutes(app, hub)

	// API v1 routes
	v1.RegisterMonitoringRoutes(app)
	v1.RegisterComplianceRoutes(app)
	v1.RegisterAgenticRoutes(app)
	v1.RegisterRiskRoutes(app, riskSvc, hub.Broadcast)
	v1.RegisterEvidenceRoutes(app, evidenceSvc)
	v1.RegisterChatRoutes(app, slackBot)
	v1.RegisterCaseRoutes(app, caseSvc, hub.Broadcast)
	v1.RegisterPlaybookRoutes(app, playbookEngine)
	v1.RegisterSIEMRoutes(app, siemSvc, hub.Broadcast)
	v1.RegisterAssetRoutes(app, assetSvc, hub.Broadcast)
	v1.RegisterThreatIntelRoutes(app, tiSvc, hub.Broadcast)
	v1.RegisterTenantRoutes(app, tenantSvc, hub.Broadcast)
	v1.RegisterReportingRoutes(app, reporting.New(caseSvc, riskSvc, assetSvc))
	v1.RegisterIntegrationRoutes(app, serviceNow, jira)
	v1.RegisterISO27001Routes(app, iso27001Svc, hub.Broadcast)
	v1.RegisterNISTRoutes(app, nistSvc, hub.Broadcast)
	v1.RegisterGDPRRoutes(app, gdprSvc, hub.Broadcast)
	v1.RegisterPolicyRoutes(app, policySvc, hub.Broadcast)
	v1.RegisterVendorRoutes(app, vendorSvc, hub.Broadcast)
	v1.RegisterSOC2Routes(app, soc2Svc, hub.Broadcast)
	v1.RegisterPCIDSSRoutes(app, pciDssSvc, hub.Broadcast)

	// Health check (no auth required)
	app.Get("/api/system/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
