The following instructions specify how to add(re-base) a new version of heka

1. Download the desired heka SOURCE CODE tarball from https://github.com/mozilla-services/heka/releases
2. Extract the tarball to a temporary directory using tar -xzvf
3. Rename the extracted heka directory to just "heka"
3. Copy the metadata.txt, ps_build.sh, license-manifext.txt, _build.py, build-readme.txt, build.sh and env.sh files from the current heka sandbox to the temporary directory
4. Copy the plugins/http/http_filter.go file to the temporary directory in the same location
5. Copy the .if_top directory to the temporary directory(the .if_top directory might be hidden try ll -a)
6. Ensure the temporary directory has your normal user ownership using chown if necessary
7. Copy the entire temporary directory to a NEW heka sandbox and test before committing
8. The file heka/plugins/http/http_filter.go, must be manually copied from the perfect search httpFilter github repo in order to get the latest changes

8. When rebasing ensure that the cmake/externals.cmake files git tags are either 
    a: unchanged
    b: are changed to a revision that you are sure you want.

9. The licenses MUST for the following repos must be checked against what is in the license-manifest file. The Licenses MUST match or else verified with ip-verifcation@perfectsearchcorp.com
    github.com/streadway/amqp
    github.com/bitly/go-simplejson
    github.com/mozilla-services/heka
    github.com/mozilla-services/heka-docs
    github.com/mozilla-services/heka-mozsvc-plugins
    github.com/thoj/go-ircevent
    github.com/crowdmob/goamz
    github.com/feyeleanor/slices
    github.com/feyeleanor/sets
    github.com/feyeleanor/raw
    github.com/bbangert/toml
    github.com/crankycoder/g2s
    github.com/crankycoder/xmlpath
    github.com/getsentry/raven-go
    github.com/rafrombrc/gospec
    github.com/rafrombrc/go-notify
    github.com/rafrombrc/go-dockerclient
    github.com/rafrombrc/gomock
    github.com/rafrombrc/whisper-go
    code.google.com/p/go-uuid
    code.google.com/p/gogoprotobuf
