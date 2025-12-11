let currentImageList = [];
let currentImageIndex = -1;
let currentZoom = 1.0;
let isDragging = false;
let startX = 0;
let startY = 0;
let translateX = 0;
let translateY = 0;

// Initialize
window.onload = () => {
    // Add wheel event for zooming
    const container = document.getElementById('image-container');
    container.addEventListener('wheel', handleWheel);

    // Dragging
    const img = document.getElementById('main-image');
    img.addEventListener('mousedown', startDrag);
    window.addEventListener('mousemove', drag);
    window.addEventListener('mouseup', endDrag);
};

// Functions exposed to HTML
async function selectImage() {
    try {
        const path = await window.go.main.App.OpenImageDialog();
        if (path) {
            await displayImage(path);

            // Get directory
            // Assuming Windows backslashes for now as user is on Windows
            let dir = path.substring(0, path.lastIndexOf('\\'));
            // If empty, might be root or linux style (unlikely given metadata)
            if (!dir && path.includes('/')) {
                dir = path.substring(0, path.lastIndexOf('/'));
            }

            await loadGallery(dir, path);
        }
    } catch (e) {
        console.error("Error selecting image:", e);
    }
}

async function displayImage(path) {
    if (!path) return;
    try {
        const base64 = await window.go.main.App.LoadImage(path);
        const img = document.getElementById('main-image');
        const placeholder = document.getElementById('placeholder');

        img.src = base64;
        img.style.display = 'block';
        placeholder.style.display = 'none';

        resetTransform();
    } catch (e) {
        console.error("Error loading image:", e);
    }
}

async function loadGallery(dir, currentPath) {
    try {
        const images = await window.go.main.App.GetImagesInFolder(dir);
        currentImageList = images;
        currentImageIndex = currentImageList.indexOf(currentPath);
        renderGallery();
    } catch (e) {
        console.error("Error loading gallery:", e);
    }
}

function renderGallery() {
    const gallery = document.getElementById('gallery-strip');
    gallery.innerHTML = '';

    currentImageList.forEach((path, index) => {
        const div = document.createElement('div');
        div.className = 'thumbnail ' + (index === currentImageIndex ? 'active' : '');
        div.onclick = () => jumpToImage(index);

        // We can lazy load thumbnails or just set src to full image (heavy)
        // For efficiency, we really should have a thumbnail generator, but for now we'll rely on browser caching.
        // To avoid freezing, we might want to just show names or use the LoadImage method asynchronously.
        // Using full base64 for all thumbnails immediately might be slow.
        // Let's just create the div first, and maybe load images lazily?
        // Or simpler: Just rely on on-demand loading for the main view, and minimal styling for gallery items.
        // The user wanted a gallery "strip". 
        // Let's try to load them.

        // Using a promise to load thumbnail separately to not block UI
        window.go.main.App.LoadImage(path).then(base64 => {
            const img = document.createElement('img');
            img.src = base64;
            div.appendChild(img);
        });

        gallery.appendChild(div);
    });
}

function jumpToImage(index) {
    if (index >= 0 && index < currentImageList.length) {
        currentImageIndex = index;
        const path = currentImageList[index];
        displayImage(path);
        updateActiveThumbnail();
    }
}

function nextImage() {
    if (currentImageList.length > 0) {
        let newIndex = currentImageIndex + 1;
        if (newIndex >= currentImageList.length) newIndex = 0;
        jumpToImage(newIndex);
    }
}

function prevImage() {
    if (currentImageList.length > 0) {
        let newIndex = currentImageIndex - 1;
        if (newIndex < 0) newIndex = currentImageList.length - 1;
        jumpToImage(newIndex);
    }
}

function updateActiveThumbnail() {
    const thumbnails = document.querySelectorAll('.thumbnail');
    thumbnails.forEach((t, i) => {
        if (i === currentImageIndex) t.classList.add('active');
        else t.classList.remove('active');
    });
    // Scroll thumbnail into view
    thumbnails[currentImageIndex]?.scrollIntoView({ behavior: 'smooth', inline: 'center' });
}

// Zoom / Pan Logic
function resetTransform() {
    currentZoom = 1.0;
    translateX = 0;
    translateY = 0;
    updateTransform();
}

function updateTransform() {
    const img = document.getElementById('main-image');
    img.style.transform = `translate(${translateX}px, ${translateY}px) scale(${currentZoom})`;
    document.getElementById('zoom-level').innerText = Math.round(currentZoom * 100) + '%';
}

function zoomIn() {
    currentZoom *= 1.1;
    updateTransform();
}

function zoomOut() {
    currentZoom /= 1.1;
    if (currentZoom < 0.1) currentZoom = 0.1;
    updateTransform();
}

function handleWheel(e) {
    e.preventDefault();
    const delta = Math.sign(e.deltaY) * -1;
    if (delta > 0) zoomIn();
    else zoomOut();
}

// Drag logic
function startDrag(e) {
    if (e.button !== 0) return; // Only left click
    isDragging = true;
    startX = e.clientX - translateX;
    startY = e.clientY - translateY;
    e.preventDefault();
}

function drag(e) {
    if (!isDragging) return;
    translateX = e.clientX - startX;
    translateY = e.clientY - startY;
    updateTransform();
}

function endDrag() {
    isDragging = false;
}
