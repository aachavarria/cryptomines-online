# Galaxy Online 2 - Grid, Harvest & UI Validation Research

> Compiled: 2026-02-05
> Sources: Galaxy Online II Wiki (Fandom), GameYum guides, NewRPG, F2PG, MMORPG.com, WordPress blogs, previous research (go2-visual-ui-research.md)
> Game by IGG. Flash-based browser MMORTS. Shut down January 12, 2021.

---

## 1. Sistema de colocacion de edificios en el mapa

### Respuesta: NO es grid-based. Son posiciones fijas/slots predeterminados.

**La investigacion anterior dijo "NOT grid-based, fixed positions" -- esto es CORRECTO.**

Galaxy Online 2 NO tiene un sistema de grid donde el jugador elige donde poner los edificios. Los edificios se construyen desde un **menu de construccion**, no arrastrando y soltando en un mapa.

**Evidencia:**
- El proceso de construccion es: `Click "Build" (lower right corner) > Select "Construct Buildings" > Choose building category > Select building to construct`
  - Fuente: Development Quests wiki - "Click on (Build) located in the lower right corner, and select (Construct Buildings)."
- NO hay mencion en NINGUNA fuente (wiki, guides, reviews, FAQs) de:
  - Elegir posicion en un grid
  - Arrastrar edificios a un tile
  - Mover edificios despues de colocarlos
  - Restricciones de adyacencia o colocacion espacial
  - Coordenadas de edificios en el mapa
- Los edificios tienen **limites de cantidad fija** (max 8 Metal Collectors, max 8 HE3 Extractors, max 12 Residential Areas, max 1 Resource Warehouse), lo que refuerza el modelo de "slots" no de "posiciones libres"
  - Fuente: Metal Collector wiki - "a maximum of 8 collectors per planet"
  - Fuente: Resource Warehouse wiki - "You can have a max of 1 resource warehouse"
- El juego muestra los edificios como **sprite icons clickeables en una vista estatica del planeta**
  - Fuente: Observacion consolidada de multiples fuentes + previous research

**Sobre mover edificios:**
- [NO ENCONTRADO - PREGUNTAR AL LEAD] No hay NINGUNA mencion en las fuentes de poder mover edificios de posicion. La ausencia total de esta mencion, combinada con el sistema de slots fijos, sugiere fuertemente que NO se pueden mover.

**Sobre restricciones de colocacion:**
- La unica restriccion encontrada es el **nivel de Civic Center** requerido para desbloquear y upgradear edificios
  - Fuente: Civic Center wiki - "constructing a level 3 Metal Collector requires possessing a Civic Center at minimum level 2"
- NO hay restricciones de tiles, adyacencia, o posicion espacial

---

## 2. Perspectiva del mapa

### Respuesta: Vista 2D estatica, top-down con sprites pre-renderizados. NO isometrica, NO 3D.

**Evidencia:**
- El juego es Flash-based, con **2D rendered buildings con sprites pre-renderizados** (no 3D, no isometrico)
  - Fuente: Previous research confirmado por multiples fuentes - "2D rendered buildings with pre-rendered sprite art (not 3D, not isometric)"
- Es un **window/panel-based UI**, NO un mundo 3D con camara libre
  - Fuente: Previous research - "Window/panel-based UI - NOT a free-roaming 3D world"

**Sobre rotacion del mapa:**
- NO se puede rotar. Es una vista fija 2D.
  - Fuente: Ninguna fuente menciona rotacion de camara. El juego es Flash 2D.

**Sobre zoom:**
- [NO ENCONTRADO - PREGUNTAR AL LEAD] No hay mencion de zoom en la vista del planeta. Siendo Flash 2D con resolucion fija dentro del frame del browser, es muy probable que NO tenga zoom.

**Sobre cambiar pitch/angulo:**
- NO aplica. Es 2D estatico.

**Sobre como se mueve la camara/vista:**
- [NO ENCONTRADO - PREGUNTAR AL LEAD] No hay mencion de scroll o movimiento de camara en el planet base view. La vista parece ser completamente estatica (todos los edificios visibles en una sola pantalla).

---

## 3. Menu contextual al clickear edificio

### Respuesta: Clickear un edificio abre su PANEL DE GESTION directamente. NO hay menu contextual intermedio.

