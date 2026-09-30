#!/bin/bash

GOROOT=$(go env GOROOT);
ROOT=$(pwd);
DEST="${ROOT}/../source/ui/web/public";
DESIGN="${ROOT}/../../gooey/design";

# Build the WebASM frontend
env GOOS=js GOARCH=wasm go build -o "${DEST}/main.wasm" main.go;

if [[ "$?" != "0" ]]; then
	echo "!! Build failed";
	exit 1;
fi;

# Import Go WASM Adapter
rm -f "${DEST}/wasm_exec.js";
cp "${GOROOT}/lib/wasm/wasm_exec.js" "${DEST}/wasm_exec.js";

# Import Gooey Theme
if [[ -d "${DEST}/design" ]]; then
	rm -rf "${DEST}/design";
fi;
cp -R "${DESIGN}" "${DEST}/design";

# Import App Assets
rm -f "${DEST}/index.html";
rm -f "${DEST}/wasm_init.js";
cp "${ROOT}/public/index.html" "${DEST}/index.html";
cp "${ROOT}/public/wasm_init.js" "${DEST}/wasm_init.js";

if [[ -d "${DEST}/app" ]]; then
	rm -rf "${DEST}/app";
fi;
cp -R "${ROOT}/public/app" "${DEST}/app";

echo "== App built into ${DEST} ==";
