# Galaxy Online 2 - Visual Style Research Report

**Date:** 2026-02-06
**Researcher:** Explore Agent
**Purpose:** Determine the visual presentation style of Galaxy Online 2's planet/base building view and compare to similar browser strategy games

---

## Executive Summary

Galaxy Online 2 (GO2) was a **Facebook-based Flash game** (2010-2015 era) that is **no longer in service**. Due to the game's closure and limited archival documentation, **direct visual evidence of GO2's planet building interface is scarce**. However, based on research into similar browser strategy games from the same era and technical context, we can infer the likely visual approach.

### Key Findings:

1. **GO2 likely used 2D pre-rendered sprites or Flash graphics**, NOT true 3D rendering
2. **Most browser strategy games from this era (2010-2015) used 2D sprites** with either list-based or simple graphical interfaces
3. **True isometric 3D rendering with Three.js is NOT typical** for this game genre or era
4. **"Isometric" in browser games typically means 2D sprites drawn at isometric angles**, not 3D models

---

## 1. Galaxy Online 2 Visual Interface

### What We Know:
- **Platform:** Facebook browser game (Flash-based)
- **Era:** 2010-2015 (before HTML5/WebGL became standard)
- **Graphics Quality:** Described as "decent graphics for a flash game" with a "relatively static battle system"
- **Focus:** Ship design and combat were the primary visual elements, not base building
- **Status:** Game is no longer in service, making direct verification impossible

