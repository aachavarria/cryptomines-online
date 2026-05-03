# Cryptomines — Design System

> Soft sci-fi, clean, legible. Teal as identity, orange only for action.
> Adapted for a client with a 3D scene rendered behind the entire UI.

---

## 1. Philosophy

1. **Clean first** — no unnecessary visual noise. Every border, glow, or gradient must justify itself.
2. **Color with purpose** — teal = identity/info, orange = action, green = success/resources, red = danger. Everything else is gray.
3. **Space = clarity** — more padding, less density. When in doubt, add air.
4. **Soft sci-fi** — futuristic but accessible: no saturated neon, no military frames, no pure black.

### 3D-client adaptation
The scene (`<Canvas>`) fills the viewport. Panels float on top as **opaque light surfaces** (no aggressive glassmorphism). This requires:
- Panel background `#FFFFFF` or `#F9FAFB` at **opacity 0.96–1.0** (not see-through).
- `backdrop-filter: blur(12px)` only on HUD and context menus that genuinely need to bleed the scene through.
- Soft shadows, not bright glows (which dissolve against starfields).

---

## 2. Color tokens

### Base / surfaces
```css
--ds-bg:            #F4F6F8;   /* page fallback (rarely visible — 3D is behind) */
--ds-surface:       #FFFFFF;   /* main panels */
--ds-surface-2:     #F9FAFB;   /* cards, inner sections */
--ds-surface-3:     #F3F4F6;   /* hovered list items, secondary buttons */
--ds-overlay:       rgba(15, 23, 42, 0.55); /* modal backdrop */
--ds-border:        #E5E7EB;
--ds-border-strong: #D1D5DB;
```

### Text
```css
--ds-text:          #1F2937;   /* primary */
--ds-text-muted:    #6B7280;   /* secondary, labels */
--ds-text-soft:     #9CA3AF;   /* hints, captions, disabled */
--ds-text-inverse:  #FFFFFF;   /* on top of teal/orange/dark fills */
```

### Brand
```css
--ds-teal:          #1FA39A;   /* primary identity */
--ds-teal-dark:     #167F78;   /* hover/pressed */
--ds-teal-light:    #6FD3CC;   /* soft backgrounds, charts */
--ds-teal-tint:     rgba(31, 163, 154, 0.10); /* active tab fills, selection */

--ds-orange:        #F59E0B;   /* primary CTA */
--ds-orange-strong: #F97316;   /* CTA hover */
--ds-orange-tint:   rgba(245, 158, 11, 0.10);
```

### State
```css
--ds-success:       #22C55E;
--ds-success-tint:  rgba(34, 197, 94, 0.12);
--ds-warning:       #F59E0B;   /* same as orange — disambiguate by context */
--ds-warning-tint:  rgba(245, 158, 11, 0.12);
--ds-danger:        #EF4444;
--ds-danger-tint:   rgba(239, 68, 68, 0.12);
--ds-info:          #0EA5E9;
--ds-info-tint:     rgba(14, 165, 233, 0.12);
```

### Game resource mapping
We keep parity with the existing in-game resources:
```css
--ds-metal:         #6B7280;   /* neutral metallic gray */
--ds-he3:           #0EA5E9;   /* sober cyan, calmer than current */
--ds-gold:          #F59E0B;   /* matches orange/warning */
--ds-energy:        #8B5CF6;   /* purple for energy/science */
```

### Rarity (commanders, blueprints, items)
```css
--ds-rarity-common:    #9CA3AF;
--ds-rarity-uncommon:  #22C55E;
--ds-rarity-rare:      #3B82F6;
--ds-rarity-epic:      #8B5CF6;
--ds-rarity-legendary: #F59E0B;
```

---

## 3. Typography

**Font:** Inter (variable). Load via Google Fonts or `@fontsource/inter`. Fallback: `system-ui`.
Keep `JetBrains Mono` only for timers and tabular numbers.

> `Rajdhani` and `Orbitron` (current) are **removed** — they don't fit the soft sci-fi language.

