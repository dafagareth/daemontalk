import os

posts_dir = "/home/dd/Academic/Projects/Portfolio/content/posts"
for root, _, files in os.walk(posts_dir):
    for f in files:
        if f.endswith(".md"):
            filepath = os.path.join(root, f)
            with open(filepath, "r", encoding="utf-8") as file:
                content = file.read()
            
            if "## Referensi" in content and "[^1]:" in content:
                if content.count("[^1]") == 1:
                    # Insert inline footnotes at the end of the last paragraph before the heading
                    new_content = content.replace("\n\n## Referensi", "[^1][^2]\n\n## Referensi", 1)
                    with open(filepath, "w", encoding="utf-8") as file:
                        file.write(new_content)
                    print(f"Fixed {f}")