### Planet Base Evidence:
From the [Galaxy Online II Wiki - Planet Base](https://galaxyonlineii.fandom.com/wiki/Planet_Base):
- Wiki references two images: "Home.PNG" and "Home_base.PNG"
- Minimal description: "It contains all the land based buildings that you will have"
- **No mention of grid systems, isometric views, or 3D models**
- Article marked as a "stub" with limited information

### Inference:
Given the Flash-based Facebook platform and era, GO2 most likely used:
- **2D pre-rendered sprites** for buildings
- **Simple graphical representation** of the base (possibly overhead/angled view)
- **List-based or menu-driven interface** for construction (common in this genre)
- **NOT a true isometric grid** with diamond tiles
- **NOT real-time 3D rendering**

---

## 2. Comparison: Similar Browser Strategy Games

### OGame (2002-present)
**Source:** [OGame Wiki - Interface](https://ogame.fandom.com/wiki/Interface) | [Building Screen](https://ogame.fandom.com/wiki/Building_Screen)

**Visual Style:**
- **Text-based with minimal graphics**
- Buildings shown as **list of icons with descriptions**
- **No grid placement system** - buildings are abstracted
- **No visual planet representation** - purely menu-driven
- Planet image shown as decoration, not interactive

**Key Quote:** "The Building Screen contains a list of available buildings with pictures and descriptions. It will say after the buildings name what level the building has currently been built to."

**Conclusion:** OGame is **purely list-based**, not graphical or spatial.

---

### Ikariam (2008-present)
**Source:** [Ikariam Wiki - Town View](https://ikariam.fandom.com/wiki/Town_view)

**Visual Style:**
- **2D pre-rendered graphics** with an **isometric or angled perspective**
- Buildings shown as **2D sprites** arranged in a **free-form layout** (NOT a strict grid)
- **11 fixed building spots** marked with flags
- Background evolves as Town Hall levels up (roads, parks, fountains appear)
- **Cosmetic visual progression** - purely aesthetic

**Key Quote:** "The background of your town is changed and new features are introduced as your Town Hall levels increase... roads, parks, fountains and gazebos progressively appear."

**Conclusion:** Ikariam uses **2D pre-rendered isometric sprites** with **fixed building positions**, NOT a free-placement grid system.

---

### Travian (2004-present)
**Source:** [Games Like Travian](https://gameslikefinder.com/games-like-travian/)

**Visual Style:**
- **2D browser-based** military strategy
- **Resource management** and **empire building** focus
- Described as having **minimal graphics** similar to OGame
- **Not visually elaborate** - function over form

**Conclusion:** Travian is **list/menu-driven** with minimal graphics.

---

### Classic Browser Strategy Games (General Findings)
**Sources:**
- [Browser Based Strategy Games](https://mmos.com/review/browser-games/strategy)
- [Text-based vs Graphical Comparison](https://www.mmobomb.com/browsergames/strategy)

**Spectrum of Visual Approaches:**

1. **Text-Based (Minimal Graphics):**
   - OGame, Astro Empires
   - List-based interfaces
   - No spatial building placement

2. **2D Pre-Rendered Graphics:**
   - Ikariam, Forge of Empires
   - Fixed or free-form layouts
   - 2D sprites with isometric perspective
   - **Most common approach for browser strategy games**

3. **Flash Graphics (2010-2015 Facebook Era):**
   - Social games like Galaxy Online 2
   - Moderate graphics quality
   - 2D sprites or Flash animations
   - Focus on ship design, not base building

4. **Modern HTML5/WebGL (2015+):**
   - Real-time 3D rendering
   - Three.js, Babylon.js engines
   - **NOT typical for 2010-2015 Facebook games**

**Conclusion:** Browser strategy games from GO2's era (2010-2015) **overwhelmingly used 2D sprites or text-based interfaces**, NOT real-time 3D rendering.

---

## 3. Isometric Graphics: 2D Sprites vs 3D Models

**Source:** [Isometric Video Game Graphics - Wikipedia](https://en.wikipedia.org/wiki/Isometric_video_game_graphics)

### Important Distinction:

**"Isometric" in browser games typically means:**
- **2D sprites drawn at isometric angles** (pseudo-3D)
- Buildings pre-rendered from a 3D model, then saved as 2D images
- **NOT real-time 3D rendering**
- "The main feature of a two-dimension isometric game is that sprites have a 3D appearance while being still 2D. This is done by rotating the character or object at a 45-degree angle."

**True 3D Isometric:**
- Real-time 3D models rendered with an orthographic camera
- Requires WebGL/Three.js/Babylon.js
- **NOT common in browser strategy games** until ~2015+

**Conclusion:** When browser strategy games claim "isometric graphics," they almost always mean **2D pre-rendered sprites**, not 3D models.

---

## 4. Building Placement "Ghost Preview" Patterns

**Sources:**
- [Manor Lords Discussion](https://steamcommunity.com/app/1363080/discussions/0/4362374970374094183/)
- [Placement Preview Mod](https://www.curseforge.com/minecraft/mc-mods/placement-preview)

### Common UI Pattern:
When placing buildings in strategy/building games, the standard visual feedback is:

1. **Transparent/translucent "ghost" sprite** of the building
2. **Color-coded validity:**
   - **Green tint:** Valid placement
   - **Red tint:** Invalid placement (blocked, out of resources, etc.)
3. **Snaps to grid** (if grid-based) or free-form positioning
4. **Click to confirm placement**

### Implementation in 2D Sprite Games:
- Ghost preview is the **same 2D sprite** rendered with transparency
- Color tint applied via shader or overlay
- **Simple and performant**

### Implementation in 3D Games (Three.js):
- Ghost preview is the **3D model** rendered with transparent material
- Color tint via material properties
- Requires **more complex scene management**

**Conclusion:** Ghost preview pattern is **consistent across 2D and 3D games**, but implementation differs significantly.

---

## 5. Three.js vs 2D Pre-Rendered: Technical Comparison

**Sources:**
- [Best JavaScript Game Engines 2025](https://blog.logrocket.com/best-javascript-html5-game-engines-2025/)
- [Three.js vs Babylon.js](https://blog.logrocket.com/three-js-vs-babylon-js/)
- [3D Games on the Web - MDN](https://developer.mozilla.org/en-US/docs/Games/Techniques/3D_on_the_web)

### 2D Pre-Rendered Sprites (Ikariam Style)

**Pros:**
- **Simpler to implement** - no 3D scene management
- **Better performance** - just rendering 2D images
- **Smaller file sizes** - sprite sheets vs 3D models
- **Easier art pipeline** - 2D artists can work directly
- **Consistent with genre conventions**

**Cons:**
- **Less flexible** - limited camera angles
- **Scaling/zoom is harder** - sprites pixelate
- **Less "premium" feel** - may look dated

### Three.js Real-Time 3D

**Pros:**
- **Fully dynamic camera** - pan, zoom, rotate freely
- **Modern, premium appearance**
- **Flexible lighting and effects**
- **Scalable assets** - no pixelation

**Cons:**
- **Complex implementation** - scene graph, materials, lighting
- **Performance overhead** - GPU rendering for WebGL
- **Larger file sizes** - 3D models + textures
- **Requires 3D art skills** - Blender/3D modeling
- **NOT typical for browser strategy games**

### Hybrid Approach (2.5D)
- Render 2D sprites in a 3D scene (Three.js Plane geometries)
- Benefits: Better camera control than pure 2D
- Still uses sprite assets, not true 3D models

**Conclusion:** Three.js is **overkill for a GO2-style planet view** unless you want a significantly different (more modern) visual style than the original game.

---

## 6. Recommendations for Cryptomines Online

### Based on Research Findings:

**If Goal is 1:1 Copy of Galaxy Online 2:**
- **Use 2D pre-rendered sprites** for buildings
- **Simpler grid or fixed building slots** (like Ikariam)
- **Overhead or slight angle view**, not true isometric
- **Focus on ship design/combat** as primary visual elements (per GO2 design)
- **De-scope Three.js** - it's not period-accurate for a 2010-2015 Flash game clone

**If Goal is "Modern Remake" of Galaxy Online 2:**
- **Three.js can be justified** as a modern enhancement
- But acknowledge this is a **significant visual departure** from the original
- Planet view was **not the main focus** of GO2 - ships were

### Current Cryptomines Online Implementation:
Based on the Phase 1 UI Redo research document (if available), the current implementation uses:
- **Three.js isometric 12x12 grid**
- **3D procedural models** for buildings
- **Fully 3D scene** with camera controls

**Is This Accurate to GO2?**
- **NO** - this is significantly more advanced than GO2's likely 2D sprite-based approach
- This is a **modern remake style**, not a 1:1 clone

**Should We Change It?**
- **Depends on project goals:**
  - If "1:1 copy" is paramount → switch to 2D sprites
  - If "modern remake inspired by GO2" → keep Three.js

**Key Question:** Does "Copia 1:1" apply to **visual style** or just **mechanics/systems**?

---

## 7. Visual Evidence Gaps

### What We Could NOT Find:
- **Actual screenshots** of GO2's planet building interface
- **Gameplay videos** showing base construction
- **Detailed wiki descriptions** of the visual interface
- **Whether GO2 had a grid system** or fixed building slots
- **Whether buildings were placed spatially** or abstracted (like OGame)

### Why This is Challenging:
- **Game is shut down** (no longer playable)
- **Limited archival documentation** - Flash games from this era are poorly preserved
- **Community wikis focused on mechanics**, not visuals
- **No official design documents** available publicly

### What This Means:
Any reconstruction of GO2's visual style requires **educated guessing** based on:
- Technical constraints of Flash/Facebook games (2010-2015)
- Common patterns in similar browser strategy games
- General descriptions of GO2 as a "Flash game with decent graphics"

**Conclusion:** **Perfect 1:1 visual accuracy is impossible** without direct access to the original game or archived screenshots.

---

## 8. Final Summary

### Key Takeaways:

1. **Galaxy Online 2 likely used 2D pre-rendered sprites or Flash graphics**, not true 3D rendering

2. **Most browser strategy games from 2010-2015 used:**
   - 2D sprites (Ikariam, Forge of Empires)
   - OR text/list-based interfaces (OGame, Travian)
   - NOT real-time 3D (Three.js, WebGL)

3. **"Isometric" in this genre means 2D sprites drawn at isometric angles**, not 3D models

4. **Three.js is a modern enhancement**, not period-accurate for a 2010-2015 Flash game

5. **Direct visual evidence of GO2's planet interface is scarce** due to game shutdown and poor archival

### Decision Points:

**For Cryptomines Online Project:**

- **Clarify "Copia 1:1" scope:** Does it apply to visual style or just mechanics?
- **If 1:1 visual accuracy is required:** Consider switching from Three.js to 2D sprites
- **If modern remake is acceptable:** Keep Three.js but acknowledge it's an enhancement
- **Planet building was NOT the focus of GO2** - ship design and combat were

### Recommended Next Steps:

1. **Discuss with team:** Is the current Three.js implementation aligned with project goals?
2. **Search for archived GO2 gameplay videos** (YouTube, Internet Archive)
3. **Prioritize ship design/combat visuals** over planet building (per GO2's design philosophy)
4. **If keeping Three.js:** Frame it as "GO2-inspired" rather than "GO2 clone"

---

## Sources

### Primary Research:
- [Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki)
- [Galaxy Online II - Planet Base](https://galaxyonlineii.fandom.com/wiki/Planet_Base)
- [Galaxy Online II - Construction](https://galaxyonlineii.fandom.com/wiki/Construction)
- [OGame Wiki - Interface](https://ogame.fandom.com/wiki/Interface)
- [OGame Wiki - Building Screen](https://ogame.fandom.com/wiki/Building_Screen)
- [Ikariam Wiki - Town View](https://ikariam.fandom.com/wiki/Town_view)

### Comparative Analysis:
- [Games Like Galaxy Online 2](https://www.igdb.com/games/galaxy-online-2/similar)
- [Browser Based Strategy Games](https://mmos.com/review/browser-games/strategy)
- [Isometric Video Game Graphics - Wikipedia](https://en.wikipedia.org/wiki/Isometric_video_game_graphics)
- [Isometric Sprite Sheet Guide](https://retrostylegames.com/blog/isometric-sprite-sheet/)

### Technical Comparison:
- [Best JavaScript Game Engines 2025](https://blog.logrocket.com/best-javascript-html5-game-engines-2025/)
- [Three.js vs Babylon.js](https://blog.logrocket.com/three-js-vs-babylon-js/)
- [3D Games on the Web - MDN](https://developer.mozilla.org/en-US/docs/Games/Techniques/3D_on_the_web)

### UI Patterns:
- [Manor Lords - Ghost Preview Discussion](https://steamcommunity.com/app/1363080/discussions/0/4362374970374094183/)
- [Placement Preview Mod](https://www.curseforge.com/minecraft/mc-mods/placement-preview)

---

**End of Report**