### Scale (rem, base 16px)
| Token | px | weight | use |
|---|---|---|---|
| `--fs-h1` | 28 | 600 | large modal titles |
| `--fs-h2` | 20 | 600 | panel titles |
| `--fs-h3` | 16 | 600 | sections |
| `--fs-body` | 14 | 400 | general text |
| `--fs-sm` | 13 | 400 | meta, descriptions |
| `--fs-caption` | 12 | 500 | uppercase labels |
| `--fs-button` | 13 | 600 | buttons (UPPERCASE, letter-spacing 0.06em) |

```css
--ff-sans: 'Inter', system-ui, -apple-system, sans-serif;
--ff-mono: 'JetBrains Mono', 'Fira Code', monospace;
--lh-tight: 1.2;
--lh-normal: 1.5;
--lh-loose: 1.7;
```

---

## 4. Spacing (4px grid)

```css
--sp-0: 0;
--sp-1: 4px;
--sp-2: 8px;
--sp-3: 12px;
--sp-4: 16px;
--sp-5: 20px;
--sp-6: 24px;
--sp-8: 32px;
--sp-10: 40px;
--sp-12: 48px;
```

Default panel padding: `--sp-5` (cards) or `--sp-6` (modals).

---

## 5. Radii, shadows, glows

```css
--r-sm:  6px;    /* badges, chips, tooltips */
--r-md:  10px;   /* buttons, inputs, list items */
--r-lg:  12px;   /* cards, inner panels */
--r-xl:  16px;   /* main panels, modals */
--r-pill: 999px;

--shadow-1: 0 1px 2px rgba(15, 23, 42, 0.04), 0 1px 3px rgba(15, 23, 42, 0.06);
--shadow-2: 0 2px 8px rgba(15, 23, 42, 0.06), 0 4px 12px rgba(15, 23, 42, 0.04);
--shadow-3: 0 8px 24px rgba(15, 23, 42, 0.10), 0 2px 6px rgba(15, 23, 42, 0.06);
--shadow-modal: 0 24px 60px rgba(15, 23, 42, 0.20);

--glow-teal:   0 0 0 3px rgba(31, 163, 154, 0.18);   /* focus / selected */
--glow-orange: 0 0 0 3px rgba(245, 158, 11, 0.22);   /* CTA focus */
--glow-danger: 0 0 0 3px rgba(239, 68, 68, 0.20);
```

Global motion:
```css
--motion-fast:   120ms ease-out;
--motion-base:   180ms ease-out;
--motion-slow:   260ms cubic-bezier(0.2, 0.8, 0.2, 1);
```

---

## 6. Components

### 6.1 Panel / Card
```css
.ds-panel {
  background: var(--ds-surface);
  border: 1px solid var(--ds-border);
  border-radius: var(--r-xl);
  box-shadow: var(--shadow-2);
  padding: var(--sp-5);
}
.ds-card {
  background: var(--ds-surface-2);
  border: 1px solid var(--ds-border);
  border-radius: var(--r-lg);
  padding: var(--sp-4);
}
```

### 6.2 Buttons
**Primary (CTA — orange, outline → fill on hover):**
```css
.ds-btn-primary {
  background: transparent;
  color: var(--ds-orange);
  border: 2px solid var(--ds-orange);
  border-radius: var(--r-pill);
  padding: 9px 18px;
  font: 600 var(--fs-button)/1 var(--ff-sans);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  transition: var(--motion-base);
}
.ds-btn-primary:hover { background: var(--ds-orange); color: #fff; }
.ds-btn-primary:focus-visible { box-shadow: var(--glow-orange); outline: none; }
.ds-btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
```

**Secondary (teal solid — common positive action: build/upgrade):**
```css
.ds-btn-secondary {
  background: var(--ds-teal);
  color: #fff;
  border: none;
  border-radius: var(--r-md);
  padding: 9px 16px;
}
.ds-btn-secondary:hover { background: var(--ds-teal-dark); }
```