**Evidencia:**
- "Click a building to open its management panel"
  - Fuente: Consolidado de multiples guias (GameYum Structures Guide, Development Quests wiki)
- El panel muestra: building level, upgrade button, upgrade requirements, stats
  - Fuente: GameYum Structures Guide - "click the upgrade button to see the various requirements for each level"
- Hay un **boton de upgrade DENTRO del panel** del edificio
  - Fuente: Accelerate wiki - "Building construction can be sped up by clicking the 'up arrow' icon beside the progress timer bar on the right side of the screen"
- Construccion en progreso muestra una **progress timer bar** en el lado derecho de la pantalla con icono de flecha para acelerar

**Opciones encontradas en buildings especificos:**
- **Ship Factory**: Abre panel con lista de blueprints (izquierda) y boton de disenar (derecha). Boton de flecha abajo para construir naves.
  - Fuente: GameYum Ship Guide - "hold up to twenty ship designs", "add your weapons and equipment to each ship design"
- **Recycling Plant**: Menu de seleccion donde se eligen naves y boton "proceed to destroy"
  - Fuente: GameYum Structures Guide
- **Compound Center**: Interfaz drag-and-drop para "merge and upgrade your commander cards"
  - Fuente: GameYum Structures Guide
- **Command Center**: Muestra fleet statistics en la parte inferior con opciones de formacion y rango
  - Fuente: GameYum Structures Guide
- **Resource Warehouse**: Muestra recursos acumulados + boton de Harvest
  - Fuente: Resource Warehouse wiki + GameYum Structures Guide

**NO se encontraron las siguientes opciones en el menu de edificios:**
- Move (mover edificio)
- Demolish (demoler edificio)
- [NO ENCONTRADO - PREGUNTAR AL LEAD] No hay evidencia de opciones de View/Move/Demolish en el menu de edificios.

---

## 4. Mecanica de Harvest/Collect

### Respuesta: El harvest se hace DESDE el Resource Warehouse. Los edificios de produccion generan recursos que se acumulan en el warehouse, y el jugador debe clickear "Harvest" en el warehouse para cobrarlos.

**Evidencia:**
- "The Resource Warehouse is what will collect your resources until you harvest from it."
  - Fuente: Resource Warehouse wiki (cita directa)
- "you must click the 'harvest' button to move your collected resources into your inventory for use"
  - Fuente: GameYum Structures Guide (cita directa)
- Metal Collector wiki confirma: "increases metal production, which is harvestable from the Resource Warehouse"
  - Fuente: Metal Collector wiki (cita directa)
- HE3 Extractor wiki confirma: "collects He3 for you"
  - Fuente: HE3 Extractor wiki
- "All resources are stored in the Resource Warehouse, which has limited capacity requiring the player to collect from it every 12 hours or less"
  - Fuente: Web search result summary

### Flujo de recursos:
```
Metal Collector ──┐
HE3 Extractor  ───┤──> Resource Warehouse (acumula) ──> Harvest (click) ──> Player Inventory
Residential Area ─┘
```

**Sobre harvest individual vs global:**
- Los edificios de produccion (Metal Collector, HE3 Extractor, Residential Area) NO tienen su propio boton de harvest individual
- TODO se acumula en el Resource Warehouse
- El harvest es UN solo boton en el Resource Warehouse que cobra TODOS los recursos acumulados de una vez
  - Fuente: Development Quests wiki, primera quest dice "Harvest your Resource Warehouse" (no "harvest your metal collector")

**Sobre "collect all" global:**
- [NO ENCONTRADO - PREGUNTAR AL LEAD] No se encontro evidencia de un boton global de "collect all" fuera del Resource Warehouse. El unico punto de harvest encontrado es el boton dentro del Resource Warehouse.

**Harvest de amigos:**
- Se pueden cosechar recursos de los space bases de amigos: "click on their harvest button and you'll get a certain amount of resources depending upon the mines they selected"
  - Fuente: GameYum Economy Guide (cita directa)

---

## 5. Storage / Resource Warehouse

### Respuesta: Max 1 por planeta. Limites INDIVIDUALES por recurso. Upgradeable hasta nivel 24. Cap aumenta con el nivel.

