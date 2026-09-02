// KaTeX Math Rendering & Scroll Overflow Indicators
function dtAttachMathScrollFade() {
	requestAnimationFrame(function() {
		document.querySelectorAll(".katex-display, .math-display").forEach(function(el) {
			if (el.parentElement && el.parentElement.classList.contains("katex-scroll-wrapper")) {
				var existingFade = el.parentElement.querySelector(".katex-fade-right");
				if (existingFade) {
					var isOver = el.scrollWidth > (el.clientWidth + 4);
					existingFade.style.display = isOver ? "block" : "none";
				}
				return;
			}
			var wrapper = document.createElement("div");
			wrapper.className = "katex-scroll-wrapper relative my-4";
			el.parentNode.insertBefore(wrapper, el);
			wrapper.appendChild(el);

			var fade = document.createElement("div");
			fade.className = "katex-fade-right pointer-events-none absolute right-0 top-0 bottom-0 w-8 sm:w-12 bg-gradient-to-l from-[var(--c-bg)] to-transparent transition-opacity duration-200";
			wrapper.appendChild(fade);

			function checkOverflow() {
				var isOverflow = el.scrollWidth > (el.clientWidth + 4);
				if (!isOverflow) {
					fade.style.display = "none";
				} else {
					fade.style.display = "block";
					var atEnd = (el.scrollWidth - el.scrollLeft - el.clientWidth) < 6;
					fade.style.opacity = atEnd ? "0" : "1";
				}
			}

			el.addEventListener("scroll", checkOverflow, { passive: true });
			checkOverflow();
			window.addEventListener("resize", checkOverflow, { passive: true });
		});
	});
}

function dtRenderMath() {
	if (typeof renderMathInElement === "function") {
		renderMathInElement(document.body, {
			delimiters: [
				{left: "$$", right: "$$", display: true},
				{left: "$", right: "$", display: false},
				{left: "\\(", right: "\\)", display: false},
				{left: "\\[", right: "\\]", display: true}
			],
			ignoredTags: ["script", "noscript", "style", "textarea", "pre", "code"],
			throwOnError: false
		});
		dtAttachMathScrollFade();
	} else {
		setTimeout(dtRenderMath, 60);
	}
}

window.dtRenderMath = dtRenderMath;
window.dtAttachMathScrollFade = dtAttachMathScrollFade;

if (document.readyState === "loading") {
	document.addEventListener("DOMContentLoaded", dtRenderMath);
} else {
	dtRenderMath();
}

document.addEventListener("htmx:afterSwap", function() {
	if (typeof dtRenderMath === "function") {
		dtRenderMath();
	}
});