**Ghost (neutral — cancel, close):**
```css
.ds-btn-ghost {
  background: var(--ds-surface-3);
  color: var(--ds-text);
  border: 1px solid var(--ds-border);
  border-radius: var(--r-md);
  padding: 8px 14px;
}
.ds-btn-ghost:hover { background: var(--ds-border); }
```

**Danger (destructive, demolish):**
```css
.ds-btn-danger {
  background: transparent;
  color: var(--ds-danger);
  border: 1px solid var(--ds-danger);
}
.ds-btn-danger:hover { background: var(--ds-danger); color: #fff; }
```

**Icon button (close X, settings):**
```css
.ds-btn-icon {
  width: 32px; height: 32px;
  background: var(--ds-surface-3);
  border: none;
  border-radius: var(--r-md);
  color: var(--ds-text-muted);
  display: grid; place-items: center;
}
.ds-btn-icon:hover { background: var(--ds-border); color: var(--ds-text); }
```

### 6.3 Tabs
```css
.ds-tabs { display: flex; gap: 4px; border-bottom: 1px solid var(--ds-border); }
.ds-tab {
  padding: 10px 14px;
  color: var(--ds-text-muted);
  font: 600 var(--fs-sm)/1 var(--ff-sans);
  text-transform: uppercase; letter-spacing: 0.05em;
  border: none; background: transparent;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}
.ds-tab[aria-selected="true"] {
  color: var(--ds-teal);
  border-bottom-color: var(--ds-teal);
}
```

### 6.4 List item (selectable)
```css
.ds-list-item {
  background: var(--ds-surface);
  border: 2px solid transparent;
  border-radius: var(--r-lg);
  padding: var(--sp-3);
  transition: var(--motion-fast);
}
.ds-list-item:hover { background: var(--ds-surface-2); }
.ds-list-item[aria-selected="true"] {
  border-color: var(--ds-orange);    /* selection = orange */
  background: var(--ds-surface);
}
```

### 6.5 Stat bar (with optional segments)
```css
.ds-bar { height: 6px; background: var(--ds-border); border-radius: var(--r-pill); overflow: hidden; }
.ds-bar-fill { height: 100%; background: var(--ds-teal); transition: width var(--motion-slow); }
/* Segmented: use background with repeating-linear-gradient with 2px gaps */
```

Color variants: `.ds-bar-fill--success` (green), `.ds-bar-fill--warning` (orange), `.ds-bar-fill--danger` (red).

### 6.6 Badge / chip
```css
.ds-badge {
  display: inline-flex; align-items: center; gap: 4px;
  padding: 2px 8px;
  border-radius: var(--r-sm);
  font: 600 var(--fs-caption)/1.4 var(--ff-sans);
  text-transform: uppercase; letter-spacing: 0.04em;
}
.ds-badge--teal    { background: var(--ds-teal-tint);    color: var(--ds-teal-dark); }
.ds-badge--orange  { background: var(--ds-orange-tint);  color: var(--ds-orange-strong); }
.ds-badge--success { background: var(--ds-success-tint); color: #15803d; }
.ds-badge--danger  { background: var(--ds-danger-tint);  color: #b91c1c; }
.ds-badge--neutral { background: var(--ds-surface-3);    color: var(--ds-text-muted); }
```

### 6.7 Modal
```css
.ds-modal-backdrop {
  position: fixed; inset: 0; z-index: 300;
  background: var(--ds-overlay);
  backdrop-filter: blur(2px);
  display: grid; place-items: center;
  padding: var(--sp-6);
}
.ds-modal {
  background: var(--ds-surface);
  border-radius: var(--r-xl);
  box-shadow: var(--shadow-modal);
  max-width: min(720px, 95vw);
  max-height: 85vh;
  display: flex; flex-direction: column;
  overflow: hidden;
}
.ds-modal-header {
  padding: var(--sp-5) var(--sp-6);
  border-bottom: 1px solid var(--ds-border);
  display: flex; justify-content: space-between; align-items: center;
}
.ds-modal-body { padding: var(--sp-6); overflow-y: auto; }
.ds-modal-footer {
  padding: var(--sp-4) var(--sp-6);
  border-top: 1px solid var(--ds-border);
  display: flex; justify-content: flex-end; gap: var(--sp-3);
  background: var(--ds-surface-2);
}
```

