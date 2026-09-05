#!/bin/bash

ROOT="${PWD}";

quickstart_exocomp() {

	local role="$1";
	local model="$2";

	if [[ -d "${ROOT}/source/.exocomp" ]]; then
		echo "Clearing previous .exocomp session ...";
		rm -rf "${ROOT}/source/.exocomp";
	fi;

	cd "${ROOT}/source";

	echo "Starting exocomp ...";

	exec go run cmds/exocomp/main.go web --role="${role}" --model="${model}";

}

quickstart_exocomp "planner" "huihui_ai/Qwen3.6-abliterated:35b";

