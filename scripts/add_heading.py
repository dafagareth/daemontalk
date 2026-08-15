import os

posts_dir = "/home/dd/Academic/Projects/Portfolio/content/posts"
for root, _, files in os.walk(posts_dir):
    for f in files:
        if f.endswith(".md"):
            filepath = os.path.join(root, f)
            with open(filepath, "r", encoding="utf-8") as file:
                content = file.read()
            
            if "[^1]:" in content and "## Referensi" not in content:
                # Add heading before the first footnote
                new_content = content.replace("[^1]:", "## Referensi\n\n[^1]:", 1)
                with open(filepath, "w", encoding="utf-8") as file:
                    file.write(new_content)
                print(f"Added 'Referensi' heading to {f}")
