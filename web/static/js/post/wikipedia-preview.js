document.addEventListener('DOMContentLoaded', () => {
    const wikiLinks = document.querySelectorAll('.prose a[href*="wikipedia.org/wiki/"]');
    if (wikiLinks.length === 0) return;

    const wikiCache = new Map();

    async function getWikiSummary(lang, title) {
        const safeTitle = encodeURIComponent(decodeURIComponent(title));
        const key = `${lang}:${safeTitle}`;
        if (wikiCache.has(key)) return wikiCache.get(key);

        const apiUrl = `https://${lang}.wikipedia.org/api/rest_v1/page/summary/${safeTitle}?redirect=true`;

        const res = await fetch(apiUrl);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        wikiCache.set(key, data);
        return data;
    }

    function parseWikiUrl(href) {
        try {
            const url = new URL(href);
            const langMatch = url.hostname.match(/^([a-z\-]+)\.wikipedia\.org/);
            const lang = langMatch ? langMatch[1] : 'en';
            const parts = url.pathname.split('/wiki/');
            if (parts.length < 2 || !parts[1]) return null;
            return { lang, title: parts[1], href: url.href };
        } catch {
            return null;
        }
    }

    function isMobileOrTouch() {
        return window.innerWidth < 768 || window.matchMedia('(hover: none)').matches || window.matchMedia('(pointer: coarse)').matches;
    }

    const desktopTooltip = document.createElement('div');
    desktopTooltip.className = 'wiki-tooltip';
    document.body.appendChild(desktopTooltip);

    let hoverTimeout;
    let currentDesktopLink = null;

    function showDesktopLoading(targetEl) {
        positionDesktopTooltip(targetEl);
        desktopTooltip.innerHTML = '<div class="wiki-loading">Memuat dari Wikipedia...</div>';
        desktopTooltip.classList.add('visible');
    }

    function hideDesktopTooltip() {
        desktopTooltip.classList.remove('visible');
    }

    function positionDesktopTooltip(targetEl) {
        const rect = targetEl.getBoundingClientRect();
        const tooltipWidth = 320;

        let top = rect.bottom + window.scrollY + 8;
        let left = rect.left + window.scrollX - (tooltipWidth / 2) + (rect.width / 2);

        if (left < 16) left = 16;
        if (left + tooltipWidth > window.innerWidth - 16) {
            left = window.innerWidth - tooltipWidth - 16;
        }

        desktopTooltip.style.top = `${top}px`;
        desktopTooltip.style.left = `${left}px`;
    }

    function renderDesktopTooltip(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<img src="${data.thumbnail.source}" class="wiki-thumb" alt="">`;
        }
        html += `<div class="wiki-content">`;
        html += `<div class="wiki-title">${data.title}</div>`;

        let extract = data.extract || '';
        if (extract.length > 200) extract = extract.substring(0, 200) + '...';

        html += `<p class="wiki-extract">${extract}</p>`;
        html += `<div class="wiki-footer">W — Disediakan oleh Wikipedia</div>`;
        html += `</div>`;

        desktopTooltip.innerHTML = html;
    }

    desktopTooltip.addEventListener('mouseenter', () => clearTimeout(hoverTimeout));
    desktopTooltip.addEventListener('mouseleave', () => {
        hoverTimeout = setTimeout(hideDesktopTooltip, 300);
    });

    const sheetModal = document.createElement('div');
    sheetModal.id = 'wiki-bottom-sheet';
    sheetModal.className = 'wiki-sheet-modal';
    sheetModal.setAttribute('role', 'dialog');
    sheetModal.setAttribute('aria-modal', 'true');
    sheetModal.setAttribute('aria-hidden', 'true');
    sheetModal.setAttribute('aria-label', 'Wikipedia Preview');
    sheetModal.innerHTML = `
        <div class="wiki-sheet-backdrop"></div>
        <div class="wiki-sheet-container">
            <div class="wiki-sheet-drag-area">
                <div class="wiki-sheet-handle"></div>
            </div>
            <div class="wiki-sheet-header">
                <div class="wiki-sheet-badge">
                    <svg class="wiki-sheet-logo" viewBox="0 0 24 24" width="15" height="15" fill="currentColor">
                        <path d="M12.09 13.124l2.77-7.858h2.02l-4.148 10.978h-1.606l-2.614-7.398-2.615 7.398H4.29L.142 5.266h2.02l2.77 7.858 2.383-6.732h1.666l2.383 6.732 2.726-7.724z"/>
                    </svg>
                    <span>Wikipedia</span>
                </div>
                <button type="button" class="wiki-sheet-close" aria-label="Tutup pratinjau">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>
            <div class="wiki-sheet-body"></div>
            <div class="wiki-sheet-footer">
                <a href="#" target="_blank" rel="noopener noreferrer" class="wiki-sheet-btn">
                    <span>Buka di Wikipedia</span>
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                        <polyline points="15 3 21 3 21 9"></polyline>
                        <line x1="10" y1="14" x2="21" y2="3"></line>
                    </svg>
                </a>
            </div>
        </div>
    `;
    document.body.appendChild(sheetModal);

    const sheetBackdrop = sheetModal.querySelector('.wiki-sheet-backdrop');
    const sheetContainer = sheetModal.querySelector('.wiki-sheet-container');
    const sheetDragArea = sheetModal.querySelector('.wiki-sheet-drag-area');
    const sheetCloseBtn = sheetModal.querySelector('.wiki-sheet-close');
    const sheetBody = sheetModal.querySelector('.wiki-sheet-body');
    const sheetExtBtn = sheetModal.querySelector('.wiki-sheet-btn');

    function openBottomSheet(lang, title, href) {
        hideDesktopTooltip();
        document.body.classList.add('wiki-sheet-lock');
        sheetExtBtn.href = href;
        sheetContainer.style.transform = '';
        sheetBody.innerHTML = `
            <div class="wiki-sheet-loading">
                <div class="wiki-sheet-spinner"></div>
                <span>Memuat ringkasan Wikipedia...</span>
            </div>
        `;
        sheetModal.classList.add('open');
        sheetModal.setAttribute('aria-hidden', 'false');

        getWikiSummary(lang, title).then(data => {
            if (!sheetModal.classList.contains('open')) return;
            renderBottomSheet(data);
        }).catch(err => {
            if (!sheetModal.classList.contains('open')) return;
            sheetBody.innerHTML = `
                <div class="wiki-sheet-error">
                    <p>Gagal memuat pratinjau Wikipedia.</p>
                </div>
            `;
        });
    }

    function renderBottomSheet(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<div class="wiki-sheet-thumb-wrap"><img src="${data.thumbnail.source}" class="wiki-sheet-thumb" alt="${data.title || ''}"></div>`;
        }
        html += `<div class="wiki-sheet-text">`;
        html += `<h3 class="wiki-sheet-title">${data.title}</h3>`;
        if (data.description) {
            html += `<div class="wiki-sheet-meta">${data.description}</div>`;
        }
        let extract = data.extract || '';
        if (extract.length > 280) extract = extract.substring(0, 280) + '...';
        html += `<p class="wiki-sheet-extract">${extract}</p>`;
        html += `</div>`;
        sheetBody.innerHTML = html;
    }

    function closeBottomSheet() {
        sheetModal.classList.remove('open');
        sheetModal.setAttribute('aria-hidden', 'true');
        document.body.classList.remove('wiki-sheet-lock');
        sheetContainer.style.transform = '';
    }

    sheetBackdrop.addEventListener('click', closeBottomSheet);
    sheetCloseBtn.addEventListener('click', closeBottomSheet);

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && sheetModal.classList.contains('open')) {
            closeBottomSheet();
        }
    });

    let startY = 0;
    let isDragging = false;
    let currentDelta = 0;

    sheetDragArea.addEventListener('touchstart', (e) => {
        startY = e.touches[0].clientY;
        isDragging = true;
        currentDelta = 0;
        sheetContainer.style.transition = 'none';
    }, { passive: true });

    sheetDragArea.addEventListener('touchmove', (e) => {
        if (!isDragging) return;
        const touchY = e.touches[0].clientY;
        currentDelta = touchY - startY;
        if (currentDelta > 0) {
            sheetContainer.style.transform = `translateY(${currentDelta}px)`;
        }
    }, { passive: true });

    sheetDragArea.addEventListener('touchend', () => {
        if (!isDragging) return;
        isDragging = false;
        sheetContainer.style.transition = '';
        if (currentDelta > 70) {
            closeBottomSheet();
        } else {
            sheetContainer.style.transform = '';
        }
        startY = 0;
        currentDelta = 0;
    });

    wikiLinks.forEach(link => {
        const parsed = parseWikiUrl(link.href);
        if (!parsed) return;

        link.addEventListener('mouseenter', () => {
            if (isMobileOrTouch()) return;
            clearTimeout(hoverTimeout);
            currentDesktopLink = link;

            hoverTimeout = setTimeout(async () => {
                showDesktopLoading(link);
                try {
                    const data = await getWikiSummary(parsed.lang, parsed.title);
                    if (currentDesktopLink === link && desktopTooltip.classList.contains('visible')) {
                        renderDesktopTooltip(data);
                        positionDesktopTooltip(link);
                    }
                } catch (e) {
                    if (currentDesktopLink === link) {
                        desktopTooltip.innerHTML = `<div class="wiki-error">Gagal: ${e.message}</div>`;
                    }
                }
            }, 300);
        });

        link.addEventListener('mouseleave', () => {
            if (isMobileOrTouch()) return;
            clearTimeout(hoverTimeout);
            if (currentDesktopLink === link) {
                hoverTimeout = setTimeout(hideDesktopTooltip, 300);
            }
        });

        link.addEventListener('click', (e) => {
            if (isMobileOrTouch()) {
                e.preventDefault();
                openBottomSheet(parsed.lang, parsed.title, link.href);
            }
        });
    });
});
