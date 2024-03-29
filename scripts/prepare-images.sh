#!/usr/bin/env bash

# Copyright 2020 The Dataomnis Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

CurrentDIR=$(cd "$(dirname "$0")" || exit;pwd)
ImagesDirDefault=${CurrentDIR}/dataomnis-images
action=""
save=""
registryurl=""

func() {
    echo "Usage:"
    echo
    echo "  $0 -a pull/push [-d IMAGES-DIR] [-r PRIVATE-REGISTRY] [-s]"
    echo
    echo "Description:"
    echo "  -d IMAGES-DIR          : the dir of files (tar.gz) which generated by \`docker save\`. default: ${ImagesDirDefault}"
    echo "  -r PRIVATE-REGISTRY    : target private registry:port for push-image."
    echo "  -a ACTION              : [pull/push] pull the images in the IMAGES-LIST, or push image to PRIVATE REGISTRY"
    echo "  -s                     : if save images in the IMAGES-LIST as a tar.gz file, or push image from image tar.gz"
    echo "  -h                     : usage message"
    echo
    echo "Examples:"
    echo
    echo "# pull and save the images"
    echo "$0 -a pull [-d ./] [-s]"
    echo
    echo "# Push the images to private docker registry."
    echo "$0 -a push -r dockerhub.testing.io [-d ./] [-s]"
    exit
}

while getopts 'sa:r:d:h' OPT; do
    case $OPT in
        d) ImagesDir="$OPTARG";;
        r) Registry="$OPTARG";;
        a) action="$OPTARG";;
        s) save="true";;
        h) func;;
        ?) func;;
        *) func;;
    esac
done

if [ -z "${ImagesDir}" ]; then
    ImagesDir=${ImagesDirDefault}
fi

if [ -n "${Registry}" ]; then
   registryurl=${Registry}
fi

ImagesList="./images-list.txt"

if [[ ${action} == "pull" ]] && [[ -n "${ImagesList}" ]]; then
    if [ ! -d ${ImagesDir} ]; then
       mkdir -p ${ImagesDir}
    fi
    ImagesListLen=$(cat ${ImagesList} | wc -l)
    name=""
    images=""
    index=0
    for image in $(<${ImagesList}); do
        if [[ ${image} =~ ^\#\#.* ]]; then
           if [[ -n ${images} && ${save} == "true" ]]; then
              echo ""
              echo "Save images: "${name}" to "${ImagesDir}"/"${name}".tar.gz  <<<"
              docker save ${images} | gzip -c > ${ImagesDir}"/"${name}.tar.gz
              echo ""
           fi
           images=""
           name=$(echo "${image}" | sed 's/#//g' | sed -e 's/[[:space:]]//g')
           ((index++))
           continue
        fi

        echo "Pull image: ${image} .."
        docker pull "${image}"
        images=${images}" "${image}

        if [[ ${index} -eq ${ImagesListLen}-1 ]]; then
           if [[ -n ${images} ]]; then
              docker save ${images} | gzip -c > ${ImagesDir}"/"${name}.tar.gz
           fi
        fi
        ((index++))
    done
elif [[ ${action} == "push" ]] && [ -n "${ImagesList}" ]; then
    if [[ "${registryurl}X" == "X" ]]; then
      # shellcheck disable=SC2028
      echo "The registry must be set by -r when action is PUSH!\n"
      func;
      exit 1
    fi
    # shellcheck disable=SC2045
    if [[ ${save} == "true" ]]; then
        for image in $(ls ${ImagesDir}/*.tar.gz); do
            echo "Load images: "${image}"  <<<"
            docker load  < $image
        done
    fi

    for image in $(<${ImagesList}); do
        if [[ ${image} =~ ^\#\#.* ]]; then
           continue
        fi

        imgArray=(${image//\// })
        if [ ${#imgArray[@]} == 3 ]; then
          if [ "${imgArray[0]}X" == "${registryurl}X" ]; then
           continue
          fi
          imageurl="${registryurl}/${imgArray[1]}/${imgArray[2]}"
        elif [ ${#imgArray[@]} == 2 ]; then
          imageurl="${registryurl}/${image}"
        else
          imageurl="${registryurl}/library/${image}"
        fi

        ## push image
        echo $imageurl
        docker tag $image $imageurl
        docker push $imageurl
    done
else
  func;
fi