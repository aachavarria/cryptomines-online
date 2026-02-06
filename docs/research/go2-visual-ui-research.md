# Galaxy Online 2 - Visual UI Research

> Research compiled from: Galaxy Online II Wiki (Fandom), GameYum guides, Gamezebo review, GamerTrust review, MMORPG.com, NewRPG.com, F2PG.com, WordPress blogs, and various web sources.
> Game by IGG. Shut down January 12, 2021. Was a Flash-based browser MMORTS (also on iOS/Android with different UI).

---

## 1. Overall Visual Theme & Aesthetic

### Color Scheme & Art Style
- **Dark sci-fi space theme** - deep blues, blacks, purples as dominant background colors
- **Pink and purple** used heavily for ships and accents (noted by reviewers as getting repetitive)
- **Turquoise/cyan** color highlights visible in space station interfaces
- **2D rendered buildings** with pre-rendered sprite art on an **isometric grid** (not 3D, but uses isometric perspective with diamond-shaped tiles)
- **High-end sci-fi graphics** for a Flash browser game of its era (2010-2013)
- **Dark themed UI panels** with organized rows, data tables, and clickable elements
- Flash-based rendering, fixed resolution within browser frame

### General Layout Paradigm
- **Window/panel-based UI** - NOT a free-roaming 3D world
- Buildings displayed as **clickable sprite icons on a static base view**
- Separate **tabs/screens** for different game systems (base, fleet, research, galaxy, social)
- Heavy use of **data tables**, **stat panels**, and **pop-up menus**
- Described by reviewers as "complex and counter-intuitive" with a "steep learning curve"
- Multiple nested menus required to complete actions (e.g., building ships requires: Ship Factory > Design > Add modules > Save > Build)

---

## 2. Main Navigation Structure

### Primary View Modes (3 main views)
The game has **three distinct spatial views** that players switch between:

1. **Ground Base** (Planet Surface) - Default view, where most buildings are constructed
2. **Space Base** (Orbital Station) - Click the "green icon that looks like planet Saturn" to access
3. **Galaxy View** - Accessed via dedicated icon, shows planets and fleet movements

### Navigation Bar / HUD Buttons
Based on wiki instructions referencing specific clicks:

- **Build button** - Located in the **lower right corner**, opens construction menu
- **Military button** - Located in the **lower right corner**, accesses fleet/military functions
- **Space Station button** - Switches to orbital/space base view (the Saturn-like green icon)
- **Galaxy button** - Switches to galaxy map view
- **Social button** (green and blue person icon) - Opens friends, corps, rankings
- **Tools button** (antenna icon) - Opens mail, settings, utilities
- **Friends button** - Under Social, shows friends list (5 per page)

### Bottom Bar
- Friends tab dropdown appears at the bottom when hovering over friends button
- Build and Military buttons confirmed at **lower right**

### Quest Interface
- Quests appear to have a dedicated side panel or popup
- Main quest, side quest, and daily quest tabs
- Quest completion tracking with reward breakdowns
- Pop-up help boxes provide in-game guidance

---

## 3. Resource HUD / Top Bar

### Three Main Resources
1. **Metal** - Produced by Metal Collectors (up to 8, max level 24)
2. **He3 (Helium-3)** - Produced by HE3 Extractors (up to 8, max level 24)
3. **Gold** - Produced by Residential Areas (income-based)

### Display Location
- Resources displayed in a **top bar area** of the interface (typical for browser strategy games of this era)
- Resources show current amount with icons
- Resource production rates visible through building interfaces
- Resource Warehouse shows accumulated resources ready to harvest
- Players manually harvest from the Resource Warehouse when they want to collect

### Additional Currencies
- **Mall Points (MP)** - Premium currency for Web Mall purchases
- **SP (Siege Points)** - Used for repairs and boosts
- **Championship/Honor Badges** - Arena rewards

---

## 4. Ground Base View (Planet Surface)

### Layout
- Buildings displayed as **clickable sprite icons** arranged on the planet surface
- Each building has a unique sprite/image (e.g., CivicCenter.png, Metal_collector.PNG, He3_extractor.PNG, Warehouse.PNG, Tech_center.PNG, Factory.png)
- Buildings appear as **industrial/futuristic structures** appropriate to their function
- **Grid-based placement system** - the planet surface uses an **isometric grid of 1x1 diamond tiles** where buildings can be freely placed by the player (not fixed positions/slots)
- Background is the planet surface with a space/dark theme

### Building Categories on Ground Base

