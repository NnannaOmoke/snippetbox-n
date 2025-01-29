#! /usr/bin/env bash

TARGET_DIR="./target"  # Default target directory
DEFAULT_MAIN_GO="./cmd/web"
DEFAULT_OUTPUT_NAME="$(basename "$PWD")"

compile() {
    local main_go="$1"
    local output_name="$2"
    local binary_path="$TARGET_DIR/$output_name"
    
    echo "Compiling $main_go..."
    mkdir -p "$TARGET_DIR"
    
    if ! go build -o "$binary_path" "$main_go"; then
        echo "Compilation failed"
        return 1
    fi
    
    echo "Successfully built $binary_path"
    return 0
}

run() {
    local output_name="$1"
    shift  # Remove output name from arguments
    local binary_path="$TARGET_DIR/$output_name"
    
    if [[ ! -f "$binary_path" ]]; then
        echo "Binary $binary_path not found - compile first"
        return 1
    fi
    
    "$binary_path" "$@"
}

escape_analysis() {
    local main_go="$1"
    echo "Performing escape analysis on $main_go"
    go build -gcflags="-m" "$main_go" 2>&1 | grep -v '^#' | grep -v '^./'
}

vet() {
    local main_go="$1"
    echo "Running go vet on $(dirname "$main_go")"
    go vet "$(dirname "$main_go")"
}

# Main command router
case "$1" in
    compile)
        shift  # Remove 'compile' from arguments
        main_go="$DEFAULT_MAIN_GO"
        output_name="$DEFAULT_OUTPUT_NAME"
        
        # Parse flags
        while getopts ":f:o:" opt; do
            case $opt in
                f) main_go="$OPTARG" ;;
                o) output_name="$OPTARG" ;;
                \?) echo "Invalid option -$OPTARG" >&2; exit 1 ;;
            esac
        done
        shift $((OPTIND -1))
        
        compile "$main_go" "$output_name"
        ;;
        
    run)
        shift
        output_name="$DEFAULT_OUTPUT_NAME"
        
        # Parse flags
        while getopts ":o:" opt; do
            case $opt in
                o) output_name="$OPTARG" ;;
                \?) echo "Invalid option -$OPTARG" >&2; exit 1 ;;
            esac
        done
        shift $((OPTIND -1))
        
        run "$output_name" "$@"
        ;;
        
    hcheck)
        shift
        main_go="$DEFAULT_MAIN_GO"
        [[ -n "$1" ]] && main_go="$1"
        escape_analysis "$main_go"
        ;;
        
    vet)
        shift
        main_go="$DEFAULT_MAIN_GO"
        [[ -n "$1" ]] && main_go="$1"
        vet "$main_go"
        ;;
        
    *)
        echo "Usage: $0 [command] [options]"
        echo "
Commands:
  compile [-f path/to/main.go] [-o output_name]   Build executable
  run [-o output_name] [args]                     Run compiled binary
  escape [path/to/main.go]                        Show escape analysis
  vet [path/to/main.go]                           Run go vet checks
"
        exit 1
        ;;
esac
