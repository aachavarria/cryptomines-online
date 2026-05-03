/**
 * Color tokens for the Three.js scene.
 *
 * Mirrors the Cryptomines design system §7.6 (ground grid) and §7.7 (galaxy / owner colors).
 * Use the numeric `colors3d` constants for Three.js materials (`color={colors3d.tileValid}`)
 * and the string `colors3dHex` constants when a hex string is required (e.g. CSS-in-JS,
 * SVG attributes, or `BUILDING_MODELS` hologram color helpers that expect "#RRGGBB").
 *
 * Spec: docs/reference/design-system.md §7.6 §7.7
 */

export const colors3d = {
  // Placement / grid (§7.6)
  tileValid: 0x22c55e, // alpha 0.35, border 0.6
  tileInvalid: 0xef4444, // alpha 0.35
  tileHover: 0x1fa39a, // alpha 0.45
  tileSelected: 0xffffff, // 2px outline
  tileSelectedFill: 0x6fd3cc, // alpha 0.25

  // Owner colors (§7.7)
  ownerSelf: 0x1fa39a,
  ownerAlly: 0x3b82f6,
  ownerEnemy: 0xef4444,
  ownerNeutral: 0x9ca3af,
} as const

export const colors3dHex = {
  tileValid: '#22C55E',
  tileInvalid: '#EF4444',
  tileHover: '#1FA39A',
  tileSelected: '#FFFFFF',
  tileSelectedFill: '#6FD3CC',
  ownerSelf: '#1FA39A',
  ownerAlly: '#3B82F6',
  ownerEnemy: '#EF4444',
  ownerNeutral: '#9CA3AF',
} as const

/** Opacity values per §7.6. */
export const colors3dOpacity = {
  tileValid: 0.35,
  tileInvalid: 0.35,
  tileHover: 0.45,
  tileSelectedFill: 0.25,
} as const

export type Colors3dKey = keyof typeof colors3d
