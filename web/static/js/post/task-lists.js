(function() {
    // Interactive task lists / checklists with local persistence
    var path = window.location.pathname;
    var match = path.match(/\/blog\/([^\/]+)$/);
    var postSlug = match ? match[1] : (window.location.pathname.split("/").filter(Boolean).pop() || "default");
    var storageKey = "dt_checklist_" + postSlug;
    var state = {};
    try {
        state = JSON.parse(localStorage.getItem(storageKey) || "{}");
    } catch(e) {}

    var checkboxes = document.querySelectorAll("#prose-body input[type='checkbox']");
    checkboxes.forEach(function(cb, idx) {
        cb.disabled = false;
        cb.removeAttribute("disabled");
        cb.style.cursor = "pointer";
        cb.style.pointerEvents = "auto";
        cb.style.touchAction = "manipulation";

        var parentLi = cb.closest("li");
        if (parentLi) {
            parentLi.style.listStyle = "none";
            var parentUl = parentLi.parentElement;
            if (parentUl) parentUl.style.listStyle = "none";
        }

        // Restore saved state
        if (state[idx]) {
            cb.checked = true;
            if (parentLi) parentLi.classList.add("task-checked");
        }

        cb.addEventListener("change", function() {
            state[idx] = cb.checked;
            if (parentLi) {
                if (cb.checked) parentLi.classList.add("task-checked");
                else parentLi.classList.remove("task-checked");
            }
            try {
                localStorage.setItem(storageKey, JSON.stringify(state));
            } catch(e) {}
        });
    });
})();
