(function(){
				// Copy button on each code block
				document.querySelectorAll("#prose-body pre").forEach(function(pre){
					var wrap=document.createElement("div");
					wrap.className="code-wrap";
					pre.parentNode.insertBefore(wrap,pre);
					wrap.appendChild(pre);
					var btn=document.createElement("button");
					btn.textContent="copy";
					btn.className="copy-btn";
					btn.setAttribute("aria-label","Copy code");
					btn.addEventListener("click",function(){
						var code=pre.querySelector("code");
						var text=code?code.innerText:pre.innerText;
						function done(){btn.textContent="copied!";setTimeout(function(){btn.textContent="copy";},2000);}
						function fb(){var ta=document.createElement("textarea");ta.value=text;ta.style.cssText="position:fixed;top:-9999px;opacity:0";document.body.appendChild(ta);ta.focus();ta.select();try{document.execCommand("copy");done();}catch(e){}document.body.removeChild(ta);}
						if(navigator.clipboard&&navigator.clipboard.writeText){navigator.clipboard.writeText(text).then(done).catch(fb);}else{fb();}
					});
					wrap.appendChild(btn);
				});
				// ToC scroll-spy: highlight the active section
				var tocLinks=document.querySelectorAll(".toc-link");
				if(tocLinks.length>0){
					var headings=Array.from(tocLinks).map(function(a){return document.getElementById(a.getAttribute("href").slice(1));}).filter(Boolean);
					function setActive(){
						var y=window.scrollY+100,active=headings[0];
						for(var i=0;i<headings.length;i++){if(headings[i].offsetTop<=y)active=headings[i];}
						tocLinks.forEach(function(a){a.classList.toggle("toc-active",active&&a.getAttribute("href")==="#"+active.id);});
					}
					window.addEventListener("scroll",setActive,{passive:true});
					setActive();
				}
				// Mobile ToC chevron rotation
				var mt=document.getElementById("mobile-toc");
				if(mt){mt.addEventListener("toggle",function(){var c=mt.querySelector(".toc-chevron");if(c)c.style.transform=mt.open?"rotate(180deg)":"";});}
			})();