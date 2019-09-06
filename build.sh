#!/usr/bin/env sh
package_name="https-proxy"
platforms="linux/amd64"
rm -rf build
mkdir build
VER="v1.2.0"
for platform in $platforms
do
    var=$(echo $platform | awk -F"/" '{print $1,$2}')   
    set -- $var
    # platform_split=(${platform//\// })
    GOOS=$1
    GOARCH=$2
    output_name=$package_name'-'$VER'-'$GOOS'-'$GOARCH
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build --ldflags "-extldflags -static" -tags netgo -o $output_name .
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution'
        exit 1
    fi

    tar cvzf ./build/$output_name.tgz $output_name
    rm -rf $output_name
done
