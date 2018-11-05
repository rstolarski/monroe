echo $'Cleaning previously created movie'
rm out.mp4
echo $'Generating images'
go run main.go 
echo $'Generating mp4'
ffmpeg -loglevel panic -r 2 -i 'img/%06d.jpg' -r 10 -pix_fmt yuv420p out.mp4
echo $'Cleaning temporary files'
rm -r img/*.jpg