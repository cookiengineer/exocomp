#!/bin/bash

ROOT="${PWD}";

serve_exocomp() {

	local role="$1";
	local model="$2";

	if [[ -d "${ROOT}/source/.exocomp" ]]; then
		echo "Clearing previous .exocomp session ...";
		rm -rf "${ROOT}/source/.exocomp";
	fi;

	cd "${ROOT}/source";

	echo "Starting exocomp ...";
	go run cmds/exocomp/main.go web --role="${role}" --model="${model}";

}

serve_exocomp "planner" "huihui_ai/Qwen3.6-abliterated:35b";

