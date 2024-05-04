#!/usr/bin/env bash


./netrunner-alt-gen netspace hedge fund
                    
flavor_sure_gamble='"I would suggest variable text box size.<BR>Wasted space on stuff like Sure Gambles is something that'"'"'s always bothered me for alt arts."'

./netrunner-alt-gen netspace sure gamble \
                    --make-back \
                    --text-box-height 26 \
                    --flavor "${flavor_sure_gamble}" \
                    --flavor-attribution "Suipe" \
                    --base-color 472F49 \
                    --min-walkers 5000 \
                    --grid-percent 0 \
                    --color-bg 472F49 \
                    --frame-color-influence-bg 3f3f3f \
                    --walker-color-1 AA9290 \
                    --walker-color-2 784E5C \
                    --walker-color-3 FBF4D7 \
                    --walker-color-4 BFB68F


# skip flavor on FFG cards as a courtesy
./netrunner-alt-gen netspace diversion of funds \
                    --skip-flavor

# get a cropped preview
(
    cd output
    mkdir -p cropped
    for f in *
    do
        echo "converting $f"
        convert $f -crop 2976x4152+144+149 cropped/$f
    done
)
