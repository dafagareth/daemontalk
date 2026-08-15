#!/usr/bin/env bash
# Script untuk mengelola postingan markdown (archive, restore, delete) lewat terminal.
set -euo pipefail

CMD=${1:-}
SLUG=${2:-}

if [ -z "$CMD" ] || [ -z "$SLUG" ]; then
    echo "Penggunaan: ./manage-post.sh [archive|restore|delete] [slug]"
    exit 1
fi

# Cari file di content/posts/ yang memiliki slug tersebut di frontmatter
FILE=""
for f in content/posts/*.md content/posts/*.md.archive; do
    if [ -f "$f" ]; then
        if grep -q "^slug: \?\(['\"]\?\)$SLUG\1$" "$f"; then
            FILE="$f"
            break
        fi
    fi
done

if [ -z "$FILE" ]; then
    echo "Error: Postingan dengan slug '$SLUG' tidak ditemukan."
    exit 1
fi

case "$CMD" in
    archive)
        if [[ "$FILE" == *.md.archive ]]; then
            echo "Postingan '$SLUG' sudah di-archive sebelumnya."
            exit 0
        fi
        NEW_FILE="${FILE}.archive"
        mv "$FILE" "$NEW_FILE"
        echo "Postingan '$SLUG' berhasil di-archive (ditarik dari publik)."
        echo "Nama file diubah menjadi: $NEW_FILE"
        ;;
    restore)
        if [[ "$FILE" != *.md.archive ]]; then
            echo "Postingan '$SLUG' tidak dalam status archive (sudah aktif/publik)."
            exit 0
        fi
        NEW_FILE="${FILE%.archive}"
        mv "$FILE" "$NEW_FILE"
        echo "Postingan '$SLUG' berhasil di-restore (dipublikasikan kembali)."
        echo "Nama file diubah menjadi: $NEW_FILE"
        ;;
    delete)
        read -p "Apakah Anda yakin ingin menghapus postingan '$SLUG' secara PERMANEN beserta aset gambarnya? (y/N) " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            rm "$FILE"
            echo "Postingan '$SLUG' ($FILE) telah dihapus secara permanen."
            IMG_DIR="web/static/images/posts/$SLUG"
            if [ -d "$IMG_DIR" ]; then
                rm -rf "$IMG_DIR"
                echo "Folder aset gambar ($IMG_DIR) telah dihapus."
            fi
        else
            echo "Penghapusan dibatalkan."
        fi
        ;;
    *)
        echo "Error: Perintah tidak dikenal: $CMD. Gunakan archive, restore, atau delete."
        exit 1
        ;;
esac
