echo $'Cleaning previously created files'
rm -r img/*.jpg
rm out.mp4
echo $'Generating images'
go run main.go 
echo $'Generating mp4'
ffmpeg -loglevel panic -i 'img/%06d.jpg' -vcodec h264 out.mp4