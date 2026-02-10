# Blueprint Research System - Implementation Status

**Date:** 2026-02-07
**Sistema:** Upgrade de hulls/módulos de tier 1 → tier 2 → tier 3

---

## ✅ COMPONENTES IMPLEMENTADOS

### **1. Base de Datos (100%)**
- ✅ Tabla `blueprints` con `research_level` (0-3)
- ✅ Tabla `player_blueprints` con `research_level` (1-3)
- ✅ Tabla `blueprint_research` para tracking de investigaciones activas
- ✅ Seed data: 75 hulls (10 líneas frigate, 10 cruiser, 5 battleship × 3 tiers)
- ✅ Seed data: 97 módulos con tiers
- ✅ Naming pattern:
  - Hulls: `{base}_i`, `{base}_ii`, `{base}_iii`
  - Módulos: mismo `name`, tier 1/2/3

### **2. Backend - Endpoints (70%)**
- ✅ `GET /api/blueprints` - Listar todos los blueprints del juego
- ✅ `GET /api/blueprints/mine` - Listar blueprints del player
- ✅ `POST /api/blueprints/{id}/activate` - Activar blueprint
- ✅ `POST /api/blueprints/{id}/research` - Iniciar research
  - ✅ Validaciones: blueprint activado, WRC level, solo 1 research activo
  - ✅ Costos: baseCost = 10000 × targetLevel
  - ✅ Tiempo: 3600 segundos × targetLevel (1 hora por nivel)
  - ✅ Crea row en `blueprint_research` con timer
- ✅ `blueprint_service.go` - Helper functions (RECIÉN CREADO):
  - ✅ `GetAvailableHullTiers()` - Parsea `{base}_i/ii/iii`
  - ✅ `GetAvailableModuleTiers()` - Filtra por name + tier
  - ✅ `CanPlayerUseHull()` - Valida tier desbloqueado
  - ✅ `CanPlayerUseModule()` - Valida tier desbloqueado

---

## ❌ COMPONENTES FALTANTES

### **3. Backend - Auto-Complete Worker (0%)**

**Falta:** Worker/cron que revisa `blueprint_research.research_finish_at` y completa automáticamente.

**Tareas:**
- [ ] Crear `blueprint_worker.go`
- [ ] Función `CompleteFinishedResearch()`:
  ```go
  // Pseudo-código:
  // 1. SELECT * FROM blueprint_research WHERE is_researching = true AND research_finish_at <= now()
  // 2. Para cada row:
  //    - UPDATE player_blueprints SET research_level = target_level WHERE id = player_blueprint_id
  //    - UPDATE blueprint_research SET is_researching = false
  //    - Quest progress: services.UpdateQuestProgress(player_id, "research_blueprint", blueprint_name, 1)
  ```
- [ ] Llamar desde worker principal cada 1 minuto
- [ ] Test: Iniciar research, avanzar tiempo manualmente, verificar auto-complete

**Tiempo estimado:** 2-3 horas

---

### **4. Backend - Validación en Ship Design (0%)**

**Falta:** Validar que el player tenga el tier desbloqueado al crear/editar ship design.

**Tareas:**
- [ ] En `CreateShipDesign()`:
  ```go
  // Validar hull tier
  canUse, err := services.CanPlayerUseHull(playerID, req.HullTypeID)
  if err != nil || !canUse {
      http.Error(w, `{"error":"hull tier not unlocked"}`, http.StatusForbidden)
      return
  }

  // Validar cada módulo tier
  for _, mod := range req.Modules {
      canUse, err := services.CanPlayerUseModule(playerID, mod.ModuleTypeID)
      if err != nil || !canUse {
          http.Error(w, `{"error":"module tier not unlocked"}`, http.StatusForbidden)
          return
      }
  }
  ```
- [ ] Mismo en `UpdateShipDesign()`
- [ ] Test: Intentar crear design con tier 2 sin haberlo investigado → debe fallar

**Tiempo estimado:** 1-2 horas

---

### **5. Backend - Endpoints Adicionales (0%)**

**Opcional pero útil:**
- [ ] `GET /api/blueprints/{id}/available-tiers` - Listar qué tiers están desbloqueados
  - Útil para el frontend para mostrar hulls/módulos disponibles
  - Usa `GetAvailableHullTiers()` / `GetAvailableModuleTiers()`
