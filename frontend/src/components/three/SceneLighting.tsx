export default function SceneLighting() {
  return (
    <>
      <ambientLight color="#1a1a3a" intensity={0.4} />
      <directionalLight
        color="#aabbff"
        intensity={1.2}
        position={[20, 30, 10]}
        castShadow
        shadow-mapSize-width={1024}
        shadow-mapSize-height={1024}
        shadow-camera-left={-50}
        shadow-camera-right={50}
        shadow-camera-top={50}
        shadow-camera-bottom={-50}
      />
      <directionalLight
        color="#4466aa"
        intensity={0.3}
        position={[-15, 10, -20]}
      />
    </>
  )
}
