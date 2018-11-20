echo $'Cleaning previously created movie'
rm out.mp4

DEFAULT_SPEED=0.75
SPEED=${1:-$DEFAULT_SPEED}
echo $SPEED  

echo $'Generating images'
go run main.go 

echo $'Generating whole film'
./ffmpeg -loglevel panic -r $SPEED  -i 'temp/f1/%06d.jpg' -r 25 -pix_fmt yuv420p temp/f1.mp4
./ffmpeg -loglevel panic -r $SPEED  -i 'temp/f2/%06d.jpg' -r 25 -pix_fmt yuv420p temp/f2.mp4
./ffmpeg -loglevel panic -i 'temp/f3/%06d.jpg' -r 25 -pix_fmt yuv420p temp/f3.mp4
./ffmpeg -loglevel panic -r $SPEED -i 'temp/col/%06d.jpg' -r 25 -pix_fmt yuv420p temp/col.mp4
./ffmpeg -loglevel panic -f concat -i temp/mylist.txt -r 25 -pix_fmt yuv420p out.mp4

echo $'Cleaning temporary files'
rm -r temp/*