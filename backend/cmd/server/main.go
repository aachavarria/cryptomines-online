package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/handlers"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/workers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := database.InitSupabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start background workers
	go workers.StartBlueprintWorker()
	go workers.StartResourceWorker()
	go workers.StartCorpWorker()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /api/health", handlers.Health)
	mux.HandleFunc("GET /api/building-types", handlers.ListBuildingTypes)

	// Public reference data (Phase 2)
	mux.HandleFunc("GET /api/hull-types", handlers.ListHullTypes)
	mux.HandleFunc("GET /api/module-types", handlers.ListModuleTypes)
	mux.HandleFunc("GET /api/blueprints", handlers.ListBlueprints)

	// Protected routes (require JWT)
	protected := http.NewServeMux()

	// Phase 1 - Planet Management
	protected.HandleFunc("GET /api/player/me", handlers.PlayerMe)
	protected.HandleFunc("GET /api/planets", handlers.ListPlanets)
	protected.HandleFunc("GET /api/planets/{id}", handlers.GetPlanet)
	protected.HandleFunc("GET /api/planets/{id}/buildings", handlers.ListBuildings)
	protected.HandleFunc("POST /api/planets/{id}/buildings", handlers.ConstructBuilding)
	protected.HandleFunc("POST /api/planets/{id}/buildings/{buildingId}/upgrade", handlers.UpgradeBuilding)
	protected.HandleFunc("POST /api/planets/{id}/buildings/{buildingId}/cancel", handlers.CancelUpgrade)
	protected.HandleFunc("PUT /api/planets/{id}/buildings/{buildingId}/move", handlers.MoveBuilding)
	protected.HandleFunc("GET /api/planets/{id}/resources", handlers.GetResources)
	protected.HandleFunc("POST /api/planets/{id}/resources/collect", handlers.CollectResources)
	protected.HandleFunc("POST /api/resources/collect-warehouse", handlers.CollectWarehouse)

	// Phase 2 - Ship Factory
	protected.HandleFunc("GET /api/ship-factory", handlers.GetShipFactory)
	protected.HandleFunc("GET /api/ship-factory/slots", handlers.GetShipFactorySlots)
	protected.HandleFunc("POST /api/ship-factory/build", handlers.BuildShips)
	protected.HandleFunc("POST /api/ship-factory/cancel/{slot}", handlers.CancelShipBuild)

	// Phase 2 - Ship Designs
	protected.HandleFunc("GET /api/ship-designs", handlers.ListShipDesigns)
	protected.HandleFunc("POST /api/ship-designs", handlers.CreateShipDesign)
	protected.HandleFunc("PUT /api/ship-designs/{id}", handlers.UpdateShipDesign)
	protected.HandleFunc("DELETE /api/ship-designs/{id}", handlers.DeleteShipDesign)
	protected.HandleFunc("GET /api/ship-designs/{id}/stats", handlers.GetDesignStats)

	// Phase 2 - Blueprints (player-specific)
	protected.HandleFunc("GET /api/blueprints/mine", handlers.ListMyBlueprints)
	protected.HandleFunc("POST /api/blueprints/{id}/activate", handlers.ActivateBlueprint)
	protected.HandleFunc("POST /api/blueprints/{id}/research", handlers.ResearchBlueprint)
	protected.HandleFunc("GET /api/blueprint-research/active", handlers.GetActiveBlueprintResearch)

	// Phase 2 - Fleets
	protected.HandleFunc("GET /api/fleets", handlers.ListFleets)
	protected.HandleFunc("POST /api/fleets", handlers.CreateFleet)
	protected.HandleFunc("PUT /api/fleets/{id}", handlers.UpdateFleet)
	protected.HandleFunc("DELETE /api/fleets/{id}", handlers.DeleteFleet)
	protected.HandleFunc("POST /api/fleets/{id}/assign-stack", handlers.AssignStack)
	protected.HandleFunc("POST /api/fleets/{id}/remove-stack", handlers.RemoveStack)
	protected.HandleFunc("POST /api/fleets/{id}/move", handlers.MoveFleet)
	protected.HandleFunc("POST /api/fleets/{id}/recall", handlers.RecallFleet)
	protected.HandleFunc("POST /api/fleets/{id}/dismiss", handlers.DismissFleet)

	// Phase 2 - Instances
	protected.HandleFunc("GET /api/instances", handlers.ListInstances)
	protected.HandleFunc("GET /api/instances/progress", handlers.GetInstanceProgress)
	protected.HandleFunc("GET /api/instances/{id}", handlers.GetInstance)
	protected.HandleFunc("POST /api/instances/{id}/attempt", handlers.AttemptInstance)

	// Phase 2 - Spacedock
	protected.HandleFunc("GET /api/spacedock", handlers.GetSpacedock)
	protected.HandleFunc("GET /api/spacedock/repairs", handlers.ListRepairs)
	protected.HandleFunc("POST /api/spacedock/repair", handlers.StartRepair)
	protected.HandleFunc("POST /api/spacedock/accelerate", handlers.AccelerateRepair)

	// Quests
	protected.HandleFunc("GET /api/quests", handlers.ListQuests)
	protected.HandleFunc("GET /api/quests/daily", handlers.GetDailyQuests)
	protected.HandleFunc("POST /api/quests/{id}/claim", handlers.ClaimQuest)
	protected.HandleFunc("POST /api/quests/daily/claim-tier", handlers.ClaimDailyTier)
	protected.HandleFunc("POST /api/quests/sync", handlers.SyncQuests)

	// Research
	protected.HandleFunc("GET /api/research", handlers.ListResearch)
	protected.HandleFunc("GET /api/research/trees/{tree}", handlers.GetResearchTree)
	protected.HandleFunc("POST /api/research/start", handlers.StartResearch)
	protected.HandleFunc("POST /api/research/cancel", handlers.CancelResearch)
	protected.HandleFunc("POST /api/research/speedup", handlers.SpeedupResearch)
	protected.HandleFunc("GET /api/research/active", handlers.GetActiveResearch)

	// Inventory
	protected.HandleFunc("GET /api/inventory", handlers.GetInventory)
	protected.HandleFunc("POST /api/inventory/{id}/use", handlers.UseItem)

	// Commanders
	protected.HandleFunc("POST /api/commanders/recruit", handlers.RecruitCommander)
	protected.HandleFunc("GET /api/commanders", handlers.ListCommanders)
	protected.HandleFunc("POST /api/commanders/merge", handlers.MergeCommanders)
	protected.HandleFunc("POST /api/fleets/{id}/assign-commander", handlers.AssignCommander)
	protected.HandleFunc("POST /api/fleets/{id}/unassign-commander", handlers.UnassignCommander)
	protected.HandleFunc("DELETE /api/commanders/{id}", handlers.DismissCommander)

	// Combat Reports
	protected.HandleFunc("GET /api/combat-reports", handlers.ListCombatReports)
	protected.HandleFunc("GET /api/combat-reports/{id}", handlers.GetCombatReport)

	// Recycling Plant
	protected.HandleFunc("POST /api/recycling-plant/recycle", handlers.StartRecycle)
	protected.HandleFunc("GET /api/recycling-plant/jobs", handlers.ListRecyclingJobs)
	protected.HandleFunc("POST /api/recycling-plant/collect/{id}", handlers.CollectRecycle)
	protected.HandleFunc("DELETE /api/recycling-plant/jobs/{id}", handlers.CancelRecycle)
	protected.HandleFunc("GET /api/ship-instances/available", handlers.ListAvailableShips)

	// PvP Combat
	protected.HandleFunc("POST /api/pvp/attack", handlers.AttackPlanet)
	protected.HandleFunc("GET /api/pvp/search", handlers.SearchPlanets)

	// World Chat
	protected.HandleFunc("POST /api/chat/send", handlers.SendChatMessage)
	protected.HandleFunc("GET /api/chat/messages", handlers.GetChatMessages)

	// Corps & Galaxy
	protected.HandleFunc("GET /api/corp", handlers.GetCorp)
	protected.HandleFunc("POST /api/corp", handlers.CreateCorp)
	protected.HandleFunc("POST /api/corp/join", handlers.JoinCorp)
	protected.HandleFunc("POST /api/corp/leave", handlers.LeaveCorp)
	protected.HandleFunc("GET /api/corp/members", handlers.ListCorpMembers)
	protected.HandleFunc("POST /api/corp/donate", handlers.DonateResources)
	protected.HandleFunc("PUT /api/corp/members/{id}/role", handlers.UpdateMemberRole)
	protected.HandleFunc("GET /api/corp/search", handlers.SearchCorps)
	protected.HandleFunc("POST /api/corp/rbp/{id}/attack", handlers.AttackRBP)
	protected.HandleFunc("GET /api/galaxy/map", handlers.GetGalaxyMap)

	// Dev Tools (require X-Dev-Mode: true header)
	protected.HandleFunc("POST /api/dev/reset", handlers.DevResetPlayer)
	protected.HandleFunc("POST /api/dev/give-resources", handlers.DevGiveResources)
	protected.HandleFunc("POST /api/dev/give-item", handlers.DevGiveItem)
	protected.HandleFunc("POST /api/dev/complete-constructions", handlers.DevCompleteConstructions)
	protected.HandleFunc("POST /api/dev/complete-research", handlers.DevCompleteResearch)
	protected.HandleFunc("POST /api/dev/complete-ship-builds", handlers.DevCompleteShipBuilds)
	protected.HandleFunc("POST /api/dev/give-blueprint", handlers.DevGiveBlueprint)
	protected.HandleFunc("POST /api/dev/give-commander", handlers.DevGiveCommander)
	protected.HandleFunc("POST /api/dev/give-ships", handlers.DevGiveShips)
	protected.HandleFunc("POST /api/dev/set-building-level", handlers.DevSetBuildingLevel)
	protected.HandleFunc("POST /api/dev/complete-recycling", handlers.DevCompleteRecycling)
	protected.HandleFunc("POST /api/dev/create-corp", handlers.DevCreateCorp)
	protected.HandleFunc("POST /api/dev/join-corp", handlers.DevJoinCorp)
	protected.HandleFunc("POST /api/dev/give-corp-wealth", handlers.DevGiveCorpWealth)

	mux.Handle("/api/", middleware.Auth(protected))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := middleware.CORS(mux)

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
