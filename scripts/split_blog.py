import os

with open('internal/handler/blog.go', 'r') as f:
    lines = f.readlines()

comment_funcs = ["func (h *Handler) DeleteComment", "func (h *Handler) PostComment", "func (h *Handler) renderCommentList", "func (h *Handler) sendCommentNotification"]

in_comment_func = False
comment_lines = []
blog_lines = []

for line in lines:
    is_start = any(line.startswith(cf) for cf in comment_funcs)
    
    if is_start:
        in_comment_func = True
    
    if in_comment_func:
        comment_lines.append(line)
        # simplistic block detection
        if line.startswith("}"):
            in_comment_func = False
    else:
        blog_lines.append(line)

with open('internal/handler/blog.go', 'w') as f:
    f.writelines(blog_lines)

comment_go = """package handler

import (
\t"bytes"
\t"encoding/json"
\t"fmt"
\t"log"
\t"net/http"
\t"os"
\t"strings"
\t"time"

\t"portfolio/internal/comment"
\t"portfolio/internal/i18n"
\t"portfolio/internal/postdb"

\t"github.com/go-chi/chi/v5"
)

""" + "".join(comment_lines)

with open('internal/handler/comment.go', 'w') as f:
    f.write(comment_go)
