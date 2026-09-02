// Tab switching and hash routing for the contribute page
function switchContributeTab(tabName) {
	document.querySelectorAll('.contribute-panel').forEach(function(p) {
		p.classList.add('hidden');
	});
	var activePanel = document.getElementById('panel-' + tabName);
	if (activePanel) {
		activePanel.classList.remove('hidden');
	}
	document.querySelectorAll('.contribute-tab-btn').forEach(function(btn) {
		btn.classList.remove('bg-surface', 'text-text', 'border-b-link');
		btn.classList.add('text-muted', 'bg-chip/30', 'border-b-transparent');
	});
	var activeBtn = document.getElementById('tab-btn-' + tabName);
	if (activeBtn) {
		activeBtn.classList.remove('text-muted', 'bg-chip/30', 'border-b-transparent');
		activeBtn.classList.add('bg-surface', 'text-text', 'border-b-link');
	}
	if (history.pushState) {
		history.pushState(null, null, '#' + tabName);
	}
}

window.switchContributeTab = switchContributeTab;

document.addEventListener('DOMContentLoaded', function() {
	var hash = window.location.hash.replace('#', '');
	if (hash && ['dispatches', 'engine', 'corrections', 'i18n', 'contributors'].indexOf(hash) !== -1) {
		switchContributeTab(hash);
	}
});