**Evidencia:**

**Cuantos se pueden construir:**
- "You can have a max of 1 resource warehouse that is fully upgradable to level 24, this building is already available to you when you first enter the game."
  - Fuente: Resource Warehouse wiki (cita directa)
- Confirmado: **Solo 1 por planeta**, ya esta construido al empezar

**Limites de storage por recurso:**
- Los limites son **INDIVIDUALES por cada recurso** (Metal, He3, Gold por separado), NO combinados
- Tabla de capacidad por nivel (del Resource Warehouse wiki):

| Level | Storage per Resource |
|-------|---------------------|
| 1     | 10,000              |
| 10    | 350,000             |
| 20    | 6,000,000           |
| 24    | 20,000,000          |

- Fuente: Resource Warehouse wiki - "Level 1: 10,000 each resource (Metal, He3, Gold)"

**El warehouse level aumenta el cap:**
- SI. Cada nivel aumenta significativamente la capacidad
  - Fuente: Resource Warehouse wiki - "Upgrade it to collect more resources over time"
- Ademas, la investigacion **Expand Capacity** en el tech tree de Logistics Construction Science permite exceder los limites base
  - Fuente: Resource Warehouse wiki - "Expand Capacity research under the Logistics Construction Science tech tree"

**Que pasa cuando un recurso llega al cap:**
- [NO ENCONTRADO - PREGUNTAR AL LEAD] Ninguna fuente consultada especifica que sucede cuando el storage esta lleno. Las opciones probables son: (a) la produccion se pierde/para, o (b) se desborda y se pierde el exceso. Pero no hay confirmacion en las fuentes.

**Build times y costs:**
- "Build times and construction costs may vary depending on your Construction Boost and Quality Materials research"
  - Fuente: Resource Warehouse wiki

---

## 6. Panel de info de edificio

### Respuesta: Es un panel/modal dentro de la misma pantalla. Muestra info del edificio + boton de upgrade.

**Evidencia:**
- Al clickear un edificio, se abre un **panel modal** (estilo ventana oscura) superpuesto sobre la vista del planeta
  - Fuente: Previous research - "Modal-style pop-ups for building details", "Dark-themed panels with organized data rows"
- El panel muestra:
  - Nombre del edificio
  - Nivel actual
  - Stats de produccion/funcion
  - Costos de upgrade (Metal, He3, Gold)
  - Tiempo de upgrade
  - Boton de Upgrade
  - Fuente: GameYum Structures Guide - "click the upgrade button to see the various requirements for each level"
- Cuando un upgrade esta en progreso:
  - Se muestra una **progress timer bar** en el lado derecho de la pantalla
  - Icono de **flecha azul/amarilla** para acelerar la construccion
  - Fuente: Accelerate wiki - "Building construction can be sped up by clicking the 'up arrow' icon beside the progress timer bar on the right side of the screen"
- **5 construction slots** disponibles simultaneamente
  - Fuente: Previous research consolidado + Guide To Advancing Quickly - "five construction slots"
- **Flashing yellow wrench icon** indica un edificio que necesita atencion
  - Fuente: Previous research consolidado

**Es un panel lateral? Un popup? Una nueva pantalla?**
- Es un **popup/modal** que se abre sobre la vista actual. NO cambia de pantalla, NO es una nueva pagina.
- Estilo: Dark-themed panel con filas de datos organizadas
  - Fuente: Previous research - "Window-based layout (not full-screen transitions)"

**Hay boton de upgrade DENTRO del panel?**
- SI. El boton de upgrade esta dentro del panel del edificio.
  - Fuente: GameYum Structures Guide, Accelerate wiki, Development Quests wiki (multiples fuentes confirman)

---

## Resumen de Confianza

