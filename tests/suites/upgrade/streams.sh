run_simplestream_metadata_last_stable() {
    local jujud_version previous_version

    jujud_version=$(jujud_version)
    previous_version=$(last_stable_version "${jujud_version}")

    exec_simplestream_metadata "stable" "${jujud_version}" "${previous_version}"
}

run_simplestream_metadata_prior_stable() {
    local jujud_version previous_version

    jujud_version=$(jujud_version)
    previous_version=$(prior_stable_version "${jujud_version}")

    exec_simplestream_metadata "prior" "${jujud_version}" "${previous_version}"
}

exec_simplestream_metadata() {
    local test_name version jujud_version stable_version

    version=$(jujud version)

    test_name=${1}
    jujud_version=${2}
    stable_version=${3}

    echo "===> Using jujud version ${version}"
    echo "===> Testing against stable version ${stable_version}"

    local focal_version bionic_version

    focal_version=$(series_version "${version}" "focal")
    bionic_version=$(series_version "${version}" "bionic")

    OUT=$(sudo snap install juju --classic --channel="${stable_version}/stable" 2>&1 || echo "FALLBACK")
    if [ "${OUT}" == "FALLBACK" ] || [[ "${OUT}" =~ (.*)is\ already\ installed(.*)$ ]]; then
        echo "===> Using snap refresh juju ${stable_version}/stable"
        sudo snap refresh juju --channel="${stable_version}/stable"
    fi

    add_clean_func "remove_upgrade_tools"
    add_upgrade_tools "${focal_version}"
    add_upgrade_tools "${bionic_version}"

    add_clean_func "remove_upgrade_metadata"
    juju metadata generate-agents \
        --clean \
        --prevent-fallback \
        -d "./tests/suites/upgrade/streams/"

    add_clean_func "kill_server"
    start_server "./tests/suites/upgrade/streams/tools"

    ip_address=$(ip -4 -o addr show scope global | awk '{gsub(/\/.*/,"",$4); print $4}' | head -n 1)

    name="test-upgrade-${test_name}-stream"

    file="${TEST_DIR}/test-upgrade-${test_name}-stream.log"
    /snap/bin/juju bootstrap "lxd" "${name}" \
        --config agent-metadata-url="http://${ip_address}:8000/" \
        --config test-mode=true 2>&1 | OUTPUT "${file}"
    echo "${name}" >> "${TEST_DIR}/jujus"

    juju deploy ./tests/suites/upgrade/charms/ubuntu
    wait_for "ubuntu" "$(idle_condition "ubuntu")"

    local CURRENT UPDATED

    CURRENT=$(juju machines -m controller --format=json | jq -r '.machines | .["0"] | .["juju-status"] | .version')
    echo "==> Current juju version ${CURRENT}"

    juju upgrade-controller --agent-version="${jujud_version}"

    attempt=0
    while true; do
        UPDATED=$(juju machines -m controller --format=json | jq -r '.machines | .["0"] | .["juju-status"] | .version' || echo "${CURRENT}")
        if [ "$CURRENT" != "$UPDATED" ]; then
            break
        fi
        echo "[+] (attempt ${attempt}) polling machines"
        sleep 10
        attempt=$((attempt+1))
        if [ "$attempt" -eq 48 ]; then
            echo "Upgrade controller timed out"
            exit 1
        fi
    done

    sleep 10
    juju upgrade-charm ubuntu --path=./tests/suites/upgrade/charms/ubuntu

    sleep 10
    wait_for "ubuntu" "$(idle_condition "ubuntu")"
}

test_upgrade_simplestream() {
    if [ -n "$(skip 'test_upgrade_simplestream')" ]; then
        echo "==> SKIP: Asked to skip stream tests"
        return
    fi

    (
        set_verbosity

        cd .. || exit

        run "run_simplestream_metadata_last_stable"
        run "run_simplestream_metadata_prior_stable"
    )
}

last_stable_version() {
    local version major minor patch parts

    version="${1}"

    # shellcheck disable=SC2116
    version=$(echo "${version%-*}")

    major=$(echo "${version}" | cut -d '.' -f 1)
    minor=$(echo "${version}" | cut -d '.' -f 2)
    patch=$(echo "${version}" | cut -d '.' -f 3)

    parts=$(echo "${version}" | awk -F. '{print NF-1}')
    if [ "${parts}" -eq 2 ]; then
        if [ "${patch}" -eq 0 ]; then
            minor=$((minor-1))
        fi
        echo "${major}.${minor}"
        return
    fi

    minor=$((minor-1))
    echo "${major}.${minor}"
}

prior_stable_version() {
    local version major minor patch parts

    version="${1}"

    # shellcheck disable=SC2116
    version=$(echo "${version%-*}")

    major=$(echo "${version}" | cut -d '.' -f 1)
    minor=$(echo "${version}" | cut -d '.' -f 2)

    minor=$((minor-1))
    echo "${major}.${minor}"
}

series_version() {
    local version series arch

    version="${1}"
    series="${2}"

    arch=$(echo "${version}" | sed 's:.*-::')

    # shellcheck disable=SC2116
    version=$(echo "${version%-*}")
    # shellcheck disable=SC2116
    version=$(echo "${version%-*}")

    echo "${version}-${series}-${arch}"
}

add_upgrade_tools() {
    local version jujud_path

    version=${1}

    jujud_path=$(which jujud)
    cp "${jujud_path}" "${TEST_DIR}"
    cd "${TEST_DIR}" || exit

    tar -zcvf "juju-${version}.tgz" jujud >/dev/null
    cd "${CURRENT_DIR}/.." || exit

    mkdir -p "./tests/suites/upgrade/streams/tools/released/"
    mv "${TEST_DIR}/juju-${version}.tgz" "./tests/suites/upgrade/streams/tools/released"
}

remove_upgrade_tools() {
    cd "${CURRENT_DIR}/.." || exit

    echo "==> Removing tools"
    rm -rf ./tests/suites/upgrade/streams/tools/released || true
    echo "==> Removed tools"
}

remove_upgrade_metadata() {
    cd "${CURRENT_DIR}/.." || exit

    echo "==> Removing metadata"
    rm -rf ./tests/suites/upgrade/streams/tools/streams || true
    echo "==> Removed metadata"
}