**Main Structures:**
- Space Station (gateway structure)
- Civic Center (administrative hub, central to unlocking other buildings)

**Resource Buildings (up to 8 each):**
- Metal Collector - mining/industrial structure
- HE3 Extractor - extraction facility
- Residential Area - housing/settlement
- Resource Warehouse (1 per planet) - grey industrial warehouse

**City Services:**
- Alliance Center - multiplayer corps building
- Compound Center - commander card management
- Technology Center - research hub
- Trading Center - marketplace
- Galaxy Transporter - PvP/league access

**Landscaping (decorative, provide morale/bonuses):**
- Casino Resort, Beacon, Monument, Fountain, Library
- Theater, Park, College, Hospital, Shopping Center
- Statue, Santa Sculpture

**Military Buildings:**
- Command Center - fortress-like, recruits commanders
- Ship Factory - industrial shipyard
- Weapon Research Center - weapons lab
- Recycling Plant - ship salvage
- Radar - radar dish structure
- Spacedock - ship repair facility

### Building Interactions
- **Left-click a building** to open a **context menu** with 3 blue buttons stacked vertically: **View / Move / Upgrade**
  - **View** opens the building's management panel (building level, stats, details)
  - **Move** enters move mode to reposition the building on the grid
  - **Upgrade** initiates the building upgrade
  - For the **Warehouse**, a 4th button appears: **Harvest** (to collect stored resources)
- Management panel (via View) shows: building level, upgrade requirements, stats
- **Upgrade button** shows requirements for each level
- **Blue/yellow "up" arrow** beside construction in progress for acceleration
- **Flashing yellow wrench icon** indicates a building needing attention
- **5 construction slots** available simultaneously for upgrades
- Construction timers visible showing remaining build time
- Listed on the right side of the screen when visiting

### Grid System
- The planet surface uses an **isometric grid of 1x1 diamond tiles**
- When constructing or moving a building, the grid tiles become visible and color-coded:
  - **Green tiles** = valid placement position
  - **Red tiles** = invalid placement position (occupied or restricted)
- Buildings can be freely placed anywhere on valid grid tiles

### Hover Tooltip
- Hovering the mouse over a building displays a **dark tooltip** showing: `Lv: X [Building Name]`
- Example: "Lv: 1 Technology Center"

### Construction Visual Feedback
- Buildings under construction show a **yellow progress bar** displayed over the building sprite
- In the **lower-right corner** of the screen, a **construction info panel** appears showing:
  - Building name
  - Building level
  - Countdown timer (e.g., "Metal Collector Lv: 1 00:00:32")

### Camera Controls
- The planet surface view supports **camera pan** (click and drag to move the viewport)
- Allows players to scroll around the planet surface to view all buildings

### Warehouse Details
- Only **1 Warehouse per planet** is allowed
- Stores each resource type (Metal, He3, Gold) separately
- Its context menu has **4 options** instead of the standard 3: **View / Move / Upgrade / Harvest**
- The Harvest button collects accumulated resources from the warehouse

---

## 5. Space Base View (Orbital)

### Access
- Click the **green Saturn-like icon** in the navigation
- Shows your Space Station in the center

### Layout
- Space Station is the **central structure** surrounded by orbiting defense platforms
- Defensive structures arranged around the station
- Fleet arrangement interface available here
- Dark space background with the station floating

### Space Base Buildings (Defenses)
- **Meteor Star** - space defense
- **Particle Cannon** - energy weapon platform
- **Anti-aircraft Gun** - projectile defense
- **Thor's Cannon** - heavy weapons platform

### Defense Slots
- Number of defense positions determined by Space Station level:
  - Level 1: 2 defenses
  - Higher levels: progressively more (up to 40+ at max level)
- Defenses won't stop attackers but slow them down

### Additional Space Features
- **Celestial Base** - subsidiary resource operations
- **Fleet management** - arrange and deploy fleets from here
- **Instance Map** - access to PvE instances/dungeons
- **Constellation instances** - accessed via Space Station > Instances > Constellations

---

## 6. Galaxy View

### Access
- Via Galaxy button in the navigation HUD

### Features
- Shows planetary systems and fleet movements
- **Radar symbols** indicate ships in transit
- **White flag** on a planet indicates it was recently attacked (prevents leveling)
- Player planets visible with coordinate system
- Can check other players' planets and their space bases
- Fleet transit tracking visible
- Corps territory management

---

## 7. Key Game Screens & Panels

