/* ------------------------------------------------------------------
   /static/js/modules/auxiliary.js
   ------------------------------------------------------------------
   Unified Three.js scene that can render with or without globe
------------------------------------------------------------------ */

let scene = null;  // singleton holder for the unified scene

// Initialize scene with globe
export function initNetworkSphere() {
    const targetCanvas = document.getElementById('network-sphere');
    if (!targetCanvas) return;

    if (!scene) {
        scene = buildScene(targetCanvas, { showGlobe: true });
        return;
    }

    if (scene.canvas.isConnected) return;
    targetCanvas.replaceWith(scene.canvas);
    scene.setGlobeVisibility(true);
    scene.fit();
}

// Initialize scene without globe (ambient background)
export function initAmbientBackground() {
    const targetCanvas = document.getElementById('ambient-background');
    if (!targetCanvas) return;

    if (!scene) {
        scene = buildScene(targetCanvas, { showGlobe: false });
        return;
    }

    if (scene.canvas.isConnected) return;
    targetCanvas.replaceWith(scene.canvas);
    scene.setGlobeVisibility(false);
    scene.fit();
}

/* ====================================================================
   Unified Scene Builder
==================================================================== */
function buildScene(canvas, options = {}) {
    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
    renderer.setPixelRatio(devicePixelRatio);
    renderer.physicallyCorrectLights = true;

    const threeScene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 100);
    camera.position.set(0, 1.5, 11);

    createLights(threeScene);
    const globe = createGlobe(threeScene);
    const dust = createDust(threeScene);

    // Set initial globe visibility
    globe.visible = options.showGlobe !== false;

    function fit() {
        renderer.setSize(canvas.clientWidth, canvas.clientHeight, false);
        camera.aspect = canvas.clientWidth / canvas.clientHeight;
        camera.updateProjectionMatrix();
    }

    function setGlobeVisibility(visible) {
        globe.visible = visible;
        // Adjust camera position based on mode
        if (visible) {
            camera.position.set(0, 1.5, 11);  // Globe view
        } else {
            camera.position.set(0, 0, 8);     // Ambient view (closer for better particle view)
        }
    }

    window.addEventListener('resize', fit);
    fit();

    animate({ renderer, scene: threeScene, camera, canvas, globe, dust });

    return {
        canvas,
        fit,
        setGlobeVisibility,
        globe,
        dust
    };
}

/* ====================================================================
   Scene Components (unchanged)
==================================================================== */

function createLights(scene) {
    scene.add(new THREE.AmbientLight(0x0e1b27, 0.65));

    const orbLight = new THREE.PointLight(0xffffff, 1.7, 30, 2);
    scene.add(orbLight);

    scene.userData.orbLight = orbLight;
}

function createGlobe(scene) {
    const glowTex = makeSprite(
        'rgba(255,255,255,1)',
        'rgba(80,220,255,0.9)',
        'rgba(80,220,255,0)'
    );

    const group = new THREE.Group();
    group.position.y = -2.4;
    group.scale.set(1.25, 1.25, 1.25);
    scene.add(group);

    const RADIUS = 3, SUBDIV = 3;
    const baseGeo = new THREE.IcosahedronGeometry(RADIUS, SUBDIV);

    const edgeBank = [];
    {
        const edges = new THREE.EdgesGeometry(baseGeo);
        const p = edges.attributes.position.array;
        for (let i = 0; i < p.length; i += 6) {
            edgeBank.push({
                a: new THREE.Vector3(p[i],   p[i+1], p[i+2]),
                b: new THREE.Vector3(p[i+3], p[i+4], p[i+5])
            });
        }
    }
    const flareArr = new Float32Array(baseGeo.attributes.position.count);

    group.userData.wireMat = new THREE.LineBasicMaterial({ color: 0x42c7ff });
    group.add(new THREE.LineSegments(new THREE.WireframeGeometry(baseGeo), group.userData.wireMat));

    group.userData.plasmaMat = new THREE.LineBasicMaterial({
        color: 0x42c7ff,
        transparent: true,
        opacity: 0.18,
        blending: THREE.AdditiveBlending
    });
    group.add(new THREE.LineSegments(new THREE.WireframeGeometry(baseGeo), group.userData.plasmaMat));

    const vCount   = baseGeo.attributes.position.count;
    const phaseArr = Float32Array.from({ length: vCount }, () => Math.random() * Math.PI * 2);
    const colourArr = new Float32Array(vCount * 3);
    baseGeo.setAttribute('color', new THREE.BufferAttribute(colourArr, 3));

    const nodes = new THREE.Points(
        baseGeo.clone(),
        new THREE.PointsMaterial({
            map: glowTex,
            size: 0.13,
            transparent: true,
            vertexColors: true,
            blending: THREE.AdditiveBlending,
            depthWrite: false
        })
    );
    group.add(nodes);

    const halos = new THREE.Points(
        baseGeo.clone(),
        new THREE.PointsMaterial({
            map: glowTex,
            size: 0.28,
            transparent: true,
            opacity: 0.25,
            color: 0xffffff,
            blending: THREE.AdditiveBlending,
            depthWrite: false
        })
    );
    group.add(halos);

    group.add(
        new THREE.Mesh(
            new THREE.SphereGeometry(RADIUS - 0.1, 32, 32),
            new THREE.MeshStandardMaterial({
                color: 0x0e2433,
                metalness: 0.25,
                roughness: 0.55,
                transparent: true,
                opacity: 0.38
            })
        )
    );

    group.userData.nodes       = nodes;
    group.userData.halos       = halos;
    group.userData.nodePhase   = phaseArr;
    group.userData.nodeColours = colourArr;
    group.userData.vCount      = vCount;
    group.userData.edgeBank    = edgeBank;
    group.userData.nodeFlare   = flareArr;

    return group;
}

