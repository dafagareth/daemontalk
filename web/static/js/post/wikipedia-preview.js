document.addEventListener('DOMContentLoaded', () => {
    const wikiLinks = document.querySelectorAll('.prose a[href*="wikipedia.org/wiki/"]');
    if (wikiLinks.length === 0) return;

    const tooltip = document.createElement('div');
    tooltip.className = 'wiki-tooltip';
    document.body.appendChild(tooltip);

    let hoverTimeout;
    let currentLink = null;

    wikiLinks.forEach(link => {
        link.addEventListener('mouseenter', () => {
            clearTimeout(hoverTimeout);
            currentLink = link;

            const url = new URL(link.href);
            const langMatch = url.hostname.match(/^([a-z\-]+)\.wikipedia\.org/);
            const lang = langMatch ? langMatch[1] : 'en';
            const title = url.pathname.split('/wiki/')[1];
            
            if (!title) return;

            hoverTimeout = setTimeout(() => {
                showLoading(link);
                fetchWikiData(lang, title, link);
            }, 300);
        });

        link.addEventListener('mouseleave', () => {
            clearTimeout(hoverTimeout);
            if (currentLink === link) {
                hoverTimeout = setTimeout(hideTooltip, 300);
            }
        });
    });

    tooltip.addEventListener('mouseenter', () => clearTimeout(hoverTimeout));
    tooltip.addEventListener('mouseleave', () => {
        hoverTimeout = setTimeout(hideTooltip, 300);
    });

    function showLoading(targetEl) {
        positionTooltip(targetEl);
        tooltip.innerHTML = '<div class="wiki-loading">Memuat dari Wikipedia...</div>';
        tooltip.classList.add('visible');
    }

    function hideTooltip() {
        tooltip.classList.remove('visible');
    }

    function positionTooltip(targetEl) {
        const rect = targetEl.getBoundingClientRect();
        const tooltipWidth = 320;
        
        let top = rect.bottom + window.scrollY + 8;
        let left = rect.left + window.scrollX - (tooltipWidth / 2) + (rect.width / 2);

        if (left < 16) left = 16;
        if (left + tooltipWidth > window.innerWidth - 16) {
            left = window.innerWidth - tooltipWidth - 16;
        }

        tooltip.style.top = `${top}px`;
        tooltip.style.left = `${left}px`;
    }

    async function fetchWikiData(lang, title, targetEl) {
        try {
            const safeTitle = encodeURIComponent(decodeURIComponent(title));
            const apiUrl = `https://${lang}.wikipedia.org/api/rest_v1/page/summary/${safeTitle}?redirect=true`;
            
            // Mengirim request polos (tanpa header khusus) agar menjadi CORS "Simple Request".
            // Jika ada header tambahan, browser akan mengirim preflight (OPTIONS).
            // Preflight yang berujung pada Redirect (301) dari Wikipedia akan diblokir browser (NetworkError).
            const res = await fetch(apiUrl);
            
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            const data = await res.json();
            
            if (currentLink !== targetEl || !tooltip.classList.contains('visible')) return;

            renderTooltip(data);
            positionTooltip(targetEl);
        } catch (e) {
            if (currentLink === targetEl) {
                tooltip.innerHTML = `<div class="wiki-error">Gagal: ${e.message}</div>`;
            }
        }
    }

    function renderTooltip(data) {
        let html = '';
        if (data.thumbnail && data.thumbnail.source) {
            html += `<img src="${data.thumbnail.source}" class="wiki-thumb" alt="">`;
        }
        html += `<div class="wiki-content">`;
        html += `<div class="wiki-title">${data.title}</div>`;
        
        let extract = data.extract;
        if (extract.length > 200) extract = extract.substring(0, 200) + '...';
        
        html += `<p class="wiki-extract">${extract}</p>`;
        html += `<div class="wiki-footer">W — Disediakan oleh Wikipedia</div>`;
        html += `</div>`;
        
        tooltip.innerHTML = html;
    }
});
