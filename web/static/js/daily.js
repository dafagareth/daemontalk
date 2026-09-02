// Local date formatting for the daily edition
(function() {
	function initDailyDate() {
		var el = document.getElementById("daily-date");
		if (!el) return;
		var now = new Date();
		var lang = document.documentElement.lang || "en";
		var locale = lang === "id" ? "id-ID" : (lang === "es" ? "es-ES" : "en-US");
		try {
			var formatted = now.toLocaleDateString(locale, {
				weekday: "long",
				year: "numeric",
				month: "long",
				day: "numeric"
			});
			if (formatted) {
				el.textContent = formatted;
			}
		} catch (e) {}
	}

	if (document.readyState === "loading") {
		document.addEventListener("DOMContentLoaded", initDailyDate);
	} else {
		initDailyDate();
	}
})();
