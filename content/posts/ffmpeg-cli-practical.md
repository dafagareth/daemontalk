---
title: "Practical ffmpeg on the Command Line"
slug: c6d1a8f3
aliases: ffmpeg-cli-practical
date: 2026-03-30
tags: [tools, linux, shell]
lang: en
draft: false
---

ffmpeg is a complete audio and video processing toolkit contained in a single binary. It handles transcoding, remuxing, filtering, streaming, and metadata editing, and its command-line interface exposes all of this through a consistent flag structure. This post covers five common tasks that come up in practice, with commands that you can adapt directly.

## Fun Facts

**Fact 1.** ffmpeg was started by Fabrice Bellard in 2000. Bellard also wrote QEMU, the Tiny C Compiler, and held a world record for computing digits of pi on commodity hardware.

**Fact 2.** The `-c copy` flag instructs ffmpeg to copy encoded streams without decoding or re-encoding them. This is called remuxing and is orders of magnitude faster than transcoding because no CPU-intensive codec work occurs.

**Fact 3.** H.265/HEVC typically achieves the same perceptual quality as H.264/AVC at roughly half the bitrate. The tradeoff is slower encoding and slightly higher decoder CPU usage.

---

## Tips and Tricks

### 1. Transcode to H.265/HEVC for Storage Savings

Converting a library of H.264 files to H.265 with `libx265` at a reasonable CRF (Constant Rate Factor) produces smaller files with no visible quality loss at normal viewing distances. CRF 28 for H.265 is roughly equivalent to CRF 23 for H.264.

```bash
# Transcode a single file to H.265 with AAC audio
ffmpeg -i input.mp4 \
  -c:v libx265 -crf 28 -preset slow \
  -c:a aac -b:a 128k \
  output_h265.mp4

# Use hardware-accelerated encoding on Intel Quick Sync
ffmpeg -i input.mp4 \
  -c:v hevc_qsv -global_quality 28 \
  -c:a copy \
  output_h265_hw.mp4

# Check the resulting file size and stream info
ls -lh output_h265.mp4
ffprobe -v quiet -show_streams output_h265.mp4 | grep codec_name
```

The `preset` value controls the speed/compression tradeoff. `slow` gives better compression than the default `medium` at the cost of more CPU time. For a large batch, `fast` is a practical compromise.

### 2. Batch Process Files with a Shell Loop

Processing an entire directory of files requires nothing more than a shell loop. The pattern substitution `${f%.mp4}` strips the `.mp4` extension for constructing the output filename.

```bash
#!/usr/bin/env bash
# transcode_all.sh: convert all .mp4 files in current dir to H.265

set -euo pipefail

for f in *.mp4; do
  output="${f%.mp4}_h265.mp4"
  if [[ -f "$output" ]]; then
    echo "Skipping $f (output already exists)"
    continue
  fi
  echo "Processing: $f -> $output"
  ffmpeg -i "$f" \
    -c:v libx265 -crf 28 -preset fast \
    -c:a aac -b:a 128k \
    "$output" \
    -loglevel warning -stats
done

echo "Done."
```

Run it with `bash transcode_all.sh`. The `-loglevel warning -stats` flags suppress verbose output while still printing encoding progress.

### 3. Generate Video Thumbnails at Specific Timestamps

Extracting a single frame at a given timestamp is useful for generating preview images or checking a specific scene.

```bash
# Extract one frame at 1 minute 30 seconds
ffmpeg -ss 00:01:30 -i input.mp4 -vframes 1 -q:v 2 thumbnail.jpg

# Extract a frame every 60 seconds into a numbered sequence
ffmpeg -i input.mp4 -vf fps=1/60 thumb_%04d.jpg

# Extract frames at specific timestamps using a filter
ffmpeg -i input.mp4 \
  -vf "select='eq(t,10)+eq(t,60)+eq(t,120)'" \
  -vsync vfr \
  frame_%d.jpg
```

Placing `-ss` before `-i` uses the fast seek mode, which seeks before decoding and is much faster for long files. Placing `-ss` after `-i` is slower but more accurate to the exact frame.

### 4. Strip All Metadata from a Video File

Video files recorded by phones and cameras contain embedded GPS coordinates, device model information, timestamps, and other metadata. Removing it before sharing is straightforward.

```bash
# Remove all metadata from a file
ffmpeg -i input.mp4 -map_metadata -1 -c copy output_clean.mp4

# Verify that metadata is gone
ffprobe -v quiet -show_format output_clean.mp4 | grep -E "TAG:|title|artist|location"

# Strip metadata from all files in a directory
for f in *.mp4; do
  ffmpeg -i "$f" -map_metadata -1 -c copy "clean_${f}" -loglevel error
done
```

The `-map_metadata -1` flag discards all global, stream, and chapter metadata. The `-c copy` flag ensures no re-encoding occurs, so this runs at the speed of a file copy.

### 5. Remux Without Re-encoding

Remuxing changes the container format (for example, from MKV to MP4) without touching the encoded streams. This is useful when a player does not support MKV but the H.264/AAC content inside is perfectly valid MP4.

```bash
# Remux MKV to MP4 (no re-encoding)
ffmpeg -i input.mkv -c copy output.mp4

# Remux while selecting specific streams
# (video stream 0, audio stream 1, drop subtitles)
ffmpeg -i input.mkv \
  -map 0:v:0 -map 0:a:1 \
  -c copy \
  output_selected.mp4

# Remux TS to MKV and preserve all subtitle tracks
ffmpeg -i broadcast.ts -c copy -map 0 output.mkv

# Check the output container and stream layout
ffprobe -v quiet -show_streams -show_format output.mp4 \
  | grep -E "codec_name|codec_type|duration"
```

When remuxing from MKV to MP4, be aware that some subtitle formats (ASS/SSA) are not supported in the MP4 container. ffmpeg will warn about this and skip unsupported streams if you do not include `-map 0` explicitly.