- [ ] `GET /api/blueprints/research/active` - Listar investigaciones activas del player
- [ ] `POST /api/blueprints/research/cancel` - Cancelar research activo (refund 50%)
- [ ] `POST /api/blueprints/research/speedup` - Gastar vouchers/gold para acelerar

**Tiempo estimado:** 3-4 horas

---

### **6. Frontend - UI Completa (0%)**

**Falta:** Todo el UI del sistema de blueprint research.

**Tareas:**

#### **6.1. BlueprintPanel Updates (2-3 horas)**
- [ ] Indicador visual de `research_level`:
  ```tsx
  // Mostrar estrellas: ★★★ (filled) vs ☆☆☆ (empty)
  // Ejemplo: research_level = 2 → ★★☆
  ```
- [ ] Botón "Research" (solo visible si `research_level < 3`)
- [ ] Al click → mostrar modal de confirmación con costos + tiempo
- [ ] Llamar `POST /api/blueprints/{id}/research`
- [ ] Mostrar timer countdown si hay research activo
- [ ] Tooltip: "Tier 1 unlocked, Tier 2 locked (requires research level 2)"

#### **6.2. Ship Design Panel Updates (1-2 horas)**
- [ ] Filtrar hulls disponibles por research_level:
  ```tsx
  // Llamar GET /api/blueprints/{id}/available-tiers
  // Mostrar solo hulls desbloqueados
  // Hulls bloqueados en gris con candado + tooltip "Research blueprint to unlock Tier 2"
  ```
- [ ] Mismo para módulos
- [ ] Error handling: Si backend rechaza design por tier bloqueado, mostrar mensaje claro

#### **6.3. Research Progress Display (1 hora)**
- [ ] Panel "Active Research" en la UI (sidebar o modal)
- [ ] Mostrar:
  - Blueprint name
  - Target tier
  - Progress bar (tiempo restante)
  - Button "Speed Up" (si implementamos speedup endpoint)
- [ ] Auto-refresh cada 10 segundos para actualizar timer

**Tiempo estimado total:** 4-6 horas

---

### **7. Quest Integration (0%)**

**Falta:** Auto-progress de quests relacionados con blueprint research.

**Tareas:**
- [ ] En `CompleteFinishedResearch()` (worker):
  ```go
  services.UpdateQuestProgress(playerID, "research_blueprint", blueprintName, 1)
  ```
- [ ] Verificar que existan quests de tipo `research_blueprint` en seed data
- [ ] Si no existen, agregar a daily quests o side quests

**Tiempo estimado:** 30 minutos

---

## 📊 RESUMEN POR MÓDULO

| Módulo | Completado | Faltante | Tiempo Estimado |
|--------|------------|----------|-----------------|
| **1. Base de Datos** | 100% | - | - |
| **2. Backend - Endpoints** | 70% | Validación en ship design | 1-2 horas |
| **3. Backend - Worker** | 0% | Auto-complete research | 2-3 horas |
| **4. Backend - Endpoints Extra** | 0% | available-tiers, cancel, speedup | 3-4 horas (opcional) |
| **5. Service Layer** | 100% | - | - |
| **6. Frontend - UI** | 0% | BlueprintPanel + ShipDesign updates + Research display | 4-6 horas |
| **7. Quest Integration** | 0% | Auto-progress on research complete | 30 min |

**TOTAL TIEMPO ESTIMADO:** 11-16 horas (sin endpoints opcionales) o 14-20 horas (con opcionales)

---

## 🎯 PRIORIDAD DE IMPLEMENTACIÓN

### **Fase 1: MVP Funcional (6-8 horas)**
1. Backend Worker (auto-complete) - 2-3h
2. Validación ship design - 1-2h
3. Frontend BlueprintPanel básico - 2-3h

### **Fase 2: UX Completo (5-8 horas)**
4. Frontend ShipDesign filters - 1-2h
5. Research Progress Display - 1h
6. Endpoints adicionales (available-tiers, cancel, speedup) - 3-4h
7. Quest integration - 30min

---

## ✅ LISTO PARA IMPLEMENTAR

- [x] Service layer con parseo de nombres (RECIÉN CREADO)
- [x] Endpoints de research start (YA EXISTE)
- [ ] Worker para auto-complete (SIGUIENTE PASO)
- [ ] Validación en ship design (SIGUIENTE PASO)
- [ ] Frontend UI (SIGUIENTE PASO)

---

**Last Updated:** 2026-02-07
**Status:** Service layer completado, listo para continuar con Worker + Validation
