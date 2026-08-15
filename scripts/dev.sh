#!/usr/bin/env bash
# Helper buat ngembangin daemontalk. Jalankan: ./dev.sh <perintah>
#
#   ./dev.sh up        build penuh (templ + css + go) lalu jalankan server
#   ./dev.sh run       sama dengan 'up'
#   ./dev.sh watch     build penuh lalu watch perubahan file secara otomatis
#   ./dev.sh build     build penuh tanpa menjalankan server
#   ./dev.sh css       rebuild Tailwind CSS saja
#   ./dev.sh templ     generate templ saja
#   ./dev.sh restart   jalankan ulang server (refresh cache GitHub, tanpa build)
#   ./dev.sh stop      matikan server
#   ./dev.sh logs      lihat log server (live)
#   ./dev.sh help      tampilkan bantuan ini

set -euo pipefail

# Selalu kerja dari root project (parent dari folder scripts/).
cd "$(dirname "$(dirname "$0")")"

PORT=8080
BIN=/tmp/daemontalk_run
LOG=/tmp/daemontalk.log
CSS_IN=web/static/css/input.css
CSS_OUT=web/static/css/main.css

say()  { printf '\033[1;36m=> %s\033[0m\n' "$*"; }
ok()   { printf '\033[1;32m=> %s\033[0m\n' "$*"; }
warn() { printf '\033[1;33m=> %s\033[0m\n' "$*"; }

gen_templ() {
	say "templ generate"
	templ generate
}

build_css() {
	say "rebuild Tailwind CSS"
	npx @tailwindcss/cli -i "$CSS_IN" -o "$CSS_OUT" --minify 2>/dev/null
}

build_go() {
	say "go build -> $BIN"
	go build -o "$BIN" .
}

stop_server() {
	fuser -k "${PORT}/tcp" 2>/dev/null || true
	sleep 0.3
}

start_server() {
	say "jalankan server di http://localhost:${PORT}"
	if [ ! -f .env ]; then
		warn "peringatan: .env tidak ada, GITHUB_TOKEN tidak dimuat"
	fi
	set -a
	# shellcheck disable=SC1091
	[ -f .env ] && . ./.env
	set +a
	nohup "$BIN" >"$LOG" 2>&1 &
	sleep 0.5
	ok "server jalan  (log: $LOG)"
}

# ---------------------------------------------------------------------------
# watch: hot-rebuild saat file berubah
# Strategi: pakai inotifywait jika ada, fallback ke polling checksum sha1sum.
# Klasifikasi perubahan:
#   .templ          -> templ generate + go build + restart
#   .go             -> go build + restart
#   input.css / web/static/**  -> css rebuild saja, restart
#   content/posts/** (markdown) -> restart saja (tidak perlu rebuild)
# ---------------------------------------------------------------------------
watch_loop() {
	local changed=""
	local last_hash=""
	local cur_hash=""

	# Kumpulkan semua file yang ingin dipantau.
	_collect_hash() {
		find web internal main.go \
			-type f \( -name "*.go" -o -name "*.templ" -o -name "*.css" \) \
			2>/dev/null | sort | xargs sha1sum 2>/dev/null
		find content/posts -type f -name "*.md" 2>/dev/null | sort | xargs sha1sum 2>/dev/null
	}

	if command -v inotifywait &>/dev/null; then
		say "menggunakan inotifywait untuk deteksi perubahan"
		_watch_inotify
	else
		warn "inotifywait tidak ditemukan, memakai polling setiap 1 detik"
		warn "(install inotify-tools untuk respons yang lebih cepat)"
		_watch_poll
	fi
}

_watch_inotify() {
	say "memantau perubahan... (Ctrl+C untuk berhenti)"
	while true; do
		# Tunggu event dari inotifywait, tangkap nama file yang berubah.
		changed=$(inotifywait -q -r \
			--event modify,create,delete,move \
			--format '%w%f' \
			web internal main.go content/posts 2>/dev/null | head -1) || true

		[ -z "$changed" ] && continue

		_handle_change "$changed"
	done
}

_watch_poll() {
	say "memantau perubahan... (Ctrl+C untuk berhenti)"
	local last_hash=""
	local cur_hash=""

	last_hash=$(_collect_hash)

	while true; do
		sleep 1
		cur_hash=$(_collect_hash)

		if [ "$cur_hash" != "$last_hash" ]; then
			# Temukan file mana yang berubah.
			changed=$(diff <(echo "$last_hash") <(echo "$cur_hash") \
				| grep '^[<>]' | awk '{print $NF}' | head -1)
			last_hash="$cur_hash"
			_handle_change "$changed"
		fi
	done
}

_handle_change() {
	local file="$1"
	say "perubahan terdeteksi: $file"

	local rebuild_templ=0
	local rebuild_go=0
	local rebuild_css=0
	local restart_only=0

	case "$file" in
		*.templ)
			rebuild_templ=1
			rebuild_go=1
			;;
		*.go)
			rebuild_go=1
			;;
		*.css | web/static/*)
			rebuild_css=1
			;;
		content/posts/*)
			restart_only=1
			;;
		*)
			rebuild_go=1
			;;
	esac

	local failed=0

	if [ "$rebuild_templ" = 1 ]; then
		gen_templ || { warn "templ generate gagal, skip restart"; failed=1; }
	fi

	if [ "$failed" = 0 ] && [ "$rebuild_go" = 1 ]; then
		build_go || { warn "go build gagal, server tidak direstart"; failed=1; }
	fi

	if [ "$rebuild_css" = 1 ]; then
		build_css || warn "css build gagal"
	fi

	if [ "$failed" = 0 ]; then
		stop_server
		start_server
	fi
}

_collect_hash() {
	find web internal main.go \
		-type f \( -name "*.go" -o -name "*.templ" -o -name "*.css" \) \
		2>/dev/null | sort | xargs sha1sum 2>/dev/null
	find content/posts -type f -name "*.md" 2>/dev/null | sort | xargs sha1sum 2>/dev/null
}

# ---------------------------------------------------------------------------

cmd="${1:-up}"
case "$cmd" in
	up | run)
		gen_templ
		build_css
		build_go
		stop_server
		start_server
		;;
	watch)
		gen_templ
		build_css
		build_go
		stop_server
		start_server
		watch_loop
		;;
	build)
		gen_templ
		build_css
		build_go
		say "selesai build (server tidak dijalankan)"
		;;
	css)
		build_css
		ok "CSS selesai. Hard refresh browser (Ctrl+Shift+R)"
		;;
	templ)
		gen_templ
		;;
	restart)
		stop_server
		start_server
		;;
	stop)
		stop_server
		ok "server dimatikan"
		;;
	logs)
		say "log live (Ctrl+C untuk keluar)"
		tail -f "$LOG"
		;;
	help | -h | --help)
		grep '^#' "$0" | grep -v '^#!' | sed 's/^# \{0,1\}//'
		;;
	*)
		echo "perintah tidak dikenal: $cmd" >&2
		echo "jalankan: ./dev.sh help" >&2
		exit 1
		;;
esac
