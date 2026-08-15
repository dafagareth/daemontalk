import os

with open('web/templates/blog_portal.templ', 'r') as f:
    content = f.read()

# Define headers
header_base = """package templates

import (
\t"fmt"

\t"portfolio/internal/i18n"
\t"portfolio/internal/post"
)

"""
header_minimal = """package templates

import (
\t"portfolio/internal/i18n"
\t"portfolio/internal/post"
)

"""

# We'll split the content using the templ keyword
components_text = content.split('// ---------------------------------------------------------')[1:]

# components_text[0] is STRICT COMPONENTS block
# components_text[1] is portalStrictCategoryGrid block

# Save main layout to blog_portal.templ
main_layout = content.split('// ---------------------------------------------------------')[0]
with open('web/templates/blog_portal.templ', 'w') as f:
    f.write(main_layout.strip() + "\n")

# Save strict components
strict_components = header_minimal + components_text[1]
# But wait, components_text[1] has the components, but wait, components_text is split by the comment line
# Let's just find the index of templ functions

import re

def extract_templ(name, text):
    pattern = r'(templ ' + name + r'\(.*?\)\s*\{.*?\n\})'
    match = re.search(pattern, text, re.DOTALL)
    if match:
        return match.group(1)
    return ""

def remove_templ(name, text):
    pattern = r'(templ ' + name + r'\(.*?\)\s*\{.*?\n\})'
    return re.sub(pattern, '', text, flags=re.DOTALL)

# Let's just do it manually with split since regex with balanced braces is hard in Python
# Actually, I can just use line numbers or find "templ "
lines = content.split('\n')
blocks = []
current_block = []
in_templ = False

for line in lines:
    if line.startswith("templ "):
        if current_block and not in_templ:
            pass # ignore text outside templ
        in_templ = True
    
    if in_templ:
        current_block.append(line)
        if line == "}":
            blocks.append("\n".join(current_block))
            current_block = []
            in_templ = False

# blocks[0] = portalNewsLayout
# blocks[1] = portalStrictLeadStory
# blocks[2] = portalStrictSubStory
# blocks[3] = portalStrictFeatureStory
# blocks[4] = portalStrictListItem
# blocks[5] = portalStrictCategoryGrid
# blocks[6] = portalSidebarWithCategories

with open('web/templates/blog_portal.templ', 'w') as f:
    f.write(header_base + blocks[0] + "\n")

with open('web/templates/portal_components.templ', 'w') as f:
    f.write(header_minimal + blocks[1] + "\n\n" + blocks[2] + "\n\n" + blocks[3] + "\n\n" + blocks[4] + "\n")

with open('web/templates/portal_categories.templ', 'w') as f:
    f.write(header_minimal + blocks[5] + "\n")

with open('web/templates/portal_sidebar.templ', 'w') as f:
    f.write(header_base + blocks[6] + "\n")

print("Templates split successfully")