### 6.8 Tooltip (3D and UI)
```css
.ds-tooltip {
  background: rgba(31, 41, 55, 0.95);  /* dark over 3D for contrast */
  color: #fff;
  padding: 6px 10px;
  border-radius: var(--r-sm);
  font: 500 var(--fs-caption)/1.3 var(--ff-sans);
  box-shadow: var(--shadow-2);
  pointer-events: none;
}
```
> Exception to "no dark backgrounds": tooltips appearing over the 3D scene need high contrast. We use a translucent `slate-800`.

### 6.9 Input / select
```css
.ds-input {
  background: var(--ds-surface);
  border: 1px solid var(--ds-border-strong);
  border-radius: var(--r-md);
  padding: 8px 12px;
  font: 400 var(--fs-body)/1.4 var(--ff-sans);
  color: var(--ds-text);
}
.ds-input:focus-visible { border-color: var(--ds-teal); box-shadow: var(--glow-teal); outline: none; }
```

---

## 7. Game-specific patterns

### 7.1 Resource HUD (top bar, 56px)
- Background: `var(--ds-surface)` with `box-shadow: var(--shadow-2)`, **no** blue glow.
- Logo on the left, resources centered (round icon + value in mono + subtle teal delta), Collect button on the right (`ds-btn-secondary`).
- Each resource: 18px outline icon, number in `--ff-mono`, rate `+12/h` in `--ds-text-muted`.
- Warehouse warning: `ds-badge--danger` only when ≥90%.

### 7.2 Side Nav (lateral 60px → 72px)
- Background `var(--ds-surface)`, border-right `var(--ds-border)`.
- Each item: 48×48 with 22px outline icon + 9px uppercase label.
- Active state: teal-tint fill, icon and label in `--ds-teal`, **no** border (the difference is color, not chrome).
- Side tooltip via `.ds-tooltip`.
- Ground/Space switcher: two horizontal pills at the top, not mixed-in nav items.

### 7.3 Building context menu (radial/popover)
- `ds-panel` with `padding: var(--sp-2)`, ~140px wide.
- Items as full-width `.ds-btn-ghost`, left-aligned, with 16px outline icon.
- "Demolish" uses `.ds-btn-danger` (red text on transparent).
- Entry animation: `scale(0.96) → 1`, `--motion-fast`.

### 7.4 Building detail panel (right slide-in, 380px)
- `.ds-panel` with no border-radius on right/top/bottom (flush with the viewport).
- Header: 56×56 icon with `ds-badge` category beneath, H2 name, "Level X · Y/max" muted.
- Sections separated by `border-bottom: 1px solid var(--ds-border)`, padding `--sp-5`.
- Costs: each resource with 8px dot in resource color + amount + green check / red X.
- Upgrade button: full-width `.ds-btn-secondary` (teal). At max level → `.ds-badge--orange` "MAX LEVEL", no button.

### 7.5 Construction panel (modal)
- `.ds-modal`, 760px wide.
- Category tabs (`.ds-tabs`) with `(3/5)` count in muted.
- 2-column grid, gap `--sp-4`, each card:
  - Building image/render (16:10 ratio) on top, content below.
  - Name + category chip, costs inline, build button (`.ds-btn-secondary` or disabled if insufficient).
- Search filter in header with `.ds-input` (search icon).

### 7.6 Ground grid (3D overlay — IsometricGrid)
The grid is rendered in Three.js, not HTML. Tokens to propagate to the shader material:
- Valid tile (placement OK): `#22C55E` alpha 0.35, border 0.6.
- Invalid tile: `#EF4444` alpha 0.35.
- Hovered tile: `#1FA39A` alpha 0.45 + stronger outline.
- Selected tile: 2px white outline + teal-light fill 0.25.
- No more colors. No permanent pulsing animations (only during active placement mode).

