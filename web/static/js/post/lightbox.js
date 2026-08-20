(function() {
    // Image Lightbox / Full Photo Viewer for Articles & Carousels
    var lightbox = null;
    var imgEl = null;
    var captionEl = null;
    var counterEl = null;
    var bottomCountEl = null;
    var bottomBarEl = null;
    var prevBtn = null;
    var nextBtn = null;

    var galleryItems = [];
    var currentIndex = 0;

    function createLightbox() {
        if (lightbox) return;

        lightbox = document.createElement("div");
        lightbox.className = "lightbox-overlay";
        lightbox.innerHTML = `
            <div class="lightbox-backdrop"></div>
            <div class="lightbox-toolbar">
                <div class="lightbox-caption"></div>
                <div class="lightbox-actions">
                    <span class="lightbox-counter"></span>
                    <button type="button" class="lightbox-close" aria-label="Close photo viewer">
                        <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>

            <button type="button" class="lightbox-nav-btn lightbox-prev" aria-label="Previous photo">
                <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
                </svg>
            </button>

            <div class="lightbox-content">
                <img class="lightbox-img" src="" alt="" />
            </div>

            <button type="button" class="lightbox-nav-btn lightbox-next" aria-label="Next photo">
                <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                </svg>
            </button>

            <div class="lightbox-bottom-bar">
                <button type="button" class="lightbox-bottom-btn lightbox-bottom-prev" aria-label="Previous slide">
                    <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
                    </svg>
                    <span>PREV</span>
                </button>
                <span class="lightbox-bottom-count"></span>
                <button type="button" class="lightbox-bottom-btn lightbox-bottom-next" aria-label="Next slide">
                    <span>NEXT</span>
                    <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                    </svg>
                </button>
            </div>
        `;
        document.body.appendChild(lightbox);

        imgEl = lightbox.querySelector(".lightbox-img");
        captionEl = lightbox.querySelector(".lightbox-caption");
        counterEl = lightbox.querySelector(".lightbox-counter");
        bottomCountEl = lightbox.querySelector(".lightbox-bottom-count");
        bottomBarEl = lightbox.querySelector(".lightbox-bottom-bar");

        prevBtn = lightbox.querySelector(".lightbox-prev");
        nextBtn = lightbox.querySelector(".lightbox-next");
        var bPrevBtn = lightbox.querySelector(".lightbox-bottom-prev");
        var bNextBtn = lightbox.querySelector(".lightbox-bottom-next");

        // Click event listeners
        lightbox.querySelector(".lightbox-close").addEventListener("click", closeLightbox);
        lightbox.querySelector(".lightbox-backdrop").addEventListener("click", closeLightbox);

        prevBtn.addEventListener("click", prevPhoto);
        nextBtn.addEventListener("click", nextPhoto);
        bPrevBtn.addEventListener("click", prevPhoto);
        bNextBtn.addEventListener("click", nextPhoto);

        lightbox.querySelector(".lightbox-content").addEventListener("click", function(e) {
            if (e.target === this) {
                closeLightbox();
            }
        });

        // Touch swipe support on mobile
        var touchStartX = 0;
        var touchStartY = 0;
        var contentEl = lightbox.querySelector(".lightbox-content");

        contentEl.addEventListener("touchstart", function(e) {
            touchStartX = e.touches[0].clientX;
            touchStartY = e.touches[0].clientY;
        }, { passive: true });

        contentEl.addEventListener("touchend", function(e) {
            if (!touchStartX || !touchStartY) return;
            var diffX = e.changedTouches[0].clientX - touchStartX;
            var diffY = e.changedTouches[0].clientY - touchStartY;

            if (Math.abs(diffX) > 40 && Math.abs(diffX) > Math.abs(diffY)) {
                if (diffX > 0) {
                    prevPhoto();
                } else {
                    nextPhoto();
                }
            } else if (Math.abs(diffY) > 80 && Math.abs(diffY) > Math.abs(diffX)) {
                closeLightbox();
            }

            touchStartX = 0;
            touchStartY = 0;
        }, { passive: true });

        // Keyboard controls
        document.addEventListener("keydown", function(e) {
            if (!lightbox || !lightbox.classList.contains("open")) return;
            if (e.key === "Escape") {
                closeLightbox();
            } else if (e.key === "ArrowLeft" || e.key === "a" || e.key === "A") {
                prevPhoto();
            } else if (e.key === "ArrowRight" || e.key === "d" || e.key === "D") {
                nextPhoto();
            }
        });
    }

    function renderActivePhoto() {
        if (!galleryItems.length || currentIndex < 0 || currentIndex >= galleryItems.length) return;

        var item = galleryItems[currentIndex];
        imgEl.style.opacity = "0.4";
        imgEl.style.transform = "scale(0.96)";

        setTimeout(function() {
            imgEl.src = item.src;
            imgEl.alt = item.alt || "";

            if (item.caption) {
                captionEl.textContent = item.caption;
                captionEl.style.display = "block";
            } else {
                captionEl.textContent = "";
                captionEl.style.display = "none";
            }

            var showNav = galleryItems.length > 1;
            var countStr = showNav ? ((currentIndex + 1) + " / " + galleryItems.length) : "";
            counterEl.textContent = countStr;
            counterEl.style.display = showNav ? "block" : "none";
            bottomCountEl.textContent = countStr;

            prevBtn.style.display = showNav ? "flex" : "none";
            nextBtn.style.display = showNav ? "flex" : "none";
            bottomBarEl.style.display = showNav ? "flex" : "none";

            imgEl.style.opacity = "1";
            imgEl.style.transform = "scale(1)";
        }, 80);
    }

    function nextPhoto() {
        if (galleryItems.length <= 1) return;
        currentIndex = (currentIndex + 1) % galleryItems.length;
        renderActivePhoto();
    }

    function prevPhoto() {
        if (galleryItems.length <= 1) return;
        currentIndex = (currentIndex - 1 + galleryItems.length) % galleryItems.length;
        renderActivePhoto();
    }

    function openGallery(items, startIdx) {
        createLightbox();
        galleryItems = items;
        currentIndex = startIdx || 0;

        renderActivePhoto();
        lightbox.classList.add("open");
        document.body.style.overflow = "hidden";
    }

    function closeLightbox() {
        if (!lightbox) return;
        lightbox.classList.remove("open");
        document.body.style.overflow = "";
        setTimeout(function() {
            if (imgEl && !lightbox.classList.contains("open")) {
                imgEl.src = "";
            }
        }, 250);
    }

    function initLightbox() {
        // Collect all images in articles, covers, galleries, and carousels
        var containerSelectors = [
            ".post-carousel-wrap",
            ".post-gallery-wrap",
            "figure",
            ".post-cover",
            "#prose-body"
        ];

        // Also bind standalone cover & content images
        var images = document.querySelectorAll("#prose-body img, .post-cover img, figure img, .post-gallery-wrap img, .post-carousel-wrap img, article header img");

        images.forEach(function(img) {
            if (img.dataset.lightboxBound) return;
            img.dataset.lightboxBound = "true";
            img.classList.add("lightbox-trigger");

            img.addEventListener("click", function(e) {
                var parentLink = img.closest("a");
                if (parentLink && parentLink.href && !parentLink.href.match(/\.(jpg|jpeg|png|webp|gif|svg)(\?.*)?$/i)) {
                    return; // Standard web links navigate normally
                }

                e.preventDefault();

                // Determine gallery scope: ONLY carousel or gallery containers
                var carouselScope = img.closest(".post-carousel-wrap") || img.closest(".post-gallery-wrap");
                var scopeImages = [];

                if (carouselScope) {
                    scopeImages = Array.from(carouselScope.querySelectorAll("img")).filter(function(i) {
                        return i.src && !i.closest(".post-author-card") && !i.closest(".no-lightbox");
                    });
                }

                if (!scopeImages.length) {
                    scopeImages = [img];
                }

                var items = scopeImages.map(function(i) {
                    var cap = "";
                    var fig = i.closest("figure");
                    if (fig) {
                        var fc = fig.querySelector("figcaption");
                        if (fc) cap = fc.innerText;
                    }
                    if (!cap && i.alt) cap = i.alt;
                    if (!cap && i.title) cap = i.title;

                    var pLink = i.closest("a");
                    var src = (pLink && pLink.href && pLink.href.match(/\.(jpg|jpeg|png|webp|gif|svg)(\?.*)?$/i)) ? pLink.href : i.src;
                    return { src: src, alt: i.alt || "", caption: cap };
                });

                var activeIndex = scopeImages.indexOf(img);
                if (activeIndex < 0) activeIndex = 0;

                openGallery(items, activeIndex);
            });
        });
    }

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", initLightbox);
    } else {
        initLightbox();
    }

    document.addEventListener("htmx:afterSwap", initLightbox);
})();
