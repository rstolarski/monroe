echo $'Cleaning previously created movie'
rm out.mp4


echo $'Generating images'
./monroe

echo $'Generating whole film'
ffmpeg -loglevel panic -r 0.75 -i "temp/f1/%06d.jpg" -r 25 -pix_fmt yuv420p temp/f1.mp4
ffmpeg -loglevel panic -r 2.0 -i "temp/f2/%06d.jpg" -r 25 -pix_fmt yuv420p temp/f2.mp4
ffmpeg -loglevel panic -i "temp/f3/%06d.jpg" -r 25 -pix_fmt yuv420p temp/f3.mp4
ffmpeg -loglevel panic -r 3.5 -i "temp/col/%06d.jpg" -r 25 -pix_fmt yuv420p temp/col.mp4
ffmpeg -loglevel panic -f concat -i temp/mylist.txt -r 25 -pix_fmt yuv420p out.mp4

echo $'Cleaning temporary files'
rm -r temp/*