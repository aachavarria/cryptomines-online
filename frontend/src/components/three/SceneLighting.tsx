export default function SceneLighting() {
  return (
    <>
      {/* Ambient light - increased for better base illumination */}
      <ambientLight color="#2a2a4a" intensity={0.7} />

      {/* Main top-down light for buildings */}
      <directionalLight
        color="#ffffff"
        intensity={1.8}
        position={[0, 50, 0]}
        castShadow
        shadow-mapSize-width={2048}
        shadow-mapSize-height={2048}
        shadow-camera-left={-60}
        shadow-camera-right={60}
        shadow-camera-top={60}
        shadow-camera-bottom={-60}
      />

      {/* Side fill light */}
      <directionalLight
        color="#aabbff"
        intensity={0.8}
        position={[20, 30, 10]}
      />

      {/* Back rim light */}
      <directionalLight
        color="#6688cc"
        intensity={0.5}
        position={[-15, 10, -20]}
      />

      {/* Hemisphere light for natural sky/ground lighting */}
      <hemisphereLight
        color="#b8c8ff"
        groundColor="#2a2a4a"
        intensity={0.6}
      />
    </>
  )
}
