@echo off

echo Cleaning previously created movie
del out.mp4

set DEFAULT_SPEED=0.75
set SPEED=%DEFAULT_SPEED%

if NOT [%1] EQU [] set SPEED=%1

echo Generating images
monroe.exe

echo Generating whole film
ffmpeg.exe -loglevel panic -r 0.75 -i "temp\f1\%%06d.jpg" -r 25 -pix_fmt yuv420p temp/f1.mp4
ffmpeg.exe -loglevel panic -r 2.0 -i "temp\f2\%%06d.jpg" -r 25 -pix_fmt yuv420p temp/f2.mp4
ffmpeg.exe -loglevel panic -i "temp\f3\%%06d.jpg" -r 25 -pix_fmt yuv420p temp/f3.mp4
ffmpeg.exe -loglevel panic -r 3.5 -i "temp\col\%%06d.jpg" -r 25 -pix_fmt yuv420p temp/col.mp4
ffmpeg.exe -loglevel panic -f concat -i temp/mylist.txt -r 25 -pix_fmt yuv420p out.mp4

echo Cleaning temporary files
RMDIR temp\ /S /Q