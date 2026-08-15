import os

os.makedirs('internal/runner', exist_ok=True)

with open('internal/handler/runner.go', 'r') as f:
    lines = f.readlines()

# Find where executeGoCode starts
exec_start = -1
for i, line in enumerate(lines):
    if "func formatGoCode(code string)" in line or "// formatGoCode ensures" in line:
        exec_start = i
        break

if exec_start != -1:
    handler_code = "".join(lines[:exec_start])
    runner_code = "".join(lines[exec_start:])

    # Add package declaration and imports to runner_code
    runner_imports = """package runner

import (
\t"bytes"
\t"context"
\t"fmt"
\t"os"
\t"os/exec"
\t"path/filepath"
\t"strings"
\t"time"
)

// RunResponse holds the output of the executed code
type RunResponse struct {
\tStdout string `json:"stdout"`
\tStderr string `json:"stderr"`
\tError  string `json:"error,omitempty"`
}

"""
    runner_code = runner_code.replace("type runResponse struct {\n\tStdout string `json:\"stdout\"`\n\tStderr string `json:\"stderr\"`\n\tError  string `json:\"error,omitempty\"`\n}\n\n", "")
    runner_code = runner_code.replace("runResponse", "RunResponse")
    runner_code = runner_imports + runner_code

    with open('internal/runner/runner.go', 'w') as f:
        f.write(runner_code)

    # Now update handler/runner.go to use the new runner package
    handler_code = handler_code.replace("runResponse", "runner.RunResponse")
    handler_code = handler_code.replace("executeGoCode", "runner.ExecuteGoCode")
    handler_code = handler_code.replace("executePythonCode", "runner.ExecutePythonCode")
    handler_code = handler_code.replace("executeNodeCode", "runner.ExecuteNodeCode")
    handler_code = handler_code.replace("executeShellSim", "runner.ExecuteShellSim")
    
    # Capitalize the functions in runner package
    with open('internal/runner/runner.go', 'r') as f:
        rc = f.read()
    rc = rc.replace("func executeGoCode", "func ExecuteGoCode")
    rc = rc.replace("func executePythonCode", "func ExecutePythonCode")
    rc = rc.replace("func executeNodeCode", "func ExecuteNodeCode")
    rc = rc.replace("func executeShellSim", "func ExecuteShellSim")
    with open('internal/runner/runner.go', 'w') as f:
        f.write(rc)
        
    with open('internal/handler/runner.go', 'w') as f:
        f.write(handler_code)
    
    print("Runner split successful.")
else:
    print("Could not find execution functions.")