### 7a. Ship Factory / Ship Design Screen
- **Two-panel layout:**
  - **Left panel**: Shows saved blueprints (ship designs you've already created)
  - **Right panel**: Design button that opens the customization area
- **Design area** allows adding weapons, shields, engines, and components
- **Module types**: Ballistic weapons, Directional weapons, Missiles, Shield modules, Structure modules
- **Stats displayed at bottom**: Attack, Shields, Range, Storage, Hull Structure, Movement, Build Time, Transmission Time
- **Ship name** can be entered
- Up to **20 different ship designs** can be stored
- **5 production slots** for simultaneous construction
- Maximum 2 million ships in production simultaneously
- Blueprint research tab links to weapons factory for upgrades
- Arrow button below blueprints to start construction

### 7b. Technology Center / Research Screen
- Access to **7 scientific research fields**
- Each field has multiple technologies in a tree structure
- Stats table showing: research time reduction per level, upgrade costs
- Research time reduction: 3% at Level 1 up to 36% at Level 12

### 7c. Command Center / Commander Screen
- Commander recruitment through **random draw system**
- Commander rarities: Common, Skill, Super, Legendary
- Commander Cards can be drawn using Mall Points
- Recruitment cooldown timer displayed
- Commander stats and abilities visible
- Reviewers note: commanders vary visually, with female commanders notably "risque"

### 7d. Compound Center
- Commander card compounding (fusion) interface
- Gem management
- Bionic Chip management
- Commander Ranks and Levels progression display

### 7e. Trading Center / Marketplace
- Buy/sell interface for: ships, blueprints, gems, boosts, commander cards
- Prices set by sellers
- Payment via MP (Mall Points) or Gold
- Commission system: 3-10% based on level and listing duration (12/24/48 hours)
- **Black Market** tab: resource conversion (Metal to He3 and vice versa)

### 7f. Instance / Dungeon Selection
- **Instance Map** showing geographic layout of instances
- **30 normal instances** displayed in a table:
  - Instance icons (unique per instance)
  - Instance names (e.g., "Ancestral Recall", "Deadzone")
  - Level numbers (1-30)
  - Max fleet capacity
  - Checkpoints
  - Experience rewards
  - Treasure box reward icons
- **Constellation Instances**: 12 zodiac-themed instances
  - Grid-based selection screen with zodiac symbols
  - Blue-illuminated planets/signs for available instances
  - Dark theme with constellation imagery

### 7g. Arena (PvP)
- Accessed via: Space Station > Arena button
- Room-based system for 1v1 battles
- Password protection option
- Spectator mode available
- Unlimited fleet deployment
- Practice mode (no gains/losses)

### 7h. Combat View
- **Graphic-based combat** that is "fun to watch"
- Fleets of **9 "stacks"** of ships firing at opponents
- Players cannot control ship movement during combat
- Combat outcomes determined by ship design (weapons, components, modules)
- Multiple rounds to complete a battle

### 7i. Chat System
- **World chat** - white text, visible to all
- **Corps chat** - green text, visible only to corps members
- Multiple specified chat channels

### 7j. Friends Screen
- 5 friends per page
- Online status indicators
- Hover over name reveals menu: Check, Invite, Compose, PM, Delete
- Can visit friend's planet and space base

### 7k. Web Mall / Shop
- Premium currency (Mall Points) purchases
- Categories: Development items, Battle items, Collectibles
- Currency exchange interface

---

## 8. Common UI Elements

### Buttons & Controls
- **Upgrade buttons** on all building panels
- **Blue/yellow directional arrows** for construction acceleration
- **Build button** (lower right) to access construction menu
- **Military button** (lower right) for fleet operations
- **Design button** in Ship Factory for ship customization
- **Proceed to destroy** button in Recycling Plant

### Visual Indicators
- **Construction progress bars** with timers showing remaining time
- **Flashing yellow wrench** icon for buildings needing attention
- **White flag** on attacked planets
- **Blue illumination** for available/selectable items
- **Online/offline status** indicators for friends

### Data Display Patterns
- Heavy use of **stat tables** showing level progression, costs, bonuses
- **Tooltips/pop-up help boxes** throughout the interface
- **Resource cost breakdowns** (Metal, He3, Gold amounts)
- **Timer displays** for construction and research
- **Percentage indicators** (research bonuses, repair rates)

### Panel/Window Style
- Dark-themed panels with organized data rows
- **Modal-style** pop-ups for building details
- Tabbed interfaces (quest tabs: main/side/daily)
- Dropdown menus for social features
- Window-based layout (not full-screen transitions)

---

## 9. Mobile Version (iPad/Android)

- "Playing on the phone or iPad gives the players the feeling of playing a brand new game as the UI is totally different from the Facebook version"
- Runs on separate servers (no account migration from Facebook)
- Touch-optimized interface
- Same game mechanics but completely different visual presentation

---

## 10. Summary: Key Design Patterns for Replication

### What Cryptomines Online needs to replicate:

1. **Dark sci-fi space theme** - Deep blue/black backgrounds, cyan/teal accents, futuristic metallic UI elements
2. **Three-view navigation** - Ground Base, Space Base, Galaxy View with clear toggle buttons
3. **Resource bar at top** - Metal, He3, Gold with icons and current amounts
4. **Bottom navigation buttons** - Build, Military, and view-switching icons
5. **Building sprites on isometric grid** - Each building type has a unique visual representation freely placed on an isometric 1x1 tile grid, with green/red placement indicators
6. **Context menu on click** - Left-clicking a building opens a context menu (View/Move/Upgrade), with View opening the management panel showing stats and costs
7. **Construction queue** - 5 slots with progress timers visible
8. **Ship design dual-panel** - Left: saved blueprints, Right: design editor with module slots
9. **Data-heavy tables** - Stat tables for research, upgrades, costs throughout
10. **Dark modal panels** - For building details, ship design, research, etc.
11. **Chat system** - World (white) and Corps (green) channels
12. **Instance map** - Visual map with numbered instances and rewards
13. **Commander cards** - Visual cards with rarity tiers and draw system
14. **Space station centered** - Defense platforms orbiting around central station

### What to NOT replicate:
- Flash-based rendering limitations
- The overly complex multi-menu navigation (can be streamlined slightly)
- Pink/purple color dominance (use our own accent colors)

---

## Sources

- Galaxy Online II Wiki (Fandom): https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki
- Walkthrough: https://galaxyonlineii.fandom.com/wiki/Walkthrough
- Beginner FAQ: https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ
- Beginner's Guide: https://galaxyonlineii.fandom.com/wiki/Beginner's_Guide
- Guide To Advancing Quickly: https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly
- Building Categories: https://galaxyonlineii.fandom.com/wiki/Category:Buildings
- GameYum Structures Guide: https://www.gameyum.com/galaxy-online/115041-galaxy-online-ii-guides-building-structures/
- GameYum Overview: https://www.gameyum.com/galaxy-online/109124-facebookgame-guides-galaxy-online-ii/
- GameYum Ship Guide: https://www.gameyum.com/galaxy-online/109247-galaxy-online-ii-guides-how-to-build-ships/
- GameYum Economy Guide: https://www.gameyum.com/galaxy-online/109252-galaxy-online-ii-guides-economy-guide/
- Gamezebo Review (2011): https://www.gamezebo.com/2011/03/14/galaxy-online-2-review/
- GamerTrust Review (2012): https://gamertrust.wordpress.com/2012/01/05/galaxy-online-2-an-intricate-strategy-mmo-on-facebook/
- MMORPG.com: https://www.mmorpg.com/galaxy-online-2
- NewRPG: https://newrpg.com/browser-games/galaxy-online-2/
- F2PG: https://www.f2pg.com/galaxy-online-2/
- 4GameGround: https://4gameground.com/galaxy-online-2/
- Blog with Client Screenshots: https://m3rilix.wordpress.com/2011/03/04/galaxy-online-iiclient/
- GO2 Shutdown Article: https://deajae.co.uk/posts/the-end-of-galaxy-online-2/
- Ship Creation Guide: https://galaxyonline2cheats.wordpress.com/2012/04/13/galaxy-online-ii-ship-creation-guide-strategy-games/
- Individual Building Pages: Civic Center, Ship Factory, Technology Center, Command Center, Radar, Metal Collector, Resource Warehouse, HE3 Extractor, Spacedock, Trading Center, Alliance Center, Galaxy Transporter, Space Station, Compound Center
- Normal Instances: https://galaxyonlineii.fandom.com/wiki/Normal_Instances
- Constellation Instances: https://galaxyonlineii.fandom.com/wiki/Constellation_Instances
- The Arena: https://galaxyonlineii.fandom.com/wiki/The_Arena
- Friends: https://galaxyonlineii.fandom.com/wiki/Friends
- Development Quests: https://galaxyonlineii.fandom.com/wiki/Development_Quests