function makeSignalSprite(colour) {
    const S = 64, cvs = Object.assign(document.createElement('canvas'), { width: S, height: S });
    const ctx = cvs.getContext('2d');
    const g   = ctx.createRadialGradient(S/2, S/2, 0, S/2, S/2, S/2);
    g.addColorStop(0,   'rgba(255,255,255,1)');
    g.addColorStop(0.2, colour);
    g.addColorStop(1,   'rgba(255,255,255,0)');
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, S, S);
    const tex = new THREE.CanvasTexture(cvs);
    tex.minFilter = THREE.LinearFilter;
    return new THREE.Sprite(
        new THREE.SpriteMaterial({ map: tex, transparent: true, depthWrite: false, opacity: 1 })
    );
}

function createDust(scene) {
    const dustTex = makeSprite(
        'rgba(255,255,255,0.85)',
        'rgba(150,220,255,0.35)',
        'rgba(150,220,255,0)'
    );

    const DUST_POINTS = 2700;
    const positions   = new Float32Array(DUST_POINTS * 3);
    for (let i = 0; i < DUST_POINTS; i++) {
        const r = 7 + Math.random() * 4;
        const u = Math.random();
        const v = Math.random();
        const theta = Math.acos(2 * u - 1);
        const phi   = 2 * Math.PI * v;
        positions[i*3]     = r * Math.sin(theta) * Math.cos(phi);
        positions[i*3 + 1] = r * Math.sin(theta) * Math.sin(phi);
        positions[i*3 + 2] = r * Math.cos(theta);
    }

    const dustGeo = new THREE.BufferGeometry();
    dustGeo.setAttribute('position', new THREE.BufferAttribute(positions, 3));

    const dust = new THREE.Points(
        dustGeo,
        new THREE.PointsMaterial({
            map: dustTex,
            size: 0.07,
            transparent: true,
            opacity: 0.8,
            blending: THREE.AdditiveBlending,
            depthWrite: false
        })
    );
    dust.rotation.order = 'YXZ';
    scene.add(dust);

    return dust;
}

function makeSprite(inner, mid, outer) {
    const S   = 64;
    const cvs = Object.assign(document.createElement('canvas'), { width: S, height: S });
    const ctx = cvs.getContext('2d');
    const g   = ctx.createRadialGradient(S / 2, S / 2, 0, S / 2, S / 2, S / 2);
    g.addColorStop(0.0, inner);
    g.addColorStop(0.2, mid);
    g.addColorStop(1.0, outer);
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, S, S);
    const tex = new THREE.CanvasTexture(cvs);
    tex.minFilter = THREE.LinearFilter;
    return tex;
}

