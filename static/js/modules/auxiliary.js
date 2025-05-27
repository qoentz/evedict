/* ------------------------------------------------------------------
   /static/js/modules/auxiliary.js
   ------------------------------------------------------------------
   Renders an animated turquoise “network sphere” plus drifting particles.
   Call initBall() once on DOMContentLoaded (or after an HTMX swap) and
   it will initialise if the required <canvas id="network-sphere"> exists.
------------------------------------------------------------------ */

/* ------------------------------------------------------------------
   /static/js/modules/auxiliary.js
   ------------------------------------------------------------------ */

let sphere = null;        // singleton holder

export function initNetworkSphere() {
    const targetCanvas = document.getElementById('network-sphere');
    if (!targetCanvas) return;                       // fragment has no sphere

    /* 1 ▸ first time ever? build the globe ------------------------------ */
    if (!sphere) {
        sphere = buildSphere(targetCanvas);            // ← your original builder
        return;
    }

    /* 2 ▸ globe already exists ------------------------------------------ */
    // If the running canvas is already in the DOM, nothing to do.
    if (sphere.canvas.isConnected) return;

    // Otherwise replace the newly-inserted placeholder with the live canvas.
    targetCanvas.replaceWith(sphere.canvas);
    sphere.fit();                                    // update size/aspect
}

/* ====================================================================
   The original builder (slightly renamed) -----------------------------
   Returns an object with { canvas, fit } so the singleton can manage
   re-attachment and resizing.
==================================================================== */
function buildSphere(canvas) {
    /* ---- SETUP (your previous code) ---------------------------------- */
    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
    renderer.setPixelRatio(devicePixelRatio);
    renderer.physicallyCorrectLights = true;

    const scene  = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 100);
    camera.position.set(0, 1.5, 11);

    createLights(scene);
    const globe = createGlobe(scene);
    const dust  = createDust(scene);

    /* ---- fit helper --------------------------------------------------- */
    function fit() {
        renderer.setSize(canvas.clientWidth, canvas.clientHeight, false);
        camera.aspect = canvas.clientWidth / canvas.clientHeight;
        camera.updateProjectionMatrix();
    }
    window.addEventListener('resize', fit);
    fit();

    /* ---- animation loop (unchanged) ----------------------------------- */
    animate({ renderer, scene, camera, canvas, globe, dust });

    /* ---- expose canvas + fit so we can re-use them later -------------- */
    return { canvas, fit };
}

/* ======================================================================== */
/*  Helper creation functions                                               */
/* ======================================================================== */

function createLights(scene) {
    scene.add(new THREE.AmbientLight(0x0e1b27, 0.65));

    const orbLight = new THREE.PointLight(0xffffff, 1.7, 30, 2);
    scene.add(orbLight);

    // store on scene userData so animate() can find it
    scene.userData.orbLight = orbLight;
}

function createGlobe(scene) {
    const glowTex = makeSprite(
        'rgba(255,255,255,1)',
        'rgba(80,220,255,0.9)',
        'rgba(80,220,255,0)'
    );

    const group = new THREE.Group();
    group.position.y = -2.4;                         // lower the globe
    group.scale.set(1.25, 1.25, 1.25);
    scene.add(group);

    /* -- geometry ---------------------------------------------------------- */
    const RADIUS = 3, SUBDIV = 3;
    const baseGeo = new THREE.IcosahedronGeometry(RADIUS, SUBDIV);

    const edgeBank = [];                             // array of {a, b} vectors
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
    const flareArr = new Float32Array(baseGeo.attributes.position.count); // start 0

    // solid wireframe
    group.userData.wireMat = new THREE.LineBasicMaterial({ color: 0x42c7ff });
    group.add(new THREE.LineSegments(new THREE.WireframeGeometry(baseGeo), group.userData.wireMat));

    // additive plasma skin
    group.userData.plasmaMat = new THREE.LineBasicMaterial({
        color: 0x42c7ff,
        transparent: true,
        opacity: 0.18,
        blending: THREE.AdditiveBlending
    });
    group.add(new THREE.LineSegments(new THREE.WireframeGeometry(baseGeo), group.userData.plasmaMat));

    // nodes with per-vertex colours
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
            size: 0.28,                 // roughly 2× the core halo
            transparent: true,
            opacity: 0.25,
            color: 0xffffff,            // will be tinted each frame
            blending: THREE.AdditiveBlending,
            depthWrite: false
        })
    );
    group.add(halos);

    // translucent core for depth
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

    // stash animated bits for easy access
    group.userData.nodes       = nodes;
    group.userData.halos       = halos;
    group.userData.nodePhase   = phaseArr;
    group.userData.nodeColours = colourArr;
    group.userData.vCount      = vCount;

    group.userData.edgeBank  = edgeBank;
    group.userData.nodeFlare = flareArr;


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
        const r = 7 + Math.random() * 4;               // shell radius
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

