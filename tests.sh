#!/usr/bin/env bash

# update-readme.sh
# Updates README.md test results sections with actual test output
#
# Usage: ./update-readme.sh [readme_path]
#
# The script looks for marker pairs like:
#   <!-- TST8080.COM -->
#   <!-- /TST8080.COM -->
# and replaces content between them with the test output.

set -euo pipefail

README="${1:-README.md}"
EMULATOR="./tmp/main"
TEST_DIR="./sub/8080/cpu_tests"

# Build the emulator first
echo "Building emulator..."
go build -o "$EMULATOR" .

# Array of test files to run
TEST_FILES=(
    "TST8080.COM"
    "8080PRE.COM"
    "CPUTEST.COM"
    "8080EXM.COM"
)

# Function to update a section in the README
update_section() {
    local marker="$1"
    local content="$2"
    local readme="$3"

    local start_marker="<!-- ${marker} -->"
    local end_marker="<!-- /${marker} -->"

    # Check if markers exist
    if ! grep -q "$start_marker" "$readme" || ! grep -q "$end_marker" "$readme"; then
        echo "Warning: Markers for $marker not found in $readme, skipping..."
        return
    fi

    # Create temp file
    local tmpfile
    tmpfile=$(mktemp)

    # Use awk to replace content between markers
    awk -v start="$start_marker" -v end="$end_marker" -v content="$content" '
        $0 == start {
            print
            print "```txt"
            print content
            print "```"
            skip = 1
            next
        }
        $0 == end {
            skip = 0
            print
            next
        }
        !skip { print }
    ' "$readme" > "$tmpfile"

    mv "$tmpfile" "$readme"
    echo "Updated section: $marker"
}

# Run each test and update README
for test_file in "${TEST_FILES[@]}"; do
    test_path="${TEST_DIR}/${test_file}"

    if [[ ! -f "$test_path" ]]; then
        echo "Warning: Test file $test_path not found, skipping..."
        continue
    fi

    echo "Running $test_file..."

    # Run the test and capture output (with timeout)
    output=$(timeout 120s "$EMULATOR" --cpm --headless --unthrottle "$test_path" 2>&1 | tr -d '\r\0' | cat -s || true)

    # Update the README section
    update_section "$test_file" "$output" "$README"
done

echo "Done! README updated."