`PlacementIndicator` (HTML, bottom-center): compact `.ds-panel` with H3 "Placing on Ground Base" + muted hint "ESC to cancel". Teal border.

### 7.7 Galaxy map / Space grid
- 2D map in modal panel or full-screen overlay.
- Background: radial gradient `#0F172A → #020617` (dark exception — galaxy is deep space, not UI).
- Sectors as hex/squares with `#1E293B` border, teal hover, orange selected.
- Right-side info panel for the selected sector: standard **light** theme (white panel over dark background).
- Fleet markers: 6px dots, color by owner (you: teal, ally: blue, enemy: red, neutral: muted).
- Tooltips over sectors: `.ds-tooltip` (dark, already defined).

> Galaxy is the **only** view that breaks the light theme, by thematic necessity. Floating UI panels stay light.

### 7.8 Quest / Research / Inventory panels
- Two-column layout: list on the left (selectable items) + detail on the right.
- List items per §6.4 — selection with orange border.
- Quest progress: segmented `.ds-bar` (one segment per objective).
- Tech tree (research): nodes as `.ds-card` connected by `--ds-border-strong` 2px lines; completed in teal, available with orange outline, locked muted.

### 7.9 Combat reports / battle playback
- Reports as `.ds-list-item` rows with outcome chip (`ds-badge--success/danger`).
- Playback panel: teal progress bar, controls (`.ds-btn-icon`), mono labels for timestamps.

---

## 8. Iconography

- **Lucide React** as the single icon library (install). Replaces ad-hoc emojis and SVGs.
- Stroke width: `1.75` (default) for UI, `2` only on small buttons.
- Sizes: 14, 16, 18, 20, 22, 24.
- Color: inherits from `currentColor` — control via text class (`color: var(--ds-text-muted)` by default, teal when active).
- **Don't mix** emoji + SVG + outline in the same view. Today there are flag/hammer/etc. emojis: replace.

---

## 9. Motion

| Pattern | Duration | Curve |
|---|---|---|
| Hover (color, bg) | 120ms | ease-out |
| Button pressed scale | 80ms | ease-out |
| Modal in/out | 220ms | cubic-bezier(0.2, 0.8, 0.2, 1) |
| Panel slide-in | 260ms | same |
| Tooltip fade | 100ms | ease-out (with 300ms appear delay) |
| Bar fill | 400ms | ease-out |

No infinite animations except: loading spinner, warehouse-full pulse (critical only).

---

## 10. Accessibility

- Minimum contrast 4.5:1 for body text. Verify `--ds-text-muted` (#6B7280) on `--ds-surface-2` (#F9FAFB) → 5.4:1 ✓.
- Focus is **always** visible via `outline: none; box-shadow: var(--glow-teal/orange)`.
- Tap targets ≥ 32×32 (44×44 on mobile if supported).
- `aria-selected`, `aria-pressed`, `role="dialog"` on modals.
- Reduce-motion: honor `@media (prefers-reduced-motion)` by disabling all transitions except opacity.

---

## 11. Rules

✅ Teal = identity / info / positive progress.
✅ Orange = call-to-action and selection.
✅ Plenty of whitespace; generous padding.
✅ Consistent rounded corners (`--r-md` for everything interactive, `--r-xl` for surfaces).
✅ One icon set (Lucide).
✅ Soft shadows over glows.

❌ No pure `#000` (use `--ds-text` = `#1F2937`).
❌ No saturated gradients (except galaxy map and 3D particles).
❌ No mixing Rajdhani/Orbitron/Inter — Inter only (+ JetBrains Mono for numerics).
❌ No decorative emojis in UI (chat user content is fine).
❌ No `border-glow` or `box-shadow: 0 0 12px <color>` on base panels — they wash out against the scene.
❌ No aggressive transparency (< 0.85) on main panels.

---

## 12. Full tokens (CSS)

See §2–§5. Target file: `frontend/src/styles/tokens.css` (replaces `variables.css`).