/* ======================================================================== */
/*  Animation loop                                                          */
/* ======================================================================== */

function animate({ renderer, scene, camera, canvas, globe, dust }) {
    const edgeBank = globe.userData.edgeBank;
    const flare    = globe.userData.nodeFlare;
    const signals  = [];               // { mesh, a, b, t, speed }
    const maxSignals = 6;

    const clock      = new THREE.Clock();
    const nodeLow    = 0.25;
    const nodeHigh   = 1.4;
    const hueBase    = 0.55;   // ≈ 200° — blue-leaning turquoise
    const hueSpread  = 0.04;   // 185–215° swing
    const satBase    = 0.70;   // start a bit desaturated
    const satPulse   = 0.07;   // tiny saturation sway
    const lightBase  = 0.53;   // about the same as current
    const lightPulse = 0.06;   // subtler brightness breathing
    const hueSpeed   = 0.06;   // slower global colour cycle
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

        /* -- palette --------------------------------------------------------- */
        const hue  = hueBase + hueSpread * Math.sin(t * hueSpeed);
        const sat   = satBase  + satPulse  * Math.sin(t * 1.0);
        const light = lightBase + lightPulse * Math.sin(t * 0.8);
        const base  = new THREE.Color().setHSL(hue, sat, light);
        wireMat.color.copy(base);
        plasmaMat.color.copy(base);
        orbLight.color.copy(base);
        getContainer().style.backgroundColor = `hsl(${hue * 360}deg 38% 9%)`;

        /* -- orbiting light -------------------------------------------------- */
        const a = t * 0.45;
        orbLight.position.set(Math.sin(a) * 6, Math.cos(a * 0.7) * 5, Math.cos(a) * 6);

        /* -- rotations ------------------------------------------------------- */
        globe.rotation.set(-t * spinSpeed * 0.8, t * spinSpeed, 0);
        dust.rotation.y = t * 0.004;
        dust.rotation.x = t * 0.002;

        /* ------- SIGNAL DOTS ---------------------------------------------- */
        if (signals.length < maxSignals && Math.random() < 0.025) {
            const { a, b } = edgeBank[Math.random() * edgeBank.length | 0];
            const spr = makeSignalSprite(`#${base.getHexString()}`);
            spr.scale.setScalar(0.12);                 // dot size
            globe.add(spr);
            signals.push({ mesh: spr, a, b, t: 0, speed: 0.5 + Math.random() * 0.8 });
        }

        for (let i = signals.length - 1; i >= 0; i--) {
            const s = signals[i];
            s.t += s.speed * 0.016;                    // ~0.016 ≈ 60 fps
            if (s.t >= 1) {
                globe.remove(s.mesh);
                s.mesh.material.map.dispose();
                s.mesh.material.dispose();
                signals.splice(i, 1);
                continue;
            }
            s.mesh.position.lerpVectors(s.a, s.b, s.t);
            s.mesh.material.opacity = 1 - s.t;         // fade as it travels
        }


        /* -- node crescendo -------------------------------------------------- */
        for (let i = 0, j = 0; i < vCount; i++, j += 3) {
            // occasional spike
            if (Math.random() < 0.002) flare[i] = 1.0;
            flare[i] *= 0.93;                          // decay

            const k = nodeLow + (nodeHigh - nodeLow) * (0.5 + 0.5 * Math.sin(t * 0.6 + phase[i]))
                + flare[i];                      // add the flare
            colours[j]     = base.r * k;
            colours[j + 1] = base.g * k;
            colours[j + 2] = base.b * k;
        }
        nodes.geometry.attributes.color.needsUpdate = true;



        /* -- halo glow (size + opacity breathe together) -------------------- */
        halos.material.color.copy(base);                       // tint to scene hue
        halos.material.size =
            0.25 + 0.06 * (0.5 + 0.5 * Math.sin(t * 1.0));     // 0.25 – 0.31
        halos.material.opacity =
            0.20 + 0.10 * (0.5 + 0.5 * Math.sin(t * 1.0));     // 0.20 – 0.30
         // keep aligned

        /* -- plasma shimmer -------------------------------------------------- */
        plasmaMat.opacity = 0.17 + 0.10 * (0.5 + 0.5 * Math.sin(t * 1.3));

        renderer.render(scene, camera);
    })();
}
