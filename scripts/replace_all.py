import os
import re

directories = ['internal', 'web', 'content', 'cmd']
files_to_check = ['README.md', 'LICENSE', 'Makefile', 'docker-compose.dev.yml']

def replace_in_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return

    if not re.search(r'sudolog', content, re.IGNORECASE):
        return

    # Replace specific casing
    new_content = re.sub(r'sudolog', 'daemontalk', content)
    new_content = re.sub(r'Sudolog', 'DaemonTalk', new_content)
    new_content = re.sub(r'SUDOLOG', 'DAEMONTALK', new_content)

    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Updated {filepath}")

for d in directories:
    if not os.path.exists(d):
        continue
    for root, _, files in os.walk(d):
        for file in files:
            ext = os.path.splitext(file)[1]
            if ext in ['.go', '.templ', '.md', '.js', '.css', '.html', '.txt', '.yml', '.yaml', '.toml', '.json', '.xml']:
                replace_in_file(os.path.join(root, file))

for f in files_to_check:
    if os.path.exists(f):
        replace_in_file(f)

print("Done replacing.")
