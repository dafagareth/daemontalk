import re

with open('web/static/css/input.css', 'r') as f:
    content = f.read()

# 1. Themes & Variables
themes = re.search(r'(@theme.*?)\n@layer base', content, re.DOTALL).group(1)
with open('web/static/css/theme.css', 'w') as f:
    f.write(themes)

# 2. Base
base = re.search(r'(@layer base \{.*?\n\})\n\n/\* Sub-nav', content, re.DOTALL).group(1)
with open('web/static/css/base.css', 'w') as f:
    f.write(base)

# 3. Components (Prose, Admin, etc.)
components_raw = content.split('@layer components {')[1].split('/* HTMX loading states')[0]
components = '@layer components {\n' + components_raw
with open('web/static/css/components.css', 'w') as f:
    f.write(components)

# 4. Utilities and others
utilities = '/* HTMX loading states' + content.split('/* HTMX loading states')[1]
with open('web/static/css/utilities.css', 'w') as f:
    f.write(utilities)

# 5. New input.css
new_input = """@import "tailwindcss";
@source "../../templates";
@source "../../templates/**/*.templ";
@source "../../**/*.templ";
@source "../../../*.go";
@source "../../../internal/**/*.go";

@import "./theme.css";
@import "./base.css";
@import "./components.css";
@import "./utilities.css";
"""
with open('web/static/css/input.css', 'w') as f:
    f.write(new_input)

print("CSS Split completed.")