/* ====================================================================
   Unified Animation (handles both modes)
==================================================================== */
function animate({ renderer, scene, camera, canvas, globe, dust }) {
    const edgeBank = globe.userData.edgeBank;
    const flare    = globe.userData.nodeFlare;
    const signals  = [];
    const maxSignals = 6;

    const clock      = new THREE.Clock();
    const nodeLow    = 0.25;
    const nodeHigh   = 1.4;

    // Original network scene colors
    const hueBase    = 0.57;  // Shifted away from turquoise
    const hueSpread  = 0.025; // Reduced range to avoid turquoise tones
    const satBase    = 0.70;
    const satPulse   = 0.07;
    const lightBase  = 0.53;
    const lightPulse = 0.06;

    // Fixed blue for ambient background - no color shifting
    const ambientHueBase    = 0.57;  // Fixed to the good blue hue
    const ambientHueSpread  = 0;     // NO hue variation - completely consistent
    const ambientSatBase    = 0.68;  // Fixed saturation level
    const ambientSatPulse   = 0;     // NO saturation variation
    const ambientLightBase  = 0.50;  // Fixed lightness level
    const ambientLightPulse = 0;     // NO lightness variation

    const hueSpeed   = 0.06;
    const spinSpeed  = 0.018;

    const getContainer = () => canvas.parentElement || document.body;

    const orbLight   = scene.userData.orbLight;
    const wireMat    = globe.userData.wireMat;
    const plasmaMat  = globe.userData.plasmaMat;
    const nodes      = globe.userData.nodes;
    const halos      = globe.userData.halos;
    const phase      = globe.userData.nodePhase;
    const colours    = globe.userData.nodeColours;
    const vCount     = globe.userData.vCount;

    (function loop() {
        requestAnimationFrame(loop);
        const t = clock.getElapsedTime();

        // Choose color scheme based on globe visibility
        let hue, sat, light, base;

        if (globe.visible) {
            // Original network scene colors
            hue = hueBase + hueSpread * Math.sin(t * hueSpeed);
            sat = satBase + satPulse * Math.sin(t * 1.0);
            light = lightBase + lightPulse * Math.sin(t * 0.8);
            base = new THREE.Color().setHSL(hue, sat, light);

            // Original background color for network scene
            getContainer().style.backgroundColor = `hsl(${hue * 360}deg 38% 9%)`;
        } else {
            // Light blue ambient background colors
            hue = ambientHueBase + ambientHueSpread * Math.sin(t * hueSpeed);
            sat = ambientSatBase + ambientSatPulse * Math.sin(t * 1.0);
            light = ambientLightBase + ambientLightPulse * Math.sin(t * 0.8);
            base = new THREE.Color().setHSL(hue, sat, light);

            // NetworkSphere background style for ambient scene
            getContainer().style.backgroundColor = `hsl(${hue * 360}deg 38% 9%)`;
        }

        // Adjust orbLight intensity based on mode
        if (globe.visible) {
            orbLight.color.copy(base);
        } else {
            // Dimmer light for ambient mode to reduce dust flickering
            const dimmedBase = base.clone().multiplyScalar(0.4);
            orbLight.color.copy(dimmedBase);
        }

        const a = t * 0.45;
        orbLight.position.set(Math.sin(a) * 6, Math.cos(a * 0.7) * 5, Math.cos(a) * 6);

        // Always animate dust
        dust.rotation.y = t * 0.004;
        dust.rotation.x = t * 0.002;

        // Only animate globe if it's visible
        if (globe.visible) {
            wireMat.color.copy(base);
            plasmaMat.color.copy(base);

            globe.rotation.set(-t * spinSpeed * 0.8, t * spinSpeed, 0);

            // Handle signals
            if (signals.length < maxSignals && Math.random() < 0.025) {
                const { a, b } = edgeBank[Math.random() * edgeBank.length | 0];
                const spr = makeSignalSprite(`#${base.getHexString()}`);
                spr.scale.setScalar(0.12);
                globe.add(spr);
                signals.push({ mesh: spr, a, b, t: 0, speed: 0.5 + Math.random() * 0.8 });
            }

            for (let i = signals.length - 1; i >= 0; i--) {
                const s = signals[i];
                s.t += s.speed * 0.016;
                if (s.t >= 1) {
                    globe.remove(s.mesh);
                    s.mesh.material.map.dispose();
                    s.mesh.material.dispose();
                    signals.splice(i, 1);
                    continue;
                }
                s.mesh.position.lerpVectors(s.a, s.b, s.t);
                s.mesh.material.opacity = 1 - s.t;
            }

            // Animate node colors
            for (let i = 0, j = 0; i < vCount; i++, j += 3) {
                if (Math.random() < 0.002) flare[i] = 1.0;
                flare[i] *= 0.93;

                const k = nodeLow + (nodeHigh - nodeLow) * (0.5 + 0.5 * Math.sin(t * 0.6 + phase[i]))
                    + flare[i];
                colours[j]     = base.r * k;
                colours[j + 1] = base.g * k;
                colours[j + 2] = base.b * k;
            }
            nodes.geometry.attributes.color.needsUpdate = true;

            halos.material.color.copy(base);
            halos.material.size = 0.25 + 0.06 * (0.5 + 0.5 * Math.sin(t * 1.0));
            halos.material.opacity = 0.20 + 0.10 * (0.5 + 0.5 * Math.sin(t * 1.0));

            plasmaMat.opacity = 0.17 + 0.10 * (0.5 + 0.5 * Math.sin(t * 1.3));
        }

        renderer.render(scene, camera);
    })();
}