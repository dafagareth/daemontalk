import os
import re

# We will replace all rounded utility classes with rounded-none
# Regex to match tailwind rounded classes
rounded_pattern = re.compile(r'\brounded(?:-(?:sm|md|lg|xl|2xl|3xl|full|t|b|l|r|tl|tr|bl|br)(?:-[a-z0-9]+)?)?\b')

templates_dir = 'web/templates'
for root, dirs, files in os.walk(templates_dir):
    for file in files:
        if file.endswith('.templ'):
            filepath = os.path.join(root, file)
            
            # Exceptions!
            if file == 'terminal.templ':
                continue # Skip terminal border entirely
                
            with open(filepath, 'r') as f:
                lines = f.readlines()
                
            new_lines = []
            for line in lines:
                # Skip Buy Me a Coffee button
                if 'buymeacoffee.com' in line or 'Buy me a coffee' in line or '@IconCoffee' in line:
                    new_lines.append(line)
                    continue
                # For layout_footer.templ, the <a> tag wrapping the buy me a coffee needs to be skipped
                if file == 'layout_footer.templ' and 'href="https://buymeacoffee.com' in line:
                    new_lines.append(line)
                    continue
                # Actually, the <a> tag is multi-line in layout_footer.templ
                # Let's just do a simpler string replacement for layout_footer.templ later
                
                # Replace rounded classes
                new_line = rounded_pattern.sub('rounded-none', line)
                
                # Cleanup if we created duplicates like "rounded-none rounded-none"
                new_line = re.sub(r'(rounded-none\s+)+rounded-none', 'rounded-none', new_line)
                new_lines.append(new_line)
                
            with open(filepath, 'w') as f:
                f.writelines(new_lines)

# Fix layout_footer.templ specifically for the Buy Me A Coffee button if it got modified on a different line
with open('web/templates/layout_footer.templ', 'r') as f:
    content = f.read()

# We need to ensure the Buy Me a Coffee button retains 'rounded-full'
# It looks like: class="inline-flex items-center gap-2 text-xs px-3.5 py-1.5 rounded-none bg-[#FFDD00] ...
content = content.replace('rounded-none bg-[#FFDD00]', 'rounded-full bg-[#FFDD00]')

with open('web/templates/layout_footer.templ', 'w') as f:
    f.write(content)

# Now handle CSS files
css_dir = 'web/static/css'
for root, dirs, files in os.walk(css_dir):
    for file in files:
        if file.endswith('.css'):
            filepath = os.path.join(root, file)
            with open(filepath, 'r') as f:
                content = f.read()
            
            # We want to change border-radius: <anything>; to border-radius: 0;
            # EXCEPT in utilities.css for .reaction-btn
            
            if file == 'utilities.css':
                # don't touch utilities.css because reaction-btn is there and it's the only one using border-radius
                continue
                
            if file == 'components.css':
                # Replace all border-radius: ...; with border-radius: 0;
                content = re.sub(r'border-radius:\s*[^;]+;', 'border-radius: 0;', content)
                with open(filepath, 'w') as f:
                    f.write(content)
