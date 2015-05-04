#!/bin/bash
source build.sh
cpack --config ../../built.linux_x86-64/heka/CPackConfig.cmake
rm -f ../../built.linux_x86-64/heka/heka/src/github.com/carlanton/go-dockerclient/testing/data/symlink