| Punto | Confianza | Notas |
|-------|-----------|-------|
| 1. Grid system | ALTA | Multiples fuentes confirman menu-based construction, NO grid |
| 1. Mover edificios | MEDIA-BAJA | No encontrado, pero ausencia de mencion sugiere que NO |
| 2. Perspectiva | ALTA | 2D estatica, Flash, confirmado por multiples fuentes |
| 2. Zoom/scroll | BAJA | No encontrado en fuentes |
| 3. Menu al click | ALTA | Abre panel directo, confirmado por guias y wiki |
| 3. Options (move/demolish) | BAJA | No encontrado |
| 4. Harvest desde warehouse | MUY ALTA | Multiples citas directas lo confirman |
| 4. Collect all | BAJA | No encontrado |
| 5. Max 1 warehouse | MUY ALTA | Cita directa del wiki |
| 5. Limites individuales | ALTA | Wiki muestra "each resource" |
| 5. Que pasa al cap | BAJA | No encontrado |
| 6. Panel tipo modal | ALTA | Confirmado por multiples fuentes |
| 6. Upgrade dentro del panel | MUY ALTA | Multiples fuentes confirman |

---

## Items [NO ENCONTRADO - PREGUNTAR AL LEAD]

1. **Mover edificios**: No hay mencion de poder mover edificios. Probablemente NO se pueden mover.
2. **Zoom en planet view**: No hay mencion de zoom. Probablemente la vista es fija.
3. **Scroll/camara en planet view**: No hay mencion de scroll. Probablemente la vista es completamente estatica.
4. **Opciones de Move/Demolish en menu de edificio**: No encontrado. Solo se confirman: upgrade, harvest (warehouse), y funciones especificas por edificio.
5. **Collect all global**: No encontrado. Solo el harvest desde el Resource Warehouse.
6. **Que pasa cuando storage llega al cap**: No especificado en las fuentes.

---

## Fuentes Consultadas

### Wiki Pages (Fandom)
- Galaxy Online II Wiki: https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki
- Walkthrough: https://galaxyonlineii.fandom.com/wiki/Walkthrough
- Beginner FAQ: https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ
- Beginner's Guide: https://galaxyonlineii.fandom.com/wiki/Beginner%27s_Guide
- Resource Warehouse: https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse
- Metal Collector: https://galaxyonlineii.fandom.com/wiki/Metal_Collector
- HE3 Extractor: https://galaxyonlineii.fandom.com/wiki/HE3_Extractor
- Residential Area: https://galaxyonlineii.fandom.com/wiki/Residential_Area
- Civic Center: https://galaxyonlineii.fandom.com/wiki/Civic_Center
- Development Quests: https://galaxyonlineii.fandom.com/wiki/Development_Quests
- Accelerate (Buildings and Research): https://galaxyonlineii.fandom.com/wiki/Accelerate_(Buildings_and_Research)
- Planet Base: https://galaxyonlineii.fandom.com/wiki/Planet_Base
- Category: Planet Base Structure: https://galaxyonlineii.fandom.com/wiki/Category:Planet_Base_Structure
- Module Placement: https://galaxyonlineii.fandom.com/wiki/Module_Placement (ship modules, not building placement)
- Guide To Advancing Quickly: https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly
- Galaxy Online II for Beginners: https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_for_beginners!

### External Guides
- GameYum Structures Guide: https://www.gameyum.com/galaxy-online/115041-galaxy-online-ii-guides-building-structures/
- GameYum Overview: https://www.gameyum.com/galaxy-online/109124-facebookgame-guides-galaxy-online-ii/
- GameYum Economy Guide: https://www.gameyum.com/galaxy-online/109252-galaxy-online-ii-guides-economy-guide/
- Galaxy Online 2 Facebook Wiki - Buildings: https://galaxyonline2facebook.fandom.com/wiki/Buildings

### Reviews & General Info
- Gamezebo Review (2011): https://www.gamezebo.com/2011/03/14/galaxy-online-2-review/
- MMORPG.com: https://www.mmorpg.com/galaxy-online-2
- F2PG: https://www.f2pg.com/galaxy-online-2/
- NewRPG: https://newrpg.com/browser-games/galaxy-online-2/
- 4GameGround: https://4gameground.com/galaxy-online-2/
- Blog con screenshots del client: https://m3rilix.wordpress.com/2011/03/04/galaxy-online-iiclient/
- GO2 Shutdown Article: https://deajae.co.uk/posts/the-end-of-galaxy-online-2/

### Private Server
- SuperGO2 (private server): https://github.com/SuperGO2/supergo2-issues

### Previous Research
- /Users/yurei/cryptomines-online/docs/research/go2-visual-ui-research.md (compiled earlier)
