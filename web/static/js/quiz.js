(function() {
					var currentQuizIndex = 0;
					function getSlides() {
						return document.querySelectorAll('.quiz-slide');
					}

					window.goToQuizStep = function(idx) {
						var slides = getSlides();
						var totalSlides = slides.length;
						if (totalSlides === 0 || idx < 0 || idx >= totalSlides) return;
						currentQuizIndex = idx;
						slides.forEach(function(slide, i) {
							if (i === idx) {
								slide.classList.remove('hidden');
								slide.classList.add('block');
							} else {
								slide.classList.remove('block');
								slide.classList.add('hidden');
							}
						});

						var dots = document.querySelectorAll('.quiz-dot');
						dots.forEach(function(dot, i) {
							if (i === idx) {
								dot.className = 'quiz-dot rounded-full transition-all cursor-pointer w-4 h-1.5 bg-[var(--c-link)]';
							} else {
								dot.className = 'quiz-dot rounded-full transition-all cursor-pointer w-1.5 h-1.5 bg-[var(--c-border)] hover:bg-[var(--c-muted)]';
							}
						});

						var stepLabel = document.getElementById('quiz-step-label');
						if (stepLabel) {
							var isId = stepLabel.getAttribute('data-lang') === 'id';
							stepLabel.textContent = isId ? ((idx + 1) + ' dari ' + totalSlides) : ((idx + 1) + ' of ' + totalSlides);
						}

						var prevBtn = document.getElementById('quiz-prev-btn');
						var nextBtn = document.getElementById('quiz-next-btn');
						if (prevBtn) prevBtn.disabled = (idx === 0);
						if (nextBtn) nextBtn.disabled = (idx === totalSlides - 1);
					};

					window.prevQuizStep = function() {
						if (currentQuizIndex > 0) {
							window.goToQuizStep(currentQuizIndex - 1);
						}
					};

					window.nextQuizStep = function() {
						var slides = getSlides();
						if (currentQuizIndex < slides.length - 1) {
							window.goToQuizStep(currentQuizIndex + 1);
						}
					};

					window.handleQuizClick = function(btn) {
						var slide = btn.closest('.quiz-slide');
						if (!slide) return;
						var allBtns = slide.querySelectorAll('.quiz-options button');
						var isCorrect = btn.getAttribute('data-correct') === 'true';
						var explain = btn.getAttribute('data-explain') || '';
						var feedback = slide.querySelector('.quiz-feedback');

						allBtns.forEach(function(b) {
							b.disabled = true;
							b.classList.remove('hover:bg-hover', 'cursor-pointer');
							b.classList.add('opacity-60', 'cursor-default');
							if (b.getAttribute('data-correct') === 'true') {
								b.classList.remove('opacity-60', 'border-border');
								b.classList.add('border-emerald-500/60', 'bg-emerald-500/10', 'text-emerald-400');
							}
						});

						if (isCorrect) {
							btn.classList.remove('opacity-60', 'border-border');
							btn.classList.add('border-emerald-500', 'bg-emerald-500/15', 'text-emerald-400');
						} else {
							btn.classList.remove('opacity-60', 'border-border');
							btn.classList.add('border-rose-500', 'bg-rose-500/15', 'text-rose-400');
						}

						if (feedback) {
							var isId = (document.getElementById('quiz-step-label')?.getAttribute('data-lang') === 'id');
							feedback.className = 'quiz-feedback mt-4 p-3.5 rounded-lg border text-xs font-mono leading-relaxed ' + 
								(isCorrect ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-300' : 'border-rose-500/40 bg-rose-500/10 text-rose-300');
							var statusText = isCorrect ? (isId ? '✓ Tepat Sekali' : '✓ Correct') : (isId ? '✗ Kurang Tepat' : '✗ Incorrect');
							feedback.innerHTML = '<span class="font-bold uppercase tracking-wider block mb-1">' + statusText + '</span>' + explain;
							feedback.classList.remove('hidden');
						}
					};
				})();